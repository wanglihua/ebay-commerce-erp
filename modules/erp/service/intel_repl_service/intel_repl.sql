# 生成智能补货条目

# 数据库修改

/*

ALTER TABLE erp_intel_repl
    DROP stock_id;
ALTER TABLE erp_intel_repl
    ADD mel_qty BIGINT NOT NULL DEFAULT 0;
ALTER TABLE erp_intel_repl
    ADD bne_qty BIGINT NOT NULL DEFAULT 0;

*/

DROP PROCEDURE IF EXISTS intel_repl;

CREATE PROCEDURE intel_repl()
    BEGIN
        DECLARE v_product_id BIGINT;
        DECLARE v_cn_stock_id BIGINT;
        DECLARE v_mel_stock_id BIGINT;
        DECLARE v_bne_stock_id BIGINT;
        DECLARE v_mel_zt_stock_id BIGINT;
        DECLARE v_bne_zt_stock_id BIGINT;
        DECLARE v_mel_daily_mean_sales DECIMAL(18, 6); # MEL 日平均销量
        DECLARE v_bne_daily_mean_sales DECIMAL(18, 6); # BNE 日平均销量
        DECLARE v_purchase_days BIGINT;
        DECLARE v_stock_alert_days BIGINT;
        DECLARE v_mel_nvt_qty BIGINT;
        DECLARE v_mel_nvt_pending_qty BIGINT;
        DECLARE v_bne_nvt_qty BIGINT;
        DECLARE v_bne_nvt_pending_qty BIGINT;
        DECLARE v_mel_zt_nvt_qty BIGINT;
        DECLARE v_bne_zt_nvt_qty BIGINT;
        DECLARE v_mel_pc_qty BIGINT; # MEL 采购中数量
        DECLARE v_bne_pc_qty BIGINT; # BNE 采购中数量
        DECLARE v_pc_qty BIGINT; # 采购中的总数量，不是从上面的两个相加而来。直接从采购单的总数量来。

        DECLARE v_count INT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        # 查找所有应该补货的 product，  psa.type = 200 统计理论值
        DECLARE cur CURSOR FOR
            SELECT
                DISTINCT psa.product_id
            FROM erp_product_sale_analysis psa
            WHERE psa.type = 200; # 统计理论值

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        # 确定 CN 仓库的 stock_id
        SELECT s.id
        INTO v_cn_stock_id
        FROM erp_stock s
        WHERE s.stock_type_id = 100 # 100 中国仓库
        LIMIT 1;

        # 确定 MEL 仓库的 stock_id
        SELECT s.id
        INTO v_mel_stock_id
        FROM erp_stock s
        WHERE s.short_name = 'MEL' AND s.name LIKE '%墨尔本%' AND s.stock_type_id = 300; # 300 澳洲仓库

        # 确定 BNE 仓库的 stock_id
        SELECT s.id
        INTO v_bne_stock_id
        FROM erp_stock s
        WHERE s.short_name = 'BNE' AND s.name LIKE '%布里斯班%' AND s.stock_type_id = 300; # 300 澳洲仓库

        # 确定 MEL 在途库 的 stock_id
        SELECT s.id
        INTO v_mel_zt_stock_id
        FROM erp_stock s
        WHERE s.target_stock_id = v_mel_stock_id AND s.stock_type_id = 200; # 200 在途仓库

        # 确定 BNE 在途库 的 stock_id
        SELECT s.id
        INTO v_bne_zt_stock_id
        FROM erp_stock s
        WHERE s.target_stock_id = v_bne_stock_id AND s.stock_type_id = 200; # 200 在途仓库


        TRUNCATE TABLE erp_intel_repl;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP
            -- 提取游标里的数据
            FETCH cur
            INTO v_product_id;

            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;

            -- 这里做你想做的循环的事件

            # 计算 澳仓部分 psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty
            # 计算 在途仓部分
            # 计算 正在进行的采购单部分， MEL 和 BNE 要分开
            # 计算 中国仓部分， MEL 和 BNE 要按比例分开

            # 获取商品的相关配置
            SET v_purchase_days = 0, v_stock_alert_days = 0;
            SELECT
                p.purchase_days,
                p.stock_alert_days
            INTO v_purchase_days, v_stock_alert_days
            FROM erp_product p
            WHERE p.id = v_product_id;

            # MEL

            # 获取 MEL 仓库的 日平均销量
            SET v_mel_daily_mean_sales = 0;
            SELECT psa.daily_mean_sales
            INTO v_mel_daily_mean_sales
            FROM erp_product_sale_analysis psa
            WHERE psa.product_id = v_product_id AND psa.stock_id = v_mel_stock_id AND psa.type = 200; # 统计理论值

            # 获取 MEL 仓库的 数量和待发货数量
            SET v_mel_nvt_qty = 0, v_mel_nvt_pending_qty = 0;
            SELECT
                nvt.qty,
                nvt.pending_qty
            INTO v_mel_nvt_qty, v_mel_nvt_pending_qty
            FROM erp_inventory nvt
            WHERE nvt.product_id = v_product_id AND nvt.stock_id = v_mel_stock_id;

            # 获取 MEL 在途仓库 数量
            SET v_mel_zt_nvt_qty = 0;
            SELECT nvt.qty
            INTO v_mel_zt_nvt_qty
            FROM erp_inventory nvt
            WHERE nvt.product_id = v_product_id AND nvt.stock_id = v_mel_zt_stock_id;

            # 获取 MEL 采购中的 数量
            SELECT sum(pbe.mel_qty)
            INTO v_mel_pc_qty
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = v_product_id;

            # BNE

            # 获取 BNE 仓库的 日平均销量
            SET v_bne_daily_mean_sales = 0;
            SELECT psa.daily_mean_sales
            INTO v_bne_daily_mean_sales
            FROM erp_product_sale_analysis psa
            WHERE psa.product_id = v_product_id AND psa.stock_id = v_bne_stock_id AND psa.type = 200; # 统计理论值

            # 获取 BNE 仓库的 数量和待发货数量
            SET v_bne_nvt_qty = 0, v_bne_nvt_pending_qty = 0;
            SELECT
                nvt.qty,
                nvt.pending_qty
            INTO v_bne_nvt_qty, v_bne_nvt_pending_qty
            FROM erp_inventory nvt
            WHERE nvt.product_id = v_product_id AND nvt.stock_id = v_bne_stock_id;

            # 获取 BNE 在途仓库 数量
            SET v_bne_zt_nvt_qty = 0;
            SELECT nvt.qty
            INTO v_bne_zt_nvt_qty
            FROM erp_inventory nvt
            WHERE nvt.product_id = v_product_id AND nvt.stock_id = v_bne_zt_stock_id;

            # 获取 BNE 采购中的 数量
            SELECT sum(pbe.bne_qty)
            INTO v_bne_pc_qty
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = v_product_id;

            # MEL AND BNE

            # 获取 总的 采购数量
            SELECT sum(pbe.qty)
            INTO v_pc_qty
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = v_product_id;

            # 生成或更新 智能补货单
            SET v_count = 0;
            SELECT count(*)
            INTO v_count
            FROM erp_intel_repl ir
            WHERE ir.product_id = v_product_id AND ir.status_id = 5700; # 5700 新生成的智能补货单

            IF v_count = 0
            THEN
            # 新增
            ELSE
            # 更新
            END IF;

        END LOOP;
        -- 关闭游标
        CLOSE cur;

        /*

        INSERT INTO erp_intel_repl (product_id, stock_id, qty, profit, status_id, create_time, update_time)
            SELECT
                psa.product_id                                                                                                                                                                                                               AS product_id,
                psa.stock_id                                                                                                                                                                                                                 AS stock_id,
                psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - ifnull((SELECT sum(pbe.qty)
                                                                                                                    FROM erp_purchase_bill_entry pbe
                                                                                                                    WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id), 0) AS qty,
                ifnull((SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
                        FROM erp_order_item oi
                            INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
                        WHERE oi.item_product_id = psa.product_id), 0)                                                                                                                                                                       AS profit,
                # 新生成
                5700                                                                                                                                                                                                                         AS status_id,
                now()                                                                                                                                                                                                                        AS create_time,
                now()                                                                                                                                                                                                                        AS update_time
            FROM erp_product_sale_analysis psa
                INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
                INNER JOIN erp_product p ON p.id = psa.product_id
                LEFT JOIN erp_intel_repl ir ON ir.stock_id = psa.stock_id AND ir.product_id = psa.product_id AND ir.status_id = 5700 # 5700 新生成
            WHERE ir.stock_id IS NULL AND psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - ifnull((SELECT sum(pbe.qty)
                                                                                                                                              FROM erp_purchase_bill_entry pbe
                                                                                                                                              WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id), 0) > 0;

        UPDATE erp_intel_repl ir
            INNER JOIN (
                           SELECT
                               psa.product_id                                                                                                                                                                                                               AS product_id,
                               psa.stock_id                                                                                                                                                                                                                 AS stock_id,
                               psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - ifnull((SELECT sum(pbe.qty)
                                                                                                                                   FROM erp_purchase_bill_entry pbe
                                                                                                                                   WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id), 0) AS qty,
                               ifnull((SELECT cast(avg(oi.gross_profit / (oi.transaction_price * oi.quantity)) AS DECIMAL(18, 6))
                                       FROM erp_order_item oi
                                           INNER JOIN erp_order o ON o.id = oi.order_id AND o.status_id = 4300 # 订单状态 已完成
                                       WHERE oi.item_product_id = psa.product_id), 0)                                                                                                                                                                       AS profit
                           FROM erp_product_sale_analysis psa
                               INNER JOIN erp_inventory nvt ON nvt.product_id = psa.product_id AND nvt.stock_id = psa.stock_id AND psa.type = 200
                               INNER JOIN erp_product p ON p.id = psa.product_id
                           WHERE psa.daily_mean_sales * (p.purchase_days + p.stock_alert_days) - nvt.qty - nvt.pending_qty - ifnull((SELECT sum(pbe.qty)
                                                                                                                                     FROM erp_purchase_bill_entry pbe
                                                                                                                                     WHERE pbe.status_id <> 4700 AND pbe.status_id <> 4710 AND pbe.status_id <> 4720 AND pbe.product_id = psa.product_id), 0) > 0
                       ) sa
                ON ir.status_id = 5700 AND ir.stock_id = sa.stock_id AND ir.product_id = sa.product_id
        SET ir.qty = sa.qty, ir.profit = sa.profit, ir.update_time = now();

         */

    END

