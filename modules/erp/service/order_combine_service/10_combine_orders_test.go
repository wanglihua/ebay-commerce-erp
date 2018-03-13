package order_combine_service

import (
	"testing"
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"time"
	"runtime/debug"
	"matrix/core"
	"fmt"
	"matrix/modules/erp/models"
)

func TestCombineOrders(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("合并订单出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	db_engine, err := xorm.NewEngine("mysql", "root:root@tcp(localhost:3306)/ozoerp?charset=utf8mb4")
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	err = Combine_orders(fmt.Sprintf("o.status_id = %d OR o.status_id = %d OR o.status_id = %d ", models.ORDER_STATUS_YZF, models.ORDER_STATUS_DPH, models.ORDER_STATUS_HOLD), core.LoginUser{}, db_session)
	core.HandleError(err)
}
