package serve

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/taxibeat/hypatia/search"

	"github.com/beatlabs/patron/log"
	"github.com/julienschmidt/httprouter"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/template"
)

type Handler struct {
	sync.Mutex
	apiDocDefs   []scrape.DocDef
	api2htmlDocs map[string][]byte
	ready        bool
	Searcher     search.Finder
}

func (hd *Handler) APIList(wr http.ResponseWriter, req *http.Request) {
	buffer := new(bytes.Buffer)
	template.ApiList(hd.apiDocDefs, buffer)
	wr.Write(buffer.Bytes())
}

func (hd *Handler) APISearch(wr http.ResponseWriter, req *http.Request) {
	filtered := []scrape.DocDef{}
	req.ParseForm()
	queries := req.Form["query"]
	query := strings.Join(queries, " ")
	if strings.Trim(query, " ") == "" {
		filtered = hd.apiDocDefs
	} else {
		docs, err := hd.Searcher.Find(query)
		if err != nil {
			wr.WriteHeader(400)
			return
		}

		for _, r := range docs {
			for _, f := range hd.apiDocDefs {
				if r == f.ID {
					filtered = append(filtered, f)
					continue
				}
			}
		}
	}
	buffer := new(bytes.Buffer)
	template.ApiList(filtered, buffer)
	wr.Write(buffer.Bytes())
}

func (hd *Handler) ApiRender(wr http.ResponseWriter, req *http.Request) {
	vars := extractFields(req)
	repoName := vars["repoName"]
	repoType, err := strconv.Atoi(vars["type"])
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		log.Warn(err)
		return
	}
	buffer := new(bytes.Buffer)
	for _, d := range hd.apiDocDefs {
		if d.RepoName == repoName && d.Type == scrape.DocType(repoType) {
			if scrape.DocType(repoType) == scrape.Swagger {
				wr.Header().Set("Etag", strconv.FormatInt(time.Now().UnixNano(), 16))
				wr.Header().Set("Cache-Control", "public, max-age=0")
				template.ApiRender(d, buffer)
				wr.Write(buffer.Bytes())
			} else if scrape.DocType(repoType) == scrape.Async {
				wr.Header().Set("Etag", strconv.FormatInt(time.Now().UnixNano(), 16))
				wr.Header().Set("Cache-Control", "public, max-age=0")
				wr.Header().Set("Content-Type", "text/html; charset=utf-8")
				wr.Write(hd.api2htmlDocs[d.ID])
				return
			}
			return
		}
	}
	wr.WriteHeader(http.StatusNotFound)
}

func (hd *Handler) SpecRender(wr http.ResponseWriter, req *http.Request) {
	vars := extractFields(req)
	repoName := vars["repoName"]
	repoType, err := strconv.Atoi(vars["type"])
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	buffer := new(bytes.Buffer)
	for _, d := range hd.apiDocDefs {
		if d.RepoName == repoName && d.Type == scrape.DocType(repoType) {
			wr.Header().Set("Content-Type", "application/json")
			wr.Header().Set("Cache-Control", "public, max-age=0")
			wr.Header().Set("Etag", strconv.FormatInt(time.Now().UnixNano(), 16))
			buffer.Write([]byte(d.Definition))
			wr.Write(buffer.Bytes())
			return
		}
	}
	wr.WriteHeader(http.StatusNotFound)
}

func (hd *Handler) HealthStatus(wr http.ResponseWriter, req *http.Request) {
	if hd.ready {
		wr.WriteHeader(http.StatusOK)
		return
	}

	wr.WriteHeader(http.StatusBadRequest)
}

func (hd *Handler) Update(docs []scrape.DocDef, asyncRawDocs map[string][]byte) {
	hd.Lock()
	hd.apiDocDefs = docs
	hd.ready = true
	hd.api2htmlDocs = asyncRawDocs
	hd.Unlock()
}

func extractFields(r *http.Request) map[string]string {
	f := make(map[string]string)

	for name, values := range r.URL.Query() {
		f[name] = values[0]
	}

	for k, v := range extractParams(r) {
		f[k] = v
	}
	return f
}

func extractParams(r *http.Request) map[string]string {
	par := httprouter.ParamsFromContext(r.Context())
	if len(par) == 0 {
		return make(map[string]string)
	}
	p := make(map[string]string)
	for _, v := range par {
		p[v.Key] = v.Value
	}
	return p
}
