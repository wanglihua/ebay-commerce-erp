package controllers

import (
	"github.com/revel/revel"
)

type ErpSingleProductSalesQuery struct {
	*revel.Controller
}

func (c ErpSingleProductSalesQuery) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/single_product_sales_query/single_product_sales_query_index.html")
}
