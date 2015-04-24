package model

import (
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/mattbaird/elastigo/core"
	"strings"
	"time"
)

type ValueHash map[string]interface{}

type ValueHashes []*ValueHash

//Used to return searchs see search.go
type Models []IsModel

type Model struct {
	*Connection `json:"-"`
	*ValueHash  `json:"valueHash"`
	Id          string    `json:"_id"`
	Rev         string    `json:"_rev"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Type        string    `json:"type"`
}

type IsModel interface {
	crudy
	Init() error
	BeforeValidate() error
	Validate() error
	AfterValidate() error
	BeforeSave() error
	BeforeUpdate() error
	Update() (string, string, error)
	AfterUpdate() error
	BeforeCreate() error
	Create() (string, string, error)
	AfterCreate() error
	AfterSave() error
	BeforeDelete() error
	Delete() error
	AfterDelete() error
	//others
	GetValueHash() *ValueHash
	QueryObject() *jsonq.JsonQuery

	SetType(string)
	GetType() string
	//depricated by QueryObject but needed for time still
	GetValue(string) (interface{}, bool)

	GetModel() *Model
	//for Startup
	Register() error
	DeRegister() error
}

func (c *Connection) NewModel(typeOf string) *Model {
	return &Model{
		Connection: c,
		ValueHash:  &ValueHash{},
		Id:         "",
		Rev:        "",
		Type:       typeOf,
	}
}

// This method will determine if the record should be updated or created.
// This method will run the update or create hooks and also run the before and after save hooks.
// Im not forcing this to be a reciever as the model should have reference to the database
// func (c *CouchConnection) Save(m IsModel) (string, string, error) {
func Save(m IsModel) error {
	if err := m.BeforeSave(); err != nil {
		return err
	}
	//This assumption might turn out to be to limiting.
	//If the revision id is empty it means we are creating the model.
	var err error
	if m.GetRev() == "" && m.GetId() == "" {
		err = Create(m)
	} else {
		err = Update(m)
	}
	if err != nil {
		return err
	}
	if err := m.AfterSave(); err != nil {
		return err
	}
	return nil
}

// Im not forcing this to be a reciever as the model should have reference to the database
// func (c *CouchConnection) Create(m IsModel) (string, string, error) {
func Create(m IsModel) error {
	if err := Validate(m); err != nil {
		return err
	}
	if err := m.BeforeCreate(); err != nil {
		return err
	}
	id, rev, err := m.Create()
	if err != nil {
		return err
	}
	m.SetId(id)
	m.SetRev(rev)
	if err := m.AfterCreate(); err != nil {
		return err
	}
	return nil
}

// Im not forcing this to be a reciever as the model should have reference to the database
// func (c *CouchConnection) Update(m IsModel) (string, string, error) {
func Update(m IsModel) error {
	if err := Validate(m); err != nil {
		return err
	}
	if err := m.BeforeUpdate(); err != nil {
		return err
	}
	_, rev, err := m.Update()
	if err != nil {
		return err
	}
	m.SetRev(rev)
	if err := m.AfterUpdate(); err != nil {
		return err
	}
	return nil
}

// Im not forcing this to be a reciever as the model should have reference to the database
// func (c *CouchConnection) Validate(m IsModel) error {
func Validate(m IsModel) error {
	if err := m.BeforeValidate(); err != nil {
		return err
	}
	if err := m.Validate(); err != nil {
		return err
	}
	if err := m.AfterValidate(); err != nil {
		return err
	}
	return nil
}

// Im not forcing this to be a reciever as the model should have reference to the database
// func (c *CouchConnection) Delete(m IsModel) error {
func Delete(m IsModel) error {
	if err := m.BeforeDelete(); err != nil {
		return err
	}
	if err := m.Delete(); err != nil {
		return err
	}
	if err := m.AfterDelete(); err != nil {
		return err
	}
	return nil
}

////
// These methods will no longer have access to any of the struct fields of the interface Model type
// They will only have access to the fields of the embded model as they are recievers once again
////

//Use this function to store any needed startup routines
func (m *Model) Init() error {
	return nil
} // Im using receivers here and not Methods that have interfaces as args so you can overwrite them easily.

func (m *Model) GetId() string {
	return m.Id
}

func (m *Model) SetId(id string) {
	m.SetValue("_id", id)
	m.Id = id
}

func (m *Model) GetRev() string {
	return m.Rev
}

func (m *Model) SetRev(rev string) {
	m.Rev = rev
}

func (m *Model) GetValueHash() *ValueHash {
	return m.ValueHash
}

func (m *Model) BeforeValidate() error {
	return nil
}
func (m *Model) Validate() error {
	return nil
}
func (m *Model) AfterValidate() error {
	return nil
}
func (m *Model) BeforeSave() error {
	return nil
}
func (m *Model) BeforeUpdate() error {
	return nil
}
func (m *Model) Update() (string, string, error) {
	return createOrUpdateCrudy(m)
}
func (m *Model) AfterUpdate() error {
	return nil
}
func (m *Model) BeforeCreate() error {
	return nil
}
func (m *Model) Create() (string, string, error) {
	return createOrUpdateCrudy(m)
}
func (m *Model) AfterCreate() error {
	return nil
}
func (m *Model) AfterSave() error {
	return nil
}
func (m *Model) BeforeDelete() error {
	return nil
}
func (m *Model) Delete() error {
	_type := strings.ToLower(m.GetType())
	id := m.GetId()
	err := deleteCrudy(m)
	if err != nil {
		return err
	}
	_, err = core.Delete(false, _type, _type, id, 0, "")
	return err
}

func (m *Model) AfterDelete() error {
	return nil
}

func (m *Model) SetCreatedAt(t time.Time) {
	m.SetValue("createdAt", t)
	m.CreatedAt = t
}

func (m *Model) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Model) SetUpdatedAt(t time.Time) {
	m.SetValue("updatedAt", t)
	m.UpdatedAt = t
}

func (m *Model) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

/////////////////////////////////////
//Shared objects to dry things up ///
/////////////////////////////////////

func (m *Model) GetConnection() *Connection {
	return m.Connection
}

///////////////////////////////////////////////////////////////
//Additional helpers that I think are helpful to have around //
///////////////////////////////////////////////////////////////
func (m *Model) Register() error {
	river, designDoc := m.getComponentsForModel()
	err := designDoc.EnsureUptoDate()
	if err != nil {
		return err
	}
	err = river.EnsureExists()
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) DeRegister() error {
	river, designDoc := m.getComponentsForModel()
	err := river.Delete()
	if err != nil {
		return err
	}
	err = designDoc.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) getComponentsForModel() (*River, *DesignDoc) {
	river := NewRiver(m)
	designDoc := NewDesignDoc(m)
	designDoc.Filters["all"] = fmt.Sprintf(`function(doc) { if (doc.type == "%s") { return true; } else {return false; }}`, m.GetType())
	river.Couchdb["filter"] = m.GetType() + "/all"
	return river, designDoc
}

func (m *Model) GetModel() *Model {
	return m
}

func (m *Model) SetType(name string) {
	m.SetValue("type", name)
	(*m).Type = name
}

func (m *Model) GetType() string {
	return (*m).Type
}

func (m *Model) SetValue(name string, value interface{}) {
	(*m.ValueHash)[name] = value
}

func (m *Model) GetValue(name string) (interface{}, bool) {
	value, found := (*m.ValueHash)[name]
	return value, found
}

func (m *Model) QueryObject() *jsonq.JsonQuery {
	jq := jsonq.NewQuery(map[string]interface{}(*m.GetValueHash()))
	return jq
}
