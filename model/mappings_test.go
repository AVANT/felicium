package model

import (
	"testing"
)

func TestMappingUpdate(t *testing.T) {
	testHelperRiverCreate(t, "mappingTest")
	n := NewTestModel("mappingTest")
	_, err := n.PutMapping(`{"mappingTest": {"properties": { "testProperty": "string" } } }`)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = n.GetMapping()
	if err != nil {
		t.Error(err.Error())
	}
	testHelperRiverDelete(t, "mappingTest")
}
