package controllers

import (
	"github.com/revel/revel"
)

type ErpDiarySalesReport struct {
	*revel.Controller
}

func (c ErpDiarySalesReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/diary_sales_report/diary_sales_report_index.html")
}
