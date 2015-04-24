package helpers

import (
	"errors"

	"github.com/jmoiron/jsonq"
	//"github.com/robfig/revel"
	"net/url"
	"reflect"
	"strconv"
)

//jsonqArrayOfStrings is a patch untill the patch that i submitted gets pulled.
//https://github.com/jmoiron/jsonq/pull/2
func JsonqArrayOfStrings(jq *jsonq.JsonQuery, s ...string) ([]string, error) {
	array, err := jq.Array(s...)
	if err != nil {
		return []string{}, err
	}

	stringArray := make([]string, len(array))
	for i := range array {
		found, err := jq.String(append(s, strconv.Itoa(i))...)
		if err != nil {
			return stringArray, err
		}
		stringArray[i] = found
	}
	return stringArray, nil
}

//jsonqArrayOfInts is a patch untill the patch that i submitted gets pulled.
//https://github.com/jmoiron/jsonq/pull/2
func JsonqArrayOfInts(jq *jsonq.JsonQuery, s ...string) ([]int, error) {
	array, err := jq.Array(s...)
	if err != nil {
		return []int{}, err
	}

	intArray := make([]int, len(array))
	for i := range array {
		s, err := jq.Int(append(s, strconv.Itoa(i))...)
		if err != nil {
			return intArray, err
		}
		intArray[i] = s
	}
	return intArray, nil
}

//InterfaceArrayToStringArray converts an []interface{} to a []string
func InterfaceArrayToStringArray(in interface{}) ([]string, error) {
	var itterable []interface{}
	switch in.(type) {
	case []interface{}:
		itterable = in.([]interface{})
	default:
		return []string{}, errors.New("couldn't convert to [] the argument provided wasn't type []interface{}")
	}
	toReturn := make([]string, len(itterable))
	for i := range itterable {
		switch itterable[i].(type) {
		case string:
			toReturn[i] = itterable[i].(string)
		default:
			return toReturn, errors.New("an object in the array wasn't of type string and couldn't be converted.")
		}
	}
	return toReturn, nil
}

// The Bulk Set functions could use this later on if we want but it will make the nameing convention way more strict
func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

//UrlValuesToJsonObject Take the url.Values and break it down into a MassAssin object
func UrlValuesToJsonObject(values url.Values) *map[string]interface{} {
	toReturn := map[string]interface{}{}
	for k, v := range values {
		if len(v) == 1 {
			toReturn[k] = v[0]
		} else {
			toReturn[k] = v
		}

	}
	return &toReturn
}

//func ForEach(func(k string, v map[string]interface{}, out interface{}))
