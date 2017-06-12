package tasks

import (
	"matrix/app"
	"matrix/core/lic"
	"time"
	"net/http"
	"github.com/revel/revel"
	"runtime/debug"
)

type SystemCheckAppLicTask struct {
}

const System_check_app_lic_task_name = "check_al"

func (t SystemCheckAppLicTask) Name() string {
	return System_check_app_lic_task_name
}

func (t SystemCheckAppLicTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("check lic error： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	if lic.Is_lic_expire == false {
		var current_time = get_current_time()
		revel.TRACE.Println(current_time)
		lic.App_lic_date = lic.GetAppLicDate(app.AppName)

		if current_time.After(lic.App_lic_date) {
			lic.Is_lic_expire = true
		}
	}
}

func get_current_time() time.Time {
	var time_layout = "Mon, 02 Jan 2006 15:04:05 GMT"

	response, err := http.Head("http://www.taobao.com")
	if err == nil {
		time_taobao, err := time.Parse(time_layout, response.Header.Get("Date"))
		if err == nil {
			return time_taobao.Local()
		}
	}

	response, err = http.Head("http://www.qq.com")
	if err == nil {
		time_qq, err := time.Parse(time_layout, response.Header.Get("Date"))
		if err == nil {
			return time_qq.Local()
		}
	}

	response, err = http.Head("http://www.baidu.com")
	if err == nil {
		time_baidu, err := time.Parse(time_layout, response.Header.Get("Date"))
		if err == nil {
			return time_baidu.Local()
		}
	}

	return time.Now()
}
