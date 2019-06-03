package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"hypatia/scrape"
	"hypatia/scrape/github"
	"hypatia/serve"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	ghtoken := os.Getenv("GITHUB_TOKEN")
	ghorganization := os.Getenv("GITHUB_ORGANIZATION")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9024"
		log.Println("No port set, defaulting to 9024\n")
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("Wrong port value: %q is not an integer.\n", port)
	}

	branch := os.Getenv("GITHUB_BRANCH")
	if branch == "" {
		branch = "master"
		log.Println("No branch set, defaulting to master")
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

	scraper := github.New(ghtoken, ghorganization, branch)

	hdl := &serve.Handler{}

	scrapRepos(&scraper, hdl, refreshTime)

	r := mux.NewRouter()

	staticFolder := "/static/"

	r.HandleFunc("/", hdl.ApiList)
	r.HandleFunc("/doc/{repoName}/{type}", hdl.ApiRender)
	r.HandleFunc("/spec/{repoName}/{type}", hdl.SpecRender)
	r.HandleFunc("/health", hdl.HealthStatus)
	r.PathPrefix(staticFolder).Handler(http.StripPrefix(staticFolder, http.FileServer(http.Dir("../public"))))

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

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
