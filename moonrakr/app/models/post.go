package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/emilsjolander/goson"
	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/lib/helpers"
	"github.com/AVANT/felicium/moonrakr/app/lib/results"
	"github.com/jmoiron/jsonq"
	"github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
)

/////
////  Post type and related constructs
///

const postTemplate = "posts/public/post"

//Post is the main Post struct that is connected to the db
type Post struct {
	*model.Model
}

//NewPost creates a new Post struct and initalizes it correctly
func NewPost() *Post {
	return &Post{
		Connection.NewModel("post"),
	}
}

//Post is a collection of *Post it implements goson.Collection
type Posts []*Post

func (p *Posts) Len() int                                         { return len(*p) }
func (p *Posts) Get(index int) interface{}                        { return (*p)[index] }
func (p *Posts) Add(post *Post)                                   { (*p) = append((*p), post) }
func (p *Posts) RenderElement(i int) (*results.JsonResult, error) { return (*p)[i].Render() }
func (p *Posts) Render() (*results.JsonResult, error)             { return results.RenderRenderableCollection(p) }
func (p *Posts) FromInterfaceArray(array []interface{}) {
	for i := range array {
		postObject := array[i].(map[string]interface{})
		post := NewPost()
		post.PostMassAssign(&postObject)
		post.SetId(postObject["id"].(string))
		p.Add(post)
	}
}

func (p *Post) Render() (*results.JsonResult, error) {
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	err := goson.RenderTo(w, postTemplate, goson.Args{"Post": p})
	if err != nil {
		return new(results.JsonResult), err
	}
	w.Flush()
	return &results.JsonResult{buffer}, nil
}

//This will safely set fields of the post that can be set by mass assignment
func (p *Post) PostMassAssign(bulk *map[string]interface{}) (*Post, error) {
	jq := jsonq.NewQuery(*bulk)
	for k := range *bulk {
		switch k {
		case "excerpt":
			s, err := jq.String(k)
			if err == nil {
				p.SetExcerpt(s)
			}
		case "body":
			s, err := jq.String(k)
			if err == nil {
				p.SetBody(s)
			}
		case "title":
			s, err := jq.String(k)
			if err == nil {
				p.SetTitle(s)
			}
		case "status":
			s, err := jq.String(k)
			if err == nil {
				p.SetStatus(s)
			}
		case "slug":
			s, err := jq.String(k)
			if err == nil {
				p.SetSlug(s)
			}
		case "moduleEnabled":
			s, err := jq.Bool(k)
			if err == nil {
				p.SetModuleEnabled(s)
			}
		case "javascriptData":
			s, err := jq.String(k)
			if err == nil {
				p.SetJavascriptData(s)
			}
		case "cssData":
			s, err := jq.String(k)
			if err == nil {
				p.SetCssData(s)
			}
		case "htmlData":
			s, err := jq.String(k)
			if err == nil {
				p.SetHtmlData(s)
			}
		case "headerImage":
			headerImageId, err := jq.String(k, "id")
			headerImage := NewMedia()
			headerImage.SetId(headerImageId)
			if err == nil {
				p.SetHeaderImage(headerImage)
			}
		case "meta":
			array, err := helpers.JsonqArrayOfStrings(jq, k)
			if err == nil {
				p.SetMeta(array)
			}
		case "tags":
			array, err := helpers.JsonqArrayOfStrings(jq, k)
			if err == nil {
				p.SetTags(array)
			}
		case "media":
			mediaObject, err := jq.Array(k)
			if err == nil {
				media := Medium{}
				for i := range mediaObject {
					id, err := jq.String(k, strconv.Itoa(i), "id")
					if err == nil {
						m := NewMedia()
						m.SetId(id)
						media.Add(m)
					}
				}
				p.SetMedia(media)
			}
		case "authorsArray":
			array, err := helpers.JsonqArrayOfStrings(jq, k)
			if err == nil {
				p.SetAuthorsArray(array)
			}
		case "authors":
			authorsObject, err := jq.Array(k)
			if err == nil {
				authors := Users{}
				for i := range authorsObject {
					id, err := jq.String(k, strconv.Itoa(i), "id")
					if err == nil {
						a := NewUser()
						a.SetId(id)
						authors.Add(a)
					}
				}
				p.SetAuthors(authors)
			}
		}
	}
	return p, nil
}

/////
//// Setters and getters for posts
///

///
// Excerpt
///

//this was a []byte but the seed data isn't right so moving on for now. Same with the body
func (p *Post) GetExcerpt() string {
	e, _ := p.QueryObject().String("excerpt")
	return e
}

func (p *Post) SetExcerpt(e string) {
	p.SetValue("excerpt", e)
}

///
// Body
///

func (p *Post) GetBody() string {
	b, _ := p.QueryObject().String("body")
	return b
}

func (p *Post) SetBody(b string) {
	p.SetValue("body", b)
}

///
// Status
///

