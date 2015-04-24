package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type DesignDoc struct {
	*Connection `json:"-"`
	Id          string                       `json:"_id"`
	Rev         string                       `json:"_rev"`
	Language    string                       `json:"language"`
	Views       map[string]map[string]string `json:"views"`
	Filters     map[string]string            `json:"filters"`
	CreatedAt   time.Time                    `json:"createdAt"`
	UpdatedAt   time.Time                    `json:"updatedAt"`
}

func NewDesignDoc(m IsModel) *DesignDoc {
	toReturn := DesignDoc{Language: "javascript"}
	toReturn.Connection = m.GetModel().Connection
	toReturn.Views = map[string]map[string]string{}
	toReturn.Filters = map[string]string{}
	toReturn.Id = "_design/" + m.GetType()
	return &toReturn
}

///
// connectable Interface Requirenments
///

func (d *DesignDoc) GetConnection() *Connection {
	return d.Connection
}

func (d *DesignDoc) SetCreatedAt(t time.Time) {
	d.CreatedAt = t
}

func (d *DesignDoc) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *DesignDoc) SetUpdatedAt(t time.Time) {
	d.UpdatedAt = t
}

func (d *DesignDoc) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}

///
// crudy Interface Requirenments
///

func (d *DesignDoc) GetId() string {
	return d.Id
}

func (d *DesignDoc) GetRev() string {
	return d.Rev
}

func (d *DesignDoc) SetId(id string) {
	d.Id = id
}

func (d *DesignDoc) SetRev(rev string) {
	d.Rev = rev
}

///
// basic lifecycle methods
///

func (d *DesignDoc) Create() (string, string, error) {
	return createOrUpdateCrudy(d)
}

func (d *DesignDoc) Update() (string, string, error) {
	return createOrUpdateCrudy(d)
}

func (d *DesignDoc) Delete() error {
	return deleteCrudy(d)
}

///
// idempotent helpers
///

func (d *DesignDoc) EnsureUptoDate() error {
	uptodate, err := d.IsUpToDate()
	if err != nil {
		return err
	}
	//if you need the update make sure you have the rev before you try and update
	if !uptodate {
		if d.Rev == "" {
			var trashme interface{}
			rev, err := d.Retrieve(d.Id, &trashme)
			if err != nil {
				if err.Error() != fmt.Sprintf("couldn't Retrieve %s: 404 Object Not Found", d.Id) {
					return err
				}
			}
			d.Rev = rev
		}
		//could also be create
		d.Update()
	}
	return nil
}

/*
This will check that the current design doc matches the one in the db.
Currently this uses a sloppy string comparison. Because of time and the limited amout
of times this method will be run this should be sufficient.
*/
func (d *DesignDoc) IsUpToDate() (bool, error) {
	dTest := &DesignDoc{}
	_, err := d.Retrieve(d.Id, dTest)
	if err != nil {
		if err.Error() == fmt.Sprintf("couldn't Retrieve %s: 404 Object Not Found", d.Id) {
			return false, nil
		} else {
			return false, err
		}
	}
	//ignore the rev numbers if set because that is not the point
	d.Rev = ""
	dTest.Rev = ""
	dJSON, _ := json.Marshal(d)
	if err != nil {
		return false, err
	}
	dTestJSON, err := json.Marshal(dTest)
	if err != nil {
		return false, err
	}
	return (string(dTestJSON) == string(dJSON)), nil
}
