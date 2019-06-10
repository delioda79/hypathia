package berserker

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/hypatia/search"
	"github.com/taxibeat/hypatia/search/searchfakes"
)

func TestBerserkIndexing(t *testing.T) {

	idxer := &searchfakes.FakeIndexer{}
	fdr := &searchfakes.FakeFinder{}
	ch := make(chan search.Document)
	docs := []*searchfakes.FakeDocument{
		&searchfakes.FakeDocument{},
		&searchfakes.FakeDocument{},
	}

	docs[0].GetIDReturns("1")
	docs[1].GetIDReturns("2")

	docs[0].ContentReturns("", nil)
	docs[1].ContentReturns("", errors.New("error"))

	errs := []error{}

	wg := sync.WaitGroup{}

	brsrk := Berserker{
		idxChan: ch,
		idx:     idxer,
		fdr:     fdr,
	}

	idxer.IndexStub = func(d search.Document) error {
		wg.Done()
		_, err := d.Content()
		if err != nil {
			errs = append(errs, err)
		}
		return err
	}

	brsrk.Run()
	for _, d := range docs {
		wg.Add(1)
		brsrk.Index(d)
	}
	wg.Wait()

	if idxer.IndexCallCount() != len(docs) {
		t.Error("Wrong calls count: ", idxer.IndexCallCount(), " but expected ", len(docs))
	}

	if len(errs) != 1 {
		t.Errorf("Expected only one error but found %d", len(errs))
	}
}

func TestBerserkFinding(t *testing.T) {
	idxer := &searchfakes.FakeIndexer{}
	fdr := &searchfakes.FakeFinder{}
	ch := make(chan search.Document)

	brsrk := Berserker{
		idxChan: ch,
		idx:     idxer,
		fdr:     fdr,
	}

	fdr.FindReturnsOnCall(0, []string{}, nil)
	fdr.FindReturnsOnCall(1, []string{}, errors.New("error"))

	_, e := brsrk.Find("text")
	if e != nil {
		t.Errorf("Expected no error but we've got %v", e)
	}

	_, e = brsrk.Find("text")
	if e == nil {
		t.Error("Expected an error")
	}
}

func TestNewBerserker(t *testing.T) {
	brsk, err := NewBerserker()
	assert.IsType(t, &Berserker{}, brsk)
	assert.Nil(t, err)
}

func BenchmarkBerserker_Find(b *testing.B) {
	bsrk, _ := NewBerserker()
	bsrk.Run()

	b.Run("Indexing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bsrk.Index(generateDoc(i))
		}
	})

	b.Run("Searching", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bsrk.Find(fmt.Sprintf("Name%d", i))
		}
	})
}

func generateDoc(i int) search.Document {
	d := searchfakes.FakeDocument{}
	d.GetIDReturns(fmt.Sprintf("%d", i))
	d.ContentReturns(map[string]string{"id": fmt.Sprintf("%d", i), "name": fmt.Sprintf("Name%d", i)}, nil)
	return &d
}
