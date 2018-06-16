/*

SELECT DISTINCT product_id
FROM (
         SELECT nvt.product_id FROM erp_inventory nvt
         WHERE nvt.stock_id = 1
         UNION
         SELECT pbe.product_id FROM erp_purchase_bill_entry pbe
         WHERE pbe.status_id = 1
     ) p


 */

DROP PROCEDURE IF EXISTS calculate_stock_cn_zt_qty;

CREATE PROCEDURE calculate_stock_cn_zt_qty()
    BEGIN
        DECLARE v_product_id BIGINT;
        DECLARE v_cn_zt_stock_id BIGINT;

        SELECT s.id INTO v_cn_zt_stock_id
        FROM erp_stock s
        WHERE s.stock_type_id = 80 # 中国仓在途
        LIMIT 1;

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR
            SELECT DISTINCT product_id
            FROM (
                     SELECT nvt.product_id FROM erp_inventory nvt
                     WHERE nvt.stock_id = v_cn_zt_stock_id
                     UNION
                     SELECT pbe.product_id FROM erp_purchase_bill_entry pbe
                     WHERE pbe.status_id = 4600 # 采购单状态 采购中
                 ) p;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    END;