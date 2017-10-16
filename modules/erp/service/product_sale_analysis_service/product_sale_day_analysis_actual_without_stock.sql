/*
	// 统计实际值
	PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL int64 = 100
	// 统计理论值
	PRODUCT_SALE_ANALYSIS_TYPE_LOGIC int64 = 200
	// 只统计SKU，不分仓库统计
	PRODUCT_SALE_ANALYSIS_TYPE_SKU int64 = 300
 */

DROP PROCEDURE product_sale_day_analysis_actual_without_stock;

CREATE PROCEDURE product_sale_day_analysis_actual_without_stock()
    BEGIN

        DECLARE analysis_begin_time DATETIME;
        -- 需要定义接收游标数据的变量
        DECLARE v_shipped_day DATETIME;
        DECLARE v_product_id BIGINT;
        DECLARE v_sku VARCHAR(200);
        DECLARE v_qty BIGINT;
        DECLARE da_count INT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DATE(o.shipped_time) AS shipped_day, oi.item_product_id AS product_id, p.sku AS sku, sum(oi.quantity) AS qty
            FROM erp_order_item oi
                INNER JOIN erp_order o ON o.id = oi.order_id
                INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
                INNER JOIN erp_product p ON p.id = oi.item_product_id
            WHERE o.shipped_time >= analysis_begin_time AND o.stock_id <> 0 AND ol.stock_id <> 0 AND o.status_id = 4300 # 4300 已完成
            GROUP BY DATE(o.shipped_time), oi.item_product_id
            ORDER BY DATE(o.shipped_time) ASC, oi.item_product_id ASC;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        # 获取需要分析的开始时间
        SELECT da.day
        INTO analysis_begin_time
        FROM erp_product_sale_day_analysis da
        ORDER BY da.day DESC
        LIMIT 1;

        IF analysis_begin_time IS NULL
        THEN
            SET analysis_begin_time = '2017-01-01';
        END IF;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP
            -- 提取游标里的数据
            FETCH cur
            INTO v_shipped_day, v_product_id, v_sku, v_qty;

            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;

            -- 这里做你想做的循环的事件
            SELECT count(*)
            INTO da_count
            FROM erp_product_sale_day_analysis da
            WHERE da.day = v_shipped_day AND da.product_id = v_product_id AND da.stock_id = 0 AND da.type = 300; # 300 不分仓库统计

            IF da_count = 0
            THEN
                INSERT INTO erp_product_sale_day_analysis (stock_id, product_id, sku, day, qty, type, create_time, update_time) VALUES (0, product_id, sku, v_shipped_day, qty, 300, now(), now());
            ELSE
                UPDATE erp_product_sale_day_analysis
                SET sku = v_sku, qty = v_qty, update_time = now()
                WHERE day = v_shipped_day AND product_id = v_product_id AND stock_id = 0 AND type = 300;
            END IF;

        END LOOP;
        -- 关闭游标
        CLOSE cur;

    END

