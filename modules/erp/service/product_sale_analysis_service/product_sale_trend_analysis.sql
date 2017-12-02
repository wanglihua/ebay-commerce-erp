DROP PROCEDURE IF EXISTS product_sale_trend_analysis;

CREATE PROCEDURE product_sale_trend_analysis()
    BEGIN


    END;

SELECT DISTINCT product_id FROM erp_product_sale_day_analysis;

SELECT * FROM erp_product_sale_day_analysis;

CREATE INDEX ix_product_sale_day_analysis_type ON erp_product_sale_day_analysis (type);