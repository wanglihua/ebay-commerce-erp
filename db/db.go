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

var (
	DbDriver string = ""
	DbUrl    string = ""
)

func InitDatabase() {
	DbDriver, _ = revel.Config.String("db.driver")
	DbUrl, _ = revel.Config.String("db.url")

	var err error
	Engine, err = xorm.NewEngine(DbDriver, DbUrl)
	if err != nil {
		revel.ERROR.Println(err)
		revel.ERROR.Println(string(debug.Stack()))
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
