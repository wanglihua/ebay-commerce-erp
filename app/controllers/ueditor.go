package controllers

import "github.com/revel/revel"

type UEditor struct {
	*revel.Controller
}

func (c UEditor) Index(action string) revel.Result {
	return c.RenderEmpty()
}

func (c UEditor) Test() revel.Result {
	return c.RenderTemplate("ueditor/ueditor_test.html")
}
