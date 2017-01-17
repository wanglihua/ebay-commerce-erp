package controllers

import (
	"github.com/revel/revel"
)

type ErpPersonSalesReport struct {
	*revel.Controller
}

func (c ErpPersonSalesReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/person_sales_report/person_sales_report_index.html")
}
