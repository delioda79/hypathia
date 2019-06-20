package berserker

import (
	"sync"

	"github.com/beatlabs/patron/log"

	"github.com/taxibeat/hypatia/search"

	"github.com/taxibeat/hypatia/search/engine"
)

// Berserker is a text indexing and search worker
type Berserker struct {
	*sync.Mutex
	idx     search.Indexer
	fdr     search.Finder
	idxChan chan search.Document
}

// NewBerserker returns a new Berserker
func NewBerserker() (*Berserker, error) {
	eng, err := engine.NewEngine()
	if err != nil {
		return nil, err
	}

	return &Berserker{
		Mutex:   &sync.Mutex{},
		idx:     eng,
		fdr:     eng,
		idxChan: make(chan search.Document),
	}, nil
}

func (b *Berserker) Index(d search.Document) {
	b.idxChan <- d
}

func (b *Berserker) Find(txt string) ([]string, error) {

	return b.fdr.Find(txt)
}

func (b *Berserker) Run() {
	go func() {
		for {
			d := <-b.idxChan
			b.Lock()
			err := b.idx.Index(d)
			if err != nil {
				log.Debugf("Berserker errored while indexing", err)
			}
			b.Unlock()
		}
	}()
}
