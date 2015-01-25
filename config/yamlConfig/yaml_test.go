package yamlConfig

import (
	"fmt"
	"testing"

	. "github.com/avant/felicium/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/avant/felicium/config"
)

type expectedResults []struct {
	key          string
	value        interface{}
	defaultValue interface{}
}

func lookupTestHelper(yc config.Configurator, expected expectedResults) {
	for _, v := range expected {

		found, err := yc.Lookup(v.key)
		Convey(fmt.Sprintf("lookup of %s should have no error", v.key), func() {
			So(err, ShouldBeNil)
		})
		Convey(fmt.Sprintf("%s should return %v\n", v.key, v.value), func() {
			So(found, ShouldResemble, v.value)
		})
	}
}

func lookupWithDefaultsTestHelper(yc config.Configurator, expected expectedResults) {
	for _, v := range expected {

		found, err := yc.LookupWithDefault(v.key, v.defaultValue)
		Convey(fmt.Sprintf("lookup of %s should have no error", v.key), func() {
			So(err, ShouldBeNil)
		})
		Convey(fmt.Sprintf("%s should return %v\n", v.key, v.value), func() {
			So(found, ShouldResemble, v.value)
		})
	}
}

func TestSetup(t *testing.T) {

	Convey("When Config is setup from file", t, func() {

		Convey("Requesting an undefined env should cause an error", func() {
			_, err := NewConfig("test_files/test.yaml", "undefined_env")
			So(err, ShouldEqual, config.InvalidEnvError)
		})

		Convey("You should be able to request the default env", func() {
			_, err := NewConfig("test_files/test.yaml", "default")
			So(err, ShouldBeNil)
		})

		Convey("You should be able to request a defined env", func() {
			_, err := NewConfig("test_files/test2.yaml", "production")
			So(err, ShouldBeNil)
		})
	})

	Convey("When seting up config without constructor Config methods should panic", t, func() {
		yc := &configuration{}
		So(func() { yc.Env() }, ShouldPanic)
	})

}

func TestExpectedErrors(t *testing.T) {

	Convey("Querying an undefined key should cause an error", t, func() {
		yc, _ := NewConfig("test_files/test.yaml", "production")
		_, err := yc.Lookup("dne")
		So(err, ShouldEqual, config.ValueNotFound)
	})

}

func TestConfigLookups(t *testing.T) {

	Convey("Lookups using test.yaml should return the expected values in default env", t, func() {

		var yamlTestCases = expectedResults{
			{
				key:   "key1",
				value: "default1",
			},
			{
				key:   "key2",
				value: float64(1000),
			},
			{
				key: "key3",
				value: []interface{}{
					"arrayDefault1",
					"arrayDefault2",
					"arrayDefault3",
				},
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "default")
		So(err, ShouldBeNil)
		lookupTestHelper(yc, yamlTestCases)
	})

	Convey("Lookups using test.yaml should return the expected values in test env", t, func() {
		var yamlTestCases = expectedResults{
			{
				key:   "key1",
				value: float64(2000),
			},
			{
				key:   "key2",
				value: float64(1000),
			},
			{
				key: "key3",
				value: []interface{}{
					"arrayDefault1",
					"arrayDefault2",
					"arrayDefault3",
				},
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "test")
		So(err, ShouldBeNil)
		lookupTestHelper(yc, yamlTestCases)
	})

	Convey("Lookups using test.yaml should return the expected values in production env", t, func() {
		var yamlTestCases = expectedResults{
			{
				key:   "key1",
				value: "default1",
			},
			{
				key:   "key2",
				value: float64(1000),
			},
			{
				key: "key3",
				value: []interface{}{
					"arrayProduction1",
					"arrayProduction2",
					"arrayProduction3",
				},
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "production")
		So(err, ShouldBeNil)
		lookupTestHelper(yc, yamlTestCases)
	})

	Convey("LookupWithDefaults using test.yaml should return the expected values in production env", t, func() {
		var yamlTestCases = expectedResults{
			{
				key:          "key1",
				value:        "default1",
				defaultValue: "key exits so it should not be this value",
			},
			{
				key:          "key2",
				value:        float64(1000),
				defaultValue: "key exits so it should not be this value",
			},
			{
				key: "key3",
				value: []interface{}{
					"arrayProduction1",
					"arrayProduction2",
					"arrayProduction3",
				},
				defaultValue: "key exits so it should not be this value",
			},
			{
				key:          "doesntexist",
				value:        "should be this value because the key doesn't exist",
				defaultValue: "should be this value because the key doesn't exist",
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "production")
		So(err, ShouldBeNil)
		lookupWithDefaultsTestHelper(yc, yamlTestCases)
	})

	Convey("LookupORPanic using test.yaml should return the expected values in production env", t, func() {
		var yamlTestCases = expectedResults{
			{
				key:   "key1",
				value: "default1",
			},
			{
				key:   "key2",
				value: float64(1000),
			},
			{
				key: "key3",
				value: []interface{}{
					"arrayProduction1",
					"arrayProduction2",
					"arrayProduction3",
				},
				defaultValue: "key exits so it should not be this value",
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "production")
		So(err, ShouldBeNil)
		lookupTestHelper(yc, yamlTestCases)

		panicKey := "notThereKey"
		Convey(fmt.Sprintf("lookup of %s should have no error", panicKey), func() {
			So(func() {
				yc.LookupOrPanic(panicKey)
			}, ShouldPanic)
		})
	})
}

func TestReloadFunctionality(t *testing.T) {

	Convey("Reload should update config values", t, func() {
		var yamlTestCases = expectedResults{
			{
				key:   "key1",
				value: "productionTest2Value1",
			},
		}
		yc, err := NewConfig("test_files/test.yaml", "production")
		So(err, ShouldBeNil)
		conf := yc.(*configuration)
		conf.path = "test_files/test2.yaml"
		err = yc.Reload()
		So(err, ShouldBeNil)
		lookupTestHelper(yc, yamlTestCases)
	})

	Convey("Reload should remove untracked values on update", t, func() {
		yc, err := NewConfig("test_files/test.yaml", "production")
		So(err, ShouldBeNil)
		conf := yc.(*configuration)
		conf.path = "test_files/test2.yaml"
		err = yc.Reload()
		So(err, ShouldBeNil)
		_, err = yc.Lookup("key3")
		So(err, ShouldNotBeNil)
		_, err = yc.Lookup("key2")
		So(err, ShouldNotBeNil)
	})

}
