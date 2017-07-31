package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductSaleDayAnalysis struct {
	*revel.Controller
}

func (c ErpProductSaleDayAnalysis) Index() revel.Result {
	return c.RenderTemplate("erp/spbb/product_sale_day_analysis/product_sale_day_analysis_index.html")
}

func (c ErpProductSaleDayAnalysis) ListData() revel.Result {
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

	product_sale_day_analysis_list := make([]models.ProductSaleDayAnalysisInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_sale_day_analysis_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductSaleDayAnalysisInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_sale_day_analysis_list,
		Total: count,
	})
}
