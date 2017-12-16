# CALL  create_sale_report_time();

DROP PROCEDURE IF EXISTS create_sale_report_time;

CREATE PROCEDURE create_sale_report_time()
    BEGIN
        DECLARE v_time_begin DATETIME;
        DECLARE v_time_end DATETIME;
        DECLARE v_report_time DATETIME;
        DECLARE v_report_time_count INT;

        SET v_time_begin = '2016-01-01';
        SET v_time_end = '2050-12-31';

        # 时间 天 的生成
        SET v_report_time = v_time_begin;

        WHILE v_report_time < v_time_end DO
            SELECT count(id) INTO v_report_time_count
            FROM erp_sale_report_time_day
            WHERE report_day = v_report_time;

            IF v_report_time_count = 0
            THEN
                # 时间取  report_day =< time < date_add(report_day, INTERVAL 1 DAY)
                INSERT INTO erp_sale_report_time_day (report_day) VALUES (v_report_time);
            END IF;

            SET v_report_time = date_add(v_report_time, INTERVAL 1 DAY);
        END WHILE;

        # 时间 月 的生成
        SET v_report_time = v_time_begin;

        WHILE v_report_time < v_time_end DO
            SELECT count(id) INTO v_report_time_count
            FROM erp_sale_report_time_month
            WHERE report_month = v_report_time;

            IF v_report_time_count = 0
            THEN
                # 时间取  report_month =< time < date_add(report_month, INTERVAL 1 MONTH)
                INSERT INTO erp_sale_report_time_month (report_month) VALUES (v_report_time);
            END IF;

            SET v_report_time = date_add(v_report_time, INTERVAL 1 MONTH);
        END WHILE;
    END

