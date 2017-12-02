DROP PROCEDURE IF EXISTS product_sale_trend_analysis;

CREATE PROCEDURE product_sale_trend_analysis()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE time_range_type BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_latest_begin_time DATETIME;

        # 10 天一个间隔
        # TIME_RANGE_TYPE_TEN_DAYS int64 = 100
        SET time_range_type = 100;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DISTINCT product_id FROM erp_product_sale_day_analysis;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP
            -- 提取游标里的数据
            FETCH cur INTO v_product_id;

            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;

            -- 这里做你想做的循环的事件
            SET v_latest_begin_time = NULL;
            SELECT ta.begin_time INTO v_latest_begin_time
            FROM erp_product_sale_trend_analysis ta
            WHERE ta.product_id = v_product_id AND ta.type_id = time_range_type
            ORDER BY ta.begin_time DESC
            LIMIT 1;

            IF v_latest_begin_time IS NULL
            THEN
            ELSE
            END IF;

            SET done = FALSE;
        END LOOP;
        -- 关闭游标
        CLOSE cur;

    END;