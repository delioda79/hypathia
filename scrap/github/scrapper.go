package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v25/github"
	"githubscrapper/scrap"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type docFileSpec struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Scrapper struct {
	httpCLient *http.Client
	ghclient   *github.Client
	account    string
}

func New(token, account string) Scrapper {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return Scrapper{
		httpCLient: tc,
		ghclient:   client,
		account:    account,
	}
}

func (sc *Scrapper) Scrap() []scrap.DocDef {
	result := []scrap.DocDef{}
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
		var specs []docFileSpec
		err = json.Unmarshal(bts, &specs)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, doc := range specs {
			def, err := sc.Define(rp.GetName(), doc)
			if err != nil {
				log.Println(err)
				continue
			}
			if def != nil {
				result = append(result, *def)
			}
		}
	}

	return result
}

func (sc *Scrapper) Define(sourceRepo string, doc docFileSpec) (*scrap.DocDef, error) {
	result := scrap.DocDef{}

	switch doc.Name {
	case "swagger.json":
		result.Type = scrap.Swagger
	case "async.json":
		result.Type = scrap.Async
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type: %s", doc.Name))
	}
	result.RepoName = sourceRepo
	result.URL = doc.DownloadURL
	rsp, err := sc.httpCLient.Get(doc.DownloadURL)
	if err != nil {
		return nil, err
	}
	definition, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	result.Definition = string(definition)

	return &result, nil
}
