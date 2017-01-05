package controllers

import (
	"github.com/revel/revel"
)

type ErpProductImport struct {
	*revel.Controller
}

func (c ErpProductImport) Index() revel.Result {
	return c.RenderTemplate("erp/product_import/product_import.html")
}
