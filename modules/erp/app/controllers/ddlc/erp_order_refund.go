package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/service/order_paypal_service"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/core"
)

type ErpOrderRefund struct {
	*revel.Controller
}

func (c ErpOrderRefund) Refund() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_REFUND, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "refund.order_id")

	var amount float64
	c.Params.Bind(&amount, "refund.amount")

	var memo string
	c.Params.Bind(&memo, "refund.memo")

	var result = order_paypal_service.Order_refund(order_id, amount, memo, db_session)

	return c.RenderJson(result)
}

func (c ErpOrderRefund) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var refund_list = make([]ddlc_models.OrderRefundInfo, 0)
	err := db_session.Where("order_id = ?", order_id).Find(&refund_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"refund_list": refund_list,
	})
}
