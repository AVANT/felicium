package models

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/core"
)

///
//	These are the required methods
///

type Tags struct {
	Hash map[string]bool
	Tags []string
}

func NewTags() *Tags {
	return &Tags{
		Hash: map[string]bool{},
		Tags: []string{},
	}
}

func (t *Tags) Len() int                  { return len((*t).Tags) }
func (t *Tags) Get(index int) interface{} { return (*t).Tags[index] }
func (t *Tags) Add(s string) {
	_, found := (*t).Hash[s]
	if !found {
		(*t).Hash[s] = true
		(*t).Tags = append((*t).Tags, s)
	}
}
func (t *Tags) Merge(s []string) {
	for _, v := range s {
		t.Add(v)
	}
}

//GetPostsByUpdatedAt gets posts based on the last updated time. This will error out if there isn't a mapping defined for posts yet. That is to say if no post ever made it into the system.
func GetAllTags() (*Tags, error) {
	toReturn := NewTags()

	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]string{},
		},
		"fields": []string{"valueHash.tags"},
		"filter": map[string]interface{}{
			"exists": map[string]string{
				"field": "valueHash.tags",
			},
		},
		"size": 1000,
	}

	results, err := core.SearchRequest(false, "post", "post", qry, "", 0)
	if err != nil {
		return toReturn, err
	}

	type tagPartial struct {
		Values []string `json:"valueHash.tags"`
	}

	for _, v := range results.Hits.Hits {
		tags := &tagPartial{}
		err := json.Unmarshal(v.Fields, tags)
		if err != nil {
			return toReturn, err
		}
		toReturn.Merge((*tags).Values)
	}
	return toReturn, nil
}
