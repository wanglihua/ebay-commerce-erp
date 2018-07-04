package controllers

import (
	"github.com/revel/revel"
)

type ErpRefundReport struct {
	*revel.Controller
}

func (c ErpRefundReport) Index() revel.Result {

	// var db_session = c.DbSession

	return c.RenderTemplate("erp/ppgl/refund_report/refund_report_index.html")
}
