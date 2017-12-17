package inventory_service

import (
	"fmt"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
)

type Product_inventory_view_item_t struct {
	Inventory  kcgl_models.InventoryInfo  `xorm:"extends" json:"nvt"`
	Product    kcgl_models.ProductInfo    `xorm:"extends" json:"p"`
	Stock      sppz_models.StockInfo      `xorm:"extends" json:"st"`
	StorageLoc sppz_models.StorageLocInfo `xorm:"extends" json:"sl"`
	TransferId int64                      `json:"transfer_id"` // 集装箱调拨 id，订单界面显示在途仓集装箱调拨的时候用
	BatchCode  string                     `json:"batch_code"`  // 集装箱调拨批次，订单界面显示在途仓集装箱调拨的时候用
	DeptTime   string                     `json:"dept_time"`   // 在途仓的到港时间，订单界面显示在途仓到港时间的时候用
	StockTime  string                     `json:"stock_time"`  // 在途仓的到仓时间，订单界面显示在途仓到仓时间的时候用
}

func get_inventory_list_sql_builder() *core.SqlBuilder {
	sb := core.NewSqlBuilder()
	sb.Select("nvt.*, p.*, st.*, sl.*")
	sb.From(`erp_inventory nvt
    INNER JOIN erp_product p ON p.id = nvt.product_id
    INNER JOIN erp_stock st ON st.id = nvt.stock_id
    LEFT JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id`)

	return sb
}

func Get_product_inventory_list(db_session *xorm.Session, product_id int64) []Product_inventory_view_item_t {
	sb := get_inventory_list_sql_builder()
	sb.Where(fmt.Sprintf("nvt.product_id = %d", product_id))
	sb.OrderBy("st.stock_type_id DESC, st.name ASC")

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
	ProductId   int64   `xorm:"'product_id'"`
	StockId     int64   `xorm:"'stock_id'"`
	StockTypeId int64   `xorm:"'stock_type_id'"`
	Qty         int64   `xorm:"'qty'"`
	PendingQty  int64   `xorm:"'pending_qty'"`
	CostPrice   float64 `xorm:"'cost_price'"` //这个是RMB，但匹配是用的应该是澳币
}

// 获取包含"这个商品"的所有"产品库存数量记录"
func Get_product_stock_qty_list(product_id int64, db_session *xorm.Session) []Product_stock_qty {
	/*
SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id, sum(nvt.qty) AS qty, sum(nvt.pending_qty) as pending_qty, nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.product_id = 1
GROUP BY nvt.product_id, nvt.stock_id
	 */

	sql := `SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id, sum(nvt.qty) AS qty, sum(nvt.pending_qty) as pending_qty, nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.product_id = %d
GROUP BY nvt.product_id, nvt.stock_id`
	sql = fmt.Sprintf(sql, product_id)

	var product_stock_qty_list = make([]Product_stock_qty, 0)
	err := db_session.Sql(sql).Find(&product_stock_qty_list)
	core.HandleError(err)

	return product_stock_qty_list
}

// 获取包含"这个商品列表"的所有"产品库存数量记录"
func Get_products_stock_qty_list(product_id_list []int64, db_session *xorm.Session) []Product_stock_qty {
	/*
SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id AS stock_type_id,
    sum(nvt.qty) AS qty, sum(nvt.pending_qty) as pending_qty,
    if(sum(nvt.qty) <> 0,sum(nvt.cost_price * nvt.qty) / sum(nvt.qty), sum(nvt.cost_price)/count(*)) AS cost_price
FROM erp_inventory nvt
INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.product_id IN (10000001)
GROUP BY nvt.product_id, nvt.stock_id, s.stock_type_id
	 */

	sql := `SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id AS stock_type_id,
	sum(nvt.qty) AS qty, sum(nvt.pending_qty) as pending_qty,
    if(sum(nvt.qty) <> 0, sum(nvt.cost_price * nvt.qty) / sum(nvt.qty), sum(nvt.cost_price) / count(*)) AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.product_id IN (%s)
GROUP BY nvt.product_id, nvt.stock_id, s.stock_type_id`
	if len(product_id_list) != 0 {
		var product_id_str_list = core.Int_list_to_str_list(product_id_list)
		sql = fmt.Sprintf(sql, strings.Join(product_id_str_list, ","))
	} else {
		sql = fmt.Sprintf(sql, "0") // 以免sql出错
	}

	var product_stock_qty_list = make([]Product_stock_qty, 0)
	err := db_session.Sql(sql).Find(&product_stock_qty_list)
	core.HandleError(err)

	return product_stock_qty_list
}
