package serve

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/taxibeat/hypatia/bounddata"

	http2 "github.com/beatlabs/patron/sync/http"

	"github.com/taxibeat/hypatia/search"

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
	repoID := vars["repoID"]
	buffer := new(bytes.Buffer)
	for _, d := range hd.apiDocDefs {
		if d.ID == repoID {
			if d.Type == scrape.Swagger {
				wr.Header().Set("Etag", strconv.FormatInt(time.Now().UnixNano(), 16))
				wr.Header().Set("Cache-Control", "public, max-age=0")
				template.ApiRender(d, buffer)
				wr.Write(buffer.Bytes())
			} else if d.Type == scrape.Async {
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
	repoID := vars["repoID"]
	buffer := new(bytes.Buffer)
	for _, d := range hd.apiDocDefs {
		if d.ID == repoID {
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

func (hd *Handler) StaticFiles(wr http.ResponseWriter, req *http.Request) {

	req.URL.Path = strings.TrimPrefix(req.URL.Path, "/static/")

	bts, err := bounddata.Asset(req.URL.Path)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		wr.Write([]byte(err.Error()))
	}

	if p := strings.TrimPrefix(req.URL.Path, "css"); len(p) < len(req.URL.Path) {
		wr.Header().Add("Content-Type", "text/css")
	}

	wr.Write(bts)
}

func (hd *Handler) HealthStatus() http2.HealthStatus {
	if hd.ready {
		return http2.Healthy
	}

	return http2.Initializing
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
