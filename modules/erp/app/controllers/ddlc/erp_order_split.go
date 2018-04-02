package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/core"
)

type ErpOrderSplit struct {
	*revel.Controller
}

func (c ErpOrderSplit) Show() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SPLIT, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var order_item_list = make([]ddlc_models.OrderItemInfo, 0)
	err := db_session.Where("order_id = ?", order_id).Find(&order_item_list)
	core.HandleError(err)

	c.RenderArgs["order_item_list"] = order_item_list
	return c.RenderTemplate("erp/ddlc/order_split/order_split.html")
}

func (c ErpOrderSplit) Split() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SPLIT, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单分割成功！",
	})
}
