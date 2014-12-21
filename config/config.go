package config

import (
	"errors"
)

var (
	InvalidEnvError = errors.New("Invalid Environment Error.")
	ValueNotFound   = errors.New("Value was not found.")
)

type Configurator interface {
	// lookup should return error eventhough some implementions don't need it
	// those implementaions should return errors with a package specific error
	// type when values cannot be found so the user can determine the outcome
	// of a lookup
	Lookup(key string) (interface{}, error)
	LookupOrPanic(key string) interface{}
	LookupWithDefault(key string, defaultValue interface{}) (interface{}, error)
	// some implementations could have errors when attempting the lookup even with the default
	LookupWithDefaultOrPanic(key string, defaultValue interface{}) interface{}
	Reload() error
	Env() string
}
