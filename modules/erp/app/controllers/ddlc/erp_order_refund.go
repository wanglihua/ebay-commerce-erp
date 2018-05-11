package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/service/order_paypal_service"
)

type ErpOrderRefund struct {
	*revel.Controller
}

func (c ErpOrderRefund) Refund() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_REFUND, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var amount float64
	c.Params.Bind(&amount, "amount")

	var memo string
	c.Params.Bind(&memo, "memo")

	var result = order_paypal_service.Order_refund(order_id, amount, memo, db_session)

	return c.RenderJson(result)
}
