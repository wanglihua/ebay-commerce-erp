package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
	"strconv"
	"matrix/modules/erp/models/ddpz_models"
	"fmt"
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

	c.RenderArgs["order_id"] = order_id
	c.RenderArgs["order_item_list"] = order_item_list
	return c.RenderTemplate("erp/ddlc/order_split/order_split.html")
}

func (c ErpOrderSplit) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SPLIT, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	type item_quantity_t struct {
		ItemId int64
		Qty    int
	}

	var item_quantity_list = make([]item_quantity_t, 0)
	for form_key, form_value := range c.Params.Form {
		if strings.HasSuffix(form_key, "_qty") {
			item_id, err := strconv.ParseInt(form_key[0:len(form_key)-4], 10, 64)
			core.HandleError(err)

			qty, err := strconv.ParseInt(form_value[0], 10, 64)
			if err != nil {
				return c.RenderJson(map[string]interface{}{
					"success": false,
					"message": "分割数量输入有误！",
				})
			}

			var item_quantity item_quantity_t
			item_quantity.ItemId = item_id
			item_quantity.Qty = int(qty)

			item_quantity_list = append(item_quantity_list, item_quantity)
		}
	}

	var order ddlc_models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	/*
	ORDER_STATUS_YCJ            = 3850 // 订单状态 已创建
	ORDER_STATUS_YZF            = 3900 // 订单状态 已支付
	ORDER_STATUS_HOLD           = 3950 // 订单状态 Hold
	ORDER_STATUS_PROCESSED_HOLD = 3975 // 订单状态 已处理Hold
	ORDER_STATUS_DPH            = 4000 // 订单状态 待配货
	 */

	if order.StatusId > models.ORDER_STATUS_DPH {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "只有 已创建，已支付，Hold，已处理Hold，待配货 状态的订单才能分割！",
		})
	}

	//订单类型都要设置为 分割订单类型
	order.TypeId = ddpz_models.ORDER_TYPE_CFDD
	_, err = db_session.Id(order.Id).Cols("type_id").Update(&order)
	core.HandleError(err)

	// 复制 订单 和 和订单地址

	split_count, err := db_session.Where(fmt.Sprintf("code LIKE '%sS%%'", order.Code)).Count(&ddlc_models.OrderInfo{})
	core.HandleError(err)

	// 复制订单
	var split_order ddlc_models.OrderInfo
	split_order = order
	split_order.Id = 0
	split_order.Code = fmt.Sprintf("%sS%d", order.Code, split_count+1)
	split_order.TranId = fmt.Sprintf("%sS%d", order.TranId, split_count+1)
	split_order.FullTranId = fmt.Sprintf("%sS%d", order.FullTranId, split_count+1)

	split_order.TagId = 0
	split_order.StockId = 0
	split_order.CompanyId = 0
	split_order.ModeId = 0
	split_order.EbayOrderId = ""

	_, err = db_session.Insert(&split_order)
	core.HandleError(err)

	// 复制原来订单的地址
	var order_address ddlc_models.OrderShippingAddressInfo
	_, err = db_session.Id(order.AddressId).Get(&order_address)
	core.HandleError(err)

	var split_order_address ddlc_models.OrderShippingAddressInfo
	split_order_address = order_address
	split_order_address.Id = 0
	split_order_address.OrderId = split_order.Id
	_, err = db_session.Insert(&split_order_address)
	core.HandleError(err)

	// 更新新订单的地址
	split_order.AddressId = split_order_address.Id
	_, err = db_session.Id(split_order.Id).Cols("address_id").Update(&split_order)

	for _, item_quantity := range item_quantity_list {
		var order_item ddlc_models.OrderItemInfo
		_, err = db_session.Id(item_quantity.ItemId).Get(&order_item)
		core.HandleError(err)

		order_item.Quantity = order_item.Quantity - item_quantity.Qty
		_, err = db_session.Id(order_item.Id).Cols("quantity").Update(&order_item)
		core.HandleError(err)

		var split_order_item = order_item
		split_order_item.Id = 0
		split_order_item.OrderId = split_order.Id
		split_order_item.Quantity = item_quantity.Qty
		_, err = db_session.Insert(&split_order_item)
		core.HandleError(err)

	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单分割成功！",
	})
}
