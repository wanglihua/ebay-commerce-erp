package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"fmt"
	"matrix/modules/erp/service/sale_report_service"
	"time"
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

func (c ErpSaleAmountReportDay) Chart() revel.Result {
	return c.RenderTemplate("erp/xsbb/sale_amount_report_day/sale_amount_report_day_chart.html")
}

func (c ErpSaleAmountReportDay) ChartData() revel.Result {

	var time_now = time.Now()
	var today_date = time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, time.Local)

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "query_begin_time")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	if query_begin_date_err != nil {
		query_begin_date = today_date.Add(-100 * 24 * time.Hour) // 当前日期回退100天
	}

	var query_end_time string
	c.Params.Bind(&query_end_time, "query_end_time")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)
	if query_end_date_err != nil {
		query_end_date = today_date
	}

	var db_session = c.DbSession

	var sale_amount_report_day_list = make([]xsbb_models.SaleAmountReportDayInfo, 0)
	err := db_session.Where("report_day >= ? AND report_day <= ?", query_begin_date, query_end_date).
		OrderBy("report_day ASC").
		Find(&sale_amount_report_day_list);
	core.HandleError(err)

	var time_str_list = make([]string, 0)
	for _, report := range sale_amount_report_day_list {
		time_str_list = append(time_str_list, time.Time(report.ReportDay).Format("2006-01-02"))
	}

	var series = make([]interface{}, 0)

	var amount_list = make([]float64, 0)
	for _, report := range sale_amount_report_day_list {
		amount_list = append(amount_list, report.Amount)
	}

	var serie_data = map[string]interface{}{
		"name": "商品日销量",
		"type": "line",
		"data": amount_list,
		"markPoint": map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"type": "max",
					"name": "最大值",
				},
				map[string]interface{}{
					"type": "min",
					"name": "最小值",
				},
			},
		},
		"markLine": map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"type": "average",
					"name": "平均值",
				},
			},
		},
	}

	series = append(series, serie_data)

	var options = map[string]interface{}{
		"title": map[string]interface{}{
			"text":    "商品日销量",
			"subtext": fmt.Sprintf("%s至%s", query_begin_date.Format("2006-01-02"), query_end_date.Format("2006-01-02")), // 时间范围
		},
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"legend": map[string]interface{}{
			"data": []string{"日销量"},
		},
		"toolbox": map[string]interface{}{
			"show": true,
			"feature": map[string]interface{}{
				"dataZoom": map[string]interface{}{
					"yAxisIndex": "none",
				},
				"dataView": map[string]interface{}{
					"readOnly": false,
				},
				"magicType": map[string]interface{}{
					"type": []string{"line", "bar"},
				},
				"restore":     map[string]interface{}{},
				"saveAsImage": map[string]interface{}{},
			},
		},
		"xAxis": map[string]interface{}{
			"type":        "category",
			"boundaryGap": false,
			"data":        time_str_list, // time list
		},
		"yAxis": map[string]interface{}{
			"type": "value",
			"axisLabel": map[string]interface{}{
				"formatter": "{value}",
			},
		},
		"series": series,
	}

	return c.RenderJson(options)
}

func (c ErpSaleAmountReportDay) Calculate() revel.Result {

	go sale_report_service.CalculateSaleAmountReport()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动销售情况后台分析任务成功！",
	})
}
