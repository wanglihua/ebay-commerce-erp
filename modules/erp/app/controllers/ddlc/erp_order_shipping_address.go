package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type ErpOrderShippingAddress struct {
	*revel.Controller
}

func (c ErpOrderShippingAddress) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	return c.RenderEmpty()
}

func (c ErpOrderShippingAddress) Save() revel.Result {
	return c.RenderEmpty()
}
