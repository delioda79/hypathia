package scrape

import "github.com/google/go-github/v25/github"

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
	Type       DocType
	RepoName   string
	URL        string
	Definition string
}
type Scraper interface {
	Scrape() []DocDef
}

type Filter interface {
	Apply([]*github.Repository) []*github.Repository
}
