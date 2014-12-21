package middleware

import (
	"bytes"
	"github.com/vvvntdotorg/felicium/Godeps/_workspace/src/github.com/gorilla/handlers"
	"log"
	"net/http"
)

type Logger struct {
	*log.Logger
}

func NewLogger(out *log.Logger) *Logger {
	return &Logger{out}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	b := &bytes.Buffer{}
	loggedNext := handlers.CombinedLoggingHandler(b, newWrapper(next))

	loggedNext.ServeHTTP(rw, r)
	l.Print(b)
}

// needed for the conversion between http.HandlerFunc and http.Handler
type wrapper struct {
	handler http.HandlerFunc
}

func newWrapper(h http.HandlerFunc) *wrapper {
	return &wrapper{h}
}

// satify the ServerHTTP interface
func (w *wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.handler(rw, r)
}
