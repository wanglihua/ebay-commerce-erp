package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpOrderShippingAddress struct {
	*revel.Controller
}

func (c ErpOrderShippingAddress) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	db_session := c.DbSession

	var order_shipping_address models.OrderShippingAddressInfo
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
		var update_cols = []string{"name", "street1", "street2", "city_name", "province", "country", "country_name", "postal_code", "phone"}
		affected, err = db_session.Id(order_shipping_address.Id).Cols(update_cols...).Update(&order_shipping_address)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}
