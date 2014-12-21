package yamlConfig

import (
	"encoding/json"
	"io/ioutil"

	"github.com/vvvntdotorg/felicium/config"
	"github.com/vvvntdotorg/felicium/config/hashConfig"

	"github.com/vvvntdotorg/felicium/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

type configuration struct {
	config.Configurator
	path string
}

func NewConfig(path, env string) (config.Configurator, error) {

	yc := &configuration{}
	yc.path = path

	//load config
	if err := yc.loadFromFile(env); err != nil {
		return nil, err
	}

	return yc, nil

}

func (yc *configuration) Reload() error {
	return yc.loadFromFile(yc.Env())
}

func (yc *configuration) loadFromFile(env string) error {

	var err error
	var rawconfig []byte
	if rawconfig, err = ioutil.ReadFile(yc.path); err != nil {
		return err
	}

	tmp := map[string]map[string]interface{}{}

	if err = yaml.Unmarshal([]byte(rawconfig), tmp); err != nil {
		return err
	}

	_, foundDefault := tmp["default"]
	_, foundEnv := tmp[env]
	configData := map[string]interface{}{}

	// your options here are you have defined the default and the env you requested
	// or you have only defined the requested env. This will error if you don't have
	// the requested env defined. This is to reduce confusion.
	switch {
	case foundDefault && foundEnv:
		var err error
		var defaultData, envData []byte

		// are you looking at this and thinking thats getto?
		// your right the merge lib on github didn't work and
		// yaml doesn't have a Raw type to decode into
		// TODO: do this a better way
		if defaultData, err = json.Marshal(tmp["default"]); err != nil {
			return err
		}
		if envData, err = json.Marshal(tmp[env]); err != nil {
			return err
		}

		if err = json.Unmarshal(defaultData, &configData); err != nil {
			return err
		}
		// if you asked for the default we are done
		if env != "default" {
			// override if needed
			if err = json.Unmarshal(envData, &configData); err != nil {
				return err
			}
		}
	case foundEnv:
		configData = tmp[env]
	default:
		return config.InvalidEnvError
	}
	yc.Configurator = hashConfig.NewConfiguration(configData, env)
	return nil
}
