package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/emilsjolander/goson"
	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/lib/helpers"
	"github.com/AVANT/felicium/moonrakr/app/lib/results"
	"github.com/jmoiron/jsonq"
	"github.com/kr/s3/s3util"
	"github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
)

var S3FQDN string
var S3Folder string

const mediaTemplate = "medium/public/media"

///
//	These are the required methods
///

//Media object
type Media struct {
	*model.Model
}

//NewMedia create a media object with a handle to the databases
func NewMedia() *Media {
	return &Media{
		Connection.NewModel("media"),
	}
}

//Medium is a collection of *Media it implements goson.Collection
type Medium []*Media

func (m *Medium) Len() int                                         { return len(*m) }
func (m *Medium) Get(index int) interface{}                        { return (*m)[index] }
func (m *Medium) Add(media *Media)                                 { (*m) = append((*m), media) }
func (m *Medium) RenderElement(i int) (*results.JsonResult, error) { return (*m)[i].Render() }
func (m *Medium) Render() (*results.JsonResult, error)             { return results.RenderRenderableCollection(m) }
func (m *Medium) FromInterfaceArray(array []interface{}) error {
	for i := range array {
		mediaObject := array[i].(map[string]interface{})
		media := NewMedia()
		media.MediaMassAssign(&mediaObject)
		media.SetId(mediaObject["id"].(string))
		m.Add(media)
	}
	return nil
}

func (m *Media) Render() (*results.JsonResult, error) {
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	defer w.Flush()
	err := goson.RenderTo(w, mediaTemplate, goson.Args{"Media": m})
	if err != nil {
		return new(results.JsonResult), err
	}
	return &results.JsonResult{buffer}, nil
}

//This will safely set fields of the post that can be set by mass assignment
func (m *Media) MediaMassAssign(bulk *map[string]interface{}) (*Media, error) {
	jq := jsonq.NewQuery(*bulk)
	for k := range *bulk {
		switch k {
		case "belongsTo":
			s, err := jq.String(k)
			if err == nil {
				m.SetBelongsTo(s)
			}
		case "title":
			s, err := jq.String(k)
			if err == nil {
				m.SetTitle(s)
			}
		case "altText":
			s, err := jq.String(k)
			if err == nil {
				m.SetAltText(s)
			}
		}
	}
	return m, nil
}

//MediaMassAsignFromValues is a wrapper that unpacks the url.Vaules befor calling the routine set up for json
func (m *Media) MediaMassAsignFromValues(values url.Values) (*Media, error) {
	return m.MediaMassAssign(helpers.UrlValuesToJsonObject(values))
}

func (m *Media) GetAssociation() string {
	a, _ := m.QueryObject().String("association")
	return a
}

func (m *Media) SetAssociation(e string) {
	m.SetValue("association", e)
}

func (m *Media) GetBelongsTo() string {
	b, _ := m.QueryObject().String("belongsTo")
	return b
}

func (m *Media) SetBelongsTo(e string) {
	m.SetValue("belongsTo", e)
}

func (m *Media) GetUploadedBy() string {
	ub, _ := m.QueryObject().String("uploadedBy")
	return ub
}

func (m *Media) SetUploadedBy(e string) {
	m.SetValue("uploadedBy", e)
}

func (m *Media) GetTitle() string {
	t, _ := m.QueryObject().String("title")
	return t
}

func (m *Media) SetTitle(e string) {
	m.SetValue("title", e)
}

func (m *Media) GetUrl() string {
	u, _ := m.QueryObject().String("url")
	return u
}

func (m *Media) SetUrl(e string) {
	m.SetValue("url", e)
}

func (m *Media) GetStatus() string {
	s, _ := m.QueryObject().String("status")
	return s
}

func (m *Media) SetStatus(e string) {
	m.SetValue("status", e)
}

func (m *Media) GetAltText() string {
	at, _ := m.QueryObject().String("altText")
	return at
}

func (m *Media) SetAltText(e string) {
	m.SetValue("altText", e)
}

