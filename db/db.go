package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	xorm_core "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"runtime/debug"
)

var Engine *xorm.Engine = nil

func InitDatabase() {
	db_driver, _ := revel.Config.String("db.driver")
	db_url, _ := revel.Config.String("db.url")

	var err error
	Engine, err = xorm.NewEngine(db_driver, db_url)
	if err != nil {
		revel.TRACE.Println(err)
		revel.TRACE.Println(string(debug.Stack()))
	}

	Engine.TZLocation = time.Local //设置时区

	//连接池
	Engine.SetMaxIdleConns(3)   //最小
	Engine.SetMaxOpenConns(100) //最大

	if revel.DevMode == true {
		Engine.ShowSQL()
		Engine.ShowExecTime()
		xorm_core.LogQuerySql = true
	}
}
