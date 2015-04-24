package controllers

import (
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type Tags struct {
	*revel.Controller
}

func (t *Tags) Index() revel.Result {
	toReturn, err := models.GetAllTags()
	if _err, _result := CheckStopCondition(t.Controller, 500, err); _err {
		return _result
	}
	return t.RenderJson(toReturn.Tags)
}
