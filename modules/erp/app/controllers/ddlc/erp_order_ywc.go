package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderYwc struct {
	*revel.Controller
}

func (c ErpOrderYwc) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_ywc/order_ywc_index.html")
}
