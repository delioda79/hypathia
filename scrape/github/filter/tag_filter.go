package filter

import (
	"github.com/google/go-github/v25/github"
)

type TagFilter struct {
	topics []string
}

func New(topics []string) *TagFilter {
	return &TagFilter{
		topics: topics,
	}
}

//Apply filters Github repositories based on TagFilter.topics provided.
//If TagFilter.tags is empty it returns the initial repositories
func (fl *TagFilter) Apply(repositories []*github.Repository) []*github.Repository {
	if len(fl.topics) == 0 {
		return repositories
	}
	var filRepos []*github.Repository
OUTER:
	for _, repo := range repositories {
		for _, topic := range repo.Topics {
			for _, top := range fl.topics {
				if top == topic {
					filRepos = append(filRepos, repo)
					continue OUTER
				}
			}
		}
	}
	return filRepos
}
