package main

import (
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	phttp "github.com/beatlabs/patron/sync/http"
	"github.com/joho/godotenv"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/scrape/github"
	"github.com/taxibeat/hypatia/serve"
	"os"
	"strings"
	"time"
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

	ghtoken := mustGetEnv("GITHUB_TOKEN")
	ghorganization := mustGetEnv("GITHUB_ORGANIZATION")
	ghbranch := mustGetEnvWithDefault("GITHUB_BRANCH", "master")
	ghtags := mustGetEnvArray("GITHUB_TAGS")
	refreshTime := mustGetEnvDurationWithDefault("REFRESH_TIME", "1h")

	scraper := github.New(ghtoken, ghorganization, ghbranch, ghtags)

	hdl := &serve.Handler{}

	scrapRepos(&scraper, hdl, refreshTime)

	srv, err := patron.New(
		name,
		version,
		patron.Routes(routes(hdl)),
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

func scrapRepos(scraper scrape.Scraper, handler *serve.Handler, rt time.Duration) {
	ticker := time.NewTicker(rt)
	go func() {
		handler.Update(scraper.Scrape())
		fmt.Println("Updating")
		for range ticker.C {
			handler.Update(scraper.Scrape())
			fmt.Println("Updating")
		}
	}()
}

func routes(hdl *serve.Handler) []phttp.Route {
	return []phttp.Route{
		phttp.NewRouteRaw("/", "GET", hdl.ApiList, false),
		phttp.NewRouteRaw("/doc/:repoName/:type", "GET", hdl.ApiRender, false),
		phttp.NewRouteRaw("/spec/:repoName/:type", "GET", hdl.SpecRender, false),
	}
}

func mustGetEnvArray(key string) []string {
	v, ok := os.LookupEnv(key)
	if !ok {
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
