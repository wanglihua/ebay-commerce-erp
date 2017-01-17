package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderKhb struct {
	*revel.Controller
}

func (c ErpOrderKhb) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_khb/order_khb_index.html")
}
