package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"fmt"
	"matrix/modules/erp/models/ddlc_models"
	"unicode/utf8"
)

type ErpOrderWaybillNumber struct {
	*revel.Controller
}

func (c ErpOrderWaybillNumber) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var waybill_number string
	c.Params.Bind(&waybill_number, "waybill_number")

	c.Validation.Required(order_id).Message("订单Id不能为空！");
	c.Validation.Required(waybill_number).Message("运单号不能为空！")

	if c.Validation.HasErrors() {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	if utf8.RuneCountInString(waybill_number) > 100 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "录入的运单号太长了！",
		})
	}

	// order

	var order ddlc_models.OrderInfo
	has, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单不存在！",
		})
	}

	order.WaybillNumber = waybill_number

	_, err = db_session.Id(order_id).Cols("waybill_number").Update(&order)
	core.HandleError(err)

	// order logistics

	var order_logistics ddlc_models.OrderLogisticsInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order_id)).Get(&order_logistics)
	core.HandleError(err)

	order_logistics.TrackingNum = waybill_number

	_, err = db_session.Id(order_logistics.Id).Cols("tracking_num").Update(&order_logistics)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}
