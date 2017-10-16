DROP PROCEDURE IF EXISTS product_sale_day_analysis;

CREATE PROCEDURE product_sale_day_analysis()
    BEGIN

        CALL product_sale_day_analysis_actual_without_stock();
        CALL product_sale_day_analysis_actual_with_stock();
        CALL product_sale_day_analysis_logic_with_stock();

    END

