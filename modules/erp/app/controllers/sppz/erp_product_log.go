package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductLog struct {
	*revel.Controller
}

func (c ErpProductLog) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/product_log/product_log_index.html")
}

func (c ErpProductLog) ListData() revel.Result {
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

	product_log_list := make([]models.ProductLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_log_list,
		Total: count,
	})
}
