package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/models"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// newfileUploadRequest creates a new file upload http request with optional extra params
// borrowed from https://gist.github.com/mattetti/5914158
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	//part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(paramName), quoteEscaper.Replace(filepath.Base(path))))
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	toReturn, err := http.NewRequest("POST", uri, body)
	toReturn.Header.Set("Content-Type", writer.FormDataContentType())
	return toReturn, err
}

func (t AppTest) TestDestroyAll() {
	posts, err := models.GetPostsByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	for _, v := range *posts {
		model.Delete(v)
	}

	users, err := models.GetUsersByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	for _, v := range *users {
		model.Delete(v)
	}

	medium, err := models.GetMediumByUpdatedAt()
	if err != nil {
		t.Assertf(false, "%s", err)
	}
	for _, v := range *medium {
		model.Delete(v)
	}
}
