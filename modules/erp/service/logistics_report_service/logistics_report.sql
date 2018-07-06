# CALL logistics_report()

DROP PROCEDURE IF EXISTS logistics_report;

CREATE PROCEDURE logistics_report()
    BEGIN
        DECLARE v_day_date DATE;
        DECLARE v_order_id BIGINT;
        DECLARE v_shop_id BIGINT;
        DECLARE v_stock_id BIGINT;
        DECLARE v_company_id BIGINT;
        DECLARE v_mode_id BIGINT;
        DECLARE v_postcode VARCHAR(200);
        DECLARE v_prd_count BIGINT;
        DECLARE v_volume DECIMAL(18, 6);
        DECLARE v_weight DECIMAL(18, 6);
        DECLARE v_amount DECIMAL(18, 6);
        DECLARE v_total_fee DECIMAL(18, 6);
        DECLARE v_acutal_total_fee DECIMAL(18, 6);
        DECLARE v_profit DECIMAL(18, 6);
        DECLARE v_count BIGINT;

        DECLARE done INT DEFAULT FALSE;

        DECLARE order_cur CURSOR FOR
            SELECT o.id, shop_id, amount_paid, profit_au, date(shipped_time)
            FROM erp_order o
            # WHERE o.shipped_time IS NOT NULL;
            WHERE o.shipped_time > date_add(date(now()), INTERVAL -3 DAY); # 只计算3天前的订单

        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        OPEN order_cur;

        order_cur_loop: LOOP

            SET done = FALSE;
            FETCH order_cur INTO v_order_id, v_shop_id, v_amount, v_profit, v_day_date;

            IF done
            THEN
                LEAVE order_cur_loop;
            END IF;

            SELECT ol.stock_id, ol.company_id, ol.mode_id, ol.length * ol.width * ol.height, ol.weight, total_fee, actual_fee
            INTO v_stock_id, v_company_id, v_mode_id, v_volume, v_weight, v_total_fee, v_acutal_total_fee
            FROM erp_order_logistics ol
            WHERE ol.order_id = v_order_id
            LIMIT 1;

            IF v_stock_id IS NULL
            THEN
                ITERATE order_cur_loop;
            END IF;

            SELECT osa.postal_code INTO v_postcode
            FROM erp_order_shipping_address osa
            WHERE osa.order_id = v_order_id
            LIMIT 1;

            IF v_postcode IS NULL
            THEN
                ITERATE order_cur_loop;
            END IF;

            SELECT sum(oi.quantity)
            INTO v_prd_count
            FROM erp_order_item oi
            WHERE oi.order_id = v_order_id;

            IF v_prd_count IS NULL
            THEN
                ITERATE order_cur_loop;
            END IF;

            SELECT count(lr.id) INTO v_count
            FROM erp_logistics_report lr
            WHERE lr.order_id = v_order_id;

            IF v_count = 0
            THEN
                INSERT INTO erp_logistics_report (day_date, order_id, shop_id, stock_id, company_id, mode_id, postcode, prd_count,
                                                  volume, weight, amount, total_fee, actual_total_fee, profit, create_time, update_time)
                VALUES (v_day_date, v_order_id, v_shop_id, v_stock_id, v_company_id, v_mode_id, v_postcode, v_prd_count,
                                    v_volume, v_weight, v_amount, v_total_fee, v_acutal_total_fee, v_profit, now(), now());
            ELSE
                UPDATE erp_logistics_report
                SET day_date         = v_day_date, order_id = v_order_id, shop_id = v_shop_id, stock_id = v_stock_id,
                    company_id       = v_company_id, mode_id = v_mode_id, postcode = v_postcode, prd_count = v_prd_count,
                    volume           = v_volume, weight = v_weight, amount = v_amount, total_fee = v_total_fee,
                    actual_total_fee = v_acutal_total_fee, profit = v_profit, update_time = now()
                WHERE order_id = v_order_id;
            END IF;

        END LOOP;

        CLOSE order_cur;

    END;
