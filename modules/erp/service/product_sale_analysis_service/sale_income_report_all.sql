# CALL sale_income_report_all();

DROP PROCEDURE IF EXISTS sale_income_report_all;

CREATE PROCEDURE sale_income_report_all()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_latest_reported_paid_time DATETIME; #  最后缓存起来的报表时间
        DECLARE v_paid_time DATETIME;
        DECLARE v_platform_id BIGINT;
        DECLARE v_shop_id BIGINT;
        DECLARE v_order_id BIGINT;
        DECLARE v_order_item_id BIGINT;
        DECLARE v_paid_amount DECIMAL(18, 6);
        DECLARE v_express_fee DECIMAL(18, 6);
        DECLARE v_income DECIMAL(18, 6);
        DECLARE v_cost_amount DECIMAL(18, 6);
        DECLARE v_logistics_fee DECIMAL(18, 6);
        DECLARE v_estimate_logistics_fee DECIMAL(18, 6);
        DECLARE v_actual_logistics_fee DECIMAL(18, 6);
        DECLARE v_jyfl DECIMAL(18, 6); # 后来没有使用这里面的值
        DECLARE v_ptfl DECIMAL(18, 6); # 后来没有使用这里面的值
        DECLARE v_jyfy_amount DECIMAL(18, 6);
        DECLARE v_ptfy_amount DECIMAL(18, 6);
        DECLARE v_pick_fee DECIMAL(18, 6);
        DECLARE v_gst_fee DECIMAL(18, 6);
        DECLARE v_outcome DECIMAL(18, 6);
        DECLARE v_profit DECIMAL(18, 6);
        DECLARE v_profit_rate DECIMAL(18, 6);
        DECLARE v_exchange_rate DECIMAL(18, 6);
        DECLARE v_product_count BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_sale_user_id BIGINT;
        DECLARE v_service_user_id BIGINT;
        DECLARE v_dev_user_id BIGINT;
        DECLARE v_order_item_count INT;
        DECLARE v_sale_income_report_count INT;
        DECLARE v_item_price DECIMAL(18, 6);
        DECLARE v_item_qty INT;
        DECLARE v_item_amount DECIMAL(18, 6);
        DECLARE v_item_total_amount DECIMAL(18, 6);
        DECLARE v_item_paid_amount DECIMAL(18, 6);
        DECLARE v_item_express_fee DECIMAL(18, 6);
        DECLARE v_item_income DECIMAL(18, 6);
        DECLARE v_item_cost_amount DECIMAL(18, 6);
        DECLARE v_item_estimate_logistics_fee DECIMAL(18, 6);
        DECLARE v_item_actual_logistics_fee DECIMAL(18, 6);
        DECLARE v_item_jyfy_amount DECIMAL(18, 6);
        DECLARE v_item_ptfy_amount DECIMAL(18, 6);
        DECLARE v_item_pick_fee DECIMAL(18, 6);
        DECLARE v_item_gst_fee DECIMAL(18, 6);
        DECLARE v_item_outcome DECIMAL(18, 6);

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE order_cur CURSOR FOR
            SELECT
                o.id              AS order_id,
                o.paid_time       AS paid_time,
                o.shop_id         AS shop_id,
                o.amount_paid     AS paid_amount,
                o.express_fee     AS express_fee,
                o.exchange_rate   AS exchange_rate,
                ifnull(o.ptfy, 0) AS ptfy,
                ifnull(o.jyfy, 0) AS jyfy
            FROM erp_order o;
        # 退后 15 天是等待后面的实际物流数据的录入

        DECLARE order_item_cur CURSOR FOR
            SELECT
                oi.id                AS order_item_id,
                oi.item_product_id   AS product_id,
                oi.transaction_price AS item_price,
                oi.quantity          AS item_qty
            FROM erp_order_item oi
            WHERE oi.order_id = v_order_id;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        SELECT sar.paid_time AS paid_time
        INTO v_latest_reported_paid_time
        FROM erp_sale_income_report sar
        ORDER BY sar.paid_time DESC
        LIMIT 1;

        IF v_latest_reported_paid_time IS NULL
        THEN
            SET v_latest_reported_paid_time = '2016-01-01';
        END IF;

        OPEN order_cur;

        order_loop: LOOP

            SET done = FALSE;
            FETCH order_cur INTO v_order_id, v_paid_time, v_shop_id, v_paid_amount, v_express_fee, v_exchange_rate, v_ptfy_amount, v_jyfy_amount;

            IF done
            THEN
                LEAVE order_loop;
            END IF;

            IF v_paid_time IS NULL
            THEN
                ITERATE order_loop;
            END IF;

            IF v_paid_amount = 0
            THEN
                ITERATE order_loop;
            END IF;

            SELECT sa.platform_id AS platform_id
            INTO v_platform_id
            FROM erp_shop_account sa
            WHERE sa.id = v_shop_id;

            IF v_platform_id IS NULL
            THEN
                SET v_platform_id = 0;
            END IF;

            SELECT
                ol.cost_amount AS cost_amount,
                ol.total_fee   AS estimate_logistics_fee,
                ol.actual_fee  AS actual_logistics_fee,
                ol.jyfl        AS jyfl,
                ol.ptfl        AS ptfl,
                ol.pick_fee    AS pick_fee
            INTO
                v_cost_amount,
                v_estimate_logistics_fee,
                v_actual_logistics_fee,
                v_jyfl,
                v_ptfl,
                v_pick_fee
            FROM erp_order_logistics ol
            WHERE ol.order_id = v_order_id;

            IF v_cost_amount IS NULL
            THEN
                SET v_cost_amount = 0;
            END IF;

            IF v_estimate_logistics_fee IS NULL
            THEN
                SET v_estimate_logistics_fee = 0;
            END IF;

            IF v_actual_logistics_fee IS NULL
            THEN
                SET v_actual_logistics_fee = 0;
            END IF;

            IF v_jyfl IS NULL
            THEN
                SET v_jyfl = 0;
            END IF;

            IF v_ptfl IS NULL
            THEN
                SET v_ptfl = 0;
            END IF;

            IF v_pick_fee IS NULL
            THEN
                SET v_pick_fee = 0;
            END IF;

            # 采用新的处理方式
            # SET v_jyfy_amount = v_paid_amount * v_jyfl;
            # SET v_ptfy_amount = v_paid_amount * v_ptfl;

            SET v_gst_fee = (v_paid_amount / 1.1) * 0.1;

            SELECT count(oi.id) AS item_count
            INTO v_order_item_count
            FROM erp_order_item oi
            WHERE oi.order_id = v_order_id;

            IF v_order_item_count = 0
            THEN
                ITERATE order_loop;
            ELSEIF v_order_item_count = 1
                THEN
                    SET v_income = v_paid_amount + v_express_fee;
                    # SET v_income = v_paid_amount;

                    SET v_logistics_fee = v_estimate_logistics_fee;
                    IF v_actual_logistics_fee IS NOT NULL AND v_actual_logistics_fee <> 0
                    THEN
                        SET v_logistics_fee = v_actual_logistics_fee;
                    END IF;

                    SET v_outcome = v_cost_amount + v_logistics_fee + v_pick_fee + v_jyfy_amount + v_ptfy_amount + v_gst_fee;
                    SET v_profit = v_income - v_outcome;

                    IF v_income = 0
                    THEN
                        ITERATE order_loop;
                    END IF;

                    # SET v_profit_rate = v_profit / v_paid_amount * 100;
                    SET v_profit_rate = v_profit / v_income * 100;

                    SELECT
                        oi.id              AS order_item_id,
                        oi.item_product_id AS product_id,
                        oi.quantity        AS item_qty
                    INTO
                        v_order_item_id,
                        v_product_id,
                        v_item_qty
                    FROM erp_order_item oi
                    WHERE oi.order_id = v_order_id
                    LIMIT 1;

                    SELECT count(*) AS count
                    INTO v_product_count
                    FROM erp_product p
                    WHERE p.id = v_product_id;

                    IF v_product_count = 0
                    THEN
                        ITERATE order_loop;
                    END IF;

                    SELECT
                        p.sale_user_id,
                        p.service_user_id,
                        p.dev_user_id
                    INTO
                        v_sale_user_id,
                        v_service_user_id,
                        v_dev_user_id
                    FROM erp_product p
                    WHERE p.id = v_product_id;

                    SELECT count(sar.id) AS sale_income_report_count
                    INTO v_sale_income_report_count
                    FROM erp_sale_income_report sar
                    WHERE sar.order_item_id = v_order_item_id;

                    IF v_sale_income_report_count = 0
                    THEN
                        INSERT INTO erp_sale_income_report
                        (paid_time, platform_id, shop_id, order_id, order_item_id, paid_amount, express_fee, income, cost_amount, estimate_logistics_fee, actual_logistics_fee, jyfy_amount, ptfy_amount, pick_fee, gst_fee, outcome, profit, profit_rate, exchange_rate, product_id, product_qty, sale_user_id, service_user_id, dev_user_id)
                        VALUES (v_paid_time, v_platform_id, v_shop_id, v_order_id, v_order_item_id, v_paid_amount, v_express_fee, v_income, v_cost_amount, v_estimate_logistics_fee, v_actual_logistics_fee, v_jyfy_amount, v_ptfy_amount, v_pick_fee, v_gst_fee, v_outcome, v_profit, v_profit_rate, v_exchange_rate, v_product_id, v_item_qty, v_sale_user_id, v_service_user_id, v_dev_user_id);
                    ELSE
                        UPDATE erp_sale_income_report
                        SET paid_time              = v_paid_time,
                            platform_id            = v_platform_id,
                            shop_id                = v_shop_id,
                            order_id               = v_order_id,
                            paid_amount            = v_paid_amount,
                            express_fee            = v_express_fee,
                            income                 = v_income,
                            cost_amount            = v_cost_amount,
                            estimate_logistics_fee = v_estimate_logistics_fee,
                            actual_logistics_fee   = v_actual_logistics_fee,
                            jyfy_amount            = v_jyfy_amount,
                            ptfy_amount            = v_ptfy_amount,
                            pick_fee               = v_pick_fee,
                            gst_fee                = v_gst_fee,
                            outcome                = v_outcome,
                            profit                 = v_profit,
                            profit_rate            = v_profit_rate,
                            exchange_rate          = v_exchange_rate,
                            product_id             = v_product_id,
                            product_qty            = v_item_qty,
                            sale_user_id           = v_sale_user_id,
                            service_user_id        = v_service_user_id,
                            dev_user_id            = v_dev_user_id
                        WHERE order_item_id = v_order_item_id;
                    END IF;
            ELSE

                SELECT sum(oi.transaction_price * oi.quantity) AS item_total_amount
                INTO v_item_total_amount
                FROM erp_order_item oi
                WHERE oi.order_id = v_order_id;

                IF v_item_total_amount = 0
                THEN
                    ITERATE order_loop;
                END IF;

                OPEN order_item_cur;

                order_item_loop: LOOP

                    SET done = FALSE;
                    FETCH order_item_cur INTO v_order_item_id, v_product_id, v_item_price, v_item_qty;

                    IF done
                    THEN
                        LEAVE order_item_loop;
                    END IF;

                    SELECT count(*) AS count
                    INTO v_product_count
                    FROM erp_product p
                    WHERE p.id = v_product_id;

                    IF v_product_count = 0
                    THEN
                        ITERATE order_item_loop;
                    END IF;

                    SELECT
                        p.sale_user_id,
                        p.service_user_id,
                        p.dev_user_id
                    INTO
                        v_sale_user_id,
                        v_service_user_id,
                        v_dev_user_id
                    FROM erp_product p
                    WHERE p.id = v_product_id;

                    SELECT count(sar.id) AS sale_income_report_count
                    INTO v_sale_income_report_count
                    FROM erp_sale_income_report sar
                    WHERE sar.order_item_id = v_order_item_id;

                    SET v_item_amount = v_item_price * v_item_qty;
                    IF v_item_amount = 0
                    THEN
                        ITERATE order_item_loop;
                    END IF;

                    SET v_item_paid_amount = v_paid_amount * v_item_amount / v_item_total_amount;
                    SET v_item_express_fee = v_express_fee * v_item_amount / v_item_total_amount;
                    SET v_item_income = v_item_paid_amount + v_item_express_fee;
                    SET v_item_cost_amount = v_cost_amount * v_item_amount / v_item_total_amount;
                    SET v_item_estimate_logistics_fee = v_estimate_logistics_fee * v_item_amount / v_item_total_amount;
                    SET v_item_actual_logistics_fee = v_actual_logistics_fee * v_item_amount / v_item_total_amount;
                    SET v_item_jyfy_amount = v_jyfy_amount * v_item_amount / v_item_total_amount;
                    SET v_item_ptfy_amount = v_ptfy_amount * v_item_amount / v_item_total_amount;
                    SET v_item_pick_fee = v_pick_fee * v_item_amount / v_item_total_amount;
                    SET v_item_gst_fee = v_gst_fee * v_item_amount / v_item_total_amount;

                    SET v_logistics_fee = v_item_estimate_logistics_fee;
                    IF v_item_actual_logistics_fee IS NOT NULL AND v_item_actual_logistics_fee <> 0
                    THEN
                        SET v_logistics_fee = v_actual_logistics_fee;
                    END IF;

                    SET v_item_outcome = v_item_cost_amount + v_logistics_fee + v_item_pick_fee + v_item_jyfy_amount + v_item_ptfy_amount + v_item_gst_fee;
                    SET v_profit = v_item_income - v_item_outcome;

                    IF v_item_income = 0
                    THEN
                        ITERATE order_item_loop;
                    END IF;

                    # SET v_profit_rate = v_profit / v_item_paid_amount * 100;
                    SET v_profit_rate = v_profit / v_item_income * 100;

                    IF v_sale_income_report_count = 0
                    THEN
                        INSERT INTO erp_sale_income_report
                        (paid_time, platform_id, shop_id, order_id, order_item_id, paid_amount, express_fee, income, cost_amount, estimate_logistics_fee, actual_logistics_fee, jyfy_amount, ptfy_amount, pick_fee, gst_fee, outcome, profit, profit_rate, product_id, exchange_rate, product_qty, sale_user_id, service_user_id, dev_user_id)
                        VALUES (v_paid_time, v_platform_id, v_shop_id, v_order_id, v_order_item_id, v_item_paid_amount, v_item_express_fee, v_item_income, v_item_cost_amount, v_item_estimate_logistics_fee, v_item_actual_logistics_fee, v_item_jyfy_amount, v_item_ptfy_amount, v_item_pick_fee, v_item_gst_fee, v_item_outcome, v_profit, v_profit_rate, v_product_id, v_exchange_rate, v_item_qty, v_sale_user_id, v_service_user_id, v_dev_user_id);
                    ELSE
                        UPDATE erp_sale_income_report
                        SET paid_time              = v_paid_time,
                            platform_id            = v_platform_id,
                            shop_id                = v_shop_id,
                            order_id               = v_order_id,
                            paid_amount            = v_item_paid_amount,
                            express_fee            = v_item_express_fee,
                            income                 = v_item_income,
                            cost_amount            = v_item_cost_amount,
                            estimate_logistics_fee = v_item_estimate_logistics_fee,
                            actual_logistics_fee   = v_item_actual_logistics_fee,
                            jyfy_amount            = v_item_jyfy_amount,
                            ptfy_amount            = v_item_ptfy_amount,
                            pick_fee               = v_item_pick_fee,
                            gst_fee                = v_item_gst_fee,
                            outcome                = v_item_outcome,
                            profit                 = v_profit,
                            profit_rate            = v_profit_rate,
                            exchange_rate          = v_exchange_rate,
                            product_id             = v_product_id,
                            product_qty            = v_item_qty,
                            sale_user_id           = v_sale_user_id,
                            service_user_id        = v_service_user_id,
                            dev_user_id            = v_dev_user_id
                        WHERE order_item_id = v_order_item_id;
                    END IF;

                END LOOP;

                CLOSE order_item_cur;
            END IF;

        END LOOP;

        CLOSE order_cur;

    END;
