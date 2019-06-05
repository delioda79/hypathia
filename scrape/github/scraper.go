package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v25/github"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/scrape/github/filter"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"sync"
)

type docFileSpec struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Filter interface {
	Apply([]*github.Repository) []*github.Repository
}

type Scraper struct {
	httpClient   *http.Client
	ghClient     *github.Client
	organization string
	branch       string
	filter       Filter
}

type scrapeResponse struct {
	out    []scrape.DocDef
	errOut []error
}

const (
	docBasePath = "/contents/docs"
	syncFile    = "swagger.json"
	asyncFile   = "async.json"
)

func New(token, organization, branch string, tags []string) Scraper {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	fil := filter.New(tags)

	return Scraper{
		httpClient:   tc,
		ghClient:     client,
		organization: organization,
		branch:       branch,
		filter:       fil,
	}
}

//Scrape fires up workers for each repository accumulate results with reporter
//Transforms them into []DocDef and returns them
func (sc *Scraper) Scrape() []scrape.DocDef {

	ctx := context.Background()

	resChan := make(chan scrapeResponse, 10)

	//Start reporter accumulator
	var wgReporter sync.WaitGroup
	wgReporter.Add(1)

	var accumulator []scrapeResponse
	sc.reporter(resChan, &accumulator, &wgReporter)

	var wgWorkers sync.WaitGroup
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	//GET on github's account with pagination
	for {
		reps, res, err := sc.ghClient.Repositories.ListByOrg(ctx, sc.organization, opt)

		fmt.Println(res)
		if err != nil {
			fmt.Print(err)
		}

		reps = sc.filter.Apply(reps)

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

	return docDefReportTransform(accumulator)
}

//docDefReportTransform transforms internal []scrapeResponse to external []DocDef
func docDefReportTransform(scrapeRes []scrapeResponse) []scrape.DocDef {
	var docDefs []scrape.DocDef
	for _, res := range scrapeRes {
		docDefs = append(docDefs, res.out...)
	}
	return docDefs
}

//processRepository fires up a go routine that scrape a specific repository
func (sc *Scraper) processRepository(rp github.Repository, resChan chan<- scrapeResponse, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		resChan <- sc.scrapeRepo(rp)
	}()
}

//reporter fires up a go routine that accumulates scrapeResponse results
func (sc *Scraper) reporter(resChan <-chan scrapeResponse, accumulator *[]scrapeResponse, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for res := range resChan {
			*accumulator = append(*accumulator, res)
		}
	}()
}

//scrapeRepo searches in rp github.Repository for any documentation file under the docBasePath path
func (sc *Scraper) scrapeRepo(rp github.Repository) scrapeResponse {
	fmt.Println("checking: ", rp.GetName())

	result := make([]scrape.DocDef, 0)
	rsp, err := sc.httpClient.Get(fmt.Sprintf("%s"+docBasePath+"?ref=%s", rp.GetURL(), sc.branch))
	if err != nil {
		return scrapeResponse{result, []error{err}}
	}
	if rsp.StatusCode != 200 {
		return scrapeResponse{result, []error{fmt.Errorf("Status: %d", rsp.StatusCode)}}
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return scrapeResponse{result, []error{fmt.Errorf("impossible to unmarshal")}}
	}
	var specs []docFileSpec
	err = json.Unmarshal(bts, &specs)
	if err != nil {
		return scrapeResponse{result, []error{err}}
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

	return scrapeResponse{result, nil}
}

//retrieveDocumentation retrieves and returns files of supported types
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
