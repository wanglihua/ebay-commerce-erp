package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/modules/erp/service/product_sale_analysis_service"
)

type ErpProductSaleDayAnalysis struct {
	*revel.Controller
}

func (c ErpProductSaleDayAnalysis) Index() revel.Result {
	return c.RenderTemplate("erp/spbb/product_sale_day_analysis/product_sale_day_analysis_index.html")
}

func (c ErpProductSaleDayAnalysis) PageData() revel.Result {
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

func (c ErpProductSaleDayAnalysis) ListData() revel.Result {
	db_session := c.DbSession

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")

	var analysis_type_id int64
	c.Params.Bind(&analysis_type_id, "analysis_type_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_sale_day_analysis_item_t struct {
		ProductSaleDayAnalysis models.ProductSaleDayAnalysisInfo `xorm:"extends" json:"da"`
		Stock                  models.StockInfo                  `xorm:"extends" json:"s"`
	}

	/*
SELECT da.*, s.*
FROM erp_product_sale_day_analysis da
    LEFT JOIN erp_stock s ON da.stock_id = s.id
	 */

	var sb = core.NewSqlBuilder()

	sb.Select("da.*, s.*")
	sb.Count("count(1) AS da_count")
	sb.From(`erp_product_sale_day_analysis da
    LEFT JOIN erp_stock s ON da.stock_id = s.id`)

	var sql_args = make([]interface{}, 0)

	if sku_query != "" {
		filter = core.SqlAnd(filter, "da.sku = ?")
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
		sb.OrderBy("da.id ASC")
	}

	sb.Limit(offset, limit)

	product_sale_day_analysis_list := make([]product_sale_day_analysis_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList(), sql_args...).Find(&product_sale_day_analysis_list)
	core.HandleError(err)

	for index, product_sale_day_analysis := range product_sale_day_analysis_list {
		if product_sale_day_analysis.ProductSaleDayAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL {
			product_sale_day_analysis_list[index].ProductSaleDayAnalysis.TypeName = "实际仓库发货"
		}

		if product_sale_day_analysis.ProductSaleDayAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC {
			product_sale_day_analysis_list[index].ProductSaleDayAnalysis.TypeName = "理论仓库发货"
		}

		if product_sale_day_analysis.ProductSaleDayAnalysis.Type == models.PRODUCT_SALE_ANALYSIS_TYPE_SKU {
			product_sale_day_analysis_list[index].ProductSaleDayAnalysis.TypeName = "不分仓统计"
		}

		var product models.ProductInfo
		_, err := db_session.Id(product_sale_day_analysis.ProductSaleDayAnalysis.ProductId).Get(&product)
		core.HandleError(err)

		product_sale_day_analysis_list[index].ProductSaleDayAnalysis.ProductName = product.CnName
	}

	var da_count struct {
		Count int64 `xorm:"'da_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount(), sql_args...).Get(&da_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_sale_day_analysis_list,
		Total: da_count.Count,
	})
}

func (c ErpProductSaleDayAnalysis) Analyze() revel.Result {

	go product_sale_analysis_service.ProductSaleAnalyzeSingleDay()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("立即启动 近%d天日销售情况单日分析任务成功！", product_sale_analysis_service.DaysToAnalyze),
	})
}
