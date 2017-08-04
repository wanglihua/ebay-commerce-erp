package product_sale_analysis_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"matrix/core"
	"runtime/debug"
	"github.com/revel/revel"
	"fmt"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
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
		var time_now = time.Now()
		var today_date = time.Date(time_now.Year(), time_now.Month(), time_now.Day(), 0, 0, 0, 0, time.Local)

		type product_sale_t struct {
			OrderId   int64 `xorm:"'order_id'"` //不分仓统计的时候，这个基本没用的
			ProductId int64 `xorm:"'product_id'"`
			Qty       int64 `xorm:"'qty'"`
		}

		for i := 0; i < DaysToAnalyze; i++ {
			// 从前一天开始
			var day_time_begin = today_date.AddDate(0, 0, -(i + 1))
			var day_time_end = today_date.AddDate(0, 0, -i)

			// 不分仓统计
			{

				/*
	SELECT oi.item_product_id AS product_id, sum(oi.quantity) AS qty
	FROM erp_order_item oi
		INNER JOIN erp_order o ON oi.order_id = o.id
	WHERE '2017-06-01' <= o.paid_time AND o.paid_time < '2017-06-02'
	GROUP BY oi.item_product_id
	ORDER BY oi.item_product_id ASC;
			 */

				var sql = `SELECT oi.item_product_id AS product_id, sum(oi.quantity) AS qty
FROM erp_order_item oi
    INNER JOIN erp_order o ON oi.order_id = o.id
WHERE '%s' <= o.paid_time AND o.paid_time < '%s'
GROUP BY oi.item_product_id
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
							var product models.ProductInfo
							_, err := db_session.Id(product_sale.ProductId).Get(&product)
							core.HandleError(err)

							product_sale_day_analysis.StockId = 0
							product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_SKU
							product_sale_day_analysis.ProductId = product_sale.ProductId
							product_sale_day_analysis.Sku = product.Sku
							product_sale_day_analysis.Qty = product_sale.Qty
							product_sale_day_analysis.Day = core.Date(day_time_begin)

							_, err = db_session.Insert(&product_sale_day_analysis)
							core.HandleError(err)
						}
					}
				}
			}

			// 分仓统计
			{
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
								var product models.ProductInfo
								_, err := db_session.Id(product_sale.ProductId).Get(&product)
								core.HandleError(err)

								product_sale_day_analysis.StockId = stock_id
								product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL
								product_sale_day_analysis.ProductId = product_sale.ProductId
								product_sale_day_analysis.Sku = product.Sku
								product_sale_day_analysis.Qty = product_sale.Qty
								product_sale_day_analysis.Day = core.Date(day_time_begin)

								_, err = db_session.Insert(&product_sale_day_analysis)
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
								var product models.ProductInfo
								_, err := db_session.Id(product_sale.ProductId).Get(&product)
								core.HandleError(err)

								product_sale_day_analysis.StockId = stock_id
								product_sale_day_analysis.Type = models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC
								product_sale_day_analysis.ProductId = product_sale.ProductId
								product_sale_day_analysis.Sku = product.Sku
								product_sale_day_analysis.Qty = product_sale.Qty
								product_sale_day_analysis.Day = core.Date(day_time_begin)

								_, err = db_session.Insert(&product_sale_day_analysis)
								core.HandleError(err)
							}
						}
					}
				}
			}
		}
	}

	// 按多日统计
	{
		// 按多日统计

		// 先删除多日统计的所有数据
		var sql_clear = `TRUNCATE TABLE erp_product_sale_analysis;`
		_, err := db_session.Exec(sql_clear)
		core.HandleError(err)

		// 先插入空数据(product_id, sku, stock_id, type)，之后的操作就可以统一进行 UPDATE
		var sql_insert = `INSERT INTO erp_product_sale_analysis (stock_id, product_id, sku, level1_qty, level2_qty, level3_qty, level4_qty, level5_qty, type, create_time, update_time)
    SELECT
        da.stock_id   AS stock_id,
        da.product_id AS product_id,
        da.sku        AS sku,
        0             AS level1_qty,
        0             AS level2_qty,
        0             AS level3_qty,
        0             AS level4_qty,
        0             AS level5_qty,
        da.type       AS type,
        now()         AS create_time,
        now()         AS update_time
    FROM erp_product_sale_day_analysis da
    GROUP BY da.stock_id, da.product_id, da.sku, da.type;`
		_, err = db_session.Exec(sql_insert)
		core.HandleError(err)

		var analysis_setting = product_sale_analysis_setting_service.Get_product_sale_analysis_setting(db_session)

		var days1 = analysis_setting.Level1Days
		var days2 = analysis_setting.Level2Days
		var days3 = analysis_setting.Level3Days
		var days4 = analysis_setting.Level4Days
		var days5 = analysis_setting.Level5Days

		var today_date_str = time.Now().Format("2006-01-02")

		type product_sale_analysis_t struct {
			ProductId int64  `xorm:"'product_id'"`
			Sku       string `xorm:"'sku'"`
			Qty       int64  `xorm:"'qty'"`
		}

		type product_sale_stock_analysis_t struct {
			StockId   int64  `xorm:"'stock_id'"`
			ProductId int64  `xorm:"'product_id'"`
			Sku       string `xorm:"'sku'"`
			Qty       int64  `xorm:"'qty'"`
		}

		/*
UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('2017-08-04', INTERVAL -90 DAY) <= da.day AND da.day < '2017-08-04'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level1_qty = da_sum.qty;
		 */

		if days1 != 0 {
			var sql = `UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('%s', INTERVAL -%d DAY) <= da.day AND da.day < '%s'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level1_qty = da_sum.qty;`
			sql = fmt.Sprintf(sql, today_date_str, days1, today_date_str)
			_, err = db_session.Exec(sql)
			core.HandleError(err)
		}

		if days2 != 0 {
			var sql = `UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('%s', INTERVAL -%d DAY) <= da.day AND da.day < '%s'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level2_qty = da_sum.qty;`
			sql = fmt.Sprintf(sql, today_date_str, days2, today_date_str)
			_, err = db_session.Exec(sql)
			core.HandleError(err)
		}

		if days3 != 0 {
			var sql = `UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('%s', INTERVAL -%d DAY) <= da.day AND da.day < '%s'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level3_qty = da_sum.qty;`
			sql = fmt.Sprintf(sql, today_date_str, days3, today_date_str)
			_, err = db_session.Exec(sql)
			core.HandleError(err)
		}

		if days4 != 0 {
			var sql = `UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('%s', INTERVAL -%d DAY) <= da.day AND da.day < '%s'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level4_qty = da_sum.qty;`
			sql = fmt.Sprintf(sql, today_date_str, days4, today_date_str)
			_, err = db_session.Exec(sql)
			core.HandleError(err)
		}

		if days5 != 0 {
			var sql = `UPDATE erp_product_sale_analysis sa
    INNER JOIN (
                   SELECT
                       da.stock_id   AS stock_id,
                       da.product_id AS product_id,
                       sum(da.qty)   AS qty,
                       da.type       AS type
                   FROM erp_product_sale_day_analysis da
                   WHERE date_add('%s', INTERVAL -%d DAY) <= da.day AND da.day < '%s'
                   GROUP BY da.stock_id, da.product_id, da.type
               ) da_sum
        ON sa.stock_id = da_sum.stock_id AND sa.product_id = da_sum.product_id AND sa.type = da_sum.type
SET sa.level5_qty = da_sum.qty;`
			sql = fmt.Sprintf(sql, today_date_str, days5, today_date_str)
			_, err = db_session.Exec(sql)
			core.HandleError(err)
		}
	}
}
