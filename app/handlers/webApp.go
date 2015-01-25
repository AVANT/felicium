package handlers

import (
	"fmt"
	"net/http"

	"github.com/avant/felicium/app"
	"github.com/avant/felicium/config"
	"github.com/avant/felicium/middleware"
	"github.com/avant/felicium/routes"

	"github.com/avant/felicium/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/avant/felicium/Godeps/_workspace/src/github.com/phyber/negroni-gzip/gzip"
)

func determineCompressionLevel(desiredLevel interface{}) int {
	switch desiredLevel {
	case "BestCompression", 9:
		return gzip.BestCompression
	case "BestSpeed", 1:
		return gzip.BestSpeed
	case "NoCompression", 0:
		return gzip.NoCompression
	}

	//case "DefaultCompression", -1:
	return gzip.DefaultCompression

}

// want to make the distinction here so that application doesn't need to boot
// web framework for command line utils when those happen
type WebApp struct {
	*app.App
	negroni *negroni.Negroni
}

func NewWebApp(config config.Configurator) (*WebApp, error) {
	wa := &WebApp{}
	application, err := app.NewApp(config)
	if err != nil {
		return nil, err
	}
	wa.App = application

	staticDir, err := config.LookupWithDefault("static_dir", "public")
	if err != nil {
		return nil, err
	}
	compressionLevel, err := config.LookupWithDefault("compression_level", -1)
	if err != nil {
		return nil, err
	}

	wa.negroni = negroni.New()
	// this is given its own logger by default im going to line it up with our applications logger
	recovery := negroni.NewRecovery()
	recovery.Logger = wa.App.ServerLog

	wa.negroni.Use(recovery)
	wa.negroni.Use(middleware.NewLogger(wa.App.AccessLog))
	wa.negroni.Use(gzip.Gzip(determineCompressionLevel(compressionLevel)))
	wa.negroni.Use(negroni.NewStatic(http.Dir(staticDir.(string))))
	router := routes.NewRouter()
	wa.negroni.UseHandler(router)
	return wa, nil
}

func (wa *WebApp) Run() {
	webPort := wa.App.Config.LookupWithDefaultOrPanic("web_port", "8080")
	bindAddress := wa.App.Config.LookupWithDefaultOrPanic("bind_address", "0.0.0.0")
	// negroni has its own run function but it uses its own logger so we assimilate it
	http.ListenAndServe(fmt.Sprintf("%v:%v", bindAddress, webPort), wa.negroni)
}
