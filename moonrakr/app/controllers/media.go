package controllers

import (
	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type Medium struct {
	*revel.Controller
}

func (m *Medium) Index() revel.Result {
	//use the posts query method of choice here
	medium, err := models.GetMediumByUpdatedAt()
	//check if we should abort
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	//lets get an array of posts going using the goson template
	result, err := medium.Render()
	//check if we should abort
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result
}

func (m *Medium) Show() revel.Result {
	var id string = m.Params.Get("id")

	//Check that this is a vaid request
	found, err := models.CheckMediaExists(id)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(m.Controller, 404, "")).(revel.Result)
	}

	//model exists get it and throw error if it something breaks
	//TODO can this be done in one strep?
	media, err := models.GetMediaById(id)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	//render the post and throw error on failure
	result, err := media.Render()
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (m *Medium) Create() revel.Result {
	newMedia := models.NewMedia()

	//read in the media form values checking for malformed syntax
	newMedia, err := newMedia.MediaMassAsignFromValues(m.Params.Form)
	if _err, _result := CheckStopCondition(m.Controller, 400, err); _err {
		return _result
	}

	//upload the models payload
	err = newMedia.UploadFromParams(m.Params)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	err = model.Save(newMedia)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	//render the template and throw 500 on an error
	result, err := newMedia.Render()
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result

}

func (m *Medium) Update() revel.Result {
	var id string = m.Params.Get("id")

	//Check that this request targets a valid model
	found, err := models.CheckMediaExists(id)
	//we can break or just not find it differnt responce codes are appropriate
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(m.Controller, 404, "")).(revel.Result)
	}

	//ok post exists lets get it and make the mods. Bail if there is an error
	media, err := models.GetMediaById(id)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	bulkSetter := new(map[string]interface{})

	//decode into bulk setter object and complain if the json is bad
	err = DecodeJsonRequestToMassAssign(m.Controller, bulkSetter)
	if _err, _result := CheckStopCondition(m.Controller, 400, err); _err {
		return _result
	}

	//set the allowed fields and check for Unprocessable Entities
	media, err = media.MediaMassAssign(bulkSetter)
	if _err, _result := CheckStopCondition(m.Controller, 422, err); _err {
		return _result
	}

	//update and catch errors
	err = model.Update(media)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	result, err := media.Render()
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (m *Medium) Delete() revel.Result {
	var id string = m.Params.Get("id")

	//Check that this is a vaid request
	found, err := models.CheckMediaExists(id)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}
	//if there wasn't an error and we didn't find it error with a 404
	if !found {
		return ReturnLast(CheckStopCondition(m.Controller, 404, "")).(revel.Result)
	}

	media, err := models.GetMediaById(id)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	err = model.Delete(media)
	if _err, _result := CheckStopCondition(m.Controller, 500, err); _err {
		return _result
	}

	return ReturnLast(CheckStopCondition(m.Controller, 200, "")).(revel.Result)
}

func (m *Medium) For() revel.Result {
	return m.RenderJson("")
}
