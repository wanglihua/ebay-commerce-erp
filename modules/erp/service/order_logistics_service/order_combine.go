package order_logistics_service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

// 合并订单
// order_filter: o.batch_code = '20170428120352'
func Combine_orders(order_filter string, db_session *xorm.Session) error {

	if len(order_filter) == 0 {
		return errors.New("未提供订单过滤条件！")
	}

	// 同一买家，同一付款paypal，同一收货地址，同一店铺 的订单可以合并

	type order_ids struct {
		OrderIds string `xorm:"'order_ids'"`
	}

	/*
		SELECT group_concat(o.id ORDER BY o.id ASC SEPARATOR ',') AS order_ids
		FROM erp_order o
		    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
		WHERE o.batch_code = '20170428120352'
		GROUP BY o.buyer_user_id, o.buyer_email, o.shop_id, osa.country, osa.country_name, osa.province, osa.city_name, osa.street1, osa.street2, osa.postal_code, osa.phone, osa.name
		HAVING count(o.id) > 1
	 */

	var sql_get_combin_order_ids = `
SELECT group_concat(o.id ORDER BY o.id ASC SEPARATOR ',') AS order_ids
FROM erp_order o
    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
WHERE %s
GROUP BY o.buyer_user_id, o.buyer_email, o.shop_id, osa.country, osa.country_name, osa.province, osa.city_name, osa.street1, osa.street2, osa.postal_code, osa.phone, osa.name
HAVING count(o.id) > 1`
	sql_get_combin_order_ids = fmt.Sprintf(sql_get_combin_order_ids, order_filter)

	var order_ids_list = make([]order_ids, 0)
	err := db_session.Sql(sql_get_combin_order_ids).Find(&order_ids_list)
	core.HandleError(err)

	for _, order_ids := range order_ids_list {
		var order_id_str_list = strings.Split(order_ids.OrderIds, ",")
		var order_id_list = make([]int64, 0)
		for _, order_id_str := range order_id_str_list {
			order_id, _ := strconv.ParseInt(order_id_str, 10, 64)
			order_id_list = append(order_id_list, order_id)
		}

		var first_order_id = order_id_list[0]

		// 修改 合并后的 order 各种款项
		/*
UPDATE erp_order order_first
    INNER JOIN (
                   SELECT sum(order_all.adjustment_amount) AS adjustment_amount_sum, sum(order_all.amount_paid) AS amount_paid_sum, sum(order_all.amount_saved) AS amount_saved_sum, sum(order_all.total_amount) AS total_amount_sum, sum(order_all.subtotal_amount) AS subtotal_amount_sum
                   FROM erp_order order_all
                   WHERE order_all.id IN (0, -1, -2, -3)
               ) s
        ON 1 = 1
SET order_first.adjustment_amount     = s.adjustment_amount_sum,
    order_first.adjustment_amount_rmb = s.adjustment_amount_sum * order_first.exchange_rate,
    order_first.amount_paid           = s.amount_paid_sum,
    order_first.amount_paid_rmb       = s.amount_paid_sum * order_first.exchange_rate,
    order_first.amount_saved          = s.amount_saved_sum,
    order_first.amount_saved_rmb      = s.amount_saved_sum * order_first.exchange_rate,
    order_first.total_amount          = s.total_amount_sum,
    order_first.total_amount_rmb      = s.total_amount_sum * order_first.exchange_rate,
    order_first.type_id               = 0
WHERE id = 0
		 */
		// 假设各订单的币种是一样的
		var sql_sum_amount = `
UPDATE erp_order order_first
    INNER JOIN (
                   SELECT sum(order_all.adjustment_amount) AS adjustment_amount_sum, sum(order_all.amount_paid) AS amount_paid_sum, sum(order_all.amount_saved) AS amount_saved_sum, sum(order_all.total_amount) AS total_amount_sum, sum(order_all.subtotal_amount) AS subtotal_amount_sum
                   FROM erp_order order_all
                   WHERE order_all.id IN (%s)
               ) s
        ON 1 = 1
SET order_first.adjustment_amount     = s.adjustment_amount_sum,
    order_first.adjustment_amount_rmb = s.adjustment_amount_sum * order_first.exchange_rate,
    order_first.amount_paid           = s.amount_paid_sum,
    order_first.amount_paid_rmb       = s.amount_paid_sum * order_first.exchange_rate,
    order_first.amount_saved          = s.amount_saved_sum,
    order_first.amount_saved_rmb      = s.amount_saved_sum * order_first.exchange_rate,
    order_first.total_amount          = s.total_amount_sum,
    order_first.total_amount_rmb      = s.total_amount_sum * order_first.exchange_rate,
    order_first.type_id               = %d
WHERE id = %d`
		sql_sum_amount = fmt.Sprintf(sql_sum_amount, order_ids.OrderIds, models.ORDER_TYPE_HBDD, first_order_id)
		_, err := db_session.Exec(sql_sum_amount)
		core.HandleError(err)

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
		_, err = db_session.Exec(sql_insert_order_item)
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
	}

	return nil
}
