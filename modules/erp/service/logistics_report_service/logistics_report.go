package logistics_report_service

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"time"
	"matrix/core"
	"runtime/debug"
)

func Logistics_report() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("物流报表分析出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url
	db_driver, _ := revel.Config.String("db.driver")
	db_url, _ := revel.Config.String("db.url")

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

	_, err = db_session.Exec("CALL logistics_report();")
	core.HandleError(err)

	revel.TRACE.Println("物流报表分析完成！");
}
