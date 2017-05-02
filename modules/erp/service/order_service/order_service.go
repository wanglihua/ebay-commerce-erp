package order_service

// 计算GST消费税
func Get_order_gst_fee(amount_paid float64) float64 {
	return amount_paid * 0.1
}
