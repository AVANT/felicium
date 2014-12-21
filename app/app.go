package app

import (
	"io"
	"log"
	"os"

	"github.com/vvvntdotorg/felicium/config"

	_ "github.com/vvvntdotorg/felicium/Godeps/_workspace/src/github.com/jmcvetta/neoism"
)

type App struct {
	Config    config.Configurator
	ServerLog *log.Logger
	AccessLog *log.Logger
}

func NewApp(config config.Configurator) (*App, error) {
	wc := &App{}
	wc.Config = config

	serverLogLocation, err := config.LookupWithDefault("serverLog", "-")
	if err != nil {
		return nil, err
	}
	accessLogLocation, err := config.LookupWithDefault("accessLog", "-")
	if err != nil {
		return nil, err
	}

	traceTo := func(s interface{}) (io.Writer, error) {
		if s == "-" {
			return os.Stdout, nil
		}

		f, err := os.OpenFile(s.(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return os.Stdout, err
		}
		return f, nil
	}
	accessLogOut, err := traceTo(accessLogLocation)
	if err != nil {
		return nil, err
	}

	wc.AccessLog = log.New(
		accessLogOut,
		"",
		0,
	)

	serverLogOut, err := traceTo(serverLogLocation)
	if err != nil {
		return nil, err
	}

	wc.ServerLog = log.New(
		serverLogOut,
		"",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	return wc, nil
}
