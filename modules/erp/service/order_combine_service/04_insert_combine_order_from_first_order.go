package order_combine_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
	"time"
	"strings"
)

var order_code_seq int64 = 1

func insert_combine_order_from_first_order(first_order_id int64, order_id_str_list []string, db_session *xorm.Session) int64 {
	// 获取被合并订单中的第一个订单
	var first_order ddlc_models.OrderInfo
	_, err := db_session.Id(first_order_id).Get(&first_order)
	core.HandleError(err)

	// 根据第一个订单复制生成合并的订单
	var combine_order ddlc_models.OrderInfo = first_order
	combine_order.Id = 0

	combine_order.Code = time.Now().Format("200601021504") + fmt.Sprintf("%02d", order_code_seq)
	order_code_seq = (order_code_seq + 1) % 100 // 加一个序列号

	combine_order.TypeId = ddpz_models.ORDER_TYPE_HBDD // 合并订单

	_, err = db_session.Insert(combine_order)
	core.HandleError(err)

	// 把原来的订单全部设置成 被合并订单
	/*
UPDATE erp_order SET type_id = 0
WHERE id IN (-1, -2, -3);
	 */
	var sql_update_order_type = `
UPDATE erp_order SET type_id = %d
WHERE id IN (%s);`
	sql_update_order_type = fmt.Sprintf(sql_update_order_type, ddpz_models.ORDER_TYPE_BHDD, strings.Join(order_id_str_list, ","))

	_, err = db_session.Exec(sql_update_order_type)
	core.HandleError(err)

	return combine_order.Id
}
