package tasks

import (
	"github.com/revel/revel"
)

type TestTask struct {
}

func (t TestTask) Run() {
	revel.TRACE.Println("这是一个定时任务的测试！")
}
