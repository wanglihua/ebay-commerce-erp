package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
)

type ErpOrderHold struct {
	*revel.Controller
}

func (c ErpOrderHold) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_hold/order_hold_index.html")
}

func (c ErpOrderHold) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	order_list, count := get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
