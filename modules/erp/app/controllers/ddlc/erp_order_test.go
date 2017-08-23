package controllers

import (
	"testing"
	"matrix/modules/erp/models"
	"matrix/core"
)

func TestOrderCanFaHouStatus(t *testing.T) {

	/*
	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或特殊hold状态，不能发货！",
			})
		}
	}
	 */

	var order_list = make([]ddlc_models.OrderInfo, 0)
	order_list = append(order_list, ddlc_models.OrderInfo{StatusId: models.ORDER_STATUS_PHZ})
	//order_list = append(order_list, ddlc_models.OrderInfo{StatusId: models.ORDER_STATUS_HOLD})
	order_list = append(order_list, ddlc_models.OrderInfo{StatusId: models.ORDER_STATUS_HOLD, HoldTypeId: core.MakeNullInt(ddpz_models.ORDER_HOLD_TYPE_PICK_UP, true)})

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			t.Error("不可以发货")
		}
	}
}
