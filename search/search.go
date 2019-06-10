package search

// Document is a document for indexing
type Document interface {
	GetID() string
	Content() (interface{}, error)
}

// Indexer performs an indexing of a document
type Indexer interface {
	Index(d Document) error
}

// Indexer performs an indexing of a document
type AsyncIndexer interface {
	Index(d Document)
}

// FInder performs a search and returns document ids
type Finder interface {
	Find(txt string) ([]string, error)
}
