DROP PROCEDURE IF EXISTS product_sale_rank;

CREATE PROCEDURE product_sale_rank()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_sa_id BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_daily_mean_sales DECIMAL(18, 6);
        DECLARE v_count BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        /*
            // 统计实际值
            PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL int64 = 100

            // 统计理论值
            PRODUCT_SALE_ANALYSIS_TYPE_LOGIC int64 = 200

            // 只统计SKU，不分仓库统计，实际值的不分仓统计
            PRODUCT_SALE_ANALYSIS_TYPE_SKU int64 = 300
         */

        -- 游标
        DECLARE product_id_cur CURSOR FOR
            SELECT DISTINCT
                sa.id,
                sa.product_id,
                sa.daily_mean_sales
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 300;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        -- 打开游标
        OPEN product_id_cur;

        -- 开始循环
        product_id_loop: LOOP
            -- 提取游标里的数据
            SET done = FALSE;
            FETCH product_id_cur INTO v_sa_id, v_product_id, v_daily_mean_sales;

            IF done
            THEN
                LEAVE product_id_loop;
            END IF;

            -- 这里做你想做的循环的事件

            SELECT count(sa.id) INTO v_count
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 300 AND sa.daily_mean_sales > v_daily_mean_sales;

            UPDATE erp_product_sale_analysis
            SET `rank` = v_count + 1
            WHERE id = v_sa_id;

        END LOOP;
        -- 关闭游标
        CLOSE product_id_cur;
    END;
