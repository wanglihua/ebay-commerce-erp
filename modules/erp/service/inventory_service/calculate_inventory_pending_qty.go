package inventory_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"github.com/revel/revel"
	"runtime/debug"
	"matrix/core"
	"fmt"
	"matrix/modules/erp/models"
)

// 计算更新库存的 待发货 数量
func Calculate_inventory_pending_qty() {

	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("更新库存待发货数量出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url
	db_driver, _ := revel.Config.String("db.driver")
	db_url, _ := revel.Config.String("db.url")
	//只有确实要操作数据库的时候，才去初始化数据库引擎

	db_engine, err := xorm.NewEngine(db_driver, db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	/*
# 定时进行待发货的计算
# 匹配完成后马上触发代发货的计算
# 订单搬到待配货的时候是不是也要触发

# 订单帮到已发货的时候，是不是也要触发计算

# 待配货 配货中； Pick up， 信封， 快递， 国际件
UPDATE erp_inventory nvt INNER JOIN
    (SELECT nvt_inner.stock_id AS stock_id, nvt_inner.product_id AS product_id, ifnull(pending_order_item.qty, 0) AS qty
     FROM erp_inventory nvt_inner
         LEFT JOIN (SELECT o.stock_id AS stock_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
                    FROM erp_order o
                        INNER JOIN erp_order_item oi ON (o.status_id = 4000 OR o.status_id = 4100 OR (o.status_id = 3950 AND (o.hold_type_id = 10040 OR o.hold_type_id = 10140 OR o.hold_type_id = 10150 OR o.hold_type_id = 10160))) AND oi.order_id = o.id
                    GROUP BY o.stock_id, oi.item_product_id)
                   pending_order_item
             ON nvt_inner.stock_id = pending_order_item.stock_id AND nvt_inner.product_id = pending_order_item.product_id)
    pending_nvt
        ON nvt.stock_id = pending_nvt.stock_id AND nvt.product_id = pending_nvt.product_id
SET nvt.pending_qty = pending_nvt.qty;



# 匹配时要考虑 原数据表中 待发货的 数量 和 本次匹配造成的 代发货的数量
	 */

	/*
   var sql_update_inventory_pending_qty = `
UPDATE erp_inventory nvt INNER JOIN
   (SELECT nvt_inner.stock_id AS stock_id, nvt_inner.product_id AS product_id, ifnull(pending_order_item.qty, 0) AS qty
	FROM erp_inventory nvt_inner
		LEFT JOIN (SELECT o.stock_id AS stock_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
				   FROM erp_order o
					   INNER JOIN erp_order_item oi ON (o.status_id = %d OR o.status_id = %d OR o.status_id = %d OR (o.status_id = %d AND (o.hold_type_id = %d OR o.hold_type_id = %d OR o.hold_type_id = %d OR o.hold_type_id = %d))) AND oi.order_id = o.id
				   GROUP BY o.stock_id, oi.item_product_id)
				  pending_order_item
			ON nvt_inner.stock_id = pending_order_item.stock_id AND nvt_inner.product_id = pending_order_item.product_id)
   pending_nvt
	   ON nvt.stock_id = pending_nvt.stock_id AND nvt.product_id = pending_nvt.product_id
SET nvt.pending_qty = pending_nvt.qty;`
   sql_update_inventory_pending_qty = fmt.Sprintf(sql_update_inventory_pending_qty, models.ORDER_STATUS_DPH, models.ORDER_STATUS_PHZ, models.ORDER_STATUS_PROCESSED_HOLD, models.ORDER_STATUS_HOLD, ddpz_models.ORDER_HOLD_TYPE_PICK_UP, ddpz_models.ORDER_HOLD_TYPE_ENVELOPE, ddpz_models.ORDER_HOLD_TYPE_EXPRESS, ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL)
	*/

	var sql_update_inventory_pending_qty = `
UPDATE erp_inventory nvt INNER JOIN
    (SELECT nvt_inner.stock_id AS stock_id, nvt_inner.product_id AS product_id, ifnull(pending_order_item.qty, 0) AS qty
     FROM erp_inventory nvt_inner
         LEFT JOIN (SELECT o.stock_id AS stock_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
                    FROM erp_order o
                        INNER JOIN erp_order_item oi ON (o.status_id = %d OR o.status_id = %d OR o.status_id = %d OR o.status_id = %d) AND oi.order_id = o.id
                    GROUP BY o.stock_id, oi.item_product_id)
                   pending_order_item
             ON nvt_inner.stock_id = pending_order_item.stock_id AND nvt_inner.product_id = pending_order_item.product_id)
    pending_nvt
        ON nvt.stock_id = pending_nvt.stock_id AND nvt.product_id = pending_nvt.product_id
SET nvt.pending_qty = pending_nvt.qty;`
	sql_update_inventory_pending_qty = fmt.Sprintf(sql_update_inventory_pending_qty, models.ORDER_STATUS_DPH, models.ORDER_STATUS_PHZ, models.ORDER_STATUS_PROCESSED_HOLD, models.ORDER_STATUS_HOLD)

	_, err = db_session.Exec(sql_update_inventory_pending_qty)
	core.HandleError(err)
}
