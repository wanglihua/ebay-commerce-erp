package controllers

import (
	"github.com/revel/revel"
)

type Track_issue struct {
	*revel.Controller
}

func (c Track_issue) Index() revel.Result {
	return c.RenderTemplate("track/issue/issue_index.html")
}
