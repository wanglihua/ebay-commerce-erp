package controllers

import (
	"github.com/revel/revel"
)

type ErpHotSalesProductRank struct {
	*revel.Controller
}

func (c ErpHotSalesProductRank) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/hot_sales_product_rank/hot_sales_product_rank_index.html")
}
