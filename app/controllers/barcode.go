package controllers

import (
	"github.com/revel/revel"
	"github.com/boombuler/barcode/code128"
	"image/png"
	"github.com/boombuler/barcode"
)

type Barcode struct {
	*revel.Controller
}

func (c Barcode) Waybill() revel.Result {
	var code string
	c.Params.Bind(&code, "code")

	// Create the barcode
	barcode_int_cs, _ := code128.Encode(code)

	barcode_, _ := barcode.Scale_barcode_height(barcode_int_cs, 100)

	// encode the barcode as png
	png.Encode(c.Response.Out, barcode_)

	return c.RenderEmpty()
}
