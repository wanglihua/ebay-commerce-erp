package product_inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models"
	"matrix/core"
)

type Product_inventory_view_item_t struct {
	models.InventoryInfo   `xorm:"extends" json:"nvt"`
	models.ProductInfo     `xorm:"extends" json:"p"`
	models.StockInfo       `xorm:"extends" json:"st"`
	models.StorageLocInfo  `xorm:"extends" json:"sl"`
}

// 这是一个静态变量，对他的修改会被保留下来
var product_inventory_sql = `SELECT nvt.*, p.*, st.*, sl.*
FROM erp_inventory nvt
    INNER JOIN erp_product p ON p.id = nvt.product_id
    INNER JOIN erp_stock st ON st.id = nvt.stock_id
    INNER JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id
WHERE nvt.product_id = %d`

func Get_product_inventory_list(db_session *xorm.Session, product_id int64) []Product_inventory_view_item_t {

	sql := fmt.Sprintf(product_inventory_sql, product_id)

	var product_inventory_list = make([]Product_inventory_view_item_t, 0)
	err := db_session.Sql(sql).Find(&product_inventory_list)
	core.HandleError(err)

	return product_inventory_list
}
