package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
	"matrix/modules/erp/models/sppz_models"
)

func GetProductCnStockCostPrice(product_id int64, db_session *xorm.Session) float64 {
	var cn_cost_price struct {
		CostPrice float64 `xorm:"'cost_price'"`
	}

	/*
	SELECT
		nvt.cost_price AS cost_price
	FROM erp_inventory nvt
		INNER JOIN erp_stock s ON s.id = nvt.stock_id
	WHERE s.stock_type_id = 100 AND nvt.product_id = 1;
	 */
	var sql = `SELECT
    nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE s.stock_type_id = %d AND nvt.product_id = %d;
	`
	sql = fmt.Sprintf(sql, sppz_models.STOCK_TYPE_CN, product_id)
	_, err := db_session.Sql(sql).Get(&cn_cost_price)
	core.HandleError(err)

	return cn_cost_price.CostPrice
}
