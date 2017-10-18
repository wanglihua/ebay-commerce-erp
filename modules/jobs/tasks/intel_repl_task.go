package tasks

import (
	"time"
	"github.com/go-xorm/xorm"
	"runtime/debug"
	"github.com/revel/revel"
	"matrix/core"
)

type IntelReplTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t IntelReplTask) Name() string {
	return "智能补货需求生成"
}

func (t IntelReplTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("智能补货需求生成出现错误： %s\n", err)
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

	//todo: 生成商品只能补货条目
}
