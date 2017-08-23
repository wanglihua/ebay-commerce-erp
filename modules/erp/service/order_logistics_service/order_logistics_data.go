package order_logistics_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
)

type order_logistics_data_t struct {
	exchange_rate           float64
	order_shipping_address  ddlc_models.OrderShippingAddressInfo
	order_product_info_list []order_service.Order_product_info_t
	products_stock_qty_list []inventory_service.Product_stock_qty
	stock_list              []sppz_models.StockInfo
}

// 一次性全部加载
func (p_d *order_logistics_data_t) load_data(p_order *ddlc_models.OrderInfo, p_order_logistics *OrderLogisticsService, db_session *xorm.Session) {
	p_d.load_exchange_rate(p_order, p_order_logistics)
	p_d.load_order_shipping_address(p_order.Id, db_session)
	p_d.load_order_product_info_list(p_order.Id, db_session)
	p_d.load_products_stock_qty_list(p_order.Id, db_session)
}

func (p_d *order_logistics_data_t) load_exchange_rate(p_order *ddlc_models.OrderInfo, p_order_logistics *OrderLogisticsService) {
	// 下载订单的时候就把当时的汇率记录了下来了
	var exchange_rate = p_order.ExchangeRate
	if exchange_rate == 0.0 {
		exchange_rate = p_order_logistics.aud_exchange_rate
	}
	p_d.exchange_rate = exchange_rate
}

func (p_d *order_logistics_data_t) load_order_shipping_address(order_id int64, db_session *xorm.Session) {
	var order_shipping_address = order_service.Get_order_shipping_address(order_id, db_session)
	p_d.order_shipping_address = order_shipping_address
}

func (p_d *order_logistics_data_t) load_order_product_info_list(order_id int64, db_session *xorm.Session) {
	p_d.order_product_info_list = order_service.Get_order_product_info_list(order_id, db_session)
}

func (p_d *order_logistics_data_t) load_products_stock_qty_list(order_id int64, db_session *xorm.Session) {

	// 获取订单中的所有商品的id列表
	var get_order_product_id_list = func() []int64 {
		var product_id_list = make([]int64, 0)
		for _, order_product_info := range p_d.order_product_info_list {
			product_id_list = append(product_id_list, order_product_info.ProductId)
		}
		return product_id_list
	}
	var order_product_id_list = get_order_product_id_list()
	p_d.products_stock_qty_list = inventory_service.Get_products_stock_qty_list(order_product_id_list, db_session)
}
