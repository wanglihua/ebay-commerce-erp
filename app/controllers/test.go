package controllers

import (
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/revel/revel"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {

	// Create the barcode
	barcode_int_cs, _ := code128.Encode("1234567890123444")

	// Scale the barcode to 200x200 pixels
	barcode_, _ := barcode.Scale(barcode_int_cs, 200, 30)

	// encode the barcode as png
	png.Encode(c.Response.Out, barcode_)

	return c.RenderEmpty()
}
