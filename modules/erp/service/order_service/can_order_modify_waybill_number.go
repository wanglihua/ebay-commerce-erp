package order_service

import "matrix/modules/erp/models"

func Can_order_modify_waybill_number(p_order *models.OrderInfo) bool {
	return p_order.StatusId == models.ORDER_STATUS_PHZ ||
		p_order.StatusId == models.ORDER_STATUS_YFH ||
		p_order.StatusId == models.ORDER_STATUS_YWC
}
