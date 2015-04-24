package seeder

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/lib/helpers"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/jmoiron/jsonq"
	"github.com/robfig/revel"
)

type SeedData struct {
	Medium models.Medium
	Posts  models.Posts
	Users  models.Users
}

func (s *SeedData) AddMedia(m *models.Media) {
	s.Medium.Add(m)
}
func (s *SeedData) AddUser(m *models.User) {
	s.Users.Add(m)
}
func (s *SeedData) AddPost(m *models.Post) {
	s.Posts.Add(m)
}
func (s *SeedData) ConvertIntsToIds(_type string, i []int) []interface{} {
	toReturn := make([]interface{}, len(i))
	for j := range i {
		switch _type {
		case "users":
			toReturn[j] = interface{}(map[string]interface{}{"id": (*s).Users[i[j]].GetId()})
		case "medium":
			toReturn[j] = interface{}(map[string]interface{}{"id": (*s).Medium[i[j]].GetId()})
		case "posts":
			toReturn[j] = interface{}(map[string]interface{}{"id": (*s).Posts[i[j]].GetId()})
		}
	}
	return toReturn
}

func SeedFromJson(jsonFile string) (*SeedData, error) {
	toReturn := new(SeedData)
	fileData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return toReturn, err
	}
	seedData := new(map[string]interface{})
	err = json.Unmarshal(fileData, seedData)
	if err != nil {
		return toReturn, err
	}
	jq := jsonq.NewQuery(*seedData)
	users, err := jq.Array("users")
	if err != nil {
		return toReturn, err
	}
	for i := range users {
		userObject, err := jq.Object("users", strconv.Itoa(i))
		if err != nil {
			return toReturn, err
		}
		user := models.NewUser()
		user, err = user.UserMassAssign(&userObject)
		if err != nil {
			return toReturn, err
		}
		userType, err := jq.String("users", strconv.Itoa(i), "seedData", "userType")
		if err != nil {
			return toReturn, err
		}
		user.SetUserType(userType)
		err = model.Save(user)
		if err != nil {
			return toReturn, err
		}
		toReturn.AddUser(user)
	}
	revel.INFO.Println("Seed Users created")
	medium, err := jq.Array("media")
	if err != nil {
		return toReturn, err
	}
	for i := range medium {
		mediaObject, err := jq.Object("media", strconv.Itoa(i))
		if err != nil {
			return toReturn, err
		}
		media := models.NewMedia()
		media, err = media.MediaMassAssign(&mediaObject)
		if err != nil {
			return toReturn, err
		}
		imageUrl, err := jq.String("media", strconv.Itoa(i), "seedData", "url")
		if err != nil {
			return toReturn, err
		}
		uploadedBy, err := jq.Int("media", strconv.Itoa(i), "seedData", "uploadedBy")
		if err != nil {
			return toReturn, err
		}

		media.SetUploadedBy(toReturn.Users[uploadedBy].GetId())
		resp, err := http.Get(imageUrl)
		defer resp.Body.Close()
		if err != nil {
			return toReturn, err
		}

		media.UploadFile("cat.jpg", resp.Body)
		err = model.Save(media)
		if err != nil {
			return toReturn, err
		}

		// //We alread did users at this point so we are going to associate them here.
		belongsToTypeObject, err := jq.Object("media", strconv.Itoa(i), "seedData", "belongsTo")
		if err != nil {
			return toReturn, err
		}
		var belongsToType string
		for k, _ := range belongsToTypeObject {
			belongsToType = k
			break
		}
		if belongsToType == "user" {
			belongsTo, err := jq.Int("media", strconv.Itoa(i), "seedData", "belongsTo", belongsToType)
			if err != nil {
				return toReturn, err
			}
			err = toReturn.Users[belongsTo].SetUserImage(media)
			if err != nil {
				return toReturn, err
			}
			err = model.Save(toReturn.Users[belongsTo])
			if err != nil {
				return toReturn, err
			}
		}
		toReturn.AddMedia(media)
	}
	revel.INFO.Println("Media Created")
	posts, err := jq.Array("posts")
	if err != nil {
		return toReturn, err
	}
	for i := range posts {
		postObject, err := jq.Object("posts", strconv.Itoa(i))
		if err != nil {
			return toReturn, err
		}
		post := models.NewPost()
		//look up the authors and add them to the bulk object
		authorSeedIds, err := helpers.JsonqArrayOfInts(jq, "posts", strconv.Itoa(i), "seedData", "authors")
		if err != nil {
			return toReturn, err
		}
		postObject["authors"] = toReturn.ConvertIntsToIds("users", authorSeedIds)
		relatedMediaSeedIds, err := helpers.JsonqArrayOfInts(jq, "posts", strconv.Itoa(i), "seedData", "authors")
		if err != nil {
			return toReturn, err
		}
		postObject["media"] = toReturn.ConvertIntsToIds("medium", relatedMediaSeedIds)
		headerImageSeedId, err := jq.Int("posts", strconv.Itoa(i), "seedData", "headerImage")
		postObject["headerImage"] = toReturn.ConvertIntsToIds("medium", []int{headerImageSeedId})[0]
		post, err = post.PostMassAssign(&postObject)
		if err != nil {
			return toReturn, err
		}
		err = model.Save(post)
		if err != nil {
			return toReturn, err
		}
		toReturn.AddPost(post)
	}

	return &SeedData{}, nil
}

func SeedRandomData() (*SeedData, error) {

	return &SeedData{}, nil
}

// func CreateMedia() (m *models.Media error) {
// 	file, err := os.Open(path.Join(revel.BasePath, "tests/images/test-1600.jpg"))
// 	if err != nil {
// 		t.Assertf(false, "%s", err)
// 	}
// 	defer file.Close()
// 	for i := 0; i < testMediaPerPost*testPosts; i++ {
// 		media := models.NewMedia()
// 		media.UploadFile("test-1600.jpg", file)
// 		id, _, _ := model.Save(media)
// 		media.SetId(id)
// 		TestMedium.Add(media)
// 	}
// }
