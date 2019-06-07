package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v25/github"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/hypatia/scrape"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type mockFilter struct {
}

func (fl *mockFilter) Apply(repositories []*github.Repository) []*github.Repository {
	return repositories
}

type mockGitClient struct {
}

func (mgc *mockGitClient) ListByOrg(context.Context, string, *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return []*github.Repository{{}}, &github.Response{NextPage: 0}, nil
}

func TestScraper_New(t *testing.T) {
	scraper := New("o", "b", &http.Client{}, &mockFilter{}, GitClient{})

	assert.IsType(t, Scraper{}, scraper)
	assert.Equal(t, "o", scraper.organization)
	assert.Equal(t, "b", scraper.branch)
	assert.NotNil(t, scraper.httpClient)
	assert.NotNil(t, scraper.gitHubClient)
	assert.NotNil(t, scraper.filter)
}

func TestNewGithubClient(t *testing.T) {
	u, _ := url.Parse("b")
	gitClient := NewGithubClient(&http.Client{}, u)
	assert.IsType(t, GitClient{}, gitClient)
	assert.NotNil(t, gitClient.Repositories)
}

func TestNewHTTPClient(t *testing.T) {
	httpClient := NewHTTPClient("t")
	assert.NotNil(t, httpClient)
	assert.IsType(t, http.Client{}, *httpClient)
}

type mockServerResult int

const (
	FailGetOperation mockServerResult = iota
	FailUnmarshal
	Success
)

func setupMockServer(msResult mockServerResult) (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	mux.Handle(docBasePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch msResult {
		case FailUnmarshal:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			break
		case FailGetOperation:
			w.WriteHeader(http.StatusBadRequest)
			break
		default:
			w.WriteHeader(http.StatusOK)
			docDef, _ := json.Marshal([]scrape.DocDef{{}})
			w.Write(docDef)
			break
		}
	}))

	server := httptest.NewServer(apiHandler)
	return server, mux
}

func TestRetrieveDocumentation(t *testing.T) {

	type retrieveDocumentationTest struct {
		sourceRepo string
		fileType   string
		invalidURL bool
		sBehavior  mockServerResult
		result     *scrape.DocDef
		err        error
	}
	var retrieveDocumentationTests = []retrieveDocumentationTest{
		{
			sourceRepo: "a",
			fileType:   syncFile,
			sBehavior:  Success,
			invalidURL: false,
			result:     &scrape.DocDef{Type: scrape.Swagger, RepoName: "a", Definition: "[{\"Type\":0,\"RepoName\":\"\",\"URL\":\"\",\"Definition\":\"\"}]"},
			err:        nil,
		},
		{
			sourceRepo: "a",
			fileType:   asyncFile,
			sBehavior:  Success,
			invalidURL: false,
			result:     &scrape.DocDef{Type: scrape.Async, RepoName: "a", Definition: "[{\"Type\":0,\"RepoName\":\"\",\"URL\":\"\",\"Definition\":\"\"}]"},
			err:        nil,
		},
		{
			sourceRepo: "a",
			fileType:   "someOtherType",
			sBehavior:  Success,
			invalidURL: false,
			result:     nil,
			err:        fmt.Errorf("unsupported type: someOtherType"),
		},
		{
			sourceRepo: "a",
			fileType:   syncFile,
			sBehavior:  FailGetOperation,
			invalidURL: false,
			result:     nil,
			err:        fmt.Errorf("status: " + strconv.Itoa(http.StatusBadRequest)),
		},
		{
			sourceRepo: "a",
			fileType:   syncFile,
			sBehavior:  Success,
			invalidURL: true,
			result:     nil,
			err:        fmt.Errorf("Get : unsupported protocol scheme \"\""),
		},
	}

	for _, test := range retrieveDocumentationTests {

		server, _ := setupMockServer(test.sBehavior)

		scraper := Scraper{httpClient: server.Client()}

		var docDef *scrape.DocDef
		var err error
		if test.invalidURL {
			docDef, err = scraper.retrieveDocumentation(test.sourceRepo, docFileSpec{Name: test.fileType})
		} else {
			docDef, err = scraper.retrieveDocumentation(test.sourceRepo, docFileSpec{Name: test.fileType, DownloadURL: server.URL + docBasePath})
		}

		if test.err != nil {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.err, err)
		}

		//We don't do object equal assertion here because we don't know docDef.URL beforehand
		if test.result != nil {
			assert.Equal(t, test.result.RepoName, docDef.RepoName)
			assert.Equal(t, test.result.Type, docDef.Type)
			assert.Equal(t, test.result.Definition, docDef.Definition)
		} else {
			assert.Equal(t, test.result, docDef)
		}
		server.Close()
	}
}