/*
# erp_sale_income_report
CREATE TABLE erp_sale_income_report (
    id                     BIGINT AUTO_INCREMENT PRIMARY KEY,
    paid_time              DATETIME       NOT NULL,
    platform_id            BIGINT         NOT NULL,
    shop_id                BIGINT         NOT NULL,
    order_id               BIGINT         NOT NULL,
    order_item_id          BIGINT         NOT NULL,
    paid_amount            DECIMAL(18, 6) NOT NULL,
    express_fee            DECIMAL(18, 6) NOT NULL,
    income                 DECIMAL(18, 6) NOT NULL,
    cost_amount            DECIMAL(18, 6) NOT NULL,
    estimate_logistics_fee DECIMAL(18, 6) NOT NULL,
    actual_logistics_fee   DECIMAL(18, 6) NOT NULL,
    jyfy_amount            DECIMAL(18, 6) NOT NULL,
    ptfy_amount            DECIMAL(18, 6) NOT NULL,
    pick_fee               DECIMAL(18, 6) NOT NULL,
    gst_fee                DECIMAL(18, 6) NOT NULL,
    outcome                DECIMAL(18, 6) NOT NULL,
    profit                 DECIMAL(18, 6) NOT NULL,
    profit_rate            DECIMAL(18, 6) NOT NULL,
    exchange_rate          DECIMAL(18, 6) DEFAULT 5 NOT NULL,
    product_id             BIGINT         NOT NULL,
    product_qty            BIGINT         NOT NULL,
    sale_user_id           BIGINT         NOT NULL,
    service_user_id        BIGINT         NOT NULL,
    dev_user_id            BIGINT         NOT NULL
);

CREATE INDEX IDX_erp_sale_income_report_paid_time ON erp_sale_income_report (paid_time);
CREATE INDEX IDX_erp_sale_income_report_platform_id ON erp_sale_income_report (platform_id);
CREATE INDEX IDX_erp_sale_income_report_shop_id ON erp_sale_income_report (shop_id);
CREATE INDEX IDX_erp_sale_income_report_order_id ON erp_sale_income_report (order_id);
CREATE INDEX IDX_erp_sale_income_report_order_item_id ON erp_sale_income_report (order_item_id);
CREATE INDEX IDX_erp_sale_income_report_product_id ON erp_sale_income_report (product_id);
CREATE INDEX IDX_erp_sale_income_report_sale_user_id ON erp_sale_income_report (sale_user_id);
CREATE INDEX IDX_erp_sale_income_report_service_user_id ON erp_sale_income_report (service_user_id);
CREATE INDEX IDX_erp_sale_income_report_dev_user_id ON erp_sale_income_report (dev_user_id);

# erp_order
CREATE INDEX IDX_erp_order_paid_time ON erp_order (paid_time);
 */