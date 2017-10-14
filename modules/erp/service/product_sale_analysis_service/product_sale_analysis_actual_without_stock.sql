CREATE PROCEDURE product_sale_analysis_actual_without_stock()
    BEGIN
        -- 需要定义接收游标数据的变量
        DECLARE a VARCHAR(255);

        -- 遍历数据结束标志
        DECLARE done INT DEFAULT FALSE;

        -- 游标
        DECLARE cur CURSOR FOR SELECT login_name
                               FROM auth_user;

        -- 将结束标志绑定到游标
        DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

        -- 打开游标
        OPEN cur;

        -- 开始循环
        read_loop: LOOP
            -- 提取游标里的数据，这里只有一个，多个的话也一样；
            FETCH cur
            INTO a;
            -- 声明结束的时候
            IF done
            THEN
                LEAVE read_loop;
            END IF;
            -- 这里做你想做的循环的事件

            # INSERT INTO test.t VALUES (a);

        END LOOP;
        -- 关闭游标
        CLOSE cur;

    END

