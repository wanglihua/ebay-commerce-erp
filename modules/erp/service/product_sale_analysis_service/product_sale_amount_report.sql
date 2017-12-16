DROP PROCEDURE IF EXISTS product_sale_amount_report;

CREATE PROCEDURE product_sale_amount_report()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_sa_id BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_daily_mean_sales DECIMAL(18, 6);
        DECLARE v_count BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE product_id_cur CURSOR FOR
            SELECT p.id AS product_id
            FROM erp_product p
            ORDER BY p.id;

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

            -- 这里做你想做的循环的事件

            /*
            SELECT count(sd.id) INTO v_count
            FROM erp_sale_report_time_day sd
            WHERE sa.type = 300 AND sa.daily_mean_sales > v_daily_mean_sales;

            UPDATE erp_sale_report_time_day
            SET rank = v_count + 1
            WHERE id = v_sa_id;
            */

        END LOOP;

        -- 关闭游标
        CLOSE product_id_cur;
    END;
