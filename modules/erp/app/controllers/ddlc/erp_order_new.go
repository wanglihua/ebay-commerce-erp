package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
	"matrix/core"
)

type ErpOrderNew struct {
	*revel.Controller
}

func (c ErpOrderNew) Create() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	return c.RenderTemplate("erp/ddlc/order_new/order_new_create.html")
}

func (c ErpOrderNew) CreateData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	var db_session = c.DbSession

	// models.ORDER_TYPE_SGCJ //手工创建

	var order_type_list = make([]models.OrderTypeInfo, 0)
	core.HandleError(db_session.In("id", []int64{models.ORDER_TYPE_SGCJ}).Find(&order_type_list))

	var status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.In("id", []int64{models.ORDER_STATUS_YCJ, models.ORDER_STATUS_YZF}).Where("type = ?", models.STATUS_TYPE_ORDER).Find(&status_list))

	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	core.HandleError(db_session.In("id", []int64{models.ORDER_HOLD_TYPE_PICK_UP}).Find(&hold_type_list))

	var shop_list = make([]models.ShopAccountInfo, 0)
	core.HandleError(db_session.Find(&shop_list))

	return c.RenderJson(map[string]interface{}{
		"type_list":      order_type_list,
		"status_list":    status_list,
		"hold_type_list": hold_type_list,
		"shop_list":      shop_list,
	})
}

func (c ErpOrderNew) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	var db_session = c.DbSession
	_ = db_session

	// models.ORDER_TYPE_SGCJ //手工创建

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单创建成功，请添加明细！",
	})
}
