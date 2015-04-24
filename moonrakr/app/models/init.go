package models

import (
	"github.com/robfig/revel"
	"github.com/AVANT/felicium/model"
)

//Connection holds the information about the connections to couchdb and elastic search.
var Connection *model.Connection

//Setup is called anytime we want to use the models package and it creates the global Connection object.
func Setup(CouchHost, CouchPort, CouchDatabase, EsDomain, EsPort, EsProto string) {

	var err error
	Connection, err = model.New(
		CouchHost,
		CouchPort,
		CouchDatabase,
		EsDomain,
		EsPort,
		EsProto,
	)
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	//database is setup lets make sure that everything that we need is registerd
	registerModels()
}

//registerModels calls the Register method of the model interface. This triggers the rivers and couch views to be updated or
// created if needed.
func registerModels() {
	err := NewPost().Register()
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	err = NewUser().Register()
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	err = NewComment().Register()
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	err = NewMedia().Register()
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	//going to let elastic search do the work here
	// err = NewTag().Register()
	// if err != nil {
	// 	revel.ERROR.Fatal(err)
	// }
}
