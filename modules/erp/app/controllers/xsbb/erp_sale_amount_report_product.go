package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"fmt"
	"matrix/modules/erp/service/sale_report_service"
	"time"
	"matrix/modules/erp/models/xtpz_models"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strings"
)

type ErpSaleAmountReportProduct struct {
	*revel.Controller
}

func (c ErpSaleAmountReportProduct) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT, c.Session)

	var db_session = c.DbSession

	var user_type_list = make([]user_query_type, 0)
	user_type_list = append(user_type_list, user_query_type{Id: user_query_type_sale, Name: "销售人员"})
	user_type_list = append(user_type_list, user_query_type{Id: user_query_type_service, Name: "服务人员"})
	user_type_list = append(user_type_list, user_query_type{Id: user_query_type_dev, Name: "开发人员"})

	// 先不用平台查询
	var platform_list = make([]xtpz_models.ECommercePlatformInfo, 0)
	core.HandleError(db_session.Find(&platform_list))

	var shop_list = make([]xtpz_models.ShopAccountInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&shop_list))

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	c.RenderArgs["user_type_list"] = user_type_list
	c.RenderArgs["shop_list"] = shop_list
	c.RenderArgs["user_list"] = user_list

	return c.RenderTemplate("erp/xsbb/sale_amount_report_product/sale_amount_report_product_index.html")
}

