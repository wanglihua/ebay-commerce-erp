/*

根据 6 次采购的平均价格，更新商品的采购价
没有采购记录的商品，不用进行 采购价 更新

根据 6 次采购的平均库存价格，更新商品的中国仓成本价
没有采购记录的商品，不用进行 中国仓的成本价 更新

CALL calculate_product_price;

*/

DROP PROCEDURE IF EXISTS calculate_product_price;

CREATE PROCEDURE calculate_product_price()
    BEGIN

        DECLARE v_product_id BIGINT;
        DECLARE v_cn_stock_id BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DISTINCT pbe.product_id
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id = 4700; # 已入库;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        /*
            // STOCK_TYPE_CN 国内仓库
            STOCK_TYPE_CN int64 = 100
         */

        SELECT s.id
        INTO v_cn_stock_id
        FROM erp_stock s
        WHERE s.stock_type_id = 100 # 100 国内仓库
        LIMIT 1;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP

            -- 提取游标里的数据
            FETCH cur
            INTO v_product_id;

            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;

            -- 这里做你想做的循环的事件
            UPDATE erp_product
            SET cost_price = (
                SELECT avg(pbe_six.price) AS avg_price
                FROM (SELECT *
                      FROM erp_purchase_bill_entry pbe
                      WHERE pbe.product_id = v_product_id AND pbe.status_id = 4700 # 已入库
                      ORDER BY pbe.update_time DESC
                      LIMIT 6) pbe_six)
            WHERE id = v_product_id;

            UPDATE erp_inventory
            SET cost_price = (
                SELECT avg(pbe_six.inventory_cost) AS inventory_cost
                FROM (SELECT *
                      FROM erp_purchase_bill_entry pbe
                      WHERE pbe.product_id = v_product_id AND pbe.status_id = 4700 # 已入库
                      ORDER BY pbe.update_time DESC
                      LIMIT 6) pbe_six)
            WHERE product_id = v_product_id AND stock_id = v_cn_stock_id;

            SET done = FALSE;
        END LOOP;

        -- 关闭游标
        CLOSE cur;

    END;
