package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models"
	"time"
	"github.com/revel/revel"
	"runtime/debug"
	"matrix/core"
)

// 计算更新库存的 代发货 数量
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
# 定时进行代发货的计算
# 匹配完成后马上触发代发货的计算
# 订单搬到待配货的时候是不是也要触发

# 待配货 配货中
UPDATE erp_inventory nvt INNER JOIN
    (SELECT o.stock_id AS stock_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
     FROM erp_order o
         INNER JOIN erp_order_item oi ON (o.status_id = 1 OR o.status_id = 2) AND oi.order_id = o.id
     GROUP BY o.stock_id, oi.item_product_id) padding_order_item
        ON nvt.stock_id = padding_order_item.stock_id AND nvt.product_id = padding_order_item.product_id
SET nvt.qty = padding_order_item.qty;

# 匹配时要考虑 原数据表中 代发货的 数量 和 本次匹配造成的 代发货的数量
	 */

	type order_id_t struct {
		OrderId int64 `xorm:"'order_id'"`
	}
	var sql_order_list = `SELECT o.id AS order_id FROM erp_order o WHERE o.status_id = %d OR o.status_id = %d;`
	sql_order_list = fmt.Sprintf(sql_order_list, models.ORDER_STATUS_DPH, models.ORDER_STATUS_PHZ)
	var order_id_list = make([]order_id_t, 0)
	err = db_session.Sql(sql_order_list).Find(&order_id_list)
	core.HandleError(err)


}
