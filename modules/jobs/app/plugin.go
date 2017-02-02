package app

import (
	"github.com/revel/revel"
	//"time"
	//"matrix/modules/jobs/app/jobs"
	//"matrix/modules/jobs/tasks"
)

func init() {
	revel.OnAppStart(func() {

		/*
		   jobs.Schedule("0 0 0 * * ?",  TestTask{})
		   jobs.Schedule("@midnight",    TestTask{})
		   jobs.Schedule("@every 24h",   TestTask{})
		   jobs.Every(24 * time.Hour,    TestTask{})
		*/

		//cron.test = * * * * * ?
		//jobs.Schedule("cron.test", TestTask{}) //测试
		//jobs.Schedule("cron.ebay_download_order", tasks.EBayOrderDownloadTask{})
		//jobs.Schedule("cron.ebay_download_message", tasks.EBayMessageDownloadTask{})
		//jobs.Schedule("cron.paypal_download", tasks.PaypalDownloadTask{})
	})
}
