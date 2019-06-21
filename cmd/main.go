package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/taxibeat/hypatia/scrape/fs"
	"github.com/taxibeat/hypatia/scrape/github"
	"github.com/taxibeat/hypatia/scrape/github/filter"

	"github.com/taxibeat/hypatia/search"
	"github.com/taxibeat/hypatia/search/berserker"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	phttp "github.com/beatlabs/patron/sync/http"
	"github.com/joho/godotenv"
	"github.com/taxibeat/hypatia/html/api2html"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/serve"
)

const (
	version = "dev"
	name    = "hypatia"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func init() {
	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}
}

func run() error {

	brs, err := berserker.NewBerserker()
	if err != nil {
		log.Fatalf("Error creating the berserker", err)
	}
	brs.Run()

	hdl := &serve.Handler{Searcher: brs}

	mode := mustGetEnvWithDefault("MODE", "github")
	modeFl := flag.Bool("fs", false, "Running in fs mode")
	flag.Parse()
	if modeFl != nil && *modeFl {
		mode = "fs"
	}

	switch mode {
	case "github":
		ghtoken := mustGetEnv("GITHUB_TOKEN")
		ghorganization := mustGetEnv("GITHUB_ORGANIZATION")
		ghbranch := mustGetEnvWithDefault("GITHUB_BRANCH", "")
		ghtags := mustGetEnvArray("GITHUB_TAGS")
		refreshTime := mustGetEnvDurationWithDefault("REFRESH_TIME", "1h")

		filter := filter.New(ghtags)

		httpClient := github.NewHTTPClient(ghtoken)

		gitClient := github.NewGithubClient(httpClient, nil)

		scraper := github.New(ghorganization, ghbranch, httpClient, filter, gitClient)
		runScraping(scraper, api2html.Transformer{}, hdl, refreshTime, brs)
	case "fs":
		argsWithoutProg := flag.Args()
		if len(argsWithoutProg) < 1 {
			log.Fatal("Missing path parameter, please specify one after calling hypatia")
		}

		path := argsWithoutProg[0]
		scraper, err := fs.New(path)
		if err != nil {
			log.Fatalf("impossible to create the scraper", err)
		}

		syncRepos(scraper, api2html.Transformer{}, hdl, brs)
	default:
		log.Fatal("No valid mode selected, please use github or fs")
	}
	if mode == "github" {

	} else {

	}

	srv, err := patron.New(
		name,
		version,
		patron.Routes(routes(hdl)),
		patron.HealthCheck(hdl.HealthStatus),
	)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}

	return nil
}

func syncRepos(scraper *fs.Scraper, api2html api2html.Transformer, handler *serve.Handler, idxr search.AsyncIndexer) {
	scrapeRepos(scraper, api2html, handler, idxr)
	go func() {
		for {
			select {
			case <-scraper.Updates():
				scrapeRepos(scraper, api2html, handler, idxr)
			}
		}
	}()
}

func runScraping(scraper scrape.Scraper, api2html api2html.Transformer, handler *serve.Handler, rt time.Duration, idxr search.AsyncIndexer) {
	ticker := time.NewTicker(rt)
	go func() {
		scrapeRepos(scraper, api2html, handler, idxr)
		log.Infof("Updating")
		for range ticker.C {
			scrapeRepos(scraper, api2html, handler, idxr)
			log.Infof("Updating")
		}
	}()
}

func scrapeRepos(scraper scrape.Scraper, api2html api2html.Transformer, handler *serve.Handler, idxr search.AsyncIndexer) {
	repos := scraper.Scrape()
	asyncDHtmlPages := api2html.Apply(retrieveAsyncAPIs(repos))
	for _, r := range repos {
		idxr.Index(r)
	}
	handler.Update(repos, asyncDHtmlPages)
}

func retrieveAsyncAPIs(docDefs []scrape.DocDef) []api2html.ApiDef {
	var asyncDs []api2html.ApiDef
	for _, docDef := range docDefs {
		if docDef.Type == scrape.Async {
			asyncDs = append(asyncDs, api2html.NewApiDef(docDef.ID, docDef.Definition))
		}
	}
	return asyncDs
}

func routes(hdl *serve.Handler) []phttp.Route {
	return []phttp.Route{
		phttp.NewRouteRaw("/", "GET", hdl.APIList, false),
		phttp.NewRouteRaw("/", "POST", hdl.APISearch, false),
		phttp.NewRouteRaw("/doc/:repoID", "GET", hdl.ApiRender, false),
		phttp.NewRouteRaw("/spec/:repoID", "GET", hdl.SpecRender, false),
		phttp.NewRouteRaw("/static/*path", "GET", hdl.StaticFiles, false),
	}
}

func mustGetEnvArray(key string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return nil
	}
	return strings.Split(v, ",")
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Missing configuration %s", key)
	}
	return v
}

func mustGetEnvDurationWithDefault(key, def string) time.Duration {
	dur, err := time.ParseDuration(mustGetEnvWithDefault(key, def))
	if err != nil {
		log.Fatalf("env %s is not a duration: %v", key, err)
	}

	return dur
}

func mustGetEnvWithDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		if def == "" {
			log.Fatalf("Missing configuration %s", key)
		} else {
			return def
		}
	}
	return v
}
