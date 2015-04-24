package controllers

import (
	"github.com/robfig/revel"
)

type Comments struct {
	*revel.Controller
}

func (c *Comments) Index() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) Show() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) Update() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) Delete() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) GetFor() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) CreateFor() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) VoteUp() revel.Result {
	return c.RenderJson("")
}

func (c *Comments) VoteDown() revel.Result {
	return c.RenderJson("")
}
