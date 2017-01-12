package controllers

import (
	"github.com/revel/revel"
)

type ErpMonthSalesAmountReport struct {
	*revel.Controller
}

func (c ErpMonthSalesAmountReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/month_sales_amount_report/month_sales_amount_report_index.html")
}
