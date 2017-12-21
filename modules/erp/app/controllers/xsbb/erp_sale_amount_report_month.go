package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"matrix/modules/erp/service/sale_report_service"
)

type ErpSaleAmountReportMonth struct {
	*revel.Controller
}

func (c ErpSaleAmountReportMonth) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/sale_amount_report_month/sale_amount_report_month_index.html")
}

func (c ErpSaleAmountReportMonth) ListData() revel.Result {
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

	sale_amount_report_day_list := make([]xsbb_models.SaleAmountReportMonthInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&sale_amount_report_day_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xsbb_models.SaleAmountReportMonthInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  sale_amount_report_day_list,
		Total: count,
	})
}

func (c ErpSaleAmountReportMonth) Chart() revel.Result {
	return c.RenderTemplate("erp/xsbb/sale_amount_report_month/sale_amount_report_month_chart.html")
}

func (c ErpSaleAmountReportMonth) ChartData() revel.Result {
	return c.RenderJson(true)
}

func (c ErpSaleAmountReportMonth) Calculate() revel.Result {

	go sale_report_service.CalculateSaleAmountReport()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动销售情况后台分析任务成功！",
	})
}
