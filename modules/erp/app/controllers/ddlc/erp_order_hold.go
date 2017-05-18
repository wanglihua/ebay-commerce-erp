package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/setting"
)

type ErpOrderHold struct {
	*revel.Controller
}

func (c ErpOrderHold) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HOLD, c.Session)

	var has_pass_order_quxiao = pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_QUXIAO, c.Session)

	var is_button_quxiao_visible = false
	if has_pass_order_quxiao == true {
		is_button_quxiao_visible = true
	}

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	c.RenderArgs["is_button_quxiao_visible"] = is_button_quxiao_visible
	return c.RenderTemplate("erp/ddlc/order_hold/order_hold_index.html")
}

func (c ErpOrderHold) PageData() revel.Result {
	var db_session = c.DbSession

	// hold类型列表
	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	err := db_session.OrderBy("ord ASC").Find(&hold_type_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"hold_type_id":   0,
		"hold_type_list": hold_type_list,
	})
}

func (c ErpOrderHold) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HOLD, c.Session)

	db_session := c.DbSession

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_HOLD))

	if hold_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
