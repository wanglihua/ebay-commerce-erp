package controllers

import (
	"github.com/revel/revel"
	"matrix/core/captcha"
)

type Captcha struct {
	*revel.Controller
}

func (c Captcha) Image() revel.Result {
	captcha.Generate(100, 40, c.Session).WriteTo(c.Response.Out)
	c.Response.ContentType = "image/png"
	return c.RenderEmpty()
}
