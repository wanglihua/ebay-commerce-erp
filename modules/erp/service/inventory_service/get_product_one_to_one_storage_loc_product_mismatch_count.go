package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models/sppz_models"
	"matrix/core"
)

func GetProductOneToOneStorageLocProductMismatchCount(db_session *xorm.Session) int64 {
	/*
SELECT count(*) AS count
FROM (
         SELECT nvt.storage_loc_id
         FROM erp_inventory nvt
             INNER JOIN erp_storage_loc sl ON nvt.storage_loc_id = sl.id
             INNER JOIN erp_stock s ON sl.stock_id = s.id
         WHERE s.stock_type_id = 300 AND nvt.storage_loc_id <> 0
         GROUP BY nvt.storage_loc_id
         HAVING count(nvt.product_id) > 1
     ) storage_loc_product_repeat;
	 */

	var sql = `SELECT count(*) AS count
FROM (
         SELECT nvt.storage_loc_id
         FROM erp_inventory nvt
             INNER JOIN erp_storage_loc sl ON nvt.storage_loc_id = sl.id
             INNER JOIN erp_stock s ON sl.stock_id = s.id
         WHERE s.stock_type_id = %d AND nvt.storage_loc_id <> 0
         GROUP BY nvt.storage_loc_id
         HAVING count(nvt.product_id) > 1
     ) storage_loc_product_repeat;`
	sql = fmt.Sprintf(sql, sppz_models.STOCK_TYPE_AU)

	var count core.Count
	_, err := db_session.Sql(sql).Get(&count)
	core.HandleError(err)

	return count.Count
}
