package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
)

type ErpOrderAll struct {
	*revel.Controller
}

func (c ErpOrderAll) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_all/order_all_index.html")
}

func (c ErpOrderAll) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	order_list, count := get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
