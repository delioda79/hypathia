package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v25/github"
	"github.com/taxibeat/hypatia/scrape"
	"golang.org/x/oauth2"
	"net/http"
	"sync"
)

type docFileSpec struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

type Scraper struct {
	httpClient   *http.Client
	ghClient     *github.Client
	organization string
	branch       string
	filter       scrape.Filter
}

func New(token, organization, branch string, filter scrape.Filter) Scraper {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return Scraper{
		httpClient:   tc,
		ghClient:     client,
		organization: organization,
		branch:       branch,
		filter:       filter,
	}
}

//Scrape fires up workers for each repository accumulate results with reporter
//Transforms them into []DocDef and returns them
func (sc *Scraper) Scrape() []scrape.DocDef {

	ctx := context.Background()

	resChan := make(chan searchResponse, 10)

	//Start reporter accumulator
	var wgReporter sync.WaitGroup
	wgReporter.Add(1)

	var accumulator []searchResponse
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

//DocDefReportTransform transforms internal []scrapeResponse to external []DocDef
func docDefReportTransform(scrapeRes []searchResponse) []scrape.DocDef {
	var docDefs []scrape.DocDef
	for _, res := range scrapeRes {
		docDefs = append(docDefs, res.out...)
	}
	return docDefs
}

//ProcessRepository fires up a go routine that scrape a specific repository
func (sc *Scraper) processRepository(rp github.Repository, resChan chan<- searchResponse, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		resChan <- sc.scrapeRepo(rp)
	}()
}

//Reporter fires up a go routine that accumulates scrapeResponse results
func (sc *Scraper) reporter(resChan <-chan searchResponse, accumulator *[]searchResponse, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for res := range resChan {
			*accumulator = append(*accumulator, res)
		}
	}()
}
