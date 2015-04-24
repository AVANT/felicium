package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"strings"
	//"io/ioutil"
)

type River struct {
	*Connection `json:"-"`
	Type        string            `json:"type"`
	Couchdb     map[string]string `json:"couchdb"`
	Index       *RiverIndex       `json:"index"`
}

type RiverIndex struct {
	Index       string `json:"index"`
	Type        string `json:"type"`
	BulkSize    string `json:"bulk_size"`
	BulkTimeout string `json:"bulk_timeout"`
}

func NewRiver(m IsModel) *River {
	toReturn := River{}
	toReturn.Connection = m.GetModel().Connection
	toReturn.Type = "couchdb"
	toReturn.Couchdb = map[string]string{
		"host":   toReturn.GetCouchHost(),
		"port":   toReturn.GetCouchPort(),
		"db":     toReturn.GetCouchName(),
		"filter": "",
	}
	toReturn.Index = &RiverIndex{
		Index:       strings.ToLower(m.GetType()),
		Type:        strings.ToLower(m.GetType()),
		BulkSize:    "100",
		BulkTimeout: "10ms",
	}
	return &toReturn
}

/*
This is a convience method that will prevent rivers from starting their indexes over when they are created.
*/
func (r *River) EnsureExists() error {
	check, err := r.Exists()
	if err != nil {
		return err
	}
	if check {
		return nil
	} else {
		return r.CreateOrUpdate()
	}
}

/*
This creates or rebuilds a river when called on the river struct.
*/
func (r *River) CreateOrUpdate() error {
	riverName := r.Index.Index
	request, err := api.ElasticSearchRequest("PUT", fmt.Sprintf("/_river/%s/_meta", riverName))
	if err != nil {
		return err
	}
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBodyString(string(body))
	//don't know why the set body json doesn't work but it kept putting unneeded chars in the json
	//err = request.SetBodyJson(r)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(r)
	// s, _ := ioutil.ReadAll(request.Body)
	// fmt.Printf("%s", s)
	var toFill interface{}
	code, body, err := request.Do(&toFill)
	if code > 300 {
		return err
	}
	return nil
}

/*
Given a river description this method deletes the river.
*/
func (r *River) Delete() error {
	riverName := r.Index.Index
	request, err := api.ElasticSearchRequest("DELETE", fmt.Sprintf("/_river/%s", riverName))
	if err != nil {
		return err
	}
	var toFill interface{}
	code, _, err := request.Do(&toFill)
	if code > 300 {
		return err
	}
	return nil
}

/*
Check if river exists. This returns an error if an are encountered.
*/
func (r *River) Exists() (bool, error) {
	riverName := r.Index.Index
	request, err := api.ElasticSearchRequest("GET", fmt.Sprintf("/_river/%s/_status", riverName))
	if err != nil {
		return false, err
	}
	var toFill map[string]interface{}
	code, _, err := request.Do(&toFill)
	if code == 404 {
		status, exists := toFill["exists"]
		if exists {
			return status.(bool), nil
		} else {
			status, exists = toFill["error"]
			if exists && status.(string) == "IndexMissingException[[_river] missing]" {
				//this should be the case where no rivers have been setup yet
				return false, nil
			} else {
				return false, errors.New("404 was returned from ES but the error object was not reconized.")
			}
		}
	} else {
		return true, nil
	}
}
