package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"fmt"
	"matrix/modules/erp/models/ddlc_models"
	"strings"
)

type ErpOrderHistory struct {
	*revel.Controller
}

func (c ErpOrderHistory) ListData() revel.Result {
	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var user_id int64
	c.Params.Bind(&user_id, "user_id")

	filter, order, _, _ := core.GetGridRequestParam(c.Request)

	if order_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("order_id = %d", order_id))
	}

	if user_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("user_id = %d", user_id))
	}

	query := db_session.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	order_history_list := make([]ddlc_models.OrderHistoryInfo, 0)
	//err := data_query.Limit(limit, offset).Find(&order_history_list)
	err := data_query.Find(&order_history_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ddlc_models.OrderHistoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_history_list,
		Total: count,
	})
}

func (c ErpOrderHistory) DeleteHistory() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_history/order_history_delete_history.html")
}

func (c ErpOrderHistory) DeleteHistoryListData() revel.Result {
	db_session := c.DbSession

	var order_code_query string
	c.Params.Bind(&order_code_query, "order_code_query")
	order_code_query = strings.TrimSpace(order_code_query)

	filter, order, _, _ := core.GetGridRequestParam(c.Request)

	if order_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("order_code = %s", order_code_query))
	}

	filter = core.SqlAnd(filter, "is_delete = 1")

	query := db_session.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	order_history_list := make([]ddlc_models.OrderHistoryInfo, 0)
	//err := data_query.Limit(limit, offset).Find(&order_history_list)
	err := data_query.Find(&order_history_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ddlc_models.OrderHistoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_history_list,
		Total: count,
	})
}
