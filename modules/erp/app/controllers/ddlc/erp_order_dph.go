package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderDph struct {
	*revel.Controller
}

func (c ErpOrderDph) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_dph/order_dph_index.html")
}
