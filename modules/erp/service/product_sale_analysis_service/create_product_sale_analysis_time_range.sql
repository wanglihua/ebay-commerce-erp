DROP PROCEDURE IF EXISTS create_product_sale_analysis_time_range;

# 10 天一个间隔
# TIME_RANGE_TYPE_TEN_DAYS int64 = 100

CREATE PROCEDURE create_product_sale_analysis_time_range()
    BEGIN
        DECLARE v_time_begin DATETIME;
        DECLARE v_time_end DATETIME;
        DECLARE v_time_range_begin DATETIME;
        DECLARE v_time_range_end DATETIME;
        DECLARE v_range_count INT;

        SET v_time_begin = '2016-01-01';
        SET v_time_end = '2050-12-31';

        WHILE v_time_range_begin < v_time_end DO
            SET v_time_range_end = date_add(v_time_range_begin, INTERVAL 10 DAY);

            SELECT count(id) INTO v_range_count
            FROM erp_product_sale_analysis_time_range
            WHERE type_id = 100 AND begin_time = v_time_range_begin;

            IF v_range_count = 0
            THEN
                INSERT INTO erp_product_sale_analysis_time_range (type_id, begin_time, end_time) VALUES (100, v_time_range_begin, v_time_range_end);
            END IF;

            SET v_time_range_begin = date_add(v_time_range_begin, INTERVAL 10 DAY);
        END WHILE;
    END

