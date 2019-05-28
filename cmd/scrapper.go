package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"githubscrapper/scrap/github"
	"githubscrapper/serve"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
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
		log.Println("No port set, defaulting to 9024\n")
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("Wrong port value: %q is not an integer.\n", port)
	}

	scrapper := github.New(ghtoken, ghaccount)

	docs := scrapper.Scrap()

	for _, d := range docs {
		fmt.Println(d)
	}

	hdl := &serve.Handler{}
	hdl.Update(docs)

	r := mux.NewRouter()

	r.HandleFunc("/api-docs", hdl.ApiList)
	r.HandleFunc("/doc/{repoName}/{type}", hdl.ApiRender)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}
