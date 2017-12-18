DROP PROCEDURE IF EXISTS product_sale_amount_report;

CREATE PROCEDURE product_sale_amount_report()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_time_day_id BIGINT;
        DECLARE v_time_month_id BIGINT;
        DECLARE v_latest_time_day_id BIGINT; # 用于存储销售统计的最近已经统计的时间天ID
        DECLARE v_latest_time_month_id BIGINT; # 用于存储销售统计的最近已经统计的时间月ID
        DECLARE v_product_id BIGINT;
        DECLARE v_count BIGINT;

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

                # do something

            END LOOP;

            CLOSE time_day_id_cur;

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

                # do something

            END LOOP;

            CLOSE time_month_id_cur;

        END LOOP;

        -- 关闭游标
        CLOSE product_id_cur;
    END;
