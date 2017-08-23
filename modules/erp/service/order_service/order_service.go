package order_service

import (
	"matrix/modules/erp/models/ddlc_models"
)

// 计算GST消费税，这里暂时先这样
func Get_order_gst_fee(p_order *ddlc_models.OrderInfo) float64 {
	return p_order.AmountPaid * 0.1
}
