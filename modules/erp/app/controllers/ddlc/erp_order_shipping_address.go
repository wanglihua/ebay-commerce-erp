package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/service/order_history_service"
	"fmt"
	"bytes"
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

	var order ddlc_models.OrderInfo
	_, err = db_session.Id(order_shipping_address.OrderId).Get(&order)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	var affected int64
	if order_shipping_address.Id == 0 {
		affected, err = db_session.Insert(&order_shipping_address)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单地址新增， 工号： %s", login_user.LoginName), db_session)
	} else {
		//判断订单当前状态是不是可以修改
		var order_address_in_db ddlc_models.OrderShippingAddressInfo
		_, err := db_session.Id(order_shipping_address.Id).Get(&order_address_in_db)
		core.HandleError(err)

		if order_service.Can_order_modify(&order) == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能保存订单地址！",
			})
		}

		order_shipping_address.Modified = 1

		// var update_cols = []string{"name", "street1", "street2", "city_name", "province", "country", "country_name", "postal_code", "phone", "modified"}
		var update_cols = []string{"modified"}

		var buffer bytes.Buffer

		if order_address_in_db.Name != order_shipping_address.Name {
			buffer.WriteString(fmt.Sprintf("name: %s => %s, ", order_address_in_db.Name, order_shipping_address.Name))
			update_cols = append(update_cols, "name")
		}

		if order_address_in_db.Street1 != order_shipping_address.Street1 {
			buffer.WriteString(fmt.Sprintf("street1: %s => %s, ", order_address_in_db.Street1, order_shipping_address.Street1))
			update_cols = append(update_cols, "street1")
		}

		if order_address_in_db.Street2 != order_shipping_address.Street2 {
			buffer.WriteString(fmt.Sprintf("street2: %s => %s, ", order_address_in_db.Street2, order_shipping_address.Street2))
			update_cols = append(update_cols, "street2")
		}

		if order_address_in_db.CityName != order_shipping_address.CityName {
			buffer.WriteString(fmt.Sprintf("city: %s => %s, ", order_address_in_db.CityName, order_shipping_address.CityName))
			update_cols = append(update_cols, "city_name")
		}

		if order_address_in_db.Province != order_shipping_address.Province {
			buffer.WriteString(fmt.Sprintf("state: %s => %s, ", order_address_in_db.Province, order_shipping_address.Province))
			update_cols = append(update_cols, "province")
		}

		if order_address_in_db.Country != order_shipping_address.Country {
			buffer.WriteString(fmt.Sprintf("country code: %s => %s, ", order_address_in_db.Country, order_shipping_address.Country))
			update_cols = append(update_cols, "country")
		}

		if order_address_in_db.CountryName != order_shipping_address.CountryName {
			buffer.WriteString(fmt.Sprintf("country: %s => %s, ", order_address_in_db.CountryName, order_shipping_address.CountryName))
			update_cols = append(update_cols, "country_name")
		}

		if order_address_in_db.PostalCode != order_shipping_address.PostalCode {
			buffer.WriteString(fmt.Sprintf("postalcode: %s => %s, ", order_address_in_db.PostalCode, order_shipping_address.PostalCode))
			update_cols = append(update_cols, "postal_code")
		}

		if order_address_in_db.Phone != order_shipping_address.Phone {
			buffer.WriteString(fmt.Sprintf("phone: %s => %s, ", order_address_in_db.Phone, order_shipping_address.Phone))
			update_cols = append(update_cols, "phone")
		}

		affected, err = db_session.Id(order_shipping_address.Id).Cols(update_cols...).Update(&order_shipping_address)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单地址修改：%s， 工号： %s", buffer.String(), login_user.LoginName), db_session)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}
