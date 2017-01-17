package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"matrix/core/captcha"
)

type Captcha struct {
	*revel.Controller
}

func (c Captcha) Image() revel.Result {
	captcha_bytes := make([]byte, 4)
	captcha_str := captcha.NewLen(4)
	captcha_bytes = []byte(captcha_str)
	captcha_text := ""
	for v := range captcha_bytes {
		captcha_bytes[v] %= 10
		captcha_text += strconv.FormatInt(int64(captcha_bytes[v]), 32)
	}

	revel.TRACE.Println(captcha_text)

	captcha.NewImage(captcha_bytes, 100, 40).WriteTo(c.Response.Out)
	c.Response.ContentType = "image/png"
	return c.RenderEmpty()
}
