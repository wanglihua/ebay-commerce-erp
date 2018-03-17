package order_combine_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddpz_models"
)

func get_order_ids_list_to_combine(order_filter string, db_session *xorm.Session) []string {
	// 同一买家，同一付款paypal，同一收货地址，同一店铺 的订单可以合并

	type order_ids struct {
		OrderIds string `xorm:"'order_ids'"`
	}

	/*
SELECT group_concat(DISTINCT o.id ORDER BY o.id ASC SEPARATOR ',') AS order_ids
FROM erp_order o
    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
GROUP BY o.buyer_user_id, o.buyer_email, o.shop_id, osa.country, osa.country_name, osa.province, osa.city_name, osa.street1, osa.street2, osa.postal_code, osa.phone, osa.name
HAVING count(DISTINCT o.id) > 1;
	 */

	var sql_get_combine_order_ids = `
SELECT group_concat(DISTINCT o.id ORDER BY o.id ASC SEPARATOR ',') AS order_ids
FROM erp_order o
    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
WHERE %s AND o.type_id <> %d AND (o.status_id = %d OR o.status_id = %d OR o.status_id = %d) AND o.status_id <> %d
GROUP BY o.buyer_user_id, o.buyer_email, o.shop_id, osa.country, osa.country_name, osa.province, osa.city_name, osa.street1, osa.street2, osa.postal_code, osa.phone, osa.name
HAVING count(DISTINCT o.id) > 1;`
	sql_get_combine_order_ids = fmt.Sprintf(sql_get_combine_order_ids, order_filter, ddpz_models.ORDER_TYPE_BHDD, models.ORDER_STATUS_YZF, models.ORDER_STATUS_DPH, models.ORDER_STATUS_HOLD, models.ORDER_STATUS_YQX)

	var order_ids_list = make([]order_ids, 0)
	err := db_session.Sql(sql_get_combine_order_ids).Find(&order_ids_list)
	core.HandleError(err)

	var order_ids_str_list = make([]string, 0)

	for _, order_ids := range order_ids_list {
		order_ids_str_list = append(order_ids_str_list, order_ids.OrderIds)
	}

	return order_ids_str_list
}
