package main

import (
	"fmt"
	"githubscrapper/scrap"
	"githubscrapper/serve"
	"log"
	"net/http"
)

func main() {

	scrapper := scrap.New( "10a3fb18b1caf9b45b26e5f582b2f001c09fac47", "delioda79")

	docs := scrapper.Scrap()

	for _, d := range docs {
		fmt.Println(d)
	}

	hdl := &serve.Handler{}
	hdl.Update(docs)

	http.Handle("/foo", hdl)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
