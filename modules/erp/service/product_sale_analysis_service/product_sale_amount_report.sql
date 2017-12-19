DROP PROCEDURE IF EXISTS product_sale_amount_report;

CREATE PROCEDURE product_sale_amount_report()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_time_day_id BIGINT;
        DECLARE v_time_month_id BIGINT;
        DECLARE v_time_begin DATETIME;
        DECLARE v_time_end DATETIME;
        DECLARE v_latest_time_day_id BIGINT; # 用于存储销售统计的最近已经统计的时间天ID
        DECLARE v_latest_time_month_id BIGINT; # 用于存储销售统计的最近已经统计的时间月ID
        DECLARE v_product_id BIGINT;
        DECLARE v_count BIGINT;
        DECLARE v_amount DECIMAL(18, 2);

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE product_id_cur CURSOR FOR
            SELECT p.id AS product_id
            FROM erp_product p
            ORDER BY p.id;

        DECLARE time_day_id_cur CURSOR FOR
            SELECT td.id AS time_day_id
            FROM erp_sale_report_time_day td
            WHERE td.id > v_latest_time_day_id;

        DECLARE time_month_id_cur CURSOR FOR
            SELECT *
            FROM erp_sale_report_time_month tm
            WHERE tm.id > v_latest_time_month_id;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        -- 打开游标
        OPEN product_id_cur;

        -- 开始循环
        product_id_loop: LOOP
            -- 提取游标里的数据
            SET done = FALSE;
            FETCH product_id_cur INTO v_product_id;

            IF done
            THEN
                LEAVE product_id_loop;
            END IF;

            # 日统计

            # 获取商品销售统计的最近已经统计的时间天ID
            SELECT rd.id INTO v_latest_time_day_id
            FROM erp_sale_amount_report_day rd
            WHERE rd.product_id = v_product_id
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

                SELECT oi.transaction_price * oi.quantity AS amount
                INTO v_amount
                FROM erp_order_item oi
                    INNER JOIN erp_order o ON o.id = oi.order_id
                WHERE oi.item_product_id = v_product_id
                      AND o.shipped_time >= v_time_begin
                      AND o.shipped_time < v_time_end;

                SELECT count(rd.id) INTO v_count
                FROM erp_sale_amount_report_day rd
                WHERE rd.product_id = v_product_id
                      AND rd.report_day_id = v_time_day_id;

                IF v_count = 0
                THEN
                    INSERT INTO erp_sale_amount_report_day (report_day_id, report_day, product_id, amount)
                    VALUES (v_time_day_id, v_time_begin, v_product_id, v_amount);
                ELSE
                    UPDATE erp_sale_amount_report_day
                    SET report_day = v_time_begin, amount = v_amount
                    WHERE product_id = v_product_id
                          AND report_day_id = v_time_day_id;
                END IF;

            END LOOP;

            CLOSE time_day_id_cur;

            # 月统计

            # 获取商品销售统计的最近已经统计的时间月ID
            SELECT rm.id INTO v_latest_time_month_id
            FROM erp_sale_amount_report_month rm
            WHERE rm.product_id = v_product_id
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

                SELECT oi.transaction_price * oi.quantity AS amount
                INTO v_amount
                FROM erp_order_item oi
                    INNER JOIN erp_order o ON o.id = oi.order_id
                WHERE oi.item_product_id = v_product_id
                      AND o.shipped_time >= v_time_begin
                      AND o.shipped_time < v_time_end;

                SELECT count(rd.id) INTO v_count
                FROM erp_sale_amount_report_month rd
                WHERE rd.product_id = v_product_id
                      AND rd.report_month_id = v_time_month_id;

                IF v_count = 0
                THEN
                    INSERT INTO erp_sale_amount_report_month (report_month_id, report_month, product_id, amount)
                    VALUES (v_time_month_id, v_time_begin, v_product_id, v_amount);
                ELSE
                    UPDATE erp_sale_amount_report_month
                    SET report_month = v_time_begin, amount = v_amount
                    WHERE product_id = v_product_id
                          AND report_month_id = v_time_month_id;
                END IF;

            END LOOP;

            CLOSE time_month_id_cur;

        END LOOP;

        -- 关闭游标
        CLOSE product_id_cur;
    END;
