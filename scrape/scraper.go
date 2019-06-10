package scrape

import (
	"encoding/json"

	"github.com/google/go-github/v25/github"
)

type DocType int

const (
	// Swagger definition
	Swagger DocType = 0
	// Async definition
	Async DocType = 1
)

func (docDef DocType) String() string {
	names := [...]string{
		"Swagger",
		"Async",
	}
	return names[docDef]
}

// DocDef represents a documentation definition
type DocDef struct {
	ID         string
	Type       DocType
	RepoName   string
	URL        string
	Definition string
}

func (dd DocDef) GetID() string {
	return dd.ID
}

func (dd DocDef) Content() (interface{}, error) {
	var rs interface{}
	err := json.Unmarshal([]byte(dd.Definition), &rs)
	return rs, err
}

type Scraper interface {
	Scrape() []DocDef
}

type Filter interface {
	Apply([]*github.Repository) []*github.Repository
}
