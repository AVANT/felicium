package models

import (
	"encoding/json"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/AVANT/felicium/model"
	"github.com/jmoiron/jsonq"
	"github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
	"github.com/robfig/revel"
	//"github.com/moonrakr/app/helpers"
	"bufio"
	"bytes"

	"github.com/emilsjolander/goson"
	"github.com/AVANT/felicium/moonrakr/app/lib/results"
	"github.com/robfig/revel/cache"
)

///
//	User types
///

const userTemplate = "users/public/user"

//User is the main User struct that is connected to the db
type User struct {
	*model.Model
}

//NewUser creates a new user model with the db variables hooked in
func NewUser() *User {
	return &User{
		Connection.NewModel("user"),
	}
}

//Users is an array of *Users that implements goson.Collection
type Users []*User

func (u *Users) Len() int                                         { return len(*u) }
func (u *Users) Get(index int) interface{}                        { return (*u)[index] }
func (u *Users) Add(user *User)                                   { (*u) = append((*u), user) }
func (u *Users) RenderElement(i int) (*results.JsonResult, error) { return (*u)[i].Render() }
func (u *Users) Render() (*results.JsonResult, error)             { return results.RenderRenderableCollection(u) }
func (u *Users) FromInterfaceArray(array []interface{}) {
	for i := range array {
		userObject := array[i].(map[string]interface{})
		user := NewUser()
		user.UserMassAssign(&userObject)
		user.SetId(userObject["id"].(string))
		u.Add(user)
	}
}

func (u *User) Render() (*results.JsonResult, error) {
	buffer := new(bytes.Buffer)
	w := bufio.NewWriter(buffer)
	err := goson.RenderTo(w, userTemplate, goson.Args{"User": u})
	if err != nil {
		return new(results.JsonResult), err
	}
	w.Flush()
	return &results.JsonResult{buffer}, nil
}

//This will safely set fields of the post that can be set by mass assignment
func (p *User) UserMassAssign(bulk *map[string]interface{}) (*User, error) {
	jq := jsonq.NewQuery(*bulk)
	for k, _ := range *bulk {
		switch k {
		case "password":
			s, err := jq.String(k)
			if err != nil {
				return p, err
			}
			err = p.SetPassword(s)
			if err != nil {
				return p, err
			}
		case "username":
			s, err := jq.String(k)
			if err != nil {
				return p, err
			}
			p.SetUsername(s)
		case "bio":
			s, err := jq.String(k)
			if err != nil {
				return p, err
			}
			p.SetBio(s)
		case "fullName":
			s, err := jq.String(k)
			if err != nil {
				return p, err
			}
			p.SetFullName(s)
		case "email":
			s, err := jq.String(k)
			if err != nil {
				return p, err
			}
			p.SetEmail(s)
		case "userImage":
			mediaId, err := jq.String(k, "id")
			media := NewMedia()
			media.SetId(mediaId)
			if err != nil {
				return p, err
			}
			p.SetUserImage(media)
		}
	}
	return p, nil
}

////
//	Setters and Getters
////

func (u *User) GetHashedPassword() []byte {
	el, found := u.GetValue("hashedPassword")
	if found {
		return []byte(el.(string))
	} else {
		return []byte{}
	}
}

func (u *User) SetPassword(n string) error {
	pass := []byte(n)
	hp, err := bcrypt.GenerateFromPassword(pass, 0)
	if err != nil {
		return err
	}
	u.SetValue("hashedPassword", string(hp))
	return nil
}

func (u *User) GetBio() string {
	b, _ := u.QueryObject().String("bio")
	return b
}

func (u *User) SetBio(b string) {
	u.SetValue("bio", b)
}

func (u *User) GetUsername() string {
	un, _ := u.QueryObject().String("username")
	return un
}

func (u *User) SetUsername(n string) {
	u.SetValue("username", n)
}

func (u *User) GetFullName() string {
	fn, _ := u.QueryObject().String("fullName")
	return fn
}

func (u *User) SetFullName(n string) {
	u.SetValue("fullName", n)
}

