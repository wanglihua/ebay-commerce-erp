package controllers

import (
	"github.com/revel/revel"
)

type ErpHotSalesStoreRank struct {
	*revel.Controller
}

func (c ErpHotSalesStoreRank) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/hot_sales_store_rank/hot_sales_store_rank_index.html")
}
