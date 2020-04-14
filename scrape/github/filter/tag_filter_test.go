package filter

import (
	"github.com/google/go-github/v25/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagFilter_New(t *testing.T) {
	topicTests := [][]string{{"a", "b"}, {}}

	for _, test := range topicTests {
		tagFilter := New(test)
		assert.IsType(t, *tagFilter, TagFilter{})
		assert.Equal(t, tagFilter.topics, test)
		assert.Equal(t, len(tagFilter.topics), len(test))
	}
}

func TestTagFilter_Apply(t *testing.T) {
	type tagFilterApplyTest struct {
		topic       []string
		inputRepos  []*github.Repository
		outputRepos []*github.Repository
	}

	rName := "test"

	var tagFilterApplyTests = []tagFilterApplyTest{
		{
			topic:       []string{"a", "b"},
			inputRepos:  []*github.Repository{{Topics: []string{"a"}, Name: &rName}, {Topics: []string{"b"}, Name: &rName}},
			outputRepos: []*github.Repository{{Topics: []string{"a"}, Name: &rName}, {Topics: []string{"b"}, Name: &rName}},
		},
		{
			topic:       []string{"a", "b"},
			inputRepos:  []*github.Repository{{Topics: []string{"c"}, Name: &rName}, {Topics: []string{"d"}, Name: &rName}},
			outputRepos: nil,
		},
		{
			topic:       []string{},
			inputRepos:  []*github.Repository{{Topics: []string{"a"}, Name: &rName}, {Topics: []string{"b"}, Name: &rName}},
			outputRepos: []*github.Repository{{Topics: []string{"a"}, Name: &rName}, {Topics: []string{"b"}, Name: &rName}},
		},
	}

	for _, test := range tagFilterApplyTests {
		assert.Equal(t, New(test.topic).Apply(test.inputRepos), test.outputRepos)
	}
}
