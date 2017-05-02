package order_item_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

type Order_id_sku_qty_t struct {
	OrderId int64   `xorm:"'order_id'"`
	Sku     string  `xorm:"'sku'"`
	Qty     int64   `xorm:"'qty'"`
}

func Get_order_id_products_map(order_id_list []int64, db_session *xorm.Session) map[int64][]Order_id_sku_qty_t {
	var order_id_sku_qty_map = make(map[int64][]Order_id_sku_qty_t)

	if len(order_id_list) == 0 {
		return order_id_sku_qty_map
	}

	/*
	SELECT oi.order_id AS order_id, oi.item_sku AS sku, oi.quantity AS qty FROM erp_order_item oi WHERE oi.order_id IN (1, 2, 3);
	 */

	var sql = `SELECT oi.order_id AS order_id, oi.item_sku AS sku, oi.quantity AS qty FROM erp_order_item oi WHERE oi.order_id IN (%s)`
	sql = fmt.Sprintf(sql, strings.Join(core.Int_list_to_str_list(order_id_list), ","))

	var order_id_sku_qty_all_list = make([]Order_id_sku_qty_t, 0)
	err := db_session.Sql(sql).Find(&order_id_sku_qty_all_list)
	core.HandleError(err)

	for _, order_id_sku_qty := range order_id_sku_qty_all_list {
		order_id_sku_qty_list, found := order_id_sku_qty_map[order_id_sku_qty.OrderId]
		if found == false {
			order_id_sku_qty_list = make([]Order_id_sku_qty_t, 0)
		}
		order_id_sku_qty_list = append(order_id_sku_qty_list, order_id_sku_qty)
		order_id_sku_qty_map[order_id_sku_qty.OrderId] = order_id_sku_qty_list
	}

	return order_id_sku_qty_map;
}
