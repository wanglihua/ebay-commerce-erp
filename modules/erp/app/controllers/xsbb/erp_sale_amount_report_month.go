package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"matrix/modules/erp/service/sale_report_service"
	"time"
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
	var db_session = c.DbSession

	var sale_amount_report_month_list = make([]xsbb_models.SaleAmountReportMonthInfo, 0)
	err := db_session.
		OrderBy("report_month ASC").
		Find(&sale_amount_report_month_list);
	core.HandleError(err)

	var time_str_list = make([]string, 0)
	for _, report := range sale_amount_report_month_list {
		time_str_list = append(time_str_list, time.Time(report.ReportMonth).Format("2006-01-02"))
	}

	var series = make([]interface{}, 0)

	var amount_list = make([]float64, 0)
	for _, report := range sale_amount_report_month_list {
		amount_list = append(amount_list, report.Amount)
	}

	var serie_data = map[string]interface{}{
		"name": "商品月销量",
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
			"text": "商品月销量",
			//"subtext": fmt.Sprintf("%s至%s", begin_date.Format("2006-01-02"), end_date.Format("2006-01-02")), // 时间范围
			"subtext": "",
		},
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"legend": map[string]interface{}{
			"data": []string{"月销量"},
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

func (c ErpSaleAmountReportMonth) Calculate() revel.Result {

	go sale_report_service.CalculateSaleAmountReport()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动销售情况后台分析任务成功！",
	})
}
