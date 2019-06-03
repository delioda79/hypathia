package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"githubscrapper/scrap"
	"githubscrapper/scrap/github"
	"githubscrapper/serve"
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
	ghaccount := os.Getenv("GITHUB_ACCOUNT")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9024"
		log.Println("No port set, defaulting to 9024")
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("Wrong port value: %q is not an integer.\n", port)
	}

	branch := os.Getenv("GITHUB_BRANCH")
	if branch == "" {
		branch = "master"
		log.Println("No branch set, defaulting to master")
	}

	scrapper := github.New(ghtoken, ghaccount, branch)
	hdl := &serve.Handler{}

	scrapRepos(&scrapper, hdl)

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

func scrapRepos(scrapper scrap.Scrapper, handler *serve.Handler) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		handler.Update(scrapper.Scrap())
		for range ticker.C {
			handler.Update(scrapper.Scrap())
		}
	}()
}
