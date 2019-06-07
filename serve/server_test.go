package serve

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ApiList(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	hdl := &Handler{}
	hdl.docs = []scrape.DocDef{{
		Type:       scrape.Swagger,
		Definition: "{ A swagger json}",
		URL:        "u",
		RepoName:   "rest",
	}}

	rr := httptest.NewRecorder()
	hdl.ApiList(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var expected bytes.Buffer
	template.ApiList(hdl.docs, &expected)
	assert.Equal(t, expected, *rr.Body)
}

func TestHandler_ApiRenderSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/doc/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "0")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}
	hdl.docs = []scrape.DocDef{{
		Type: scrape.Swagger,
		Definition: `{
			  "openapi": "3.0.0",
			  "info": {
				"version": "1.0.0",
				"title": "Swagger Petstore"
			  },
			  "servers": [
				{
				  "url": "http://petstore.swagger.io/api"
				}
			  ],
			  "paths": {
				"/pets": {
				  "delete": {
					"operationId": "deletePet",
					"responses": {
					  "204": {
						"description": "pet deleted"
					  }
					}
				  }
				}
			  }
			}
			`,
		URL:      "u",
		RepoName: "carrot",
	}}

	rr := httptest.NewRecorder()
	hdl.ApiRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	buffer := new(bytes.Buffer)
	template.ApiRender(hdl.docs[0], buffer)

	assert.Equal(t, buffer, rr.Body)

}

func TestHandler_ApiRenderNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/doc/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "0")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}

	rr := httptest.NewRecorder()
	hdl.ApiRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestHandler_ApiRenderInvalidType(t *testing.T) {
	req, err := http.NewRequest("GET", "/doc/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "invalid")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}

	rr := httptest.NewRecorder()
	hdl.ApiRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestHandler_SpecRenderSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/spec/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "0")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}
	hdl.docs = []scrape.DocDef{{
		Type: scrape.Swagger,
		Definition: `{
			  "openapi": "3.0.0",
			  "info": {
				"version": "1.0.0",
				"title": "Swagger Petstore"
			  },
			  "servers": [
				{
				  "url": "http://petstore.swagger.io/api"
				}
			  ],
			  "paths": {
				"/pets": {
				  "delete": {
					"operationId": "deletePet",
					"responses": {
					  "204": {
						"description": "pet deleted"
					  }
					}
				  }
				}
			  }
			}
			`,
		URL:      "u",
		RepoName: "carrot",
	}}
	rr := httptest.NewRecorder()
	hdl.SpecRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	assert.Equal(t, hdl.docs[0].Definition, rr.Body.String())
}

func TestHandler_SpecRenderNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/spec/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "0")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}

	rr := httptest.NewRecorder()
	hdl.SpecRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestHandler_SpecRenderInvalidType(t *testing.T) {
	req, err := http.NewRequest("GET", "/spec/:repoName/:type", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("type", "invalid type")
	q.Add("repoName", "carrot")
	req.URL.RawQuery = q.Encode()

	hdl := &Handler{}

	rr := httptest.NewRecorder()
	hdl.SpecRender(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestHandler_HealthStatusSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-status", nil)
	if err != nil {
		t.Fatal(err)
	}

	hdl := &Handler{}
	hdl.ready = true
	rr := httptest.NewRecorder()
	hdl.HealthStatus(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandler_HealthStatusFail(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-status", nil)
	if err != nil {
		t.Fatal(err)
	}

	hdl := &Handler{}
	rr := httptest.NewRecorder()
	hdl.HealthStatus(rr, req)

	assert.NotNil(t, rr)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestHandler_Update(t *testing.T) {
	hdl := &Handler{}

	updatedDocs := []scrape.DocDef{{
		Type:       scrape.Swagger,
		Definition: "{ A swagger json}",
		URL:        "u",
		RepoName:   "rest",
	}}

	hdl.Update(updatedDocs)

	assert.Equal(t, updatedDocs, hdl.docs)
	assert.Equal(t, true, hdl.ready)
}
