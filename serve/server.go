package serve

import (
	"bytes"
	"github.com/gorilla/mux"
	"githubscrapper/scrap"
	"githubscrapper/template"
	"net/http"
	"strconv"
	"sync"
)

type Handler struct {
	sync.Mutex
	docs []scrap.DocDef
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
		if d.RepoName == repoName && d.Type == scrap.DocType(repoType) {
			template.ApiRender(d, buffer)
			wr.Write(buffer.Bytes())
			return
		}
	}
	wr.WriteHeader(http.StatusNotFound)
}

func (hd *Handler) Update(docs []scrap.DocDef) {
	hd.Lock()
	hd.docs = docs
	hd.Unlock()
}
