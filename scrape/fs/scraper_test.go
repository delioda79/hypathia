package fs

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestScraper_Scrape(t *testing.T) {
	fl1, err := ioutil.ReadFile("./testdir1/test/docs/swagger.json")
	require.Nil(t, err)
	fl2, err := ioutil.ReadFile("./testdir2/test/doc/swagger.json")
	require.Nil(t, err)

	tests := map[string]struct {
		path string
		want []byte
	}{
		"docs folder": {
			path: "./testdir1",
			want: fl1,
		},
		"doc folder": {
			path: "./testdir2",
			want: fl2,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs, err := New(tt.path)
			assert.Nil(t, err)

			rs := fs.Scrape()

			assert.Len(t, rs, 1)
			assert.Equal(t, string(tt.want), rs[0].Definition)
		})
	}
}
