package scrape

import (
	"encoding/json"

	"github.com/google/go-github/v25/github"
)

// DocType represents the type of documentation, sync/async
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

// GetID returns the id of the documentation item
func (dd DocDef) GetID() string {
	return dd.ID
}

// Content returns the documentation
func (dd DocDef) Content() (interface{}, error) {
	var rs interface{}
	err := json.Unmarshal([]byte(dd.Definition), &rs)
	return rs, err
}

// Scraper is user to gather the documentations
type Scraper interface {
	Scrape() []DocDef
}

// Filter filters teh repos based on criteria like tags
type Filter interface {
	Apply([]*github.Repository) []*github.Repository
}
