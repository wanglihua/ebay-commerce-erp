package order_service

import (
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddlc_models"
)

// 订单可不可以修改
// 不是 已发货 已完成 已取消 状态的都能修改
func Can_order_modify(p_order *ddlc_models.OrderInfo) bool {
	return p_order.StatusId != models.ORDER_STATUS_YFH &&
		p_order.StatusId != models.ORDER_STATUS_YWC &&
		p_order.StatusId != models.ORDER_STATUS_YQX
}
