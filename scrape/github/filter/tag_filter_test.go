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

	var tagFilterApplyTests = []tagFilterApplyTest{
		{
			topic:       []string{"a", "b"},
			inputRepos:  []*github.Repository{{Topics: []string{"a"}}, {Topics: []string{"b"}}},
			outputRepos: []*github.Repository{{Topics: []string{"a"}}, {Topics: []string{"b"}}},
		},
		{
			topic:       []string{"a", "b"},
			inputRepos:  []*github.Repository{{Topics: []string{"c"}}, {Topics: []string{"d"}}},
			outputRepos: nil,
		},
		{
			topic:       []string{},
			inputRepos:  []*github.Repository{{Topics: []string{"a"}}, {Topics: []string{"b"}}},
			outputRepos: []*github.Repository{{Topics: []string{"a"}}, {Topics: []string{"b"}}},
		},
	}

	for _, test := range tagFilterApplyTests {
		assert.Equal(t, New(test.topic).Apply(test.inputRepos), test.outputRepos)
	}
}
