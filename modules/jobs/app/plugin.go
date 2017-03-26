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

		cron_ebay_download_order, found := revel.Config.String("cron.ebay_download_order")
		if found {
			jobs.Schedule(cron_ebay_download_order, tasks.EBayOrderDownloadTask{})
		}

		//jobs.Now(tasks.EBayOrderDownloadTask{}) // just for test

		cron_ebay_download_message, found := revel.Config.String("cron.ebay_download_message")
		if found {
			jobs.Schedule(cron_ebay_download_message, tasks.EBayMessageDownloadTask{})
		}

		//jobs.Now(tasks.EBayMessageDownloadTask{}) // just for test

		cron_ebay_download_product, found := revel.Config.String("cron.ebay_download_product")
		if found {
			jobs.Schedule(cron_ebay_download_product, tasks.EBayProductDownloadTask{})
		}

		//jobs.Now(tasks.EBayProductDownloadTask{})

		cron_paypal_download, found := revel.Config.String("cron.paypal_download")
		if found {
			jobs.Schedule(cron_paypal_download, tasks.PaypalDownloadTask{})
		}

		// jobs.Now(tasks.PaypalDownloadTask{})

		jobs.Every(10*time.Minute, tasks.SiteKeepHotTask{})
	})
}
