package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderYfh struct {
	*revel.Controller
}

func (c ErpOrderYfh) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_yfh/order_yfh_index.html")
}
