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
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/order_history_service"
	"time"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/service/logistics_setting_service"
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

	// ddpz_models.ORDER_TYPE_SGCJ //手工创建

	var order_type_list = make([]ddpz_models.OrderTypeInfo, 0)
	core.HandleError(db_session.In("id", []int64{ddpz_models.ORDER_TYPE_SGCJ}).Find(&order_type_list))

	var status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.In("id", []int64{models.ORDER_STATUS_YZF}).Where("type = ?", models.STATUS_TYPE_ORDER).Find(&status_list))

	var hold_type_list = make([]ddpz_models.OrderHoldTypeInfo, 0)
	core.HandleError(db_session.In("id", []int64{ddpz_models.ORDER_HOLD_TYPE_PICK_UP}).Find(&hold_type_list))

	var shop_list = make([]xtpz_models.ShopAccountInfo, 0)
	core.HandleError(db_session.Find(&shop_list))

	var pickup_pay_method_list = make([]ddpz_models.PickupPayMethodInfo, 0)
	core.HandleError(db_session.Find(&pickup_pay_method_list))

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))

	var au_stock_list = make([]sppz_models.StockInfo, 0)

	var login_user = core.GetLoginUser(c.Session)

	/*
SELECT count(*) AS gu_count
FROM auth_group_user gu
WHERE gu.group_id = 11300 AND gu.user_id = 22;
	 */

	var GuCount struct {
		Count int64 `xorm:"'gu_count'"`
	}

	var gu_count_sql = `SELECT count(*) AS gu_count
FROM auth_group_user gu
WHERE gu.group_id = %d AND gu.user_id = %d;`
	gu_count_sql = fmt.Sprintf(gu_count_sql, auth_models.GROUP_PICK_CREW, login_user.UserId)
	_, err := db_session.Sql(gu_count_sql).Get(&GuCount)
	core.HandleError(err)

	if GuCount.Count != 0 { // 如果是拣货人员
		/*
SELECT s.*
FROM erp_stock_pick_crew spc
    INNER JOIN erp_stock s ON s.id = spc.stock_id
WHERE s.stock_type_id = 300 AND spc.user_id = 22;
		 */

		var stock_list_sql = `SELECT s.*
FROM erp_stock_pick_crew spc
    INNER JOIN erp_stock s ON s.id = spc.stock_id
WHERE s.stock_type_id = %d AND spc.user_id = %d;`
		stock_list_sql = fmt.Sprintf(stock_list_sql, sppz_models.STOCK_TYPE_AU, login_user.UserId)
		core.HandleError(db_session.Sql(stock_list_sql).Find(&au_stock_list))
	} else {
		core.HandleError(db_session.Where("stock_type_id = ?", sppz_models.STOCK_TYPE_AU).Find(&au_stock_list))
	}

	var stock_id int64 = 0
	if len(au_stock_list) > 0 {
		stock_id = au_stock_list[0].Id
	}

	return c.RenderJson(map[string]interface{}{
		"type_list":              order_type_list,
		"status_list":            status_list,
		"hold_type_list":         hold_type_list,
		"shop_list":              shop_list,
		"pickup_pay_method_list": pickup_pay_method_list,
		"user_list":              user_list,
		"creator":                fmt.Sprintf("%s(%s)", login_user.LoginName, login_user.NickName),
		"stock_list":             au_stock_list,
		"stock_id":               stock_id,
		"create_time":            core.LongTime(time.Now()),
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

	if order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP && order.PaidTime.Valid == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请填写支付时间",
		})
	}

	/*
	if order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP && order.StatusId != models.ORDER_STATUS_YCJ {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请创建已支付的非PICKUP订单！",
		})
	}
	*/

	var aud_exchange_rate_setting = exchange_rate_service.Get_aud_exchange_rate(db_session)

	order.ExchangeRate = aud_exchange_rate_setting

	order.AmountPaidCurrency = "AUD"
	order.AmountPaidRmb = order.AmountPaid * aud_exchange_rate_setting

	order.AdjustmentAmount = order.AmountPaid
	order.AdjustmentAmountCurrency = "AUD"
	order.AdjustmentAmountRmb = order.AmountPaid * aud_exchange_rate_setting

	order.TotalAmount = order.AmountPaid
	order.TotalAmountCurrency = "AUD"
	order.TotalAmountRmb = order.AmountPaid * aud_exchange_rate_setting

	order.SubTotalAmount = order.AmountPaid
	order.SubTotalAmountCurrency = "AUD"
	order.SubTotalAmountRmb = order.AmountPaid * aud_exchange_rate_setting

	order.StatusId = models.ORDER_STATUS_YZF

	var order_in_db ddlc_models.OrderInfo
	has, err := db_session.Where(fmt.Sprintf("code = '%s'", order.Code)).Get(&order_in_db)
	core.HandleError(err)

	if has == true {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单编号已存在！",
		})
	}

	_, err = db_session.Insert(&order)
	core.HandleError(err)

	var shipping_address = form.ShippingAddress
	shipping_address.OrderId = order.Id
	shipping_address.AddressOwner = "eBay"
	_, err = db_session.Insert(&shipping_address)
	core.HandleError(err)

	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {

		var pickup_info = form.PickupInfo

		var login_user = core.GetLoginUser(c.Session)

		pickup_info.Creator = login_user.LoginName
		pickup_info.CreatorId = login_user.UserId

		if pickup_info.IsPaid == 1 {
			// 如果PICK UP下面选择了已收款，则建单时间默认为支付时间
			order.PaidTime = core.MakeNullLongTime(time.Time(order.CreateTime), true)
			_, err := db_session.Id(order.Id).Cols("paid_time").Update(&order)
			core.HandleError(err)

			//如果选择了已收款，那么建单人默认为收款人。

			pickup_info.Payee = pickup_info.Creator
			pickup_info.PayeeId = pickup_info.CreatorId
		}

		pickup_info.OrderId = order.Id

		_, err = db_session.Insert(&pickup_info)
		core.HandleError(err)

	}

	var stock_id = form.StockId

	var logistics_setting = logistics_setting_service.Get_logistics_setting(db_session)

	var order_logistics ddlc_models.OrderLogisticsInfo
	order_logistics.OrderId = order.Id
	order_logistics.StockId = stock_id
	order_logistics.AddressId = shipping_address.Id
	order_logistics.Ptfl = 0
	order_logistics.Jyfl = logistics_setting.Jyfl

	// pick_fee_service.Get_product_pick_fee(whole_weight, p_ol.pick_fee_setting)

	_, err = db_session.Insert(&order_logistics)
	core.HandleError(err)

	order.AddressId = shipping_address.Id
	order.StockId = stock_id
	_, err = db_session.Id(order.Id).Cols("address_id", "stock_id").Update(&order)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单手工创建， 工号： %s", login_user.LoginName), db_session)

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "订单创建成功，请添加明细！",
		"order_id": order.Id,
	})
}

type OrderCreateForm struct {
	Order           ddlc_models.OrderInfo                `json:"order"`
	ShippingAddress ddlc_models.OrderShippingAddressInfo `json:"shipping_address"`
	PickupInfo      ddlc_models.PickupOrderInfo          `json:"pickup_info"`
	StockId         int64                                `json:"stock_id"`
}

func (f *OrderCreateForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Order.Code).Message("订单编码不能为空！")
	validation.Required(f.Order.CreatedTime).Message("订单创建时间不能为空！")
	validation.Required(f.Order.StatusId).Message("订单状态不能为空！")
	validation.Required(f.Order.TypeId).Message("订单状态不能为空！")
	validation.Required(f.Order.ShopId).Message("订单来源不能为空！")

	return validation.HasErrors() == false
}
