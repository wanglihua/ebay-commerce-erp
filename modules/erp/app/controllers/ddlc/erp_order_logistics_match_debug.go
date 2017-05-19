package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderLogisticsMatchDebug struct {
	*revel.Controller
}

func (c ErpOrderLogisticsMatchDebug) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_match_debug/order_logistics_match_debug_index.html")
}

func (c ErpOrderLogisticsMatchDebug) ListData() revel.Result {
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

	order_logistics_match_debug_list := make([]models.OrderLogisticsMatchDebugInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_logistics_match_debug_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderLogisticsMatchDebugInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_logistics_match_debug_list,
		Total: count,
	})
}

func (c ErpOrderLogisticsMatchDebug) Detail(id int64) revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_match_debug/order_logistics_match_debug_detail.html")
}
