package serve

import (
	"bytes"
	"github.com/beatlabs/patron/log"
	"github.com/julienschmidt/httprouter"
	"github.com/taxibeat/hypatia/scrape"
	"github.com/taxibeat/hypatia/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	sync.Mutex
	docs []scrape.DocDef
	ready bool
}

func (hd *Handler) ApiList(wr http.ResponseWriter, req *http.Request) {
	buffer := new(bytes.Buffer)
	template.ApiList(hd.docs, buffer)
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
	for _, d := range hd.docs {
		if d.RepoName == repoName && d.Type == scrape.DocType(repoType) {
			wr.Header().Set("Etag",  strconv.FormatInt(time.Now().UnixNano(), 16))
			wr.Header().Set("Cache-Control", "public, max-age=0")
			template.ApiRender(d, buffer)
			wr.Write(buffer.Bytes())
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
	for _,d := range hd.docs {
		if d.RepoName == repoName && d.Type == scrape.DocType(repoType) {
			wr.Header().Set("Content-Type", "application/json")
			wr.Header().Set("Cache-Control", "public, max-age=0")
			wr.Header().Set("Etag",  strconv.FormatInt(time.Now().UnixNano(), 16))
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

func (hd *Handler) Update(docs []scrape.DocDef) {
	hd.Lock()
	hd.docs = docs
	hd.ready = true
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
