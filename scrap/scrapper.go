package scrap

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

// Docdef represents a documentation definition
type DocDef struct {
	Type       DocType
	RepoName   string
	URL        string
	Definition string
}
type Scrapper interface {
	Scrap() []DocDef
}