func (u *User) GetUserType() string {
	ut, _ := u.QueryObject().String("userType")
	return ut
}

func (u *User) SetUserType(n string) {
	u.SetValue("userType", n)
}

func (u *User) GetEmail() string {
	e, _ := u.QueryObject().String("email")
	return e
}

func (u *User) SetEmail(n string) {
	u.SetValue("email", n)
}

///
// UserImage
///

func (p *User) GetUserImage() *Media {
	s, _ := p.QueryObject().String("userImage")
	media, _ := GetMediaById(s)
	return media
}

func (u *User) SetUserImage(m *Media) error {
	media, err := GetMediaById(m.GetId())
	if err != err {
		return err
	}
	u.SetValue("userImage", media.GetId())
	return nil
}

///
// Queries
///

//GetPostsByUpdatedAt gets users based on their last updated time. This will error out if there isn't a mapping defined for users yet. That is to say if no user ever made it into the system.
func GetUsersByUpdatedAt() (*Users, error) {
	toReturn := Users{}
	query := search.Query().All()
	filter := search.Filter().Exists("valueHash.updatedAt")
	sort := search.Sort("valueHash.updatedAt").Desc()
	results, err := search.Search("user").Type("user").Query(query).Filter(filter).Sort(sort).Size("1000").Result()
	if err != nil {
		return &toReturn, err
	}
	for _, v := range results.Hits.Hits {
		user := NewUser()
		json.Unmarshal(v.Source, user)
		toReturn = append(toReturn, user)
	}
	return &toReturn, err
}

/*
GetUserByEmailOrUsername gets a user by one of their unique identifiers.
Currently the elastigo api doesn't support nesting boolean filters so we just used or in this case. A more complete request might look more like:
{
   "query": {
      "match_all": {}
   },
   "filter": {
      "and": [
         {
            "or": [
               {
                  "term": {
                     "valueHash.username": "bdenny68"
                  }
               },
               {
                  "term": {
                     "valueHash.email": "bdenny68"
                  }
               }
            ]
         },
         {
            "exists": {
               "field": "valueHash.email"
            }
         },
         {
            "exists": {
               "field": "valueHash.username"
            }
         }
      ]
   }
}
*/
func GetUserByEmailOrUsername(identifier string) (*User, error) {
	revel.ERROR.Printf("%s", identifier)
	toReturn := NewUser()
	query := search.Query().All()
	//existsFilter := search.Filter().Exists("valueHash.username").Exists("valueHash.email")
	usernameFilter := search.Filter().Terms("valueHash.username", identifier)
	emailFilter := search.Filter().Terms("valueHash.email", identifier)
	results, err := search.Search("user").Type("user").Query(query).Filter("or", usernameFilter, emailFilter).Size("1").Result()
	if err != nil {
		return toReturn, err
	}
	switch {
	case results.Hits.Total == 0:
		return toReturn, nil
	case results.Hits.Total > 1:
		revel.ERROR.Println("multiple users are being matched in the find user by email or username function")
	}
	json.Unmarshal(results.Hits.Hits[0].Source, toReturn)
	return toReturn, err
}

//GetUserById this gets the user by id it doesn't check that it exists. You should use CheckUserExists before calling this so you can handle that error.
func GetUserById(id string) (*User, error) {
	toReturn := NewUser()
	//check cache for new models
	if err := cache.Get(id, &toReturn); err != nil {
		if err := core.GetSource("user", "user", id, toReturn); err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

//CheckUserExists is a simple wrapper for the core.Exists function that makes it more convient for users.
func CheckUserExists(id string) (bool, error) {
	exists, _ := core.Exists(false, "user", "user", id)
	//yeah your reading this right im dumping the error because there is a bug in the elasticsearch api
	//i have submitted a bug patch lets see if we can dump this later.
	//https://github.com/mattbaird/elastigo/pull/53
	return exists, nil
}

/////
//// Callbacks
//

func (u *User) AfterSave() error {
	return tmpCache(u)
}

func (u *User) AfterUpdate() error {
	return tmpCache(u)
}
