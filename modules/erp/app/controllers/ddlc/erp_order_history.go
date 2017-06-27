package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpOrderHistory struct {
	*revel.Controller
}

func (c ErpOrderHistory) ListData() revel.Result {
	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	//query extra filter here
	filter = core.SqlAnd(filter, fmt.Sprintf("order_id = %d", order_id))

	query := db_session.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	order_history_list := make([]models.OrderHistoryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_history_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderHistoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_history_list,
		Total: count,
	})
}
