package controllers

import (
	"github.com/revel/revel"
)

type ErpLogisticsDialyReport struct {
	*revel.Controller
}

func (c ErpLogisticsDialyReport) Index() revel.Result {
	return c.RenderTemplate("erp/ddbb/logistics_report/logistics_report_index.html")
}
