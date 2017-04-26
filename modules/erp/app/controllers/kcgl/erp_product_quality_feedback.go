package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type ErpProductQualityFeedback struct {
	*revel.Controller
}

func (c ErpProductQualityFeedback) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_QUALITY_FEEDBACK, c.Session)

	return c.RenderTemplate("erp/kcgl/product_quality_feedback/product_quality_feedback_index.html")
}
