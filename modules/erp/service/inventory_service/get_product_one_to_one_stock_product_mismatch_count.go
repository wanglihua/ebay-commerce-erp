package inventory_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func GetProductOneToOneStockProductMismatchCount(db_session *xorm.Session) int64 {

	/*
SELECT count(*) AS count
FROM (
         SELECT stock_id, product_id
         FROM erp_inventory
         GROUP BY stock_id, product_id
         HAVING count(*) > 1
     ) stock_product_repeat;
	 */

	var sql = `SELECT count(*) AS count
FROM (
         SELECT stock_id, product_id
         FROM erp_inventory
         GROUP BY stock_id, product_id
         HAVING count(*) > 1
     ) stock_product_repeat;`

	var count core.Count
	_, err := db_session.Sql(sql).Get(&count)
	core.HandleError(err)

	return count.Count
}
