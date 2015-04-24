package tests

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/models"
)

func generatePosts(t *AppTest) {
	for i := 0; i < testPosts; i++ {
		post := models.NewPost()
		post.SetTitle(fmt.Sprintf("Title for test post %d", i))
		post.SetTitle(fmt.Sprintf("Title for test post %d", i))
		post.SetExcerpt(fmt.Sprint("Excerpt for test post %d", i))
		var body bytes.Buffer
		relatedMedia := []string{}
		body.WriteString(fmt.Sprint("This is the Body for post %d.", i))
		for j := 0; j < testMediaPerPost; j++ {
			media := TestMedium[i*testMediaPerPost+j]
			relatedMedia = append(relatedMedia, media.GetId())
			body.WriteString(fmt.Sprint(` This post has some associated media. <img src="%s" alt="%s" \> Some more text after it.`, media.GetUrl(), media.GetAltText()))
		}
		post.SetBody(body.String())
		post.SetTags([]string{"shared tag 1", "shared tag 2", fmt.Sprintf("unique tag %s", i)})
		//post.SetHeaderImage(TestMedium[i*testMediaPerPost].GetId())
		//post.SetRelatedMedia(relatedMedia)
		_ = model.Save(post)
		TestPosts.Add(post)
	}
}

func cleanPosts(t *AppTest) {
	for _, v := range TestPosts {
		model.Delete(v)
	}
}

//TestPostIndex checks that the post index don't return an error code
func (t AppTest) TestPostIndex() {
	t.Get(postPrefix)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestPostCreation test that post creation works as expected.
func (t AppTest) TestPostCreation() {
	buf := bytes.NewBufferString(`{
                        "title": "‘We Don’t Need No Thought Control’: Self-authorization, Bricolage, and Garageband as Critical Pedagogy",
                        "excerpt": "Over the 1957-58 Christmas holidays, the three schools were relocated",
                        "body": "This is a test post",
                        "authors": [
                                {"id": "23b822ba4cef65f230c9e49a05491588"}
                        ],
                        "promotion": 1234567,
                        "promotedComment": "487a28fag80ae1c7fb4729bf0c6e0cf2",
                        "recommendations": 34,
                        "tags": [
                                "Dubstep",
                                "Art"
                        ]
                }`)
	t.Post(postPrefix, "application/json", buf)
	t.AssertOk()
	t.AssertContentType("application/json")
}

//TestPostShow checks to see if the most recently updated post can be shown with a successful error code
func (t AppTest) TestPostShow() {
	posts, err := models.GetPostsByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if posts.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		id := (*posts)[0].GetId()
		t.Get(fmt.Sprintf("%s/%s", postPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
	}
}

//TestPostUpdate this takes the last updated post and updates it.
func (t AppTest) TestPostUpdate() {
	posts, err := models.GetPostsByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if posts.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		post := (*posts)[0]
		newTitle := "Not what the title was."
		buf := bytes.NewBufferString(fmt.Sprintf(`{"title": "%s"}`, newTitle))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", t.BaseUrl(), postPrefix, post.GetId()), buf)
		if err != nil {

			t.Assertf(false, "%s", err)
		}
		req.Header.Add("Content-Type", "application/json")
		t.MakeRequest(req)
		t.AssertOk()
		t.AssertContentType("application/json")
		updatedPost, err := models.GetPostById(post.GetId())
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if updatedPost.GetTitle() != newTitle {
			t.Assertf(false, "%s != %s", updatedPost.GetTitle(), newTitle)
		}
	}
}

//TestPostDelete this will delete a post and check the error code if there is one.
func (t AppTest) TestPostDelete() {
	posts, err := models.GetPostsByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	if posts.Len() == 0 {
		t.Assertf(false, "%s", "No Posts exist cannot check the show method.")
	} else {
		id := (*posts)[0].GetId()
		t.Delete(fmt.Sprintf("%s/%s", postPrefix, id))
		t.AssertOk()
		t.AssertContentType("application/json")
		found, err := models.CheckPostExists(id)
		if err != nil {
			t.Assertf(false, "%s", err)
		}
		if found {
			t.Assertf(false, "%s", "Delete didn't remove the post")
		}
	}
}
