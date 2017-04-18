package tasks

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_logistics_service"
)

type OrderLogisticsMatchTask struct {
}

func (t OrderLogisticsMatchTask) Name() string {
	return "匹配物流（非导入订单）"
}

func (t OrderLogisticsMatchTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	if task_db_driver == "" {
		task_db_driver, _ = revel.Config.String("db.driver")
	}

	if task_db_url == "" {
		task_db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(task_db_driver, task_db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	// 获得不是导入的且已支付的订单列表
	/*
	SELECT * FROM erp_order o WHERE o.type_id <> 0 AND o.status_id = 0
	 */

	var paid_download_order_list = make([]models.OrderInfo, 0)
	err = db_session.Sql(fmt.Sprintf("SELECT * FROM erp_order o WHERE o.type_id <> %d AND o.status_id = %d", models.ORDER_TYPE_DRDD, models.ORDER_STATUS_YZF)).Find(&paid_download_order_list)
	core.HandleError(err)

	var order_logistics order_logistics_service.OrderLogistics
	err = order_logistics.CalculateList(paid_download_order_list, db_session)
	core.HandleError(err)
}
