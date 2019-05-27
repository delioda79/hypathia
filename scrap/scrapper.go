package scrap

import (
	"encoding/json"
	"fmt"
	"githubscrapper"
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/v25/github"
	"io/ioutil"
	"log"
	"net/http"
)

type docFileSpec struct {
	Name string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Scrapper struct {
	httpCLient *http.Client
	ghclient *github.Client
	account string
}

func New(token, account string) Scrapper {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return Scrapper{
		httpCLient:tc,
		ghclient: client,
		account: account,
	}
}

func (sc *Scrapper) Scrap() []githubscrapper.DocDef {
	result := []githubscrapper.DocDef{}
	ctx := context.Background()
	reps, _, err := sc.ghclient.Repositories.List(ctx, "delioda79", nil)
	if err != nil {
		fmt.Print(err)
	}


	if len(reps) <= 0 {
		return result
	}

	for _, rp := range reps {
		fmt.Println("checking: ", rp.GetName())
		rsp, err := sc.httpCLient.Get(fmt.Sprintf("%s/contents/docs/", rp.GetURL()))
		if err != nil {
			log.Println(rp.GetName(), err)
			continue
		}
		if rsp.StatusCode != 200 {
			log.Println("Status: ", rsp.StatusCode)
			continue
		}
		bts, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println("Impossible to unmarshal")
			continue
		}
		specs := []docFileSpec{}
		err = json.Unmarshal(bts,&specs)
		if err != nil {
			log.Println(err)
			continue
		}
		for _,doc := range specs {
			fmt.Println(doc.Name)
			if doc.Name == "swagger.json" {
				result = append(result, githubscrapper.DocDef{Type: githubscrapper.Swagger, URL: doc.DownloadURL})
			}

			if doc.Name == "async.json" {
				result = append(result, githubscrapper.DocDef{Type: githubscrapper.Async, URL: doc.DownloadURL})
			}
		}
	}

	return result
}