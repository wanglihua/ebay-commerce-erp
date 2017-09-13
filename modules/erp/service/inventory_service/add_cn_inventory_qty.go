package inventory_service

import "github.com/go-xorm/xorm"

// 添加中国仓库的库存
func AddCnInventoryQty(product_id int64, qty_added int64, db_session *xorm.Session) {

	/*
SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = 100 AND nvt.product_id = 1;
	 */

}
