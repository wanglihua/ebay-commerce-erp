package product_sale_analysis_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"matrix/core"
	"runtime/debug"
	"github.com/revel/revel"
	"fmt"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
	"bytes"
	"matrix/modules/erp/models/spbb_models"
)

func product_sale_analyze_multiple_day() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("商品销售多日分析出现错误： %s\n", err)
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

	// 按多日统计

	// 先删除多日统计的所有数据
	var sql_clear = `TRUNCATE TABLE erp_product_sale_analysis;`
	_, err = db_session.Exec(sql_clear)
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

	/*
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
	*/

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

	/*

UPDATE erp_product_sale_analysis
SET daily_mean_sales = cast(level1_qty / 90 * 0.2 + level2_qty / 60 * 0.3 AS DECIMAL(18, 6));

	 */

	var update_analysis_daily_mean_sales_sql = `UPDATE erp_product_sale_analysis
SET daily_mean_sales = cast(%s AS DECIMAL(18, 6));`
	update_analysis_daily_mean_sales_sql = fmt.Sprintf(update_analysis_daily_mean_sales_sql, get_product_dialy_mean_sales_sql_expression(analysis_setting))

	_, err = db_session.Exec(update_analysis_daily_mean_sales_sql)
	core.HandleError(err)

	// 更新商品的日均销量和警戒库存
	/*
UPDATE erp_product p
    INNER JOIN erp_product_sale_analysis sa ON sa.product_id = p.id AND sa.type = 300 # 不分仓统计
SET p.daily_mean_sales = sa.daily_mean_sales, p.stock_alert_qty = p.stock_alert_days * sa.daily_mean_sales;
	 */
	var update_product_dialy_mean_sales_sql = `UPDATE erp_product p
    INNER JOIN erp_product_sale_analysis sa ON sa.product_id = p.id AND sa.type = 300 # 不分仓统计
SET p.daily_mean_sales = sa.daily_mean_sales, p.stock_alert_qty = p.stock_alert_days * sa.daily_mean_sales;`

	_, err = db_session.Exec(update_product_dialy_mean_sales_sql)
	core.HandleError(err)
}

func get_product_dialy_mean_sales_sql_expression(product_sale_analysis_setting spbb_models.ProductSaleAnalysisSettingInfo) string {
	var buffer bytes.Buffer

	if product_sale_analysis_setting.Level1Days != 0 {
		buffer.WriteString(fmt.Sprintf("level1_qty / %d * %f", product_sale_analysis_setting.Level1Days, product_sale_analysis_setting.Level1Weight))
	} else {
		return "0"
	}

	if product_sale_analysis_setting.Level2Days != 0 {
		buffer.WriteString(fmt.Sprintf("+ level2_qty / %d * %f", product_sale_analysis_setting.Level2Days, product_sale_analysis_setting.Level2Weight))
	} else {
		return buffer.String()
	}

	if product_sale_analysis_setting.Level3Days != 0 {
		buffer.WriteString(fmt.Sprintf("+ level3_qty / %d * %f", product_sale_analysis_setting.Level3Days, product_sale_analysis_setting.Level3Weight))
	} else {
		return buffer.String()
	}

	if product_sale_analysis_setting.Level4Days != 0 {
		buffer.WriteString(fmt.Sprintf("+ level4_qty / %d * %f", product_sale_analysis_setting.Level4Days, product_sale_analysis_setting.Level4Weight))
	} else {
		return buffer.String()
	}

	if product_sale_analysis_setting.Level5Days != 0 {
		buffer.WriteString(fmt.Sprintf("+ level5_qty / %d * %f", product_sale_analysis_setting.Level5Days, product_sale_analysis_setting.Level5Weight))
	} else {
		return buffer.String()
	}

	return buffer.String()
}
