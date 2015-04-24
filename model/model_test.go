package model

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//Define A New Model

type TestModel struct {
	*Model
}

//Create overrides for specific Methods defined on the model struct
func (r *TestModel) Create() (string, string, error) {
	r.Id = getTestId()
	//this is like a call to super in an object oriented language
	return r.Model.Create()
}

// Define a new method for your model
// This method should be a reciever for your ApplicationConnection
// This will just embed a new Model in your model
func NewTestModel(modelType string) *TestModel {
	return &TestModel{TestingConnection.NewModel(modelType)}
}

//Write a Custom Validation
func (r *TestModel) Validate() error {
	if val, ok := (*r.ValueHash)["name"]; ok && val == "bad value" {
		return errors.New("name cannot equal 'bad value'")
	}
	return nil
}

func TestCreate(t *testing.T) {
	t.Log("Testing Create Methods")
	ensureTestingConnection(t)
	createModelWithTests(t)
}

func TestUpdate(t *testing.T) {
	t.Log("Testing Update Methods")
	ensureTestingConnection(t)
	m := createModelWithTests(t)
	t.Log("Model Created Updating Attributes")
	(*m.ValueHash)["test"] = "something"
	err := Update(m)
	t.Log("Model Updated")
	if err != nil {
		t.Errorf("%v", err)
	}
	updateCheck := NewTestModel("TestModel")
	t.Log("Retrieving Data")
	_, err = TestingConnection.Retrieve(m.GetId(), updateCheck)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestDelete(t *testing.T) {
	t.Log("Testing Update Methods")
	ensureTestingConnection(t)
	m := createModelWithTests(t)
	id := m.Id
	t.Log("Model Created. Attempting to delete model.")
	if err := Delete(m); err != nil {
		t.Errorf("%v", err)
	}
	t.Log("Retrieving Data")
	deleteCheck := NewTestModel("TestModel")
	_, err := TestingConnection.Retrieve(id, deleteCheck)
	// There could be a better way to check this but im not sure right now
	if fmt.Sprintf("%v", err) != fmt.Sprintf("couldn't Retrieve %s: 404 Object Not Found", id) {
		t.Errorf("%v", err)
	}
}

func TestValidate(t *testing.T) {
	t.Log("Testing Validation Methods")
	ensureTestingConnection(t)
	m := createModelWithTests(t)
	id := m.Id
	t.Log("Model Created and Validation Passed. Attempting to add bad data.")
	validateCheck := NewTestModel("TestModel")
	_, err := TestingConnection.Retrieve(id, validateCheck)
	if err != nil {
		t.Errorf("%v", err)
	}
	(*validateCheck.ValueHash)["name"] = "bad value"
	err = Save(validateCheck)
	if fmt.Sprintf("%v", err) != fmt.Sprintf("name cannot equal 'bad value'") {
		t.Errorf("%v", err)
	}
}

func TestRegister(t *testing.T) {
	m := createModelByNameWithTests(t, "RegisterTest")
	err := m.Register()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer m.DeRegister()
}

// Some helpers to dry up the tests
func getTestId() string {
	return fmt.Sprintf("test-%s%b", time.Now().Format("20060102150405"), rand.New(rand.NewSource(time.Now().UnixNano())).Float64())
}

func createModelByNameWithTests(t *testing.T, name string) *TestModel {
	testModel := NewTestModel(name)
	t.Log("Creating model")
	err := Create(testModel)
	if testModel.GetId() == testModel.GetRev() && testModel.GetId() == "" {
		t.Errorf("Create does not return an id and revision.")
	}
	if err != nil {
		t.Errorf("%v", err)
	}

	return testModel
}

func createModelWithTests(t *testing.T) *TestModel {
	return createModelByNameWithTests(t, "TestModel")
}
