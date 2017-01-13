package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderYzf struct {
	*revel.Controller
}

func (c ErpOrderYzf) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_yzf/order_yzf_index.html")
}
