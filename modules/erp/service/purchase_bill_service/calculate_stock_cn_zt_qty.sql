/*

CALL calculate_stock_cn_zt_qty()

*/

DROP PROCEDURE IF EXISTS calculate_stock_cn_zt_qty;

CREATE PROCEDURE calculate_stock_cn_zt_qty()
    BEGIN
        DECLARE v_product_id BIGINT;
        DECLARE v_cn_zt_stock_id BIGINT;
        DECLARE v_product_qty BIGINT;
        DECLARE v_stock_count BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DISTINCT product_id
            FROM (
                     SELECT nvt.product_id AS product_id
                     FROM erp_inventory nvt
                     WHERE nvt.stock_id = v_cn_zt_stock_id
                     UNION
                     SELECT pbe.product_id AS product_id
                     FROM erp_purchase_bill_entry pbe
                     WHERE pbe.status_id = 4600 # 4600 采购单状态 采购中
                 ) p;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        SELECT s.id INTO v_cn_zt_stock_id
        FROM erp_stock s
        WHERE s.stock_type_id = 80 # 中国仓在途
        LIMIT 1;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP

            SET done = FALSE;

            -- 提取游标里的数据
            FETCH cur
            INTO v_product_id;

            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;

            -- 这里做你想做的循环的事件

            SELECT sum(pbe.qty)
            INTO v_product_qty
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id = 4600 AND pbe.product_id = v_product_id; # 4600 采购单状态 采购中

            IF v_product_qty IS NULL
            THEN
                SET v_product_qty = 0;
            END IF;

            SELECT count(*) INTO v_stock_count
            FROM erp_inventory nvt
            WHERE nvt.stock_id = v_cn_zt_stock_id AND nvt.product_id = v_product_id;

            IF v_stock_count = 0
            THEN
                INSERT INTO erp_inventory (stock_id, storage_loc_id, product_id, qty, create_time, update_time, cost_price)
                VALUES (v_cn_zt_stock_id, 0, v_product_id, v_product_qty, now(), now(), 0);
            ELSE
                UPDATE erp_inventory
                SET qty         = v_product_qty,
                    update_time = now()
                WHERE stock_id = v_cn_zt_stock_id AND product_id = v_product_id;
            END IF;
        END LOOP;

        -- 关闭游标
        CLOSE cur;

    END;