func mockRetrieveDocumentations(repoName string, doc docFileSpec) (*scrape.DocDef, error) {
	return nil, nil
}

func mockRetrieveDocumentationsSucc(repoName string, doc docFileSpec) (*scrape.DocDef, error) {
	docDef := scrape.DocDef{}
	return &docDef, nil
}

func mockRetrieveDocumentationsFail(repoName string, doc docFileSpec) (*scrape.DocDef, error) {
	return nil, fmt.Errorf("error")
}

func TestScrapeRepo(t *testing.T) {

	type scrapeRepoTest struct {
		rp          github.Repository
		retrieveDoc retrieveDocumentation
		invalidURL  bool
		sBehavior   mockServerResult
		expected    scrapeResponse
	}

	var scrapeRepoTests = []scrapeRepoTest{
		{
			rp:          github.Repository{},
			retrieveDoc: mockRetrieveDocumentations,
			invalidURL:  true,
			sBehavior:   Success,
			expected:    scrapeResponse{[]scrape.DocDef{}, []error{error(fmt.Errorf("Get " + docBasePath + "?ref=: unsupported protocol scheme \"\""))}},
		},
		{
			rp:          github.Repository{},
			retrieveDoc: mockRetrieveDocumentations,
			invalidURL:  false,
			sBehavior:   FailGetOperation,
			expected:    scrapeResponse{[]scrape.DocDef{}, []error{fmt.Errorf("status: " + strconv.Itoa(http.StatusBadRequest))}},
		},
		{
			rp:          github.Repository{},
			retrieveDoc: mockRetrieveDocumentations,
			invalidURL:  false,
			sBehavior:   FailUnmarshal,
			expected:    scrapeResponse{[]scrape.DocDef{}, []error{error(fmt.Errorf("json: cannot unmarshal object into Go value of type []github.docFileSpec"))}},
		},
		{
			rp:          github.Repository{},
			retrieveDoc: mockRetrieveDocumentationsSucc,
			invalidURL:  false,
			sBehavior:   Success,
			expected:    scrapeResponse{[]scrape.DocDef{{}}, nil},
		},
		{
			rp:          github.Repository{},
			retrieveDoc: mockRetrieveDocumentationsFail,
			invalidURL:  false,
			sBehavior:   Success,
			expected:    scrapeResponse{[]scrape.DocDef{}, []error{error(fmt.Errorf("error"))}},
		},
	}

	for _, test := range scrapeRepoTests {
		server, _ := setupMockServer(test.sBehavior)

		scraper := Scraper{httpClient: server.Client()}

		var actual scrapeResponse
		if test.invalidURL == true {
			actual = scraper.scrapeRepo(test.rp, test.retrieveDoc)
		} else {
			testURL := server.URL + docBasePath + "?ref="
			test.rp.URL = &testURL
			actual = scraper.scrapeRepo(test.rp, test.retrieveDoc)
		}
		assert.Equal(t, test.expected.out, actual.out)

		if actual.errOut != nil {
			assert.Equal(t, test.expected.errOut[0].Error(), actual.errOut[0].Error())
		} else {
			assert.Equal(t, test.expected.errOut, actual.errOut)
		}
		server.Close()
	}
}

func TestScrape(t *testing.T) {
	server, _ := setupMockServer(Success)

	scraper := Scraper{httpClient: server.Client(),
		gitHubClient: GitClient{Repositories: &mockGitClient{}}, organization: "thebeat",
		filter: &mockFilter{}}
	var result []scrape.DocDef

	actual := scraper.Scrape()
	assert.Equal(t, result, actual)

}
