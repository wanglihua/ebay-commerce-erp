package tasks

import (
	"runtime/debug"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/service/login_log_service"
)

type UserActiveCheckTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t UserActiveCheckTask) Name() string {
	return "用户在线检查"
}

func (t UserActiveCheckTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("用户在线检查出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url

	if t.db_driver == "" {
		t.db_driver, _ = revel.Config.String("db.driver")
	}

	if t.db_url == "" {
		t.db_url, _ = revel.Config.String("db.url")
	}

	//只有确实要操作数据库的时候，才去初始化数据库引擎

	db_engine, err := xorm.NewEngine(t.db_driver, t.db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	login_log_service.Check_user_active(db_session)
}
