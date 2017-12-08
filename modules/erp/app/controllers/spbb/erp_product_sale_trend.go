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
	"strconv"
	"time"
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

	var product_id_strs string
	c.Params.Bind(&product_id_strs, "product_ids")

	var time_now = time.Now()
	var today_date = time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, nil)

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

	var product_id_str_list = strings.Split(product_id_strs, ",")
	if len(product_id_str_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品id参数错误！",
		})
	}

	var db_session = c.DbSession

	var product_list = make([]kcgl_models.ProductInfo, 0)
	for _, product_id_str := range product_id_str_list {
		product_id, err := strconv.ParseInt(product_id_str, 10, 64)
		if err != nil {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "商品id参数错误！",
			})
		}

		var product kcgl_models.ProductInfo
		has, _ := db_session.Id(product_id).Get(&product)
		if has == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "指定的商品id不存在！",
			})
		}

		product_list = append(product_list, product)
	}

	var product_sku_list = make([]string, 0)
	for _, product := range product_list {
		product_sku_list = append(product_sku_list, product.Sku)
	}

	var time_range_list = make([]spbb_models.ProductSaleAnalysisTimeRangeInfo, 0)
	err := db_session.
		Where("begin_time >= ? AND end_time <= ?", query_begin_date, query_end_date).
		OrderBy("begin_time ASC").
		Find(&time_range_list)
	core.HandleError(err)

	var time_str_list = make([]string, 0)
	for _, time_range := range time_range_list {
		time_str_list = append(time_str_list, time.Time(time_range.BeginTime).Format("2006-01-02"))
	}

	var series = make([]interface{}, 0)

	for _, product := range product_list {
		_ = product
	}

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
			"text":    "商品每10天销量",
			"subtext": fmt.Sprintf("%s至%s", query_begin_date.Format("2006-01-02"), query_end_date.Format("2006-01-02")), // 时间范围
		},
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"legend": map[string]interface{}{
			"data": product_sku_list,
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
				"formatter": "{value} 件",
			},
		},
		"series": series,
	}

	return c.RenderJson(options)
}
