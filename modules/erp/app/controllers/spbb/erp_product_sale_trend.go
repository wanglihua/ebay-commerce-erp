package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/spbb_models"
	"strings"
	"fmt"
)

type ErpProductSaleTrend struct {
	*revel.Controller
}

func (c ErpProductSaleTrend) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	var db_session = c.DbSession

	var product_id_list = make([]int64, 0)
	c.Params.Bind(&product_id_list, "product_id_list")

	var product_name_list = make([]string, 0)
	if len(product_id_list) != 0 {
		// 显示当前趋势分析的商品名
		var product_list = make([]kcgl_models.ProductInfo, 0)
		core.HandleError(db_session.In("id", product_id_list).Find(&product_list))

		for _, product := range product_list {
			product_name_list = append(product_name_list, fmt.Sprintf("%s， %s", product.Sku, product.CnName))
		}
	}
	c.RenderArgs["product_name_list"] = product_name_list

	var product_id_str_list = make([]string, 0)
	for _, product_id := range product_id_list {
		product_id_str_list = append(product_id_str_list, fmt.Sprint(product_id))
	}
	c.RenderArgs["product_id_strs"] = strings.Join(product_id_str_list, ",")

	return c.RenderTemplate("erp/spbb/product_sale_trend/product_sale_trend_index.html")
}

type product_sale_analysis_trend_t struct {
	ProductSaleAnalysis spbb_models.ProductSaleAnalysisInfo `xorm:"extends" json:"sa"`
	Stock               sppz_models.StockInfo               `xorm:"extends" json:"s"`
}

func (c ErpProductSaleTrend) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	/*
	db_session := c.DbSession

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var sb = core.NewSqlBuilder()

	sb.Select("sa.*, s.*")
	sb.Count("count(1) AS sa_count")
	sb.From(`erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id`)

	var sql_args = make([]interface{}, 0)

	if sku_query != "" {
		filter = core.SqlAnd(filter, "sa.sku = ?")
		sql_args = append(sql_args, sku_query)
	}

	// 只取不分仓统计的数据
	filter = core.SqlAnd(filter, "sa.type = ?")
	sql_args = append(sql_args, spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_SKU)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sa.id ASC")
	}

	sb.Limit(offset, limit)

	product_sale_analysis_list := make([]product_sale_analysis_trend_t, 0, limit)
	err := db_session.Sql(sb.SqlList(), sql_args...).Find(&product_sale_analysis_list)
	core.HandleError(err)

	for index, product_sale_day_analysis := range product_sale_analysis_list {
		var product kcgl_models.ProductInfo
		_, err := db_session.Id(product_sale_day_analysis.ProductSaleAnalysis.ProductId).Get(&product)
		core.HandleError(err)

		product_sale_analysis_list[index].ProductSaleAnalysis.ProductName = product.CnName
	}

	var sa_count struct {
		Count int64 `xorm:"'sa_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount(), sql_args...).Get(&sa_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_sale_analysis_list,
		Total: sa_count.Count,
	})
	*/

	/*
			var option = {
				title: {
					text: '未来一周气温变化',
					subtext: '纯属虚构'
				},
				tooltip: {
					trigger: 'axis'
				},
				legend: {
					data: ['最高气温', '最低气温']
				},
				toolbox: {
					show: true,
					feature: {
						dataZoom: {
							yAxisIndex: 'none'
						},
						dataView: {readOnly: false},
						magicType: {type: ['line', 'bar']},
						restore: {},
						saveAsImage: {}
					}
				},
				xAxis: {
					type: 'category',
					boundaryGap: false,
					data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
				},
				yAxis: {
					type: 'value',
					axisLabel: {
						formatter: '{value} °C'
					}
				},
				series: [
					{
						name: '最高气温',
						type: 'line',
						data: [11, 11, 15, 13, 12, 13, 10],
						markPoint: {
							data: [
								{type: 'max', name: '最大值'},
								{type: 'min', name: '最小值'}
							]
						},
						markLine: {
							data: [
								{type: 'average', name: '平均值'}
							]
						}
					},
					{
						name: '最低气温',
						type: 'line',
						data: [1, -2, 2, 5, 3, 2, 0],
						markPoint: {
							data: [
								{type: 'max', name: '最大值'},
								{type: 'min', name: '最小值'}
							]
						},
						markLine: {
							data: [
								{type: 'average', name: '平均值'}
							]
						}
					}
				]
			};
	 */

	var series = make([]interface{}, 0)

	var high_temp_data = map[string]interface{}{
		"name": "最高气温",
		"type": "line",
		"data": []int64{11, 11, 15, 13, 12, 13, 10},
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

	var low_temp_data = map[string]interface{}{
		"name": "最低气温",
		"type": "line",
		"data": []int64{1, -2, 2, 5, 3, 2, 0},
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

	series = append(series, high_temp_data)
	series = append(series, low_temp_data)

	var options = map[string]interface{}{
		"title": map[string]interface{}{
			"text":    "未来一周气温变化",
			"subtext": "纯属虚构",
		},
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"legend": map[string]interface{}{
			"data": []string{"最高气温", "最低气温"},
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
			"data":        []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
		},
		"yAxis": map[string]interface{}{
			"type": "value",
			"axisLabel": map[string]interface{}{
				"formatter": "{value} °C",
			},
		},
		"series": series,
	}

	return c.RenderJson(options)
}
