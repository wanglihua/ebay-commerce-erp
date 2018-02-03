package order_service

import (
	"matrix/modules/erp/models/ddlc_models"
)

// 根据数据计算订单利润，单位是澳元
func CalculateOrderProfit(p_order *ddlc_models.OrderInfo, p_order_logistics *ddlc_models.OrderLogisticsInfo) float64 {
	var order_gst = Get_order_gst_fee(p_order)

	// 如果有实际物流费用，就按实际物流费用来计算。
	// 如果是匹配物流的时候调用这个方法，一般这时实际物流费用是0
	var order_logistics_total_fee = p_order_logistics.TotalFee
	if p_order_logistics.ActualFee > 0 {
		order_logistics_total_fee = p_order_logistics.ActualFee
	}

	// return p_order.AmountPaid + p_order.ExpressFee - p_order.AmountPaid*(p_order_logistics.Ptfl+p_order_logistics.Jyfl) - order_logistics_total_fee - p_order_logistics.PickFee - order_gst - p_order_logistics.CostAmount

	// return p_order.AmountPaid + p_order.ExpressFee - p_order.Ptfy - p_order.Jyfy - order_logistics_total_fee - p_order_logistics.PickFee - order_gst - p_order_logistics.CostAmount
	return p_order.AmountPaid - p_order.Ptfy - p_order.Jyfy - order_logistics_total_fee - p_order_logistics.PickFee - order_gst - p_order_logistics.CostAmount
}
