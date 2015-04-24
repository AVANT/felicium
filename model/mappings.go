package model

import (
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"strings"
)

type Mapping struct {
	*Properties
}

type Properties map[string]Property

type Property struct {
}

/*
This creates a mapping on models specific index.
*/
func (m *Model) PutMapping(mapping string) ([]byte, error) {
	indexName := strings.ToLower(m.GetType())
	request, err := api.ElasticSearchRequest("PUT", fmt.Sprintf("/%s/%s/_mapping", indexName, indexName))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBodyString(mapping)

	var toFill interface{}
	code, body, err := request.Do(&toFill)
	if code > 300 || err != nil {
		return body, err
	}
	return body, nil
}

/*
This gets the mapping of a model specific index.
*/
func (m *Model) GetMapping() ([]byte, error) {
	indexName := strings.ToLower(m.GetType())
	request, err := api.ElasticSearchRequest("GET", fmt.Sprintf("/%s/%s/_mapping", indexName, indexName))
	if err != nil {
		return nil, err
	}

	var toFill interface{}
	code, body, err := request.Do(&toFill)
	if code > 300 {
		return nil, err
	}
	return body, nil
}
