package product_sale_analysis_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"matrix/core"
	"runtime/debug"
	"github.com/revel/revel"
)

// 这是的所有统计都以已完成的订单为准，别的订单不统计

func product_sale_trend_analysis() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("商品销售趋势分析出现错误： %s\n", err)
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

	_, err = db_session.Exec("CALL product_sale_trend_analysis();")
	core.HandleError(err)

	revel.TRACE.Println("商品销售趋势分析完成！");
}
