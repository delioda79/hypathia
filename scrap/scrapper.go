package scrap

const (
	// Swagger definition
	Swagger = iota
	// Async definition
	Async
)

// Docdef represents a documentation definition
type DocDef struct{
	Type int
	URL string
	Definition string
}
type Scrapper interface{
	Scrap() []DocDef
}
