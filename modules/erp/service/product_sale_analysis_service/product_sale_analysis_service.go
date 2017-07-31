package product_sale_analysis_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"matrix/core"
	"runtime/debug"
	"github.com/revel/revel"
	"fmt"
	"matrix/modules/erp/models"
)

func ProductSaleAnalyze() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("商品销售分析出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url
	db_driver, _ := revel.Config.String("db.driver")
	db_url, _ := revel.Config.String("db.url")

	db_engine, err := xorm.NewEngine(db_driver, db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	// 按单日统计
	{
		// 按单日统计

		var days_to_analyze = 3;
		var time_now = time.Now()
		var today_date = time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, time.Local)

		type product_sale_t struct {
			OrderId   int64 `xorm:"'order_id'"`
			ProductId int64 `xorm:"'product_id'"`
			Qty       int64 `xorm:"'qty'"`
		}

		for i := 0; i < days_to_analyze; i++ {
			// 从前一天开始
			var day_time_begin = today_date.AddDate(0, 0, -(i + 1))
			var day_time_end = today_date.AddDate(0, 0, -i)

			/*
	SELECT oi.order_id AS order_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
	FROM erp_order_item oi
		INNER JOIN erp_order o ON oi.order_id = o.id
	WHERE '2017-06-01' <= o.paid_time AND o.paid_time < '2017-06-02'
	GROUP BY oi.order_id, oi.item_product_id
	ORDER BY oi.item_product_id ASC;
			 */

			var sql = `SELECT oi.order_id AS order_id, oi.item_product_id AS product_id, sum(oi.quantity) AS qty
FROM erp_order_item oi
    INNER JOIN erp_order o ON oi.order_id = o.id
WHERE '%s' <= o.paid_time AND o.paid_time < '%s'
GROUP BY oi.order_id, oi.item_product_id
ORDER BY oi.item_product_id ASC;`
			sql = fmt.Sprintf(sql, day_time_begin.Format("2006-01-02"), day_time_end.Format("2006-01-02"))

			var product_sale_list = make([]product_sale_t, 0)
			core.HandleError(db_session.Sql(sql).Find(&product_sale_list))

			for _, product_sale := range product_sale_list {
				{
					// 不分仓统计

					var product_sale_day_analysis models.ProductSaleDayAnalysisInfo
					has, err := db_session.Where(fmt.Sprintf("type = %d AND product_id = %d AND day = '%s'",
						models.PRODUCT_SALE_ANALYSIS_TYPE_SKU,
						product_sale.ProductId,
						day_time_begin.Format("2006-01-02"))).
						Get(&product_sale_day_analysis)
					core.HandleError(err)

					if has == true {
						product_sale_day_analysis.Qty = product_sale.Qty
						_, err := db_session.Id(product_sale_day_analysis.Id).Cols("qty").Update(&product_sale_day_analysis)
						core.HandleError(err)
					} else {
						product_sale_day_analysis.StockId = 0
						product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_SKU
						product_sale_day_analysis.ProductId = product_sale.ProductId
						product_sale_day_analysis.Qty = product_sale.Qty
						product_sale_day_analysis.Day = core.Date(day_time_begin)

						_, err := db_session.Insert(&product_sale_day_analysis)
						core.HandleError(err)
					}
				}

				{
					// 实际仓统计 stock_id != 0

					var order_logistics models.OrderLogisticsInfo
					has, err := db_session.Where(fmt.Sprintf("order_id = %d", product_sale.OrderId)).Get(&order_logistics)
					core.HandleError(err)

					var stock_id = order_logistics.StockId

					if has == true && stock_id != 0 { // 如果有匹配的物流信息，并且仓库信息不为空
						var product_sale_day_analysis models.ProductSaleDayAnalysisInfo
						has, err := db_session.Where(fmt.Sprintf("type = %d AND stock_id = %d AND product_id = %d AND day = '%s'",
							models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL,
							stock_id,
							product_sale.ProductId,
							day_time_begin.Format("2006-01-02"))).
							Get(&product_sale_day_analysis)
						core.HandleError(err)

						if has == true {
							product_sale_day_analysis.Qty = product_sale.Qty
							_, err := db_session.Id(product_sale_day_analysis.Id).Cols("qty").Update(&product_sale_day_analysis)
							core.HandleError(err)
						} else {
							product_sale_day_analysis.StockId = stock_id
							product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL
							product_sale_day_analysis.ProductId = product_sale.ProductId
							product_sale_day_analysis.Qty = product_sale.Qty
							product_sale_day_analysis.Day = core.Date(day_time_begin)

							_, err := db_session.Insert(&product_sale_day_analysis)
							core.HandleError(err)
						}
					}
				}

				{
					// 逻辑仓统计 stock_id != 0

					var order_logistics_logic models.OrderLogisticsLogicInfo
					has, err := db_session.Where(fmt.Sprintf("order_id = %d", product_sale.OrderId)).Get(&order_logistics_logic)
					core.HandleError(err)

					var stock_id = order_logistics_logic.StockId

					if has == true && stock_id != 0 { // 如果有匹配的物流信息，并且仓库信息不为空
						var product_sale_day_analysis models.ProductSaleDayAnalysisInfo
						has, err := db_session.Where(fmt.Sprintf("type = %d AND stock_id = %d AND product_id = %d AND day = '%s'",
							models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC,
							stock_id,
							product_sale.ProductId,
							day_time_begin.Format("2006-01-02"))).
							Get(&product_sale_day_analysis)
						core.HandleError(err)

						if has == true {
							product_sale_day_analysis.Qty = product_sale.Qty
							_, err := db_session.Id(product_sale_day_analysis.Id).Cols("qty").Update(&product_sale_day_analysis)
							core.HandleError(err)
						} else {
							product_sale_day_analysis.StockId = stock_id
							product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC
							product_sale_day_analysis.ProductId = product_sale.ProductId
							product_sale_day_analysis.Qty = product_sale.Qty
							product_sale_day_analysis.Day = core.Date(day_time_begin)

							_, err := db_session.Insert(&product_sale_day_analysis)
							core.HandleError(err)
						}
					}
				}
			}
		}
	}

	// 按多日统计
	{
		// 按多日统计

	}
}
