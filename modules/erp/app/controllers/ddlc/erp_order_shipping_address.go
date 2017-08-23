package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/models/ddlc_models"
)

type ErpOrderShippingAddress struct {
	*revel.Controller
}

func (c ErpOrderShippingAddress) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	db_session := c.DbSession

	var order_shipping_address ddlc_models.OrderShippingAddressInfo
	err := json.Unmarshal(c.GetRequestBody(), &order_shipping_address)
	core.HandleError(err)

	c.Validation.Required(order_shipping_address.Name).Message("Name 不能为空！")
	c.Validation.Required(order_shipping_address.Street1).Message("Street1 不能为空！")
	c.Validation.Required(order_shipping_address.CityName).Message("City Name 不能为空！")
	c.Validation.Required(order_shipping_address.Province).Message("Province 不能为空！")
	c.Validation.Required(order_shipping_address.CountryName).Message("Country Name 不能为空！")
	c.Validation.Required(order_shipping_address.PostalCode).Message("Post Code 不能为空！")

	if c.Validation.HasErrors() {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var affected int64
	if order_shipping_address.Id == 0 {
		affected, err = db_session.Insert(&order_shipping_address)
		core.HandleError(err)
	} else {
		//判断订单当前状态是不是可以修改
		var order_address_in_db ddlc_models.OrderShippingAddressInfo
		_, err := db_session.Id(order_shipping_address.Id).Get(&order_address_in_db)
		core.HandleError(err)

		var order ddlc_models.OrderInfo
		_, err = db_session.Id(order_address_in_db.OrderId).Get(&order)
		core.HandleError(err)

		if order_service.Can_order_modify(&order) == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能保存订单地址！",
			})
		}

		var update_cols = []string{"name", "street1", "street2", "city_name", "province", "country", "country_name", "postal_code", "phone"}
		affected, err = db_session.Id(order_shipping_address.Id).Cols(update_cols...).Update(&order_shipping_address)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}
