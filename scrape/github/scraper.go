package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v25/github"
	"hypatia/scrape"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"sync"
)

type docFileSpec struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Scraper struct {
	httpCLient   *http.Client
	ghclient     *github.Client
	organization string
	branch       string
}

type scrapResponse struct {
	out    []scrape.DocDef
	errOut []error
}

func New(token, organization, branch string) Scraper {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return Scraper{
		httpCLient:   tc,
		ghclient:     client,
		organization: organization,
		branch:       branch,
	}
}

//Scrap fires up workers for each repository and waits for DocDef results
func (sc *Scraper) Scrape() []scrape.DocDef {

	ctx := context.Background()

	resChan := make(chan scrapResponse, 10)

	//Start reporter accumulator
	var wgReporter sync.WaitGroup
	wgReporter.Add(1)

	var accumulator []scrapResponse
	go sc.reporter(resChan, &accumulator, &wgReporter)

	var wgWorkers sync.WaitGroup
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	//GET on github's account with pagination
	for {
		reps, res, err := sc.ghclient.Repositories.ListByOrg(ctx, sc.organization, opt)

		fmt.Println(res)
		if err != nil {
			fmt.Print(err)
		}

		//Start repository workers
		for _, rp := range reps {
			wgWorkers.Add(1)
			sc.processRepository(*rp, resChan, &wgWorkers)
		}

		if res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
	}

	//Wait on workers to finish
	wgWorkers.Wait()

	//Close channels, no more data will be passed
	close(resChan)

	//Wait on reporter to finish
	wgReporter.Wait()

	return docDefReport(accumulator)
}

func docDefReport(scrapRes []scrapResponse) []scrape.DocDef {
	var docDefs []scrape.DocDef
	for _, res := range scrapRes {
		docDefs = append(docDefs, res.out...)
	}
	return docDefs
}

func (sc *Scraper) processRepository(rp github.Repository, resChan chan<- scrapResponse, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		resChan <- sc.ScrapRepo(rp)
	}()
}

func (sc *Scraper) reporter(resChan <-chan scrapResponse, accumulator *[]scrapResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for res := range resChan {
		*accumulator = append(*accumulator, res)
	}
}

func (sc *Scraper) Define(sourceRepo string, doc docFileSpec) (*scrape.DocDef, error) {
	result := scrape.DocDef{}

	switch doc.Name {
	case "swagger.json":
		result.Type = scrape.Swagger
	case "async.json":
		result.Type = scrape.Async
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type: %s", doc.Name))
	}
	result.URL = doc.DownloadURL
	result.RepoName = sourceRepo
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

func (sc *Scraper) ScrapRepo(rp github.Repository) scrapResponse {
	fmt.Println("checking: ", rp.GetName())

	result := make([]scrape.DocDef, 0)
	rsp, err := sc.httpCLient.Get(fmt.Sprintf("%s/contents/docs?ref=%s", rp.GetURL(), sc.branch))
	if err != nil {
		return scrapResponse{result, []error{err}}
	}
	if rsp.StatusCode != 200 {
		return scrapResponse{result, []error{errors.New(fmt.Sprintf("Status: %d", rsp.StatusCode))}}
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return scrapResponse{result, []error{errors.New("impossible to unmarshal")}}
	}
	var specs []docFileSpec
	err = json.Unmarshal(bts, &specs)
	if err != nil {
		return scrapResponse{result, []error{err}}
	}
	fmt.Println("SPECS", specs)
	errs := []error{}
	for _, doc := range specs {
		def, err := sc.Define(rp.GetName(), doc)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if def != nil {
			result = append(result, *def)
		}
	}

	return scrapResponse{result, nil}
}
