package order_logistics_service

type match_logistics_mode_result struct {
	order_id    int64
	address_id  int64
	company_id  int64
	stock_id    int64
	mode_id     int64
	fee_rate_id int64
	length      float64
	width       float64
	height      float64
	dim_weight  float64
	weight      float64
	total_fee   float64
	pick_fee    float64
	cost_amount float64
	has_stock   bool
}
