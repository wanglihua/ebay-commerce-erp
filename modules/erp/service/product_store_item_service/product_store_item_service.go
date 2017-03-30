package product_store_item_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

// 检查 product_item_id 的对应的 ProductStoreItemInfo 是不是存在，如果存在就返回对应的 product_id 和 product_store_item_id
// product_item_id 就是 ebay 系统里面的 item_id
func Query_product_store_item(product_item_id string, db_session *xorm.Session) (has_item_return bool, product_id_return int64, product_store_item_id_return int64) {
	var product_store_item_query_result struct {
		ProductId          int64 `xorm:"'product_id'"`
		ProductStoreItemId int64 `xorm:"'product_store_item_id'"`
	}

	has_product_store_item, err := db_session.Sql(fmt.Sprintf("SELECT product_id, id AS product_store_item_id FROM erp_product_store_item psi WHERE item_id = '%s';", product_item_id)).Get(&product_store_item_query_result)
	core.HandleError(err)

	// return
	has_item_return = has_product_store_item
	product_id_return = product_store_item_query_result.ProductId
	product_store_item_id_return = product_store_item_query_result.ProductStoreItemId
	return
}
