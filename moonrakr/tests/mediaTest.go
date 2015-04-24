package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"path"

	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

//TestPostIndex checks that the post index don't return an error code
func (t AppTest) TestMediaIndex() {
	t.Get(postPrefix)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestMediaCreation uploads a media object using the standard multipart form
func (t AppTest) TestMediaCreation() {

	extraParams := map[string]string{
		"title":     "Test Document",
		"altText":   "Test image",
		"belongsTo": "A document with all the Go programming language secrets",
	}
	url := fmt.Sprintf("%s%s", t.BaseUrl(), mediaPrefix)
	//incase you want to see it in nc or something
	//url := fmt.Sprint("http://localhost:8000/media")
	request, err := newfileUploadRequest(url, extraParams, "file", path.Join(revel.BasePath, "tests/images/test-1600.jpg"))
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	t.MakeRequest(request)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestMediaShow shows that media is correctly displayed
func (t AppTest) TestMediaShow() {
	medium, err := models.GetMediumByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if medium.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		id := (*medium)[0].GetId()
		t.Get(fmt.Sprintf("%s/%s", mediaPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
	}
}

//TestMediaUpdate this takes the last updated media and updates it.
func (t AppTest) TestMediaUpdate() {
	medium, err := models.GetMediumByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if medium.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		media := (*medium)[0]
		newTitle := "Not what the title was."
		buf := bytes.NewBufferString(fmt.Sprintf(`{"title": "%s"}`, newTitle))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", t.BaseUrl(), mediaPrefix, media.GetId()), buf)
		if err != nil {

			t.Assertf(false, "%s", err)
		}
		req.Header.Add("Content-Type", "application/json")
		t.MakeRequest(req)
		t.AssertOk()
		t.AssertContentType("application/json")
		updatedMedia, err := models.GetMediaById(media.GetId())
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if updatedMedia.GetTitle() != newTitle {
			t.Assertf(false, "%s != %s", updatedMedia.GetTitle(), newTitle)
		}
	}
}

//TestMediaDelete this will delete media and check the error code if there is one.
func (t AppTest) TestMediaDelete() {
	medium, err := models.GetMediumByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if medium.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		id := (*medium)[0].GetId()
		t.Delete(fmt.Sprintf("%s/%s", mediaPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
		found, err := models.CheckMediaExists(id)
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if found {
			t.Assertf(false, "%s", "Delete didn't remove the post")
		}
	}
}
