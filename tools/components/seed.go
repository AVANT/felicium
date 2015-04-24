package components

import (
	"github.com/AVANT/felicium/moonrakr/app/lib/boot"
	"github.com/AVANT/felicium/moonrakr/app/lib/seeder"
	"github.com/robfig/revel"
)

//getConfigString is helper to get string value from config or halt program
func getConfigString(s string, message string) string {
	toReturn, found := revel.Config.String(s)
	if !found {
		revel.ERROR.Fatal(message)
	}
	return toReturn
}

func seed(c []CommandFlag) {
	f := findByName(c, "file").Value.(*string)
	env := findByName(TopLevelCommandFlags, "env").Value.(*string)
	boot.ToolBoot(*env)

	_, err := seeder.SeedFromJson(*f)
	if err != nil {
		panic(err)
	}

}
