package order_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Order_view_item_t struct {
	models.OrderInfo              `xorm:"extends" json:"o"`
	models.OrderTypeInfo          `xorm:"extends" json:"ot"`
	models.OrderHoldTypeInfo      `xorm:"extends" json:"ht"`
	models.ECommercePlatformInfo  `xorm:"extends" json:"pf"`
	models.ShopAccountInfo        `xorm:"extends" json:"sa"`
	models.StatusInfo             `xorm:"extends" json:"s"`
}

func Get_order_data_list(db_session *xorm.Session, filter string, order string, offset int64, limit int64) ([]Order_view_item_t, int64) {

	/*
SELECT o.*, ot.*, ht.*, pf.*, sa.*, s.*
FROM erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    inner join erp_ecommerce_platform pf on pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id
	 */

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("o.*, ot.*, ht.*, pf.*, sa.*, s.*")
	sql_builder.Count("count(o.id) as order_count")
	sql_builder.From(`
erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    inner join erp_ecommerce_platform pf on pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id
	`)
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

	return models.OrderShippingAddressInfo{}
}
