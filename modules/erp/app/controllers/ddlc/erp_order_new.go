package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
	"matrix/core"
	"fmt"
	auth_models "matrix/modules/auth/models"
	"encoding/json"
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

	var pickup_pay_method_list = make([]models.PickupPayMethodInfo, 0)
	core.HandleError(db_session.Find(&pickup_pay_method_list))

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))

	var login_user = core.GetLoginUser(c.Session)

	return c.RenderJson(map[string]interface{}{
		"type_list":              order_type_list,
		"status_list":            status_list,
		"hold_type_list":         hold_type_list,
		"shop_list":              shop_list,
		"pickup_pay_method_list": pickup_pay_method_list,
		"user_list":              user_list,
		"creator":                fmt.Sprintf("%s(%s)", login_user.LoginName, login_user.NickName),
	})
}

func (c ErpOrderNew) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	var db_session = c.DbSession

	var form OrderCreateForm
	err := json.Unmarshal(c.GetRequestBody(), &form)
	core.HandleError(err)



	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var order = form.Order
	_, err = db_session.Insert(&order)
	core.HandleError(err)

	var shipping_address = form.ShippingAddress
	shipping_address.OrderId = order.Id
	shipping_address.AddressOwner = "eBay"
	_, err = db_session.Insert(&shipping_address)
	core.HandleError(err)

	if order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PICK_UP {
		var pickup_info = form.PickupInfo
		pickup_info.OrderId = order.Id

		_, err = db_session.Insert(&pickup_info)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "订单创建成功，请添加明细！",
		"order_id": order.Id,
	})
}

type OrderCreateForm struct {
	Order           models.OrderInfo                `json:"order"`
	ShippingAddress models.OrderShippingAddressInfo `json:"shipping_address"`
	PickupInfo      models.PickupOrderInfo          `json:"pickup_info"`
}

func (f *OrderCreateForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Order.Code).Message("订单编码不能为空！")
	validation.Required(f.Order.CreatedTime).Message("订单创建时间不能为空！")
	validation.Required(f.Order.StatusId).Message("订单状态不能为空！")
	validation.Required(f.Order.TypeId).Message("订单状态不能为空！")
	validation.Required(f.Order.ShopId).Message("订单来源不能为空！")

	return validation.HasErrors() == false
}
