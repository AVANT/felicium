package hashConfig

import (
	"github.com/vvvntdotorg/felicium/config"
)

type configuration struct {
	config map[string]interface{}
	env    string
	setup  bool
}

func NewConfiguration(data map[string]interface{}, env string) config.Configurator {
	return &configuration{data, env, true}
}

func (c *configuration) Reload() error {
	return nil
}

func (c *configuration) Env() string {
	return c.env
}

func (c *configuration) Lookup(key string) (interface{}, error) {
	if value, found := c.config[key]; !found {
		return nil, config.ValueNotFound
	} else {
		return value, nil
	}
}

func (c *configuration) LookupOrPanic(key string) interface{} {
	if value, err := c.Lookup(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func (c *configuration) LookupWithDefault(key string, defaultValue interface{}) (interface{}, error) {
	if value, found := c.config[key]; !found {
		return defaultValue, nil
	} else {
		return value, nil
	}
}

func (c *configuration) LookupWithDefaultOrPanic(key string, defaultValue interface{}) interface{} {
	// nothing should make us panic here so just uphold the interface
	value, _ := c.LookupWithDefault(key, defaultValue)
	return value
}
