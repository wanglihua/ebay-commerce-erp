package tasks

import (
	"runtime/debug"
	"time"

	"github.com/revel/revel"
	"matrix/modules/erp/service/ebay_token_status_check_service"
)

type EBayTokenStatusCheckTask struct {
}

func (t EBayTokenStatusCheckTask) Name() string {
	return "店铺Token状态检查"
}

func (t EBayTokenStatusCheckTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("检查Token状态 出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	job_time_begin := time.Now()

	ebay_token_status_check_service.Check_ebay_token_status()

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("eBay GetTokenStatus time duration: ", job_duration.Seconds(), "seconds")
}
