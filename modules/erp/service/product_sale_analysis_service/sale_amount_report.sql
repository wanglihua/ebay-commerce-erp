# CALL sale_amount_report();

DROP PROCEDURE IF EXISTS sale_amount_report;

CREATE PROCEDURE sale_amount_report()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE v_latest_reported_paid_time DATETIME; #  最后缓存起来的报表时间
        DECLARE v_paid_time DATETIME;
        DECLARE v_platform_id BIGINT;
        DECLARE v_shop_id BIGINT;
        DECLARE v_order_id BIGINT;
        DECLARE v_order_item_id BIGINT;
        DECLARE v_paid_amount DECIMAL(18, 2);
        DECLARE v_express_fee DECIMAL(18, 2);
        DECLARE v_income DECIMAL(18, 2);
        DECLARE v_cost_amount DECIMAL(18, 2);
        DECLARE v_estimate_logistics_fee DECIMAL(18, 2);
        DECLARE v_actual_logistics_fee DECIMAL(18, 2);
        DECLARE v_jyfl DECIMAL(18, 2);
        DECLARE v_ptfl DECIMAL(18, 2);
        DECLARE v_jyfy_amount DECIMAL(18, 2);
        DECLARE v_ptfy_amount DECIMAL(18, 2);
        DECLARE v_pick_fee DECIMAL(18, 2);
        DECLARE v_gst_fee DECIMAL(18, 2);
        DECLARE v_outcome DECIMAL(18, 2);
        DECLARE v_profit DECIMAL(18, 2);
        DECLARE v_profit_rate DECIMAL(18, 2);
        DECLARE v_product_id BIGINT;
        DECLARE v_sale_user_id BIGINT;
        DECLARE v_service_user_id BIGINT;
        DECLARE v_dev_user_id BIGINT;
        DECLARE v_order_item_count INT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE order_cur CURSOR FOR
            SELECT
                o.id          AS order_id,
                o.paid_time   AS paid_time,
                o.shop_id     AS shop_id,
                o.amount_paid AS paid_amount,
                o.express_fee AS express_fee,
                o.profit_au   AS profit
            FROM erp_order o
            WHERE (o.status_id = 4200 OR o.status_id = 4300)
                  AND o.paid_time >= v_latest_reported_paid_time;
        # 条件里面要加上 =，防止可能的有多个相同的 支付时间 的情况

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        SELECT sar.paid_time AS paid_time
        INTO v_latest_reported_paid_time
        FROM erp_sale_amount_report sar
        ORDER BY sar.paid_time DESC
        LIMIT 1;

        IF v_latest_reported_paid_time IS NULL
        THEN
            SET v_latest_reported_paid_time = '2016-01-01';
        END IF;

        OPEN order_cur;

        order_id_loop: LOOP

            SET done = FALSE;
            FETCH order_cur INTO v_order_id, v_paid_time, v_shop_id, v_paid_amount, v_express_fee, v_profit;

            IF done
            THEN
                LEAVE order_id_loop;
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

            SELECT count(oi.id) AS item_count
            INTO v_order_item_count
            FROM erp_order_item oi
            WHERE oi.order_id = v_order_id;

            IF v_order_item_count = 1
            THEN

            END IF;

        END LOOP;

        CLOSE order_cur;

    END;
/*
# erp_sale_amount_report
CREATE TABLE erp_sale_amount_report (
    id                     BIGINT AUTO_INCREMENT PRIMARY KEY,
    paid_time              DATETIME       NOT NULL,
    platform_id            BIGINT         NOT NULL,
    shop_id                BIGINT         NOT NULL,
    order_id               BIGINT         NOT NULL,
    order_item_id          BIGINT         NOT NULL,
    paid_amount            DECIMAL(18, 2) NOT NULL,
    express_fee            DECIMAL(18, 2) NOT NULL,
    income                 DECIMAL(18, 2) NOT NULL,
    cost_amount            DECIMAL(18, 2) NOT NULL,
    estimate_logistics_fee DECIMAL(18, 2) NOT NULL,
    actual_logistics_fee   DECIMAL(18, 2) NOT NULL,
    jyfy_amount            DECIMAL(18, 2) NOT NULL,
    ptfy_amount            DECIMAL(18, 2) NOT NULL,
    pick_fee               DECIMAL(18, 2) NOT NULL,
    gst_fee                DECIMAL(18, 2) NOT NULL,
    outcome                DECIMAL(18, 2) NOT NULL,
    profit                 DECIMAL(18, 2) NOT NULL,
    profit_rate            DECIMAL(18, 2) NOT NULL,
    product_id             BIGINT         NOT NULL,
    sale_user_id           BIGINT         NOT NULL,
    service_user_id        BIGINT         NOT NULL,
    dev_user_id            BIGINT         NOT NULL
);

CREATE INDEX IDX_erp_sale_amount_report_paid_time ON erp_sale_amount_report (paid_time);
CREATE INDEX IDX_erp_sale_amount_report_platform_id ON erp_sale_amount_report (platform_id);
CREATE INDEX IDX_erp_sale_amount_report_shop_id ON erp_sale_amount_report (shop_id);
CREATE INDEX IDX_erp_sale_amount_report_order_id ON erp_sale_amount_report (order_id);
CREATE INDEX IDX_erp_sale_amount_report_order_item_id ON erp_sale_amount_report (order_item_id);
CREATE INDEX IDX_erp_sale_amount_report_product_id ON erp_sale_amount_report (product_id);
CREATE INDEX IDX_erp_sale_amount_report_sale_user_id ON erp_sale_amount_report (sale_user_id);
CREATE INDEX IDX_erp_sale_amount_report_service_user_id ON erp_sale_amount_report (service_user_id);
CREATE INDEX IDX_erp_sale_amount_report_dev_user_id ON erp_sale_amount_report (dev_user_id);

# erp_order
CREATE INDEX IDX_erp_order_paid_time ON erp_order (paid_time);
 */