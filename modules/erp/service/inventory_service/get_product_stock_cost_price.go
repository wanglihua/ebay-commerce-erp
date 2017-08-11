package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func GetProductStockCostPrice(stock_id int64, product_id int64, db_session *xorm.Session) float64 {
	var cost_price struct {
		CostPrice float64 `xorm:"'cost_price'"`
	}

	/*
SELECT nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.stock_id = 1 AND nvt.product_id = 1;
	 */

	var sql = `SELECT nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.stock_id = %d AND nvt.product_id = %d;`

	sql = fmt.Sprintf(sql, stock_id, product_id)
	_, err := db_session.Sql(sql).Get(&cost_price)
	core.HandleError(err)

	return cost_price.CostPrice
}
