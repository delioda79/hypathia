package fs

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	paths := []struct {
		path string
		err  bool
	}{
		{"./testdir", false},
		{"nonexistant", true},
	}

	for _, p := range paths {
		fs, err := New(p.path)
		if p.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.IsType(t, &Scraper{}, fs)
		}
	}
}

func TestScrape(t *testing.T) {

	path := "./testdir"

	fs, err := New(path)
	assert.Nil(t, err)

	rs := fs.Scrape()

	fl, err := ioutil.ReadFile("./testdir/test1/docs/swagger.json")
	assert.Nil(t, err)

	assert.Len(t, rs, 1)
	assert.Equal(t, string(fl), rs[0].Definition)
}
