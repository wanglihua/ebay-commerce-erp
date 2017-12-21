package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"fmt"
	"matrix/modules/erp/service/sale_report_service"
)

type ErpSaleAmountReportDay struct {
	*revel.Controller
}

func (c ErpSaleAmountReportDay) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/sale_amount_report_day/sale_amount_report_day_index.html")
}

func (c ErpSaleAmountReportDay) ListData() revel.Result {
	db_session := c.DbSession

	var begin_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")

	var end_time_query core.NullLongTime
	c.Params.Bind(&end_time_query, "end_time_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("report_day >= '%s'", begin_time_query.Time.Format("2006-01-02")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("report_day < '%s'", end_time_query.Time.AddDate(0, 0, 1).Format("2006-01-02")))
	}

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

func (c ErpSaleAmountReportDay) Calculate() revel.Result {

	go sale_report_service.CalculateSaleAmountReport()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动销售情况后台分析任务成功！",
	})
}
