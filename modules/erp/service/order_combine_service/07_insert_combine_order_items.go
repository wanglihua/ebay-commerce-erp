package order_combine_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"strings"
	"matrix/core"
)

func insert_combine_order_items(combine_order_id int64, order_id_str_list []string, db_session *xorm.Session) {
	/*
INSERT INTO erp_order_item (order_id, order_line_item_id, quantity, transaction_id, transaction_price_currency, transaction_price, transaction_price_rmb, product_store_item_id, item_id, item_product_id, item_sku, item_title, item_condition_id, item_condition_display_name, created_date, `desc`, create_time, update_time)
SELECT -1, order_line_item_id, quantity, transaction_id, transaction_price_currency, transaction_price, transaction_price_rmb, product_store_item_id, item_id, item_product_id, item_sku, item_title, item_condition_id, item_condition_display_name, created_date, `desc`, create_time, update_time
FROM erp_order_item
WHERE order_id IN (-2, -3)
	 */

	// 添加新的 order_id_item
	var sql_insert_order_item = `INSERT INTO erp_order_item (order_id, order_line_item_id, quantity, transaction_id, transaction_price_currency, transaction_price, transaction_price_rmb, product_store_item_id, item_id, item_product_id, item_sku, item_title, item_condition_id, item_condition_display_name, created_date, %s, create_time, update_time)
    SELECT %d, order_line_item_id, quantity, transaction_id, transaction_price_currency, transaction_price, transaction_price_rmb, product_store_item_id, item_id, item_product_id, item_sku, item_title, item_condition_id, item_condition_display_name, created_date, %s, create_time, update_time
    FROM erp_order_item
    WHERE order_id IN (%s)`
	sql_insert_order_item = fmt.Sprintf(sql_insert_order_item, "`desc`", combine_order_id, "`desc`", strings.Join(order_id_str_list, ","))
	_, err := db_session.Exec(sql_insert_order_item)
	core.HandleError(err)
}
