package tasks

import (
	"runtime/debug"

	"github.com/revel/revel"
)

type OrderLogisticsMatchTask struct {
}

func (t OrderLogisticsMatchTask) Name() string {
	return "匹配物流（非导入订单）"
}

func (t OrderLogisticsMatchTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()
}
