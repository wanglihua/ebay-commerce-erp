package controllers

import (
	"github.com/revel/revel"
)

type ErpDiarySalesAmountReport struct {
	*revel.Controller
}

func (c ErpDiarySalesAmountReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/diary_sales_amount_report/diary_sales_amount_report_index.html")
}
