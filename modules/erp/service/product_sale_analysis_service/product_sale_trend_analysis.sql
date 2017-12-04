DROP PROCEDURE IF EXISTS product_sale_trend_analysis;

CREATE PROCEDURE product_sale_trend_analysis()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE const_time_range_type BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_latest_begin_time DATETIME;
        DECLARE v_current_date DATETIME;
        DECLARE v_time_range_id BIGINT;
        DECLARE v_time_range_type_id BIGINT;
        DECLARE v_range_begin_time DATETIME;
        DECLARE v_range_end_time DATETIME;
        DECLARE v_range_qty BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE product_id_cur CURSOR FOR
            SELECT DISTINCT da.product_id FROM erp_product_sale_day_analysis da;

        DECLARE time_range_cur CURSOR FOR
            SELECT tr.id, tr.type_id, tr.begin_time, tr.end_time FROM erp_product_sale_analysis_time_range tr
            WHERE v_latest_begin_time < tr.begin_time AND tr.end_time < v_current_date; # todo 是不是需要取等号 还要考虑一下

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        # 10 天一个间隔
        # TIME_RANGE_TYPE_TEN_DAYS int64 = 100
        SET const_time_range_type = 100;

        # SELECT current_date() INTO v_current_date FROM dual;
        SET v_current_date = current_date();

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

            -- 这里做你想做的循环的事件
            SET v_latest_begin_time = NULL;
            SELECT ta.begin_time INTO v_latest_begin_time
            FROM erp_product_sale_trend_analysis ta
            WHERE ta.product_id = v_product_id AND ta.type_id = const_time_range_type
            ORDER BY ta.begin_time DESC
            LIMIT 1;

            IF v_latest_begin_time IS NULL
            THEN
                SET v_latest_begin_time = '2006-01-01';
            END IF;

            -- 打开游标
            OPEN time_range_cur;

            -- 开始循环
            time_range_id_loop: LOOP
                -- 提取游标里的数据
                SET done = FALSE;
                FETCH time_range_cur INTO v_time_range_id, v_time_range_type_id, v_range_begin_time, v_range_end_time;

                IF done
                THEN
                    LEAVE time_range_id_loop;
                END IF;

                -- 这里做你想做的循环的事件

                /*
                // 统计实际值
                PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL int64 = 100
                // 统计理论值
                PRODUCT_SALE_ANALYSIS_TYPE_LOGIC int64 = 200
                // 只统计SKU，不分仓库统计，实际值的不分仓统计
                PRODUCT_SALE_ANALYSIS_TYPE_SKU int64 = 300
                 */

                # 时间范围前一个是 <= 后一个是 <
                # 取按日分析的不分仓的数据
                SELECT ifnull(sum(da.qty), 0) INTO v_range_qty
                FROM erp_product_sale_day_analysis da
                WHERE da.product_id = v_product_id AND v_range_begin_time <= da.day AND da.day < v_range_end_time AND type = 300;

                # 插入数量分析结果
                INSERT INTO erp_product_sale_trend_analysis (time_range_id, type_id, begin_time, end_time, product_id, qty)
                VALUES (v_time_range_id, v_time_range_type_id, v_range_begin_time, v_range_end_time, v_product_id, v_range_qty);
            END LOOP;
            -- 关闭游标
            CLOSE time_range_cur;
        END LOOP;
        -- 关闭游标
        CLOSE product_id_cur;

    END;