func (p *Post) GetStatus() string {
	s, _ := p.QueryObject().String("status")
	return s
}

func (p *Post) SetStatus(t string) {
	if p.GetStatus() != t {
		p.SetPublishedAt(time.Now())
		p.SetValue("status", t)
	}
}

///
// Promotion
///

func (p *Post) GetPromtion() float64 {
	pro, _ := p.QueryObject().Float("promotion")
	return pro
}

func (p *Post) SetPromotion(t float64) {
	p.SetValue("promotion", t)
}

///
// HeaderImage
///

func (p *Post) GetHeaderImage() *Media {
	s, _ := p.QueryObject().String("headerImage")
	media, _ := GetMediaById(s)
	return media
}

func (p *Post) SetHeaderImage(m *Media) error {
	media, err := GetMediaById(m.GetId())
	if err != err {
		return err
	}
	p.SetValue("headerImage", media.GetId())
	return nil
}

///
// RelatedMedia
///

func (p *Post) GetMedia() Medium {
	var toReturn Medium
	sa, _ := p.QueryObject().ArrayOfStrings("media")
	for _, v := range sa {
		media, err := GetMediaById(v)
		if err == nil {
			toReturn.Add(media)
		}
	}
	return toReturn
}

func (p *Post) SetMedia(m Medium) error {
	var toSet []interface{}
	for _, v := range m {
		media, err := GetMediaById(v.GetId())
		if err != err {
			return err
		}
		toSet = append(toSet, media.GetId())
	}
	p.SetValue("media", toSet)
	return nil
}

///
// Tags
///

func (p *Post) GetTags() []string {
	tags, _ := p.QueryObject().ArrayOfStrings("tags")
	return tags
}

func (p *Post) SetTags(a []string) {
	p.SetValue("tags", a)
}

///
// Meta
///

func (p *Post) GetMeta() []string {
	meta, _ := p.QueryObject().ArrayOfStrings("meta")
	return meta
}

func (p *Post) SetMeta(a []string) {
	p.SetValue("meta", a)
}

///
// Authors
///

//nasty hack for front end because the relationships are not yet handled there.
func (p *Post) GetAuthorsArray() []string {
	saa, _ := p.QueryObject().ArrayOfStrings("authorsArray")
	return saa
}

func (p *Post) SetAuthorsArray(array []string) error {
	p.SetValue("authorsArray", array)
	return nil
}

func (p *Post) GetAuthors() Users {
	var toReturn Users
	sa, _ := p.QueryObject().ArrayOfStrings("authors")
	for _, v := range sa {
		user, err := GetUserById(v)
		if err == nil {
			toReturn.Add(user)
		}
	}
	return toReturn
}

func (p *Post) SetAuthors(m Users) error {
	var toSet []interface{}
	for _, v := range m {
		user, err := GetUserById(v.GetId())
		if err != err {
			return err
		}
		toSet = append(toSet, user.GetId())
	}
	revel.INFO.Printf("%T", toSet)
	p.SetValue("authors", toSet)
	return nil
}

///
// Title
///

func (p *Post) GetTitle() string {
	s, _ := p.QueryObject().String("title")
	return s
}

func (p *Post) SetTitle(t string) {
	p.SetValue("title", t)
}

///
// Slug
///

func (p *Post) GetSlug() string {
	s, _ := p.QueryObject().String("slug")
	return s
}

func (p *Post) SetSlug(t string) {
	t = revel.Slug(t)
	p.SetValue("slug", t)
}

///
// javascriptData
///

func (p *Post) GetJavascriptData() string {
	s, _ := p.QueryObject().String("javascriptData")
	return s
}

func (p *Post) SetJavascriptData(t string) {
	p.SetValue("javascriptData", t)
}

///
// cssData
///

func (p *Post) GetCssData() string {
	s, _ := p.QueryObject().String("cssData")
	return s
}

func (p *Post) SetCssData(t string) {
	p.SetValue("cssData", t)
}

///
// htmlData
///

func (p *Post) GetHtmlData() string {
	s, _ := p.QueryObject().String("htmlData")
	return s
}

func (p *Post) SetHtmlData(t string) {
	p.SetValue("htmlData", t)
}

///
// moduleEnabled
///

func (p *Post) GetModuleEnabled() bool {
	s, _ := p.QueryObject().Bool("moduleEnabled")
	return s
}

func (p *Post) SetModuleEnabled(t bool) {
	p.SetValue("moduleEnabled", t)
}

///
// Published At
///

func (p *Post) GetPublishedAt() time.Time {
	el, found := p.GetValue("publishedAt")
	if found {
		switch el.(type) {
		case time.Time:
			return el.(time.Time)
		case string:
			if t, err := time.Parse(time.RFC3339Nano, el.(string)); err == nil {
				return t
			} else {
				return time.Time{}
			}
		default:
			return time.Time{}
		}
	} else {
		return time.Time{}
	}
}

func (p *Post) SetPublishedAt(t time.Time) {
	p.SetValue("publishedAt", t)
}

