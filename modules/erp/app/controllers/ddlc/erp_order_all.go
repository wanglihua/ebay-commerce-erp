package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderAll struct {
	*revel.Controller
}

func (c ErpOrderAll) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_all/order_all_index.html")
}
