package controllers

import (
	"github.com/revel/revel"
)

type Log_product struct {
	*revel.Controller
}

func (c Log_product) Index() revel.Result {
	return c.RenderTemplate("log/product/product_index.html")
}
