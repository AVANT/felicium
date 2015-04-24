package components

import (
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/mattbaird/elastigo/api"
)

func drop(c []CommandFlag) {
	all := findByName(c, "all").Value.(*bool)
	if *all {
		models.Connection.DeleteDatabase()
		request, err := api.ElasticSearchRequest("DELETE", "/_all/")
		if err != nil {
			messageAndExit(err)
		}
		request.Header.Set("Content-Type", "application/json")

		var toFill interface{}
		code, _, err := request.Do(&toFill)
		if code > 300 {
			messageAndExit(err)
		}
	}
}
