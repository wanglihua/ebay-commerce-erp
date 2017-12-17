package tasks

import (
	"github.com/go-xorm/xorm"
	"time"
	"github.com/revel/revel"
	"runtime/debug"
	"matrix/core"
)

/*
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("商品销售单日分析出现错误： %s\n", err)
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

	_, err = db_session.Exec("CALL product_sale_day_analysis();")
	core.HandleError(err)
 */

type CalculateProductPriceTask struct {
}

func (t CalculateProductPriceTask) Name() string {
	return "计算商品采购价和中国仓成本价"
}

func (t CalculateProductPriceTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("计算商品采购价和中国仓成本价出现错误： %s\n", err)
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

	_, err = db_session.Exec("CALL calculate_product_price();")
	core.HandleError(err)
}
