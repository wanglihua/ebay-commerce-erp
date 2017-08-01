package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductSaleDayAnalysis struct {
	*revel.Controller
}

func (c ErpProductSaleDayAnalysis) Index() revel.Result {
	return c.RenderTemplate("erp/spbb/product_sale_day_analysis/product_sale_day_analysis_index.html")
}

func (c ErpProductSaleDayAnalysis) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_sale_day_analysis_item_t struct {
		ProductSaleDayAnalysis models.ProductSaleDayAnalysisInfo `xorm:"extends" json:"da"`
		Stock                  models.StockInfo                  `xorm:"extends" json:"s"`
	}

	/*
SELECT da.*, s.*
FROM erp_product_sale_day_analysis da
    INNER JOIN erp_stock s ON da.stock_id = s.id
	 */

	var sb = core.NewSqlBuilder()

	sb.Select("da.*, s.*")
	sb.Count("count(1) AS da_count")
	sb.From(`erp_product_sale_day_analysis da
    INNER JOIN erp_stock s ON da.stock_id = s.id`)
	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("da.id ASC")
	}

	sb.Limit(offset, limit)

	product_sale_day_analysis_list := make([]product_sale_day_analysis_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_sale_day_analysis_list)
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

	_, err = db_session.Sql(sb.SqlCount()).Get(&da_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_sale_day_analysis_list,
		Total: da_count.Count,
	})
}
