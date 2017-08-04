package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/product_sale_analysis_service"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpProductSaleAnalysis struct {
	*revel.Controller
}

func (c ErpProductSaleAnalysis) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS, c.Session)

	var db_session = c.DbSession

	var product_sale_analysis_setting = product_sale_analysis_setting_service.Get_product_sale_analysis_setting(db_session)

	c.RenderArgs["level1_days"] = product_sale_analysis_setting.Level1Days
	c.RenderArgs["level2_days"] = product_sale_analysis_setting.Level2Days
	c.RenderArgs["level3_days"] = product_sale_analysis_setting.Level3Days
	c.RenderArgs["level4_days"] = product_sale_analysis_setting.Level4Days
	c.RenderArgs["level5_days"] = product_sale_analysis_setting.Level5Days

	return c.RenderTemplate("erp/spbb/product_sale_analysis/product_sale_analysis_index.html")
}

func (c ErpProductSaleAnalysis) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS, c.Session)

	return c.RenderJson(map[string]interface{}{
		"analysis_type_id": 0,
		"analysis_type_list": []map[string]interface{}{
			map[string]interface{}{
				"id":   models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL,
				"name": "实际仓库发货",
			},
			map[string]interface{}{
				"id":   models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC,
				"name": "理论仓库发货",
			},
			map[string]interface{}{
				"id":   models.PRODUCT_SALE_ANALYSIS_TYPE_SKU,
				"name": "不分仓统计",
			},
		},
	})
}

func (c ErpProductSaleAnalysis) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS, c.Session)

	db_session := c.DbSession

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")

	var analysis_type_id int64
	c.Params.Bind(&analysis_type_id, "analysis_type_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_sale_analysis_item_t struct {
		ProductSaleAnalysis models.ProductSaleAnalysisInfo `xorm:"extends" json:"sa"`
		Stock               models.StockInfo               `xorm:"extends" json:"s"`
	}

	/*
SELECT sa.*, s.*
FROM erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id
	 */

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

	if analysis_type_id != 0 {
		filter = core.SqlAnd(filter, "type = ?")
		sql_args = append(sql_args, analysis_type_id)
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sa.id ASC")
	}

	sb.Limit(offset, limit)

	product_sale_analysis_list := make([]product_sale_analysis_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList(), sql_args...).Find(&product_sale_analysis_list)
	core.HandleError(err)

	for index, product_sale_day_analysis := range product_sale_analysis_list {
		if product_sale_day_analysis.ProductSaleAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "实际仓库发货"
		}

		if product_sale_day_analysis.ProductSaleAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "理论仓库发货"
		}

		if product_sale_day_analysis.ProductSaleAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_SKU {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "不分仓统计"
		}

		var product models.ProductInfo
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
}

func (c ErpProductSaleAnalysis) Analyze() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS, c.Session)
	
	go product_sale_analysis_service.ProductSaleAnalyzeMultipleDay()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "立即启动分档多日销售情况析任务成功！",
	})
}
