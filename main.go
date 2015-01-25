package main

import (
	"flag"

	"github.com/avantdotorg/felicium/app/handlers"
	"github.com/avantdotorg/felicium/config/yamlConfig"
)

func main() {

	env := flag.String("env", "development", "The config env to read")
	configLocation := flag.String("config", "config.yml", "Path to the config file to read from")

	flag.Parse()

	config, err := yamlConfig.NewConfig(*configLocation, *env)
	if err != nil {
		panic(err)
	}
	webApp, err := handlers.NewWebApp(config)
	if err != nil {
		panic(err)
	}
	webApp.Run()
}
