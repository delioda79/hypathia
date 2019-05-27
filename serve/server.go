package serve

import (
	"encoding/json"
	"githubscrapper"
	"net/http"
	"sync"
)

type Handler struct {
	sync.Mutex
	docs []githubscrapper.DocDef
}

func (hd *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	body, err := json.Marshal(hd.docs)
	if err != nil {

	}
	wr.Write(body)
}

func (hd *Handler) Update(docs []githubscrapper.DocDef) {
	hd.Lock()
	hd.docs = docs
	hd.Unlock()
}



