package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpRefundReport struct {
	*revel.Controller
}

const REPORT_TYPE_DAY int64 = 1

const REPORT_TYPE_MONTH int64 = 2

func (c ErpRefundReport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	// var db_session = c.DbSession

	c.RenderArgs["report_type_day"] = REPORT_TYPE_DAY
	c.RenderArgs["report_type_month"] = REPORT_TYPE_MONTH

	return c.RenderTemplate("erp/ppgl/refund_report/refund_report_index.html")
}

func (c ErpRefundReport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	return c.RenderJson(core.GridResult{
		Data:  nil,
		Total: 0,
	})
}
