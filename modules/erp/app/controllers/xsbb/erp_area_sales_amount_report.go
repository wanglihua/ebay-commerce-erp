package controllers

import (
	"github.com/revel/revel"
)

type ErpAreaSalesAmountReport struct {
	*revel.Controller
}

func (c ErpAreaSalesAmountReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/area_sales_amount_report/area_sales_amount_report_index.html")
}
