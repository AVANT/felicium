package tests

import (
	"github.com/robfig/revel"
	//"github.com/moonrakr/app/lib/seeder"
	"github.com/AVANT/felicium/moonrakr/app/models"
	//"path"
)

const (
	postPrefix  = "/posts"
	mediaPrefix = "/media"
	userPrefix  = "/users"
)

const (
	testPosts        = 4
	testMediaPerPost = 4
	testUsers        = 4
)

var (
	TestMedium models.Medium
	TestPosts  models.Posts
	TestUsers  models.Users
)

type AppTest struct {
	revel.TestSuite
}

func (t *AppTest) Before() {
	//_, err := seeder.SeedFromJson(path.Join(revel.BasePath, "app/db/seeds/test_data.v1.json"))
	//if err != nil {
	//	revel.ERROR.Println(err)
	//}
	//generateMedium(t)
	//generatePosts(t)
	//generateUsers(t)
	revel.INFO.Println("Set up")
}

func (t *AppTest) After() {
	//cleanMedium(t)
	//cleanPosts(t)
	//cleanUsers(t)
	revel.INFO.Println("Tear down")
}
