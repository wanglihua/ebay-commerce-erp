package controllers

import (
	"github.com/revel/revel"
)

type BugTrack struct {
	*revel.Controller
}

func (c BugTrack) Index() revel.Result {
	return c.RenderEmpty()
}
