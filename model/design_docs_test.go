package model

import (
	"testing"
)

func TestDesignDocCreation(t *testing.T) {
	d := NewDesignDoc(NewTestModel("design_test"))
	d.Filters["all"] = `function(doc) { if (doc.type == "design_test") { return true; } else {return false; }}`
	d.Views["all"] = map[string]string{"map": `function(doc) { if (doc.type == "design_test") { emit(null, doc) }}`}
	d.Id = d.Id + getTestId()
	id, _, err := d.Create()
	if err != nil {
		t.Errorf("%v", err)
	}
	dtest := NewDesignDoc(NewTestModel("design_test"))
	_, err = dtest.Retrieve(id, dtest)
	if err != nil {
		t.Error(err.Error())
	}
	defer d.Delete()
}

func TestEnsureUptoDate(t *testing.T) {
	d := NewDesignDoc(NewTestModel("desing_ensure"))
	d.Filters["all"] = `function(doc) { if (doc.type == "design_test") { return true; } else {return false; }}`
	d.Views["all"] = map[string]string{"map": `function(doc) { if (doc.type == "design_test") { emit(null, doc) }}`}
	d.Id = d.Id + getTestId()
	err := d.EnsureUptoDate()
	if err != nil {
		t.Error(err.Error())
	}
	dTest := NewDesignDoc(NewTestModel("design_ensure"))
	dTest.Filters["all"] = `function(doc) { if (doc.type == "somethingelse") { return true; } else {return false; }}`
	dTest.Views["all"] = map[string]string{"map": `function(doc) { if (doc.type == "design_test") { emit(null, doc) }}`}
	dTest.Id = d.Id
	err = dTest.EnsureUptoDate()
	if err != nil {
		t.Error(err.Error())
	}
	uptoDate, err := dTest.IsUpToDate()
	if err != nil {
		t.Error(err.Error())
	}
	if !uptoDate {
		t.Error("Update of design docs are not working properly")
	}
	defer dTest.Delete()
}
