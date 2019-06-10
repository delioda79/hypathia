package engine

import (
	"github.com/blevesearch/bleve"
	"github.com/taxibeat/hypatia/search"
)

// Engine is a text indexing and search engine
type Engine struct {
	idx bleve.Index
}

// NewEngine returns a new search engine
func NewEngine() (*Engine, error) {
	idmapping := bleve.NewIndexMapping()
	idmapping.StoreDynamic = true
	index, err := bleve.NewMemOnly(idmapping)
	if err != nil {
		return nil, err
	}

	return &Engine{
		idx: index,
	}, nil
}

func (e *Engine) Index(d search.Document) error {
	cnt, err := d.Content()
	if err != nil {
		return err
	}
	return e.idx.Index(d.GetID(), cnt)
}

func (e *Engine) Find(txt string) ([]string, error) {
	query := bleve.NewMatchQuery(txt)
	search := bleve.NewSearchRequest(query)
	searchResults, err := e.idx.Search(search)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, h := range searchResults.Hits {
		ids = append(ids, h.ID)
	}

	return ids, nil
}
