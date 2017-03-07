package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
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
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	order_list := make([]models.OrderInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
