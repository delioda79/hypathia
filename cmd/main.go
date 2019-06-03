package main

import (
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	phttp "github.com/beatlabs/patron/sync/http"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"hypatia/serve"
	"hypatia/scrape/github"
	"hypatia/scrape"
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

	err = godotenv.Load("../../.env")
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
		//log.Println("No branch set, defaulting to master")
	}
	var ghtags []string
	if (os.Getenv("GITHUB_TAGS")) != "" {
		ghtags = strings.Split(os.Getenv("GITHUB_TAGS"), ",")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9024"
		//log.Println("No port set, defaulting to 9024\n")
	}

	if _, err := strconv.Atoi(port); err != nil {
		//log.Printf("Wrong port value: %q is not an integer.\n", port)
	}

	refreshTime := time.Minute
	rt := os.Getenv("REFRESH_TIME")
	if rt != "" {
		parsed, err := time.ParseDuration(rt)
		if err != nil {
			//log.Fatalf("env %s is not a duration: %v", rt, err)
		}
		refreshTime = parsed
	}


	scraper := github.New(ghtoken, ghorganization, ghbranch, ghtags)

	hdl := &serve.Handler{}

	scrapRepos(&scraper, hdl, refreshTime)

	//r := mux.NewRouter()
	//
	//r.HandleFunc("/", hdl.ApiList)
	//r.HandleFunc("/doc/{repoName}/{type}", hdl.ApiRender)
	//r.HandleFunc("/spec/{repoName}/{type}", hdl.SpecRender)
	//r.HandleFunc("/health", hdl.HealthStatus)
	//
	//log.Printf("Listening on port %s\n", port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))


	if err := run(hdl); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(hdl *serve.Handler) error {

	r := phttp.NewRouteRaw("/", "GET",  Index, false)
	//r1 := phttp.NewRouteRaw("/doc/:repoName/:type", "GET",  hdl.ApiRender, false)
	//r2 := phttp.NewRouteRaw("/spec/:repoName/:type", "GET",  hdl.SpecRender, false)

	srv, err := patron.New(
		name,
		version,
		patron.Routes([]phttp.Route{r}),
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

func Index(w http.ResponseWriter, req *http.Request) {
	f, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		fmt.Println("error reading file")
	}
	w.Write(f)
}

func scrapRepos(scraper scrape.Scraper, handler *serve.Handler, rt time.Duration) {
	return
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
