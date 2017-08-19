package controllers

import (
	"github.com/revel/revel"

	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpProductFeedbackExcel struct {
	*revel.Controller
}

func (c ErpProductFeedbackExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	return c.RenderEmpty()
}

func (c ErpProductFeedbackExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	return c.RenderEmpty()
}
