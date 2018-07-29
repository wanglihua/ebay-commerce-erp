package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpOrderMessage struct {
	*revel.Controller
}

func (c ErpOrderMessage) show() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SEND_MESSAGE, c.Session)

	return c.RenderTemplate("erp/ddlc/order_message/order_message_send_to_customers.html")
}

func (c ErpOrderMessage) send() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SEND_MESSAGE, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": false,
		"message": "未知回复方式！",
	})
}
