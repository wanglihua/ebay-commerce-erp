package tasks

import (
	"runtime/debug"
	"time"

	"github.com/revel/revel"
)

type SystemLogDateCheckTask struct {
}

var last_system_log_date_check_date = time.Now()

const (
	System_log_date_check_task_name = "system_log_date_check_task"
)

func (t SystemLogDateCheckTask) Name() string {
	return System_log_date_check_task_name
}

func (t SystemLogDateCheckTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("系统日志文件日期检查出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	last_check_year, last_check_month, last_check_day := last_system_log_date_check_date.Date()

	var time_now = time.Now()
	now_check_year, now_check_month, now_check_day := time_now.Date()

	if last_check_year != now_check_year || last_check_month != now_check_month || last_check_day != now_check_day {
		revel.Recreate_all_logger()
	}
}
