package serve

import (
	"bytes"
	"github.com/gorilla/mux"
	"hypatia/scrape"
	"hypatia/template"
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
	vars := mux.Vars(req)
	repoName := vars["repoName"]
	repoType, err := strconv.Atoi(vars["type"])
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
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
	vars := mux.Vars(req)
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