func (c ErpSaleAmountReportProduct) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT, c.Session)

	var db_session = c.DbSession

	var begin_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")

	var end_time_query core.NullLongTime
	c.Params.Bind(&end_time_query, "end_time_query")

	var platform_id_query int64
	c.Params.Bind(&platform_id_query, "platform_id_query")

	var shop_id_query int64
	c.Params.Bind(&shop_id_query, "shop_id_query")

	var user_type_query int64
	c.Params.Bind(&user_type_query, "user_type_query")

	var user_id_query int64
	c.Params.Bind(&user_id_query, "user_id_query")

	_, order, offset, limit := core.GetGridRequestParam(c.Request)

	var filter = "1=1"

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.paid_time < '%s'", end_time_query.Time.AddDate(0, 0, 1).Format("2006-01-02")))
	}

	if platform_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.platform_id = %d", platform_id_query))
	}

	if shop_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.shop_id = %d", shop_id_query))
	}

	if user_type_query != 0 && user_id_query != 0 {
		if user_type_query == user_query_type_sale {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.sale_user_id = %d", user_id_query))
		}

		if user_type_query == user_query_type_service {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.service_user_id = %d", user_id_query))
		}

		if user_type_query == user_query_type_dev {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.dev_user_id = %d", user_id_query))
		}
	}

	var search = core.GetGridRequestSearch(c.Request)
	search = strings.TrimSpace(search)
	if search != "" {
		var search_filter = `(r.sale_user_id IN (SELECT u.id FROM auth_user u
          WHERE u.login_name LIKE '%%%s%%' OR u.nick_name LIKE '%%%s%%'))
          OR (r.dev_user_id IN (SELECT u.id FROM auth_user u
          WHERE u.login_name LIKE '%%%s%%' OR u.nick_name LIKE '%%%s%%'))
          OR (r.product_id IN (SELECT p.id FROM erp_product p
          WHERE p.sku LIKE '%%%s%%' OR p.cn_name LIKE '%%%s%%'))`

		filter = core.SqlAnd(filter, fmt.Sprintf(search_filter, search, search, search, search, search, search))
	}

	/*

SELECT
    r.product_id                                    AS product_id,
    (SELECT p.sku FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_sku,
    (SELECT p.cn_name FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_name,
    (SELECT p.weight FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_weight,
    (SELECT avg(nvt.cost_price) AS cost_price
     FROM erp_inventory nvt
         INNER JOIN erp_stock stk ON stk.id = nvt.stock_id
     WHERE stk.stock_type_id = 300
           AND nvt.product_id = r.product_id)       AS cost_price,
    r.sale_user_id                                  AS sale_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.sale_user_id)                   AS sale_user_name,
    r.service_user_id                               AS service_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.service_user_id)                AS service_user_name,
    r.dev_user_id                                   AS dev_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.dev_user_id)                    AS dev_user_name,
    sum(r.paid_amount)                              AS paid_amount,
    sum(r.express_fee)                              AS express_fee,
    sum(r.income)                                   AS income,
    sum(r.cost_amount)                              AS cost_amount,
    sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
    sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
    sum(r.jyfy_amount)                              AS jyfy_amount,
    sum(r.ptfy_amount)                              AS ptfy_amount,
    sum(r.pick_fee)                                 AS pick_fee,
    sum(r.gst_fee)                                  AS gst_fee,
    sum(r.outcome)                                  AS outcome,
    sum(r.income - r.outcome)                       AS profit,
    sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
    count(DISTINCT r.order_id)                      AS order_qty,
    sum(r.product_qty)                              AS product_qty
FROM erp_sale_amount_report r
WHERE r.platform_id = 0
      AND r.shop_id = 0
      AND r.paid_time >= '2006-01-01'
      AND r.paid_time < '2017-12-30'
      AND r.product_id = 0
      AND r.sale_user_id = 0
      AND r.service_user_id = 0
      AND r.dev_user_id = 0
      AND (
          (r.sale_user_id IN (SELECT u.id FROM auth_user u
          WHERE u.login_name LIKE '%abc%' OR u.nick_name LIKE '%abc%'))
          OR (r.dev_user_id IN (SELECT u.id FROM auth_user u
          WHERE u.login_name LIKE '%abc%' OR u.nick_name LIKE '%abc%'))
          OR (r.product_id IN (SELECT p.id FROM erp_product p
          WHERE p.sku LIKE '%abc%' OR p.cn_name LIKE '%abc%'))
      )
GROUP BY r.product_id
ORDER BY product_id
LIMIT 0, 10;


SELECT count(*) AS count
FROM (
         SELECT r.product_id AS product_id
         FROM erp_sale_amount_report r
         WHERE r.platform_id = 0
               AND r.shop_id = 0
               AND r.paid_time >= '2006-01-01'
               AND r.paid_time < '2017-12-30'
               AND r.product_id = 0
               AND r.sale_user_id = 0
               AND r.service_user_id = 0
               AND r.dev_user_id = 0
         GROUP BY r.product_id
     ) sale_amount_report;

	 */

	var sql_list = `
SELECT
    r.product_id                                    AS product_id,
    (SELECT p.sku FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_sku,
    (SELECT p.cn_name FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_name,
    (SELECT p.weight FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_weight,
    (SELECT avg(nvt.cost_price) AS cost_price
     FROM erp_inventory nvt
         INNER JOIN erp_stock stk ON stk.id = nvt.stock_id
     WHERE stk.stock_type_id = 300
           AND nvt.product_id = r.product_id)       AS cost_price,
    r.sale_user_id                                  AS sale_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.sale_user_id)                   AS sale_user_name,
    r.service_user_id                               AS service_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.service_user_id)                AS service_user_name,
    r.dev_user_id                                   AS dev_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.dev_user_id)                    AS dev_user_name,
    sum(r.paid_amount)                              AS paid_amount,
    sum(r.express_fee)                              AS express_fee,
    sum(r.income)                                   AS income,
    sum(r.cost_amount)                              AS cost_amount,
    sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
    sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
    sum(r.jyfy_amount)                              AS jyfy_amount,
    sum(r.ptfy_amount)                              AS ptfy_amount,
    sum(r.pick_fee)                                 AS pick_fee,
    sum(r.gst_fee)                                  AS gst_fee,
    sum(r.outcome)                                  AS outcome,
    sum(r.income - r.outcome)                       AS profit,
    sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
    count(DISTINCT r.order_id)                      AS order_qty,
    sum(r.product_qty)                              AS product_qty
FROM erp_sale_amount_report r
WHERE %s
GROUP BY r.product_id
ORDER BY %s
LIMIT %d, %d;
`
	if order == "" {
		order = "paid_time DESC"
	}

	sql_list = fmt.Sprintf(sql_list, filter, order, offset, limit)

	var sale_amount_report_list = make([]xsbb_models.SaleAmountReportProductInfo, 0, limit)
	core.HandleError(db_session.Sql(sql_list).Find(&sale_amount_report_list))

	var sql_count = `
SELECT count(*) AS count
FROM (
         SELECT r.product_id AS product_id
         FROM erp_sale_amount_report r
         WHERE %s
         GROUP BY r.product_id
     ) sale_amount_report;;
`
	sql_count = fmt.Sprintf(sql_count, filter)

	var count core.Count
	_, err := db_session.Sql(sql_count).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  sale_amount_report_list,
		Total: count.Count,
	})
}

func (c ErpSaleAmountReportProduct) Chart() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT, c.Session)

	return c.RenderTemplate("erp/xsbb/sale_amount_report_product/sale_amount_report_product_chart.html")
}

func (c ErpSaleAmountReportProduct) ChartData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT, c.Session)

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
		amount_list = append(amount_list, report.PaidAmount)
	}

	var serie_data = map[string]interface{}{
		"name": "商品日销售额",
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

func (c ErpSaleAmountReportProduct) Calculate() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT, c.Session)

	go sale_report_service.CalculateSaleAmountReport()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动销售情况后台分析任务成功！",
	})
}
