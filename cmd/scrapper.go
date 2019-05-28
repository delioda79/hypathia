package main

import (
	"fmt"
	"githubscrapper/scrap/github"
	"githubscrapper/serve"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	ghtoken := os.Getenv("GITHUB_TOKEN")
	ghaccount := os.Getenv("GITHUB_ACCOUNT")

	scrapper := github.New( ghtoken, ghaccount)

	docs := scrapper.Scrap()

	for _, d := range docs {
		fmt.Println(d)
	}

	hdl := &serve.Handler{}
	hdl.Update(docs)

	http.Handle("/foo", hdl)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
