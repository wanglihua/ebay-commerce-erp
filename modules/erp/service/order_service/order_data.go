package order_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/core/collections"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_hold_type_service"
	"matrix/modules/erp/service/order_type_service"
	"matrix/modules/erp/service/shop_service"
	"matrix/modules/erp/service/status_service"
)

type Order_view_item_t struct {
	Order        models.OrderInfo       `xorm:"extends" json:"o"`
	TypeName     string                 `json:"type_name"`
	HoldTypeName string                 `json:"hold_type_name"`
	PlatformName string                 `json:"platform_name"`
	ShopCode     string                 `json:"shop_code"`
	ShopName     string                 `json:"shop_name"`
	StatusName   string                 `json:"status_name"`
}

func Get_order_data_list(db_session *xorm.Session, filter string, order string, offset int64, limit int64) ([]Order_view_item_t, int64) {

	/*
SELECT o.*
FROM erp_order o
	 */

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("o.*")
	sql_builder.Count("count(o.id) as order_count")
	sql_builder.From(`erp_order o`)
	sql_builder.Where(filter)

	if order == "" {
		sql_builder.OrderBy("o.id DESC")
	} else {
		sql_builder.OrderBy(order)
	}

	sql_builder.Limit(offset, limit)

	order_list := make([]Order_view_item_t, 0)
	err := db_session.Sql(sql_builder.SqlList()).Find(&order_list)
	core.HandleError(err)

	// 加载订单列表额外信息
	var order_type_id_set = collections.NewIntSet()
	var hold_type_id_set = collections.NewIntSet()
	var shop_account_id_set = collections.NewIntSet()
	var status_id_set = collections.NewIntSet()

	for _, order := range order_list {
		order_type_id_set.Add(order.Order.TypeId)
		if order.Order.HoldTypeId.Valid == true {
			hold_type_id_set.Add(order.Order.HoldTypeId.Int64)
		}
		shop_account_id_set.Add(order.Order.ShopId)
		status_id_set.Add(order.Order.StatusId)
	}

	var order_type_id_name_map = order_type_service.Get_order_type_id_name_map(order_type_id_set.List(), db_session)
	var hold_type_id_name_map = order_hold_type_service.Get_order_hold_type_id_name_map(hold_type_id_set.List(), db_session)
	var platform_shop_id_name_map = shop_service.Get_shop_id_name_map(shop_account_id_set.List(), db_session)
	var status_id_name_map = status_service.Get_status_id_name_map(status_id_set.List(), db_session)

	for index, order := range order_list {
		order_list[index].TypeName = order_type_id_name_map[order.Order.TypeId]
		if order.Order.HoldTypeId.Valid == true {
			order_list[index].HoldTypeName = hold_type_id_name_map[order.Order.HoldTypeId.Int64]
		}

		var platform_shop_name = platform_shop_id_name_map[order.Order.ShopId]
		order_list[index].PlatformName = platform_shop_name.PlatformName
		order_list[index].ShopCode = platform_shop_name.ShopCode
		order_list[index].ShopName = platform_shop_name.ShopName

		order_list[index].StatusName = status_id_name_map[order.Order.StatusId]
	}

	var order_count struct {
		Count int64 `xorm:"'order_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&order_count)
	core.HandleError(err)

	return order_list, order_count.Count
}

func Get_order_id_list_by_product_sku_query(db_session *xorm.Session, product_sku_query string) []int64 {
	sql := `SELECT oi.order_id
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE p.sku like '%%%s%%'`
	sql = fmt.Sprintf(sql, product_sku_query)

	type order_id_result_t struct {
		OrderId int64 `json:"order_id"`
	}

	var order_id_result_list = make([]order_id_result_t, 0)
	err := db_session.Sql(sql).Find(&order_id_result_list)
	core.HandleError(err)

	var order_id_list = make([]int64, 0)
	for _, order_id_result := range order_id_result_list {
		order_id_list = append(order_id_list, order_id_result.OrderId)
	}

	return order_id_list
}

func Get_order_id_list_by_product_name_query(db_session *xorm.Session, product_name_query string) []int64 {
	sql := `SELECT oi.order_id
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE p.name like '%%%s%%'`
	sql = fmt.Sprintf(sql, product_name_query)

	type order_id_result_t struct {
		OrderId int64 `json:"order_id"`
	}

	var order_id_result_list = make([]order_id_result_t, 0)
	err := db_session.Sql(sql).Find(&order_id_result_list)
	core.HandleError(err)

	var order_id_list = make([]int64, 0)
	for _, order_id_result := range order_id_result_list {
		order_id_list = append(order_id_list, order_id_result.OrderId)
	}

	return order_id_list
}

func Get_paid_order_list(db_session *xorm.Session) []models.OrderInfo {
	var paid_order_list = make([]models.OrderInfo, 0)
	err := db_session.Sql(fmt.Sprintf("SELECT * FROM erp_order o WHERE o.status_id = %d", models.ORDER_STATUS_YZF)).Find(&paid_order_list)
	core.HandleError(err)

	return paid_order_list
}

type Order_product_info_t struct {
	OrderId   int64    `xorm:"'order_id'"`
	ProductId int64    `xorm:"'product_id'"`
	Price     float64  `xorm:"'price'"` // 方便后面的运费计算啥的，这里面的价格是 澳币
	Qty       int64    `xorm:"'qty'"`
	CostPrice int64    `xorm:"'cost_price'"` // 是 人民币 的成本价，由于运费都是用 澳币结算的，匹配物流时，应该把这个 成本价 换算成 澳币
	Length    float64  `xorm:"'length'"`
	Width     float64  `xorm:"'width'"`
	Height    float64  `xorm:"'height'"`
	Weight    float64  `xorm:"'weight'"`
}

func Get_order_product_info_list(order_id int64, db_session *xorm.Session) []Order_product_info_t {

	/*
SELECT
    oi.order_id,
    oi.item_product_id    AS product_id,
    oi.transaction_price  AS price,
    oi.quantity           AS qty,
    oi.product_cost_price AS cost_price,
    p.length              AS length,
    p.width               AS width,
    p.height              AS height,
    p.weight              AS weight
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE oi.order_id = 1
	 */

	sql := `SELECT
    oi.order_id,
    oi.item_product_id    AS product_id,
    oi.transaction_price  AS price,
    oi.quantity           AS qty,
    oi.product_cost_price AS cost_price,
    p.length              AS length,
    p.width               AS width,
    p.height              AS height,
    p.weight              AS weight
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE oi.order_id = %d`
	sql = fmt.Sprintf(sql, order_id)

	var order_product_qty_list = make([]Order_product_info_t, 0)
	err := db_session.Sql(sql).Find(&order_product_qty_list)
	core.HandleError(err)

	return order_product_qty_list
}

func Get_order_shipping_address(order_id int64, db_session *xorm.Session) models.OrderShippingAddressInfo {
	var order_shipping_address_list = make([]models.OrderShippingAddressInfo, 0)
	err := db_session.Where(fmt.Sprintf("order_id = %d", order_id)).Find(&order_shipping_address_list)
	core.HandleError(err)

	for _, order_shipping_address := range order_shipping_address_list {
		if order_shipping_address.AddressOwner == "PayPal" { // PayPal 收货地址
			return order_shipping_address
		}
	}

	for _, order_shipping_address := range order_shipping_address_list {
		if order_shipping_address.AddressOwner == "eBay" { // eBay 收货地址
			return order_shipping_address
		}
	}

	if len(order_shipping_address_list) != 0 {
		return order_shipping_address_list[0]
	}

	return models.OrderShippingAddressInfo{} //如果前面的都不成立，就返回一个默认的
}
