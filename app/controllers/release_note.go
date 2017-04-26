package controllers

import "github.com/revel/revel"

type ReleaseNote struct {
	*revel.Controller
}

func (c ReleaseNote) Index() revel.Result {
	return c.RenderTemplate("release_note/release_note_index.html")
}
