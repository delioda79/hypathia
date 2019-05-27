package githubscrapper

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
}
type Scrapper interface{
	Scrap() []DocDef
}
