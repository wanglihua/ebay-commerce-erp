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

	var current_time = get_current_time()
	revel.TRACE.Println(current_time)
	lic.App_lic_date = lic.GetAppLicDate(app.AppNameForLic)

	if current_time.After(lic.App_lic_date) {
		lic.Is_lic_expire = true
	} else {
		lic.Is_lic_expire = false
	}
}

func get_current_time() time.Time {
	var time_layout = "Mon, 02 Jan 2006 15:04:05 GMT"

	var time_site_list = []string{
		"http://www.taobao.com",
		"http://www.qq.com",
		"http://www.baidu.com",
		"http://www.sina.com.cn",
		"http://www.youku.com",
		"http://www.58.com",
		"https://www.jd.com",
		"https://www.toutiao.com",
		"https://www.ebay.com",
		"https://www.paypal.com",
	}

	for _, time_site := range time_site_list {
		response, err := http.Head(time_site)
		if err == nil {
			time_internet, err := time.Parse(time_layout, response.Header.Get("Date"))
			if err == nil {
				return time_internet.Local()
			}
		}
	}

	return time.Now()
}
