package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

//bad CORS implementation
func (a *App) Options() revel.Result {
	return nil
}

////
//	Global Controller Helpers
////

//DecodeJsonRequest unpacks the request body json into an interface.
func DecodeJsonRequest(r *revel.Controller, toFill interface{}) error {
	defer r.Request.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, toFill)
	if err != nil {
		return err
	}
	return nil
}

//DecodeJsonRequestToModel unpacks the request body json into a model.
func DecodeJsonRequestToMassAssign(r *revel.Controller, toFill *map[string]interface{}) error {
	err := DecodeJsonRequest(r, toFill)
	if err != nil {
		return err
	}
	return nil
}

//DecodeJsonRequestToModel unpacks the request body json into a model.
func DecodeJsonRequestToModel(r *revel.Controller, toFill model.IsModel) error {
	err := DecodeJsonRequest(r, toFill.GetValueHash())
	if err != nil {
		return err
	}
	return nil
}

//This causes a panic if there is an error. This is basically obsolete.
func CheckFatalError(err error) {
	if err != nil {
		panic(err)
	}
}

//CheckStopCondition if the app should stop execution of an action before the final return.
//All relavant callbacks are hooked in at this point.
//ToDo implement mail alert here.
func CheckStopCondition(c *revel.Controller, code int, descriptor interface{}) (bool, revel.Result) {
	if descriptor != nil {
		message := c.Message(fmt.Sprintf("general.%d", code))
		err := revel.NewErrorFromPanic(descriptor)
		err.Title = "Stop Condition Triggered."
		r := models.BuildResponse(code, message)
		if code >= 500 {
			revel.ERROR.Print(err, "\n", err.Stack)
		}
		if code >= 300 && code < 500 {
			revel.INFO.Println(err)
		}
		c.Response.Status = code
		c.Response.ContentType = "application/json"
		return true, c.RenderJson(r)
	} else {
		return false, nil
	}
}

//this returns the last argument as an interface{} type
func ReturnLast(args ...interface{}) interface{} {
	if len(args) > 0 {
		return args[(len(args) - 1)]
	} else {
		return args[0]
	}
}
