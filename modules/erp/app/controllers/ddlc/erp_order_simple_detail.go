package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderSimpleDetail struct {
	*revel.Controller
}

func (c ErpOrderSimpleDetail) Index() revel.Result {
	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	c.RenderArgs["order_id"] = order_id
	return c.RenderTemplate("erp/ddlc/order_simple_detail/order_simple_detail_index.html")
}
