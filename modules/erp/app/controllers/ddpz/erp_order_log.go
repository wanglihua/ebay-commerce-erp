package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderLog struct {
	*revel.Controller
}

func (c ErpOrderLog) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/order_log/order_log_index.html")
}

func (c ErpOrderLog) ListData() revel.Result {
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

	order_log_list := make([]models.OrderLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_log_list,
		Total: count,
	})
}
