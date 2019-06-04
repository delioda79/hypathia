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
	"strconv"
	"strings"
	"time"
)

const (
	version = "dev"
	name    = "hypatia"
)

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

func main() {

	ghtoken := os.Getenv("GITHUB_TOKEN")
	ghorganization := os.Getenv("GITHUB_ORGANIZATION")

	ghbranch := os.Getenv("GITHUB_BRANCH")
	if ghbranch == "" {
		ghbranch = "master"
		log.Warn("No branch set, defaulting to master")
	}
	var ghtags []string
	if (os.Getenv("GITHUB_TAGS")) != "" {
		ghtags = strings.Split(os.Getenv("GITHUB_TAGS"), ",")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9024"
		log.Warn("No port set, defaulting to 9024\n")
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Wrong port value: %q is not an integer.\n", port)
	}

	refreshTime := time.Minute
	rt := os.Getenv("REFRESH_TIME")
	if rt != "" {
		parsed, err := time.ParseDuration(rt)
		if err != nil {
			log.Fatalf("env %s is not a duration: %v", rt, err)
		}
		refreshTime = parsed
	}

	scraper := github.New(ghtoken, ghorganization, ghbranch, ghtags)

	hdl := &serve.Handler{}

	scrapRepos(&scraper, hdl, refreshTime)

	if err := run(hdl); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(hdl *serve.Handler) error {

	r := phttp.NewRouteRaw("/", "GET", hdl.ApiList, false)
	r1 := phttp.NewRouteRaw("/doc/:repoName/:type", "GET", hdl.ApiRender, false)
	r2 := phttp.NewRouteRaw("/spec/:repoName/:type", "GET", hdl.SpecRender, false)

	srv, err := patron.New(
		name,
		version,
		patron.Routes([]phttp.Route{r, r1, r2}),
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
