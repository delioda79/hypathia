package serve

import (
	"encoding/json"
	"githubscrapper/scrap"
	"net/http"
	"sync"
)

type Handler struct {
	sync.Mutex
	docs []scrap.DocDef
}

func (hd *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	body, err := json.Marshal(hd.docs)
	if err != nil {

	}
	wr.Write(body)
}

func (hd *Handler) Update(docs []scrap.DocDef) {
	hd.Lock()
	hd.docs = docs
	hd.Unlock()
}



