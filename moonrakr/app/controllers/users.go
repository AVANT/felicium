package controllers

import (
	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type Users struct {
	*revel.Controller
}

//Index displays an array of users to public users
func (u *Users) Index() revel.Result {

	//use the user query method of choice here
	users, err := models.GetUsersByUpdatedAt()
	//check if we should abort
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	//lets get an array of users going using the goson template
	result, err := users.Render()
	//check if we should abort
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result
}

func (u *Users) Create() revel.Result {
	newUser := models.NewUser()
	bulkSetter := new(map[string]interface{})

	//read in the post json body checking for malformed syntax
	err := DecodeJsonRequestToMassAssign(u.Controller, bulkSetter)
	if _err, _result := CheckStopCondition(u.Controller, 400, err); _err {
		return _result
	}

	//set the allowed fields and check for Unprocessable Entities
	newUser, err = newUser.UserMassAssign(bulkSetter)
	if _err, _result := CheckStopCondition(u.Controller, 422, err); _err {
		return _result
	}

	//save the post. Only bad things can happen if there is an error so trigger 500
	err = model.Save(newUser)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	//render the template and throw 500 on an error
	result, err := newUser.Render()
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result
}

func (u *Users) Show() revel.Result {
	var id string = u.Params.Get("id")

	//Check that this is a vaid request
	found, err := models.CheckUserExists(id)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(u.Controller, 404, "")).(revel.Result)
	}

	user, err := models.GetUserById(id)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	result, err := user.Render()
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (u *Users) Update() revel.Result {
	var id string = u.Params.Get("id")

	//Check that this request targets a valid model
	found, err := models.CheckUserExists(id)
	//we can break or just not find it differnt responce codes are appropriate
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(u.Controller, 404, "")).(revel.Result)
	}

	//ok post exists lets get it and make the mods. Bail if there is an error
	user, err := models.GetUserById(id)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	bulkSetter := new(map[string]interface{})

	//decode into bulk setter object and complain if the json is bad
	err = DecodeJsonRequestToMassAssign(u.Controller, bulkSetter)
	if _err, _result := CheckStopCondition(u.Controller, 400, err); _err {
		return _result
	}

	//set the allowed fields and check for Unprocessable Entities
	user, err = user.UserMassAssign(bulkSetter)
	if _err, _result := CheckStopCondition(u.Controller, 422, err); _err {
		return _result
	}

	//update and catch errors
	err = model.Update(user)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	result, err := user.Render()
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (u *Users) Delete() revel.Result {
	var id string = u.Params.Get("id")

	//Check that this is a vaid request
	found, err := models.CheckUserExists(id)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(u.Controller, 404, "")).(revel.Result)
	}

	user, err := models.GetUserById(id)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	err = model.Delete(user)
	if _err, _result := CheckStopCondition(u.Controller, 500, err); _err {
		return _result
	}

	return ReturnLast(CheckStopCondition(u.Controller, 200, "")).(revel.Result)
}
