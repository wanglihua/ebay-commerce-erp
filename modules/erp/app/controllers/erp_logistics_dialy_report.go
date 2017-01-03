
package controllers

import (
	"github.com/revel/revel"
)

type ErpLogisticsDialyReport struct {
	*revel.Controller
}

func (c ErpLogisticsDialyReport) Index() revel.Result {
	return c.RenderTemplate("erp/logistics_dialy_report/logistics_dialy_report_index.html")
}
