package model

import (
	"testing"
)

const (
	TEST_COUCH_HOST = "couchdb.local.moonrakr.co"
	TEST_COUCH_PORT = "5984"
	TEST_COUCH_NAME = "moonrakr-model-testing"
	TEST_ES_DOMAIN  = "elasticsearch.local.moonrakr.co"
	TEST_ES_PORT    = "9200"
	//for use with nc -l 9999
	//TEST_ES_DOMAIN   = "localhost"
	//TEST_ES_PORT     = "9999"
	TEST_ES_PROTOCOL = "http"
)

var TestingConnection *Connection

/*
This method will be used by testing suite to keep the testing info in one place.
*/
func ensureTestingConnection(t *testing.T) {
	if TestingConnection == nil {
		var err error
		TestingConnection, err = New(
			TEST_COUCH_HOST,
			TEST_COUCH_PORT,
			TEST_COUCH_NAME,
			TEST_ES_DOMAIN,
			TEST_ES_PORT,
			TEST_ES_PROTOCOL,
		)
		if err != nil {
			t.Error(err)
		}
	}
}

/*
This test is not really needed but I used it while I was working.
*/
func TestConnectionInfoInterface(t *testing.T) {
	ensureTestingConnection(t)
	if TEST_COUCH_HOST != TestingConnection.GetCouchHost() {
		t.Error("GetHost failed for TestCouchConnectionInfoInterface")
	}
	if TEST_COUCH_PORT != TestingConnection.GetCouchPort() {
		t.Error("GetPort failed for TestCouchConnectionInfoInterface")
	}
	if TEST_COUCH_NAME != TestingConnection.GetCouchName() {
		t.Error("GetName failed for TestCouchConnectionInfoInterface")
	}
	if TEST_ES_DOMAIN != TestingConnection.GetESDomain() {
		t.Error("GetName failed for TestCouchConnectionInfoInterface")
	}
	if TEST_ES_PORT != TestingConnection.GetESPort() {
		t.Error("GetName failed for TestCouchConnectionInfoInterface")
	}
	if TEST_ES_PROTOCOL != TestingConnection.GetESProtocol() {
		t.Error("GetName failed for TestCouchConnectionInfoInterface")
	}
}
