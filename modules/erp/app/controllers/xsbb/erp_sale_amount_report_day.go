package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
)

type ErpSaleAmountReportDay struct {
	*revel.Controller
}

func (c ErpSaleAmountReportDay) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/sale_amount_report_day/sale_amount_report_day_index.html")
}

func (c ErpSaleAmountReportDay) ListData() revel.Result {
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

	sale_amount_report_day_list := make([]xsbb_models.SaleAmountReportDayInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&sale_amount_report_day_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xsbb_models.SaleAmountReportDayInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  sale_amount_report_day_list,
		Total: count,
	})
}