func (m *Media) GetContentType() string {
	ct, _ := m.QueryObject().String("contentType ")
	if ct == "" {
		ct = "image/jpeg"
	}
	return ct
}

func (m *Media) SetContentType(e string) {
	m.SetValue("contentType", e)
}

func (m *Media) UploadFile(name string, file io.ReadCloser) error {
	defer file.Close()
	ext := filepath.Ext(name)

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	headers := http.Header{}
	headers.Add("X-Amz-Acl", "public-read")
	if m.GetContentType() != "" {
		headers.Add("Content-Type", m.GetContentType())
	} else {
		m.SetContentType("binary/octet-stream")
	}

	headers.Add("Cache-Control", "max-age=3153600")

	//lookup ext for file.
	if value, found := MimeTypes[m.GetContentType()]; found {
		checkCorrectExt := false
		//there will be multiple file exts for a single file
		for _, i := range value {
			//if it matches one we know about we are done
			if i == ext {
				checkCorrectExt = true
				break
			}
		}
		//if it doesn't grab the first ext based on the MimeType
		if !checkCorrectExt {
			ext = value[0]
		}
	}

	m.SetUrl(fmt.Sprintf("http://%s/%s/%s%s", S3FQDN, S3Folder, uuid, ext))
	w, err := s3util.Create(m.GetUrl(), headers, nil)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	defer w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *Media) UploadFromParams(files *revel.Params) error {
	for _, v := range files.Files {
		file := v[0]
		r, err := file.Open()
		if err != nil {
			return err
		}

		if value, found := file.Header["Content-Type"]; found {
			m.SetContentType(value[0])
		}

		return m.UploadFile(file.Filename, r) //only going to upload one
	}
	return nil
}

/////
//// Queries
///

//GetMediaByUpdatedAt gets media based on the last updated time. This will error out if there isn't a mapping defined for media yet. That is to say if no media ever made it into the system.
func GetMediumByUpdatedAt() (*Medium, error) {
	toReturn := Medium{}
	query := search.Query().All()
	filter := search.Filter().Exists("valueHash.updatedAt")
	sort := search.Sort("valueHash.updatedAt").Desc()
	results, err := search.Search("media").Type("media").Query(query).Filter(filter).Sort(sort).Size("1000").Result()
	if err != nil {
		return &toReturn, err
	}
	for _, v := range results.Hits.Hits {
		media := NewMedia()
		json.Unmarshal(v.Source, media)
		toReturn = append(toReturn, media)
	}
	return &toReturn, err
}

//GetMediaById this gets the media by id it doesn't check that it exists. you should use CheckMediaExists before calling this so you can handle that error.
func GetMediaById(id string) (*Media, error) {
	toReturn := NewMedia()
	//check cache for new models
	if err := cache.Get(id, &toReturn); err != nil {
		if err := core.GetSource("media", "media", id, toReturn); err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

//CheckMediaExists is a simple wrapper for the core.Exists function that makes it more convient for media.
func CheckMediaExists(id string) (bool, error) {
	exists, _ := core.Exists(false, "media", "media", id)
	//yeah your reading this right im dumping the error because there is a bug in the elasticsearch api
	//i have submitted a bug patch lets see if we can dump this later.
	//https://github.com/mattbaird/elastigo/pull/53
	return exists, nil
}

/////
//// Callbacks
//
func (m *Media) BeforeDelete() error {
	if m.GetUrl() != "" {
		r, err := http.NewRequest("DELETE", m.GetUrl(), nil)
		if err != nil {
			return err
		}
		r.ContentLength = 0
		r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
		r.Header.Set("Content-Type", "text/plain")
		s3util.DefaultConfig.Service.Sign(r, *s3util.DefaultConfig.Keys)
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return err
		}
		if resp.StatusCode != 204 {
			body, _ := ioutil.ReadAll(resp.Body)
			return errors.New(string(body))
		}
	}
	return nil
}

func (m *Media) AfterSave() error {
	return tmpCache(m)
}

func (m *Media) AfterUpdate() error {
	return tmpCache(m)
}
