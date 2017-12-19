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
        DECLARE v_amount DECIMAL(18, 2);
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

            SELECT sum(oi.transaction_price * oi.quantity) AS amount
            INTO v_amount
            FROM erp_order_item oi
                INNER JOIN erp_order o ON o.id = oi.order_id
            WHERE o.paid_time >= v_time_begin
                  AND o.paid_time < v_time_end;

            IF v_amount IS NULL
            THEN
                SET v_amount = 0;
            END IF;

            SELECT count(rd.id) INTO v_count
            FROM erp_sale_amount_report_day rd
            WHERE rd.report_day_id = v_time_day_id;

            IF v_count = 0
            THEN
                INSERT INTO erp_sale_amount_report_day (report_day_id, report_day, amount)
                VALUES (v_time_day_id, v_time_begin, v_amount);
            ELSE
                UPDATE erp_sale_amount_report_day
                SET report_day = v_time_begin, amount = v_amount
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
            INTO v_amount
            FROM erp_order_item oi
                INNER JOIN erp_order o ON o.id = oi.order_id
            WHERE o.paid_time >= v_time_begin
                  AND o.paid_time < v_time_end;

            IF v_amount IS NULL
            THEN
                SET v_amount = 0;
            END IF;

            SELECT count(rd.id) INTO v_count
            FROM erp_sale_amount_report_month rd
            WHERE rd.report_month_id = v_time_month_id;

            IF v_count = 0
            THEN
                INSERT INTO erp_sale_amount_report_month (report_month_id, report_month, amount)
                VALUES (v_time_month_id, v_time_begin, v_amount);
            ELSE
                UPDATE erp_sale_amount_report_month
                SET report_month = v_time_begin, amount = v_amount
                WHERE report_month_id = v_time_month_id;
            END IF;

        END LOOP;

        CLOSE time_month_id_cur;

    END;
