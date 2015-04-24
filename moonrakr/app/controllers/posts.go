package controllers

import (
	//"github.com/mattbaird/elastigo/api"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
	//"github.com/emilsjolander/goson"
	"github.com/AVANT/felicium/model"
)

type Posts struct {
	*revel.Controller
}

//Index displays an array of posts to public users
func (p *Posts) Index() revel.Result {

	var status string = p.Params.Get("status")
	switch status {
	case "unpublished":
	default:
		status = "published"
	}

	//use the posts query method of choice here
	posts, err := models.GetPostsByStatus(status)
	//check if we should abort
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	//lets get an array of posts going using the goson template
	result, err := posts.Render()
	//check if we should abort
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result
}

func (p *Posts) Create() revel.Result {
	newPost := models.NewPost()
	bulkSetter := new(map[string]interface{})

	//read in the post json body checking for malformed syntax
	err := DecodeJsonRequestToMassAssign(p.Controller, bulkSetter)
	if _err, _result := CheckStopCondition(p.Controller, 400, err); _err {
		return _result
	}

	//set the allowed fields and check for Unprocessable Entities
	newPost, err = newPost.PostMassAssign(bulkSetter)
	if _err, _result := CheckStopCondition(p.Controller, 422, err); _err {
		return _result
	}

	//save the post. Only bad things can happen if there is an error so trigger 500
	err = model.Save(newPost)
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	//render the template and throw 500 on an error
	result, err := newPost.Render()
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	//return the JsonResult in this case which will use its Apply method to render
	return result
}

func (p *Posts) Show() revel.Result {
	var id string = p.Params.Get("id")

	post, err := models.GetPostBySlug(id)

	if err != nil && err.Error() == "Not Found" {
		return ReturnLast(CheckStopCondition(p.Controller, 404, "")).(revel.Result)
	}

	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	//render the post and throw error on failure
	result, err := post.Render()
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (p *Posts) Update() revel.Result {
	var id string = p.Params.Get("id")

	post, err := models.GetPostBySlug(id)
	if err != nil && err.Error() == "Not Found" {
		return ReturnLast(CheckStopCondition(p.Controller, 404, "")).(revel.Result)
	}

	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	bulkSetter := new(map[string]interface{})

	//decode into bulk setter object and complain if the json is bad
	err = DecodeJsonRequestToMassAssign(p.Controller, bulkSetter)
	if _err, _result := CheckStopCondition(p.Controller, 400, err); _err {
		return _result
	}

	//set the allowed fields and check for Unprocessable Entities
	post, err = post.PostMassAssign(bulkSetter)
	if _err, _result := CheckStopCondition(p.Controller, 422, err); _err {
		return _result
	}

	//update and catch errors
	err = model.Update(post)
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	result, err := post.Render()
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	return result
}

func (p *Posts) Delete() revel.Result {
	var id string = p.Params.Get("id")

	post, err := models.GetPostBySlug(id)
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}
	if err != nil && err.Error() == "Not Found" {
		return ReturnLast(CheckStopCondition(p.Controller, 404, "")).(revel.Result)
	}

	err = model.Delete(post)
	if _err, _result := CheckStopCondition(p.Controller, 500, err); _err {
		return _result
	}

	return ReturnLast(CheckStopCondition(p.Controller, 200, "")).(revel.Result)
}

func (p *Posts) For() revel.Result {
	var id string = p.Params.Get("id")
	return p.RenderJson(id)
}

func (p *Posts) MakeRecomendation() revel.Result {
	var id string = p.Params.Get("id")
	return p.RenderJson(id)
}

func (p *Posts) UndoRecomendation() revel.Result {
	var id string = p.Params.Get("id")
	return p.RenderJson(id)
}

func (p *Posts) Comments() revel.Result {
	var id string = p.Params.Get("id")
	return p.RenderJson(id)
}
