package app

import (
	"strings"
	"time"

	"github.com/revel/revel"
	"matrix/modules/jobs/app/jobs"
	"matrix/modules/jobs/tasks"
)

func init() {
	revel.OnAppStart(func() {
		defer func() {
			if err := recover(); err != nil {
				revel.ERROR.Println(err)
			}
		}()

		/*
		   jobs.Schedule("0 0 0 * * ?",  TestTask{})
		   jobs.Schedule("@midnight",    TestTask{})
		   jobs.Schedule("@every 24h",   TestTask{})
		   jobs.Every(24 * time.Hour,    TestTask{})
		*/

		// 其实只是下载任务
		var schedule_task_enabled = revel.Config.BoolDefault("schedule.task.enabled", false) // 临时禁用或启用定时任务
		revel.TRACE.Println("schedule_task_enabled", schedule_task_enabled)

		if schedule_task_enabled == true {
			cron_ebay_download_order, found := revel.Config.String("cron.ebay_download_order")
			if found {
				check_cron_support(cron_ebay_download_order)
				jobs.Schedule(cron_ebay_download_order, tasks.EBayOrderDownloadTask{})
			}

			// cron.order_logistics_match
			// 下载的已支付的订单匹配物流
			cron_order_logistics_match, found := revel.Config.String("cron.order_logistics_match")
			if found {
				check_cron_support(cron_order_logistics_match)
				jobs.Schedule(cron_order_logistics_match, tasks.OrderLogisticsMatchTask{})
			}

			//jobs.Now(tasks.EBayOrderDownloadTask{}) // just for test

			/*
			cron_ebay_download_message, found := revel.Config.String("cron.ebay_download_message")
			if found {
				check_cron_support(cron_ebay_download_message)
				jobs.Schedule(cron_ebay_download_message, tasks.EBayMessageDownloadTask{})
			}
			*/

			//jobs.Now(tasks.EBayMessageDownloadTask{}) // just for test

			cron_ebay_download_product, found := revel.Config.String("cron.ebay_download_product")
			if found {
				check_cron_support(cron_ebay_download_product)
				jobs.Schedule(cron_ebay_download_product, tasks.EBayProductDownloadTask{})
			}

			//jobs.Now(tasks.EBayProductDownloadTask{})

			jobs.Every(4*time.Hour, tasks.EBayTokenStatusCheckTask{})
			//jobs.Every(time.Minute, tasks.EBayTokenStatusCheckTask{})
		}

		jobs.Every(5*time.Second, tasks.UserActiveCheckTask{})

		jobs.Every(time.Second, tasks.SystemLogDateCheckTask{})

		jobs.Every(10*time.Minute, tasks.SiteKeepHotTask{})

		jobs.Every(15*time.Minute, tasks.SystemCheckAppLicTask{})
		//jobs.Every(1*time.Second, tasks.SystemCheckAppLicTask{})

		jobs.Every(1*time.Hour, tasks.UpdateAuInventoryCostTask{})
		//jobs.Every(1*time.Minute, tasks.UpdateAuInventoryCostTask{})

		//jobs.Every(5*time.Minute, tasks.UpdateInventoryPendingQtyTask{})
		jobs.Every(1*time.Minute, tasks.UpdateInventoryPendingQtyTask{})

		jobs.Every(1*time.Hour, tasks.CalculateStockCnZtQtyTask{})

		// 商品销售分析
		jobs.Every(2*time.Hour, tasks.ProductSaleAnalysisTask{})

		// 销售分析
		jobs.Every(2*time.Hour, tasks.SaleReportTask{})

		// 智能补货
		jobs.Every(1*time.Hour, tasks.IntelReplTask{})

		// 计算商品采购价，中国仓成本价
		jobs.Every(2*time.Hour, tasks.CalculateProductPriceTask{})

	})
}

// 检查 cron 表达式是否支持
func check_cron_support(cron string) {
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	cron_fields := strings.Fields(cron)
	if len(cron_fields) != 5 && len(cron_fields) != 6 {
		revel.ERROR.Printf("Expected 5 or 6 fields, found %d: %s\n", len(cron_fields), cron)
	}
}
