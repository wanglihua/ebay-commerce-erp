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
        DECLARE v_ppqzlzl DECIMAL(18, 6);
        DECLARE v_ppqzbzl DECIMAL(18, 6);
        DECLARE v_pphzlzl DECIMAL(18, 6);
        DECLARE v_pphzbzl DECIMAL(18, 6);
        DECLARE v_order_address_id BIGINT;
        DECLARE v_count BIGINT;

        DECLARE done INT DEFAULT FALSE;

        DECLARE order_cur CURSOR FOR
            SELECT
                o.id,
                o.shop_id,
                o.amount_paid,
                o.profit_au,
                date(o.shipped_time),
                o.address_id
            FROM erp_order o
            # WHERE o.shipped_time IS NOT NULL;
            WHERE o.shipped_time > date_add(date(now()), INTERVAL -3 DAY); # 只计算3天前的订单

        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        OPEN order_cur;

        order_cur_loop: LOOP

            SET done = FALSE;
            FETCH order_cur INTO v_order_id, v_shop_id, v_amount, v_profit, v_day_date, v_order_address_id;

            IF done
            THEN
                LEAVE order_cur_loop;
            END IF;

            SELECT
                ol.stock_id,
                ol.company_id,
                ol.mode_id,
                ol.length * ol.width * ol.height,
                ol.weight,
                ol.total_fee,
                ol.actual_fee,
                ol.ppqzlzl,
                ol.ppqzbzl,
                ol.pphzlzl,
                ol.pphzbzl
            INTO v_stock_id, v_company_id, v_mode_id, v_volume, v_weight, v_total_fee, v_acutal_total_fee, v_ppqzlzl, v_ppqzbzl, v_pphzlzl, v_pphzbzl
            FROM erp_order_logistics ol
            WHERE ol.order_id = v_order_id
            LIMIT 1;

            IF v_stock_id IS NULL
            THEN
                ITERATE order_cur_loop;
            END IF;

            SELECT osa.postal_code INTO v_postcode
            FROM erp_order_shipping_address osa
            WHERE osa.id = v_order_address_id
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
                                                  volume, weight, amount, total_fee, actual_total_fee, profit, ppqzlzl, ppqzbzl, pphzlzl, pphzbzl, create_time, update_time, month_date)
                VALUES (v_day_date, v_order_id, v_shop_id, v_stock_id, v_company_id, v_mode_id, v_postcode, v_prd_count,
                                    v_volume, v_weight, v_amount, v_total_fee, v_acutal_total_fee, v_profit, v_ppqzlzl, v_ppqzbzl, v_pphzlzl, v_pphzbzl, now(), now(), date_add(date(day_date), INTERVAL (1 - day(day_date)) DAY));
            ELSE
                UPDATE erp_logistics_report
                SET day_date         = v_day_date, order_id = v_order_id, shop_id = v_shop_id, stock_id = v_stock_id,
                    company_id       = v_company_id, mode_id = v_mode_id, postcode = v_postcode, prd_count = v_prd_count,
                    volume           = v_volume, weight = v_weight, amount = v_amount, total_fee = v_total_fee,
                    actual_total_fee = v_acutal_total_fee, profit = v_profit,
                    ppqzlzl          = v_ppqzlzl, ppqzbzl = v_ppqzbzl, pphzlzl = v_pphzlzl, pphzbzl = v_pphzbzl,
                    update_time      = now(), month_date = date_add(date(day_date), INTERVAL (1 - day(day_date)) DAY)
                WHERE order_id = v_order_id;
            END IF;

        END LOOP;

        CLOSE order_cur;

    END;
