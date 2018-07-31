# CALL product_sale_rank();

DROP PROCEDURE IF EXISTS product_sale_rank;

CREATE PROCEDURE product_sale_rank()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_sa_id BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_daily_mean_sales DECIMAL(18, 6);
        DECLARE v_count BIGINT;
        DECLARE v_mel_stock_id BIGINT;
        DECLARE v_bne_stock_id BIGINT;

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
        DECLARE product_analysis_cur CURSOR FOR
            SELECT DISTINCT
                sa.id,
                sa.product_id,
                sa.daily_mean_sales
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 300;

        DECLARE product_analysis_mel_cur CURSOR FOR
            SELECT DISTINCT
                sa.id,
                sa.product_id,
                sa.daily_mean_sales
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 100 AND sa.stock_id = v_mel_stock_id;

        DECLARE product_analysis_bne_cur CURSOR FOR
            SELECT DISTINCT
                sa.id,
                sa.product_id,
                sa.daily_mean_sales
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 100 AND sa.stock_id = v_bne_stock_id;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        # 不分仓
        OPEN product_analysis_cur;

        product_analysis_loop: LOOP
            SET done = FALSE;
            FETCH product_analysis_cur INTO v_sa_id, v_product_id, v_daily_mean_sales;

            IF done
            THEN
                LEAVE product_analysis_loop;
            END IF;

            SELECT count(sa.id) INTO v_count
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 300 AND sa.daily_mean_sales > v_daily_mean_sales;

            UPDATE erp_product_sale_analysis
            SET `rank` = v_count + 1
            WHERE id = v_sa_id;

        END LOOP;

        CLOSE product_analysis_cur;

        # MEL 仓库
        SELECT t.id INTO v_mel_stock_id
        FROM erp_stock t
        WHERE t.short_name = 'MEL';

        OPEN product_analysis_mel_cur;

        product_analysis_mel_loop: LOOP
            SET done = FALSE;
            FETCH product_analysis_mel_cur INTO v_sa_id, v_product_id, v_daily_mean_sales;

            IF done
            THEN
                LEAVE product_analysis_mel_loop;
            END IF;

            SELECT count(sa.id) INTO v_count
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 100 AND sa.stock_id = v_mel_stock_id AND sa.daily_mean_sales > v_daily_mean_sales;

            UPDATE erp_product_sale_analysis
            SET `rank` = v_count + 1
            WHERE id = v_sa_id;

        END LOOP;

        CLOSE product_analysis_mel_cur;

        # BNE 仓库
        SELECT t.id INTO v_bne_stock_id
        FROM erp_stock t
        WHERE t.short_name = 'BNE';

        OPEN product_analysis_bne_cur;

        product_analysis_bne_loop: LOOP
            SET done = FALSE;
            FETCH product_analysis_bne_cur INTO v_sa_id, v_product_id, v_daily_mean_sales;

            IF done
            THEN
                LEAVE product_analysis_bne_loop;
            END IF;

            SELECT count(sa.id) INTO v_count
            FROM erp_product_sale_analysis sa
            WHERE sa.type = 100 AND sa.stock_id = v_bne_stock_id AND sa.daily_mean_sales > v_daily_mean_sales;

            UPDATE erp_product_sale_analysis
            SET `rank` = v_count + 1
            WHERE id = v_sa_id;

        END LOOP;

        CLOSE product_analysis_bne_cur;

    END;
