package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpOrderNew struct {
	*revel.Controller
}

func (c ErpOrderNew) Create() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	return c.RenderTemplate("erp/ddlc/order_new/order_new_create.html")
}

func (c ErpOrderNew) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)
	
	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单创建成功！",
	})
}
