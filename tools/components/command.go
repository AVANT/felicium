package components

import (
	"flag"
	"fmt"
	"math"
	"os"
	"reflect"

	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/config"
)

type SubDescription struct {
	Name string
	Help string
}

type CommandFlag struct {
	Type        string
	Help        string
	Name        string
	Value       interface{}
	Subcommands map[SubDescription][]CommandFlag
}

func NewCommandFlag(name string, defaultValue interface{}, help string, c map[SubDescription][]CommandFlag) CommandFlag {
	toReturn := CommandFlag{}
	toReturn.Type = reflect.TypeOf(defaultValue).String()
	toReturn.Help = help
	toReturn.Name = name
	toReturn.Subcommands = c
	toReturn.Value = defaultValue
	return toReturn
}

var TopLevelCommandFlags = []CommandFlag{
	NewCommandFlag("help", new(bool), "Returns the usage documentation.", nil),
	NewCommandFlag("configPath", defaultString("conf/app.conf"), "This can be used in conjunction with any command to specify the path to the revel config. Defaults to ./conf/app.conf", nil),
	NewCommandFlag("env", defaultString("dev"), "This can be used in conjunction with any command to specify the enviornment. Defaults to dev.", nil),
	NewCommandFlag("method", new(string), "This selects one of the subcommands.", map[SubDescription][]CommandFlag{
		SubDescription{Name: "search", Help: "This allows you to run queries against an index."}: []CommandFlag{
			NewCommandFlag("index", new(string), "This is the index to run the query against.", nil),
			NewCommandFlag("query", new(string), "This is the Query to run against the index.", nil),
		},
		SubDescription{Name: "seed", Help: "This will seed the database with the specified seed file."}: []CommandFlag{
			NewCommandFlag("file", new(string), "This specifies the file to use for the seed.", nil),
		},
		SubDescription{Name: "drop", Help: "This allows you to drop data in the database."}: []CommandFlag{
			NewCommandFlag("all", new(bool), "Drop all indexes.", nil),
			NewCommandFlag("just", new(string), "Drop a specific Index.", nil),
		},
	}),
}

var Config config.Config

func GenerateHelpWrapped(c []CommandFlag, indent int) {
	maxInCol := 0.0
	for _, i := range c {
		if l := float64(len(i.Name)); l > maxInCol {
			maxInCol = l
		}
	}

	for _, i := range c {
		fmt.Printf("%s--%s%s : %s\n", createTabs(indent), i.Name, createTabs(int(math.Ceil((maxInCol-float64(len(i.Name)))/5.0))), i.Help)
		generateHelpSubcommands(i.Subcommands, indent+1)
	}
}

func generateHelpSubcommands(c map[SubDescription][]CommandFlag, indent int) {
	for k, v := range c {
		fmt.Printf("%s%s : %s\n", createTabs(indent), k.Name, k.Help)
		GenerateHelpWrapped(v, indent+1)
	}
}

func createTabs(number int) string {
	toReturn := ""
	for i := 0; i < number; i++ {
		toReturn += "\t"
	}
	return toReturn
}

func defaultString(s string) *string {
	return &s
}

func GenerateHelp(c []CommandFlag) {
	GenerateHelpWrapped(c, 0)
}

func ListenForFlag(c []CommandFlag) {
	for _, v := range c {
		switch v.Type {
		case "*string":
			flag.StringVar(v.Value.(*string), v.Name, *v.Value.(*string), v.Help)
		case "*int":
			flag.IntVar(v.Value.(*int), v.Name, *v.Value.(*int), v.Help)
		case "*bool":
			flag.BoolVar(v.Value.(*bool), v.Name, *v.Value.(*bool), v.Help)
		}
		if v.Subcommands != nil {
			for _, sub := range v.Subcommands {
				ListenForFlag(sub)
			}
		}
	}
}

func findByName(c []CommandFlag, name string) CommandFlag {
	for _, v := range c {
		if v.Name == name {
			return v
		}
	}
	return CommandFlag{}
}

func helpAndExit() {
	GenerateHelp(TopLevelCommandFlags)
	os.Exit(0)
}

func messageAndExit(message interface{}) {
	fmt.Println(message)
	os.Exit(1)
}

func initDBConnection() {
	var filename, env *string
	for _, v := range TopLevelCommandFlags {
		switch v.Name {
		case "configPath":
			filename = v.Value.(*string)
		case "env":
			env = v.Value.(*string)
		}
	}
	Config, err := config.ReadDefault(*filename)
	if err != nil {
		messageAndExit(err)
	}
	if !Config.HasSection(*env) {
		messageAndExit(fmt.Sprintf("could not find env %s in config\n", *env))
	}

	////
	//	Read the config in
	////

	var CouchHost, CouchPort, CouchDatabase, EsDomain, EsPort, EsProto string
	var found error
	CouchHost, found = Config.String(*env, "couchDB.domain")
	if found != nil {
		messageAndExit("You did not specify a couch domain for current environment.")
	}
	CouchPort, found = Config.String(*env, "couchDB.port")
	if found != nil {
		messageAndExit("You did not specify a couch port for current environment.")
	}
	CouchDatabase, found = Config.String(*env, "couchDB.database")
	if found != nil {
		messageAndExit("You did not specify a couch name for current environment.")
	}
	EsDomain, found = Config.String(*env, "es.domain")
	if found != nil {
		messageAndExit("You did not specify an es domain for current environment.")
	}
	EsPort, found = Config.String(*env, "es.port")
	if found != nil {
		messageAndExit("You did not specify an es port for current environment.")
	}
	EsProto, found = Config.String(*env, "es.proto")
	if found != nil {
		messageAndExit("You did not specify a es proto for current environment.")
	}
	models.Setup(CouchHost, CouchPort, CouchDatabase, EsDomain, EsPort, EsProto)
}

func pullOutSubCommands(c map[SubDescription][]CommandFlag, name string) []CommandFlag {
	for k, v := range c {
		if k.Name == name {
			return v
		}
	}
	return []CommandFlag{}
}

func checkError(e error) {
	if e != nil {
		messageAndExit(e)
	}
}

func Execute() {
	ListenForFlag(TopLevelCommandFlags)

	flag.Parse()
	b := TopLevelCommandFlags[0].Value.(*bool)
	if *b {
		helpAndExit()
	}
	initDBConnection()

	command := findByName(TopLevelCommandFlags, "method")
	method := command.Value.(*string)

	switch *method {
	case "search":
		search(pullOutSubCommands(command.Subcommands, "search"))
	case "seed":
		seed(pullOutSubCommands(command.Subcommands, "seed"))
	case "drop":
		drop(pullOutSubCommands(command.Subcommands, "drop"))
	default:
		helpAndExit()
	}
}
