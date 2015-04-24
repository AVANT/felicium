package tests

import (
	"bytes"
)

func (t AppTest) TestThatSessionLogin() {
	buf := bytes.NewBufferString(`{"user":"ledwards", "password":"ifdecWiWyd9"}`)
	t.Post("/login", "application/json", buf)
	t.AssertOk()
	t.AssertContentType("application/json")
}

func (t AppTest) TestThatSessionLogout() {
	t.Get("/logout")
	t.AssertOk()
	t.AssertContentType("text/html")
}
