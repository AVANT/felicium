package tests

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/AVANT/felicium/moonrakr/app/models"
)

//TestUserIndex checks that the user index doesn't return an error code
func (t AppTest) TestUserIndex() {
	t.Get(userPrefix)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestUserCreation test that user creation works as expected.
func (t AppTest) TestUserCreation() {
	buf := bytes.NewBufferString(`{
                        "fullName": "Billy Dennis",
                        "email" : "test1@email.com",
                        "username": "bdenny68",
                        "userType": "user",
                        "password": "testpassword",
                        "comments": 45,
                        "promotedComments": 30,
                        "bio":"Infuriatingly humble zombie fan. Food fanatic. Pop culture specialist. Hardcore analyst.",
                        "posts": null,
                        "recommendations": null,
                        "latestPost": null,
                        "heroImage": {
                                "id": "69344e464f7e2d2231417b4a43573922",
                                "altText": "Crawlers and the blind read this.",
                                "title": "Awesome Image 1",
                                "sizes": {
                                        "mobile": "http://dummyimage.com/300",
                                        "desktop": "http://dummyimage.com/300"
                                }
                        }
                }`)
	t.Post(userPrefix, "application/json", buf)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestUserShow checks to see if the most recently updated user can be shown with a successful error code
func (t AppTest) TestUserShow() {
	users, err := models.GetUsersByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if users.Len() == 0 {
		t.Assertf(false, "%s", "No Users exist cannot check the show method.")
	} else {
		id := (*users)[0].GetId()
		t.Get(fmt.Sprintf("%s/%s", userPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
	}
}

//TestUserUpdate this takes the last updated user and updates it.
func (t AppTest) TestUserUpdate() {
	users, err := models.GetUsersByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if users.Len() == 0 {
		t.Assertf(false, "%s", "No Users exist cannot check the show method.")
	} else {
		user := (*users)[0]
		newBio := "Not what the bio was."
		buf := bytes.NewBufferString(fmt.Sprintf(`{"bio": "%s"}`, newBio))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", t.BaseUrl(), userPrefix, user.GetId()), buf)
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		req.Header.Add("Content-Type", "application/json")
		t.MakeRequest(req)
		t.AssertOk()
		t.AssertContentType("application/json")
		updatedUser, err := models.GetUserById(user.GetId())
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if updatedUser.GetBio() != newBio {
			t.Assertf(false, "%s != %s", updatedUser.GetBio(), newBio)
		}
	}
}

//TestUserDelete this will delete a user and check the error code if there is one.
func (t AppTest) TestUserDelete() {
	users, err := models.GetUsersByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if users.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		id := (*users)[0].GetId()
		t.Delete(fmt.Sprintf("%s/%s", userPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
		found, err := models.CheckUserExists(id)
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if found {
			t.Assertf(false, "%s", "Delete didn't remove the post")
		}
	}
}
