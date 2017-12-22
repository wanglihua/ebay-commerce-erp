# CALL sale_amount_report();

DROP PROCEDURE IF EXISTS sale_amount_report;

CREATE PROCEDURE sale_amount_report()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_time_day_id BIGINT;
        DECLARE v_time_month_id BIGINT;
        DECLARE v_time_begin DATETIME;
        DECLARE v_time_end DATETIME;
        DECLARE v_latest_time_day_id BIGINT; # 用于存储销售统计的最近已经统计的时间天ID
        DECLARE v_latest_time_month_id BIGINT; # 用于存储销售统计的最近已经统计的时间月ID
        DECLARE v_count BIGINT;
        DECLARE v_paid_amount DECIMAL(18, 2); # 交易金额 AUD
        DECLARE v_express_fee DECIMAL(18, 2); # 运费收入 AUD
        DECLARE v_cost_price DECIMAL(18, 2); # 商品澳仓成本 RMB
        DECLARE v_estimate_logistics_fee DECIMAL(18, 2); # 预估运费 AUD
        DECLARE v_actual_logistics_fee DECIMAL(18, 2); # 实际运费 AUD
        DECLARE v_jyfy_amount DECIMAL(18, 2); # 交易费用 AUD
        DECLARE v_ptfy_amount DECIMAL(18, 2); # 平台费用 AUD
        DECLARE v_pick_fee DECIMAL(18, 2); # 拣货费 AUD
        DECLARE v_gst_fee DECIMAL(18, 2); # GST费用 AUD
        DECLARE v_profit DECIMAL(18, 2); # 利润 AUD
        DECLARE v_product_qty INT; # 商品数量


        DECLARE v_day_time_end DATETIME;
        DECLARE v_month_time_end DATETIME;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE time_day_id_cur CURSOR FOR
            SELECT td.id AS time_day_id
            FROM erp_sale_report_time_day td
            WHERE td.id > v_latest_time_day_id AND td.report_day < v_day_time_end;

        DECLARE time_month_id_cur CURSOR FOR
            SELECT tm.id AS time_month_id
            FROM erp_sale_report_time_month tm
            WHERE tm.id > v_latest_time_month_id AND tm.report_month < v_month_time_end;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        SET v_day_time_end = date_add(now(), INTERVAL -1 DAY);
        SET v_month_time_end = date_add(now(), INTERVAL -1 MONTH);

        # 日统计

        # 获取商品销售统计的最近已经统计的时间天ID
        SELECT rd.id INTO v_latest_time_day_id
        FROM erp_sale_amount_report_day rd
        ORDER BY rd.report_day_id DESC
        LIMIT 1;

        IF v_latest_time_day_id IS NULL
        THEN
            SET v_latest_time_day_id = 0;
        END IF;

        OPEN time_day_id_cur;

        time_day_id_loop: LOOP

            SET done = FALSE;
            FETCH time_day_id_cur INTO v_time_day_id;

            IF done
            THEN
                LEAVE time_day_id_loop;
            END IF;

            # do something in loop

            SELECT td.report_day
            INTO v_time_begin
            FROM erp_sale_report_time_day td
            WHERE td.id = v_time_day_id;

            SET v_time_end = date_add(v_time_begin, INTERVAL 1 DAY);

            # 商品销售额， 运费收入, 预估运费, 实际运费, 拣货费, 交易费用, 平台费用, GST费用, 利润
            SELECT
                ifnull(sum(o.amount_paid), 0),
                ifnull(sum(o.express_fee), 0),
                ifnull(sum(ol.total_fee), 0),
                ifnull(sum(ol.actual_fee), 0),
                ifnull(sum(ol.pick_fee), 0),
                ifnull(sum(o.amount_paid * ol.jyfl), 0),
                ifnull(sum(o.amount_paid * ol.ptfl), 0),
                ifnull(sum(o.amount_paid * 0.1), 0),
                ifnull(sum(o.profit_au), 0)
            INTO v_paid_amount,
                v_express_fee,
                v_estimate_logistics_fee,
                v_actual_logistics_fee,
                v_pick_fee,
                v_jyfy_amount,
                v_ptfy_amount,
                v_gst_fee,
                v_profit
            FROM erp_order o
                INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
            WHERE o.paid_time >= v_time_begin
                  AND o.paid_time < v_time_end;

            # 商品澳仓成本, 商品数量
            SELECT
                # 把澳仓成本的 RMB 转成 AUD
                ifnull(sum(oi.cost_price * oi.quantity), 0) / o.exchange_rate,
                ifnull(sum(oi.quantity), 0)
            INTO v_cost_price,
                v_product_qty
            FROM erp_order_item oi
                INNER JOIN erp_order o ON o.id = oi.order_id
            WHERE o.paid_time >= v_time_begin
                  AND o.paid_time < v_time_end;

            SELECT count(rd.id) INTO v_count
            FROM erp_sale_amount_report_day rd
            WHERE rd.report_day_id = v_time_day_id;

            IF v_count = 0
            THEN
                INSERT INTO erp_sale_amount_report_day (report_day_id, report_day, paid_amount)
                VALUES (v_time_day_id, v_time_begin, v_paid_amount);
            ELSE
                UPDATE erp_sale_amount_report_day
                SET report_day = v_time_begin, paid_amount = v_paid_amount
                WHERE report_day_id = v_time_day_id;
            END IF;

        END LOOP;

        CLOSE time_day_id_cur;

        # 月统计

        # 获取商品销售统计的最近已经统计的时间月ID
        SELECT rm.id INTO v_latest_time_month_id
        FROM erp_sale_amount_report_month rm
        ORDER BY rm.report_month_id DESC
        LIMIT 1;

        IF v_latest_time_month_id IS NULL
        THEN
            SET v_latest_time_month_id = 0;
        END IF;

        OPEN time_month_id_cur;

        time_month_id_loop: LOOP

            SET done = FALSE;
            FETCH time_month_id_cur INTO v_time_month_id;

            IF done
            THEN
                LEAVE time_month_id_loop;
            END IF;

            # do something in loop

            SELECT td.report_month
            INTO v_time_begin
            FROM erp_sale_report_time_month td
            WHERE td.id = v_time_month_id;

            SET v_time_end = date_add(v_time_begin, INTERVAL 1 MONTH);

            SELECT sum(oi.transaction_price * oi.quantity) AS amount
            INTO v_paid_amount
            FROM erp_order_item oi
                INNER JOIN erp_order o ON o.id = oi.order_id
            WHERE o.paid_time >= v_time_begin
                  AND o.paid_time < v_time_end;

            IF v_paid_amount IS NULL
            THEN
                SET v_paid_amount = 0;
            END IF;

            SELECT count(rd.id) INTO v_count
            FROM erp_sale_amount_report_month rd
            WHERE rd.report_month_id = v_time_month_id;

            IF v_count = 0
            THEN
                INSERT INTO erp_sale_amount_report_month (report_month_id, report_month, amount)
                VALUES (v_time_month_id, v_time_begin, v_paid_amount);
            ELSE
                UPDATE erp_sale_amount_report_month
                SET report_month = v_time_begin, amount = v_paid_amount
                WHERE report_month_id = v_time_month_id;
            END IF;

        END LOOP;

        CLOSE time_month_id_cur;

    END;
