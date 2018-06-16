package purchase_bill_service

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"time"
	"runtime/debug"
	"matrix/core"
)

func Calculate_stock_cn_zt_qty() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("更新中国在途仓库存数量出现错误： %s\n", err)
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

	_, err = db_session.Exec("CALL calculate_stock_cn_zt_qty()")
	core.HandleError(err)
}
