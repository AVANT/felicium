package model

import (
	"testing"
	"time"
)

func testHelperRiverCreate(t *testing.T, name string) {
	river := NewRiver(NewTestModel(name))
	err := river.CreateOrUpdate()
	t.Log(river)
	if err != nil {
		t.Error(err)
	}
}

func testHelperRiverDelete(t *testing.T, name string) {
	river := NewRiver(NewTestModel(name))
	err := river.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestRiverLifeCycle(t *testing.T) {
	t.Log("Testing River Creation")
	ensureTestingConnection(t)
	//build the river from a model stub named RiverTestModel
	testHelperRiverCreate(t, "RiverTestModel")
	defer testHelperRiverDelete(t, "RiverTestModel")
}

func TestRiverExists(t *testing.T) {
	ensureTestingConnection(t)
	river := NewRiver(NewTestModel("RiverDoesNotExist"))
	testHelperRiverDelete(t, "RiverDoesNotExist")
	exists, err := river.Exists()
	if exists {
		t.Error("River exist and it shouldn't")
	}
	if err != nil {
		t.Error(err)
	}
	testHelperRiverCreate(t, "RiverDoesNotExist")
	//need to wait for es here should probably make this part of the core rivers functionalit later
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	<-timeout

	exists, err = river.Exists()
	if !exists {
		t.Error("River doesn't exist and it should")
	}
	if err != nil {
		t.Error(err)
	}
	defer testHelperRiverDelete(t, "RiverDoesNotExist")
}
