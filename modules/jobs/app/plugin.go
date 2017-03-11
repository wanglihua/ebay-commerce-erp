package app

import (
	"time"

	"github.com/revel/revel"
	"matrix/modules/jobs/app/jobs"
	"matrix/modules/jobs/tasks"
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

		_, found := revel.Config.String("cron.ebay_download_order")
		if found {
			jobs.Schedule("cron.ebay_download_order", tasks.EBayOrderDownloadTask{})
		}

		//jobs.Now(tasks.EBayOrderDownloadTask{}) // just for test

		_, found = revel.Config.String("cron.ebay_download_message")
		if found {
			jobs.Schedule("cron.ebay_download_message", tasks.EBayMessageDownloadTask{})
		}

		//jobs.Now(tasks.EBayMessageDownloadTask{}) // just for test

		_, found = revel.Config.String("cron.ebay_download_product")
		if found {
			jobs.Schedule("cron.ebay_download_product", tasks.EBayProductDownloadTask{})
		}

		//jobs.Now(tasks.EBayProductDownloadTask{})

		_, found = revel.Config.String("cron.paypal_download")
		if found {
			jobs.Schedule("cron.paypal_download", tasks.PaypalDownloadTask{})
		}

		jobs.Now(tasks.PaypalDownloadTask{})

		jobs.Every(10*time.Minute, tasks.SiteKeepHotTask{})
	})
}
