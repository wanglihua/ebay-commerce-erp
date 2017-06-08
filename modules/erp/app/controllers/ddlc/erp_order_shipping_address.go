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

func (c ErpOrderShippingAddress) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var address_id int64
	c.Params.Bind(&address_id, "address_id")

	var shipping_address models.OrderShippingAddressInfo
	_, err := db_session.Id(address_id).Get(&shipping_address)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success":          true,
		"shipping_address": shipping_address,
	})
}

func (c ErpOrderShippingAddress) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var shipping_address models.OrderShippingAddressInfo
	err := json.Unmarshal(c.GetRequestBody(), &shipping_address)
	core.HandleError(err)

	return c.RenderEmpty()
}
