package order_combine_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models/ddpz_models"
	"strings"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/inventory_service"
)

func cancel_combined_orders(order_id_str_list []string, db_session *xorm.Session) {

	// 把原来的订单全部设置成 被合并订单
	/*
UPDATE erp_order SET type_id = 0
WHERE id IN (-1, -2, -3);
	 */
	var sql_update_order_type = `
UPDATE erp_order SET type_id = %d
WHERE id IN (%s);`
	sql_update_order_type = fmt.Sprintf(sql_update_order_type, ddpz_models.ORDER_TYPE_BHDD, strings.Join(order_id_str_list, ","))

	_, err := db_session.Exec(sql_update_order_type)
	core.HandleError(err)

	// 把原来的订单作废
	// UPDATE erp_order SET status_id = 4350 WHERE id IN (-1, -2)
	var sql_cancel_order = `UPDATE erp_order SET status_id = %d WHERE id IN (%s)`
	sql_cancel_order = fmt.Sprintf(sql_cancel_order, models.ORDER_STATUS_YQX, strings.Join(order_id_str_list, ","))
	_, err = db_session.Exec(sql_cancel_order)
	core.HandleError(err)

	// 更新库存待配货的数量
	inventory_service.Calculate_inventory_pending_qty()
}
