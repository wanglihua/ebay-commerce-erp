package order_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

type Order_sku_qty_loc struct {
	Sku            string `xorm:"'sku'" json:"sku"`
	Quantity       int64  `xorm:"'quantity'" json:"quantity"`
	StorageLocName string `xorm:"'storage_loc_name'" json:"storage_loc_name"`
}

func Get_order_sku_qty_loc_list(order_id int64, db_session *xorm.Session) []Order_sku_qty_loc {

	/*
SELECT
    oi.item_sku AS sku,
    oi.quantity AS quantity,
    (SELECT sl.name
     FROM erp_inventory nvt
         INNER JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id
     WHERE nvt.product_id = oi.item_product_id
           AND nvt.stock_id = (SELECT ol.stock_id FROM erp_order_logistics ol WHERE ol.order_id = oi.order_id ORDER BY ol.create_time DESC LIMIT 0, 1)
     ORDER BY sl.update_time DESC
     LIMIT 0, 1) AS storage_loc_name
FROM erp_order_item oi
WHERE oi.order_id = 58880;
	 */

	var order_sku_qty_loc_list = make([]Order_sku_qty_loc, 0)
	sql := `
SELECT
    oi.item_sku AS sku,
    oi.quantity AS quantity,
    (SELECT sl.name
     FROM erp_inventory nvt
         INNER JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id
     WHERE nvt.product_id = oi.item_product_id
           AND nvt.stock_id = (SELECT ol.stock_id FROM erp_order_logistics ol WHERE ol.order_id = oi.order_id ORDER BY ol.create_time DESC LIMIT 0, 1)
     ORDER BY sl.update_time DESC
     LIMIT 0, 1) AS storage_loc_name
FROM erp_order_item oi
WHERE oi.order_id = %d;`
	sql = fmt.Sprintf(sql, order_id)
	err := db_session.Sql(sql).Find(&order_sku_qty_loc_list)
	core.HandleError(err)

	return order_sku_qty_loc_list
}
