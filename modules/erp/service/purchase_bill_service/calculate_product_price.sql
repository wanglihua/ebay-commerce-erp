/*

根据 6 次采购的平均价格，更新商品的采购价
没有采购记录的商品，不用进行 采购价 更新

*/

DROP PROCEDURE IF EXISTS calculate_product_price;

CREATE PROCEDURE calculate_product_price()
    BEGIN

        DECLARE v_product_id BIGINT;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DISTINCT pbe.product_id
            FROM erp_purchase_bill_entry pbe
            WHERE pbe.status_id = 4700; # 已入库;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

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
                SELECT avg(pbe.price) AS avg_price
                FROM erp_purchase_bill_entry pbe
                WHERE pbe.product_id = v_product_id AND pbe.status_id = 4700 # 已入库
                LIMIT 6)
            WHERE id = v_product_id;

            SET done = FALSE;
        END LOOP;

        -- 关闭游标
        CLOSE cur;

    END;
