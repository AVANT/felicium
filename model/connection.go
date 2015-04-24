package model

import (
	"github.com/mattbaird/elastigo/api"
	"github.com/peterbourgon/couch-go"
)

/*
This struct allows us to create a custom interface for the couch.Database object.
*/
type Connection struct {
	*couch.Database
}

func New(
	couchUrl string,
	couchPort string,
	couchDatabaseName string,
	esDomain string,
	esPort string,
	esProtocol string,
) (*Connection, error) {
	db, err := couch.NewDatabase(couchUrl, couchPort, couchDatabaseName)
	if err != nil {
		return nil, err
	}
	//seting the es vars
	api.Domain = esDomain
	api.Port = esPort
	api.Protocol = esProtocol
	return &Connection{&db}, nil
}

/*
This interface allows us to promice that the CouchConnection will be able to provide
connection infomation to the methods that use it later.
*/
type ConnectionInfo interface {
	GetCouchHost() string
	GetCouchPort() string
	GetCouchName() string
	GetESDomain() string
	GetESPort() string
	GetESProtocol() string
}

/*
Returns the Database.Host from the couch-go lib.
*/
func (c *Connection) GetCouchHost() string {
	return c.Database.Host
}

/*
Returns the Database.Name from the couch-go lib.
*/
func (c *Connection) GetCouchName() string {
	return c.Database.Name
}

/*
Returns the Database.Port from the couch-go lib.
*/
func (c *Connection) GetCouchPort() string {
	return c.Database.Port
}

/*
Returns the Protocol from the elasticgo lib.
*/
func (c *Connection) GetESProtocol() string {
	return api.Protocol
}

/*
Returns the Domain from the elasticgo lib.
*/
func (c *Connection) GetESDomain() string {
	return api.Domain
}

/*
Returns the Port from the elasticgo lib.
*/
func (c *Connection) GetESPort() string {
	return api.Port
}
