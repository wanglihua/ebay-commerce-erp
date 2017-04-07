package inventory_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Product_inventory_view_item_t struct {
	models.InventoryInfo   `xorm:"extends" json:"nvt"`
	models.ProductInfo     `xorm:"extends" json:"p"`
	models.StockInfo       `xorm:"extends" json:"st"`
	models.StorageLocInfo  `xorm:"extends" json:"sl"`
}

func get_inventory_list_sql_builder() *core.SqlBuilder {
	sb := core.NewSqlBuilder()
	sb.Select("nvt.*, p.*, st.*, sl.*")
	sb.From(`erp_inventory nvt
    INNER JOIN erp_product p ON p.id = nvt.product_id
    INNER JOIN erp_stock st ON st.id = nvt.stock_id
    INNER JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id`)

	return sb
}

func Get_product_inventory_list(db_session *xorm.Session, product_id int64) []Product_inventory_view_item_t {
	sb := get_inventory_list_sql_builder()
	sb.Where(fmt.Sprintf("nvt.product_id = %d", product_id))

	var product_inventory_list = make([]Product_inventory_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&product_inventory_list)
	core.HandleError(err)

	return product_inventory_list
}

func Get_inventory_by_product_list(db_session *xorm.Session, product_id_list []int64) []Product_inventory_view_item_t {
	sb := get_inventory_list_sql_builder()

	if len(product_id_list) == 0 {
		sb.Where("1=0")
	} else {
		var product_id_str_list = make([]string, 0, len(product_id_list))
		for _, product_id := range product_id_list {
			product_id_str_list = append(product_id_str_list, fmt.Sprint(product_id))
		}
		sb.Where(fmt.Sprintf("nvt.product_id IN (%s)", strings.Join(product_id_str_list, ",")))
	}

	sb.OrderBy("st.id ASC, sl.id ASC")

	var product_inventory_list = make([]Product_inventory_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&product_inventory_list)
	core.HandleError(err)

	return product_inventory_list
}

type Product_stock_qty struct {
	ProductId   int64     `xorm:"'product_id'"`
	StockId     int64     `xorm:"'stock_id'"`
	StockTypeId int64     `xorm:"'stock_type_id'"`
	Qty         int64     `xorm:"'qty'"`
}

func Get_product_stock_qty_list(product_id int64, db_session *xorm.Session) []Product_stock_qty {
	/*
SELECT ivt.product_id AS product_id, ivt.stock_id AS stock_id, s.stock_type_id, sum(ivt.qty) AS qty
FROM erp_inventory ivt
    INNER JOIN erp_stock s ON s.id = ivt.stock_id
WHERE ivt.product_id = 1
GROUP BY ivt.product_id, ivt.stock_id
	 */

	sql := `SELECT ivt.product_id AS product_id, ivt.stock_id AS stock_id, s.stock_type_id, sum(ivt.qty) AS qty
FROM erp_inventory ivt
    INNER JOIN erp_stock s ON s.id = ivt.stock_id
WHERE ivt.product_id = %d
GROUP BY ivt.product_id, ivt.stock_id`
	sql = fmt.Sprintf(sql, product_id)

	var product_stock_qty_list = make([]Product_stock_qty, 0)
	err := db_session.Sql(sql).Find(&product_stock_qty_list)
	core.HandleError(err)

	return product_stock_qty_list
}
