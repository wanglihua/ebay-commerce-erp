package order_combine_service

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_history_service"
	"matrix/modules/erp/models/ddlc_models"
)

func combine_order_list(order_ids string, user core.LoginUser, db_session *xorm.Session) {
	var order_id_list = make([]int64, 0)
	var order_id_str_list = strings.Split(order_ids, ",")
	for _, order_id_str := range order_id_str_list {
		order_id, _ := strconv.ParseInt(order_id_str, 10, 64)
		order_id_list = append(order_id_list, order_id)
	}

	var first_order_id = order_id_list[0]

	var combine_order_id = insert_combine_order_from_first_order(first_order_id, order_ids, db_session)

	insert_combine_order_address_from_first_order(combine_order_id, first_order_id, db_session)

	update_combine_order_amount(combine_order_id, order_ids, db_session)


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
	sql_insert_order_item = fmt.Sprintf(sql_insert_order_item, "`desc`", first_order_id, "`desc`", strings.Join(order_id_str_list[1:], ",")) // 修改第一个订单之外的 order item
	_, err := db_session.Exec(sql_insert_order_item)
	core.HandleError(err)

	/*
INSERT INTO erp_order_shipping_address (order_id, name, street1, street2, city_name, province, country, country_name, phone, postal_code, address_id, address_owner, external_address_id, create_time, update_time)
SELECT -3, name, street1, street2, city_name, province, country, country_name, phone, postal_code, address_id, address_owner, external_address_id, create_time, update_time
FROM erp_order_shipping_address
WHERE order_id IN (-1, -2)
	 */
	// 删除 shipping_address
	// var sql_insert_address = `INSERT INTO erp_order_shipping_address (order_id, name, street1, street2, city_name, province, country, country_name, phone, postal_code, address_id, address_owner, external_address_id, create_time, update_time)
	// SELECT %d, name, street1, street2, city_name, province, country, country_name, phone, postal_code, address_id, address_owner, external_address_id, create_time, update_time
	// FROM erp_order_shipping_address
	// WHERE order_id IN (%s)`
	// sql_insert_address = fmt.Sprintf(sql_insert_address, first_order_id, strings.Join(order_id_str_list[1:], ",")) // 修改第一个订单之外的 shipping address
	// 好像不用修改被合并订单的地址
	// _, err = db_session.Exec(sql_insert_address)
	// core.HandleError(err)

	// 把原来的订单作废
	// UPDATE erp_order SET status_id = 4350 WHERE id IN (-1, -2)
	var sql_cancel_order = `UPDATE erp_order SET status_id = %d WHERE id IN (%s)`
	sql_cancel_order = fmt.Sprintf(sql_cancel_order, models.ORDER_STATUS_YQX, strings.Join(order_id_str_list[1:], ",")) // 作废第一个之外的 order，设置为已取消状态
	_, err = db_session.Exec(sql_cancel_order)
	core.HandleError(err)

	// 订单历史
	var first_order ddlc_models.OrderInfo
	_, err = db_session.Id(order_id_list[0]).Get(&first_order)
	core.HandleError(err)

	var combined_order_list = make([]ddlc_models.OrderInfo, 0)
	core.HandleError(db_session.In("id", order_id_list[1:]).Find(&combined_order_list))

	var combined_order_code_list = make([]string, 0)
	for _, order := range combined_order_list {
		combined_order_code_list = append(combined_order_code_list, order.Code)
	}

	order_history_service.Insert_order_history(first_order, user.UserId, fmt.Sprintf("订单合并,被和并订单编码: %s, 工号：%s", strings.Join(combined_order_code_list, ","), user.LoginName), db_session)

	for _, order := range combined_order_list {
		order_history_service.Insert_order_history(order, user.UserId, fmt.Sprintf("订单被合并而取消, 并入订单编码: %s， 工号：%s", first_order.Code, user.LoginName), db_session)
	}
}
