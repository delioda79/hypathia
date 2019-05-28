package scrap

const (
	// Swagger definition
	Swagger = iota
	// Async definition
	Async
)

// Docdef represents a documentation definition
type DocDef struct {
	Type       int
	RepoName   string
	URL        string
	Definition string
}
type Scrapper interface {
	Scrap() []DocDef
}
