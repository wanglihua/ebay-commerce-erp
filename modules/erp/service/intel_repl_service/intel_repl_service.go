package intel_repl_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Calculate_intel_repl(db_session *xorm.Session) {

	/*

INSERT INTO erp_intel_repl (product_id, stock_id, qty, profit, status_id, create_time, update_time)
    SELECT
        psa.product_id                                                                                                                                                                                                    AS product_id,
        psa.stock_id                                                                                                                                                                                                      AS stock_id,
        psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - (SELECT sum(pbe.qty)
                                                                                                     FROM erp_purchase_bill_entry pbe
                                                                                                     WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id) AS qty,
        (SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
         FROM erp_order_item oi
             INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
         WHERE oi.item_product_id = psa.product_id)                                                                                                                                                                       AS profit,
        # 新生成
        100                                                                                                                                                                                                               AS status_id,
        now()                                                                                                                                                                                                             AS create_time,
        now()                                                                                                                                                                                                             AS update_time
    FROM erp_product_sale_analysis psa
        INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
        INNER JOIN erp_product p ON p.id = psa.product_id
        LEFT JOIN erp_intel_repl ir ON ir.stock_id = psa.stock_id AND ir.product_id = psa.product_id AND ir.status_id = 100 # 100 新生成
    WHERE psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty > 0 AND ir.stock_id IS NULL;

UPDATE erp_intel_repl ir
    INNER JOIN (
                   SELECT
                       psa.product_id                                                                                                                                                                                                    AS product_id,
                       psa.stock_id                                                                                                                                                                                                      AS stock_id,
                       psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - (SELECT sum(pbe.qty)
                                                                                                                    FROM erp_purchase_bill_entry pbe
                                                                                                                    WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id) AS qty,
                       (SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
                        FROM erp_order_item oi
                            INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
                        WHERE oi.item_product_id = psa.product_id)                                                                                                                                                                       AS profit
                   FROM erp_product_sale_analysis psa
                       INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
                       INNER JOIN erp_product p ON p.id = psa.product_id
                   WHERE psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty > 0
               ) sa
        ON ir.status_id = 100 AND ir.stock_id = sa.stock_id AND ir.product_id = sa.product_id
SET ir.qty = sa.qty, ir.profit = sa.profit, ir.update_time = now();

	 */

	var sql_intel_repl_insert = `INSERT INTO erp_intel_repl (product_id, stock_id, qty, profit, status_id, create_time, update_time)
    SELECT
        psa.product_id                                                                                                                                                                                                    AS product_id,
        psa.stock_id                                                                                                                                                                                                      AS stock_id,
        psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - (SELECT sum(pbe.qty)
                                                                                                     FROM erp_purchase_bill_entry pbe
                                                                                                     WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id) AS qty,
        (SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
         FROM erp_order_item oi
             INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
         WHERE oi.item_product_id = psa.product_id)                                                                                                                                                                       AS profit,
        # 新生成
        100                                                                                                                                                                                                               AS status_id,
        now()                                                                                                                                                                                                             AS create_time,
        now()                                                                                                                                                                                                             AS update_time
    FROM erp_product_sale_analysis psa
        INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
        INNER JOIN erp_product p ON p.id = psa.product_id
        LEFT JOIN erp_intel_repl ir ON ir.stock_id = psa.stock_id AND ir.product_id = psa.product_id AND ir.status_id = 100 # 100 新生成
    WHERE psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty > 0 AND ir.stock_id IS NULL;`

	_, err := db_session.Exec(sql_intel_repl_insert)
	core.HandleError(err)

	var sql_intel_repl_update = `UPDATE erp_intel_repl ir
    INNER JOIN (
                   SELECT
                       psa.product_id                                                                                                                                                                                                    AS product_id,
                       psa.stock_id                                                                                                                                                                                                      AS stock_id,
                       psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - (SELECT sum(pbe.qty)
                                                                                                                    FROM erp_purchase_bill_entry pbe
                                                                                                                    WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id) AS qty,
                       (SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
                        FROM erp_order_item oi
                            INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
                        WHERE oi.item_product_id = psa.product_id)                                                                                                                                                                       AS profit
                   FROM erp_product_sale_analysis psa
                       INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
                       INNER JOIN erp_product p ON p.id = psa.product_id
                   WHERE psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty > 0
               ) sa
        ON ir.status_id = 100 AND ir.stock_id = sa.stock_id AND ir.product_id = sa.product_id
SET ir.qty = sa.qty, ir.profit = sa.profit, ir.update_time = now();`

	_, err = db_session.Exec(sql_intel_repl_update)
	core.HandleError(err)
}