/////
//// Queries
///

//GetPostsByUpdatedAt gets posts based on the last updated time. This will error out if there isn't a mapping defined for posts yet. That is to say if no post ever made it into the system.
func GetPostsByUpdatedAt() (*Posts, error) {
	toReturn := Posts{}
	query := search.Query().All()
	filter := search.Filter().Exists("valueHash.updateAt")
	sort := search.Sort("valueHash.updatedAt").Desc()
	results, err := search.Search("post").Type("post").Query(query).Filter(filter).Sort(sort).Size("1000").Result()
	if err != nil {
		return &toReturn, err
	}
	for _, v := range results.Hits.Hits {
		post := NewPost()
		json.Unmarshal(v.Source, post)
		toReturn = append(toReturn, post)
	}
	return &toReturn, err
}

//GetPostsByUpdatedAt gets posts based on the last updated time. This will error out if there isn't a mapping defined for posts yet. That is to say if no post ever made it into the system.
func GetPostsByCreatedAt() (*Posts, error) {
	toReturn := Posts{}
	query := search.Query().All()
	filter := search.Filter().Exists("valueHash.createdAt")
	sort := search.Sort("valueHash.createdAt").Desc()
	results, err := search.Search("post").Type("post").Query(query).Filter(filter).Sort(sort).Size("1000").Result()
	if err != nil {
		return &toReturn, err
	}
	for _, v := range results.Hits.Hits {
		post := NewPost()
		json.Unmarshal(v.Source, post)
		toReturn = append(toReturn, post)
	}
	return &toReturn, err
}

//GetPostsByUpdatedAt gets posts based on the last updated time. This will error out if there isn't a mapping defined for posts yet. That is to say if no post ever made it into the system.
func GetPostsByStatus(s string) (*Posts, error) {

	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]string{},
		},
		"sort": map[string]string{
			"valueHash.publishedAt": "desc",
		},
		"filter": map[string]interface{}{
			"and": []map[string]interface{}{
				map[string]interface{}{
					"exists": map[string]interface{}{
						"field": "valueHash.status",
					},
				},
				map[string]interface{}{
					"exists": map[string]interface{}{
						"field": "valueHash.publishedAt",
					},
				},
				map[string]interface{}{
					"term": map[string]interface{}{
						"valueHash.status": s,
					},
				},
			},
		},
		"size": 1000,
	}
	toReturn := Posts{}
	results, err := core.SearchRequest(false, "post", "post", qry, "", 0)
	if err != nil {
		return &toReturn, err
	}
	for _, v := range results.Hits.Hits {
		post := NewPost()
		json.Unmarshal(v.Source, post)
		toReturn = append(toReturn, post)
	}
	return &toReturn, err
}

//GetPostsById this gets the post by id it doesn't check that it exists. you should use CheckPostExists before calling this so you can handle that error.
func GetPostById(id string) (*Post, error) {
	toReturn := NewPost()
	//check cache for new models
	if err := cache.Get(id, &toReturn); err != nil {
		if err := core.GetSource("post", "post", id, toReturn); err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

//CheckPostExists is a simple wrapper for the core.Exists function that makes it more convient for posts.
func CheckPostExists(id string) (bool, error) {
	exists, _ := core.Exists(false, "post", "post", id)
	//yeah your reading this right im dumping the error because there is a bug in the elasticsearch api
	//i have submitted a bug patch lets see if we can dump this later.
	//https://github.com/mattbaird/elastigo/pull/53
	return exists, nil
}

func GetPostBySlug(s string) (*Post, error) {
	toReturn := NewPost()

	qry := map[string]interface{}{
		"filter": map[string]interface{}{
			"script": map[string]string{
				"script": fmt.Sprintf(`if ( _source['valueHash'] != null && _source.valueHash['slug'] != null ) { return _source.valueHash.slug == %s;}`, strconv.Quote(s))},
		},
		"size": 1,
	}

	results, err := core.SearchRequest(false, "post", "post", qry, "", 0)
	if err != nil {
		return toReturn, err
	}
	if len(results.Hits.Hits) < 1 {
		return toReturn, errors.New("Not Found")
	}
	if len(results.Hits.Hits) > 1 {
		revel.ERROR.Printf("collision on the slug %s", s)
	}
	err = json.Unmarshal(results.Hits.Hits[0].Source, toReturn)
	if err != nil {
		return toReturn, err
	}
	return toReturn, nil
}

/////
//// Callbacks
//

func (p *Post) AfterSave() error {
	return tmpCache(p)
}

func (p *Post) AfterUpdate() error {
	return tmpCache(p)
}

func (p *Post) BeforeSave() error {
	if p.GetSlug() == "" {
		if p.GetTitle() == "" {
			return errors.New("Title is cannot be blank.")
		}
		p.SetSlug(p.GetTitle())
	}
	if p.GetStatus() == "" {
		p.SetStatus("unpublished")
	}
	return nil
}
