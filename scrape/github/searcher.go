package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v25/github"
	"github.com/taxibeat/hypatia/scrape"
	"io/ioutil"
)

type searchResponse struct {
	out    []scrape.DocDef
	errOut []error
}

const (
	docBasePath = "/contents/docs"
	syncFile    = "swagger.json"
	asyncFile   = "async.json"
)

//ScrapeRepo searches in rp github.Repository for any documentation file under the docBasePath path
func (sc *Scraper) scrapeRepo(rp github.Repository) searchResponse {
	fmt.Println("checking: ", rp.GetName())

	result := make([]scrape.DocDef, 0)
	rsp, err := sc.httpClient.Get(fmt.Sprintf("%s"+docBasePath+"?ref=%s", rp.GetURL(), sc.branch))
	if err != nil {
		return searchResponse{result, []error{err}}
	}
	if rsp.StatusCode != 200 {
		return searchResponse{result, []error{fmt.Errorf("Status: %d", rsp.StatusCode)}}
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return searchResponse{result, []error{errors.New("impossible to unmarshal")}}
	}
	var specs []docFileSpec
	err = json.Unmarshal(bts, &specs)
	if err != nil {
		return searchResponse{result, []error{err}}
	}
	fmt.Println("SPECS", specs)
	var errs []error
	for _, doc := range specs {
		def, err := sc.retrieveDocumentation(rp.GetName(), doc)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if def != nil {
			result = append(result, *def)
		}
	}

	return searchResponse{result, nil}
}

//RetrieveDocumentation retrieves and returns files of supported types
func (sc *Scraper) retrieveDocumentation(sourceRepo string, doc docFileSpec) (*scrape.DocDef, error) {
	result := scrape.DocDef{}

	switch doc.Name {
	case syncFile:
		result.Type = scrape.Swagger
	case asyncFile:
		result.Type = scrape.Async
	default:
		return nil, fmt.Errorf("Unsupported type: %s", doc.Name)
	}
	result.URL = doc.DownloadURL
	result.RepoName = sourceRepo
	rsp, err := sc.httpClient.Get(doc.DownloadURL)
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
