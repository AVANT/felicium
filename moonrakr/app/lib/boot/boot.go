package boot

import (
	"time"

	"github.com/emilsjolander/goson"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/kr/s3/s3util"
	"github.com/mattbaird/elastigo/core"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
)

//getConfigString is helper to get string value from config or halt program
func getConfigString(s string, message string) string {
	toReturn, found := revel.Config.String(s)
	if !found {
		revel.ERROR.Fatal(message)
	}
	return toReturn
}

//connect to database
func ConfigureDB() {
	CouchHost := getConfigString("couchDB.domain", "You did not specify a couch domain for current environment.")
	CouchPort := getConfigString("couchDB.port", "You did not specify a couch port for current environment.")
	CouchDatabase := getConfigString("couchDB.database", "You did not specify a couch name for current environment.")
	EsDomain := getConfigString("es.domain", "You did not specify an es domain for current environment.")
	EsPort := getConfigString("es.port", "You did not specify an es port for current environment.")
	EsProto := getConfigString("es.proto", "You did not specify a es proto for current environment.")
	models.Setup(CouchHost, CouchPort, CouchDatabase, EsDomain, EsPort, EsProto)
}

//set up the goson path variables
func ConfigureTemplates() {
	goson.TemplateRoot = revel.ViewsPath + "/"
}

//put elastic search core into development mode
func ConfigureES() {
	if revel.DevMode {
		core.DebugRequests = true
	}
}

//set up the s3 bucket
func ConfigureS3() {
	AccessKey := getConfigString("aws.s3.access_key", "AccessKey for s3 couldn't be found.")
	SecretKey := getConfigString("aws.s3.secret_key", "SecretKey for s3 couldn't be found.")
	S3Domain := getConfigString("aws.s3.domain", "Domain for s3 couldn't be found.")
	S3Folder := getConfigString("aws.s3.folder", "Folder for s3 couldn't be found.")

	s3util.DefaultConfig.AccessKey = AccessKey
	s3util.DefaultConfig.SecretKey = SecretKey
	models.S3FQDN = S3Domain
	models.S3Folder = S3Folder
}

//This should be called in any context where revel doesn't run fully ie the tool
func EnsureCache() {
	if cache.Instance == nil {
		cache.Instance = cache.NewInMemoryCache(30 * time.Second)
	}
}

func NormalBoot() {
	ConfigureDB()
	ConfigureES()
	ConfigureS3()
	ConfigureTemplates()
}

func ToolBoot(env string) {
	revel.Init(env, "github.com/moonrakr", "")
	NormalBoot()
	EnsureCache()
}
