package inventory_service

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/sea_container_cost_service"
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
	ProductId   int64       `xorm:"'product_id'"`
	StockId     int64       `xorm:"'stock_id'"`
	StockTypeId int64       `xorm:"'stock_type_id'"`
	Qty         int64       `xorm:"'qty'"`
	CostPrice   float64     `xorm:"'cost_price'"` //这个是RMB，但匹配是用的应该是澳币
}

// 获取包含"这个商品"的所有"产品库存数量记录"
func Get_product_stock_qty_list(product_id int64, db_session *xorm.Session) []Product_stock_qty {
	/*
SELECT ivt.product_id AS product_id, ivt.stock_id AS stock_id, s.stock_type_id, sum(ivt.qty) AS qty, ivt.cost_price AS cost_price
FROM erp_inventory ivt
    INNER JOIN erp_stock s ON s.id = ivt.stock_id
WHERE ivt.product_id = 1
GROUP BY ivt.product_id, ivt.stock_id
	 */

	sql := `SELECT ivt.product_id AS product_id, ivt.stock_id AS stock_id, s.stock_type_id, sum(ivt.qty) AS qty, ivt.cost_price AS cost_price
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

// 获取包含"这个商品列表"的所有"产品库存数量记录"
func Get_products_stock_qty_list(product_id_list []int64, db_session *xorm.Session) []Product_stock_qty {
	/*
SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id AS stock_type_id, sum(nvt.qty) AS qty, sum(nvt.cost_price * nvt.qty) / sum(nvt.qty) AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE nvt.product_id IN (10000001)
GROUP BY nvt.product_id, nvt.stock_id, s.stock_type_id
	 */

	sql := `SELECT nvt.product_id AS product_id, nvt.stock_id AS stock_id, s.stock_type_id AS stock_type_id, sum(nvt.qty) AS qty, sum(nvt.cost_price * nvt.qty) / sum(nvt.qty) AS cost_price
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

func Update_au_stock_cost() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("更新澳仓库存出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url
	db_driver, _ := revel.Config.String("db.driver")
	db_url, _ := revel.Config.String("db.url")
	//只有确实要操作数据库的时候，才去初始化数据库引擎

	db_engine, err := xorm.NewEngine(db_driver, db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	/*
SELECT nvt.* FROM erp_inventory nvt
INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE s.stock_type_id = 0;

UPDATE erp_inventory SET cost_price = 0 WHERE product_id = 0 AND (SELECT s.stock_type_id FROM erp_stock s WHERE s.id = erp_inventory.stock_id) = 0;
 */

	var inventory_cn_list = make([]models.InventoryInfo, 0)

	var get_cn_inventory_list = `SELECT nvt.* FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE s.stock_type_id = %d;`
	get_cn_inventory_list = fmt.Sprintf(get_cn_inventory_list, models.STOCK_TYPE_CN)

	err = db_session.Sql(get_cn_inventory_list).Find(&inventory_cn_list)
	core.HandleError(err)

	// 加载海运成本设置
	// 不用汇率换算了，仓库的成本都是RMB的
	var sea_container_cost = sea_container_cost_service.Get_sea_container_cost(db_session)

	// 汇率换算
	//var aud_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)

	for _, cn_inventory := range inventory_cn_list {
		// 获取中国仓商品成本
		var cn_cost_price = cn_inventory.CostPrice

		var product models.ProductInfo
		_, err = db_session.Id(cn_inventory.ProductId).Get(&product)
		core.HandleError(err)

		// 计算商品体积，计算平摊海运成本
		var avg_sea_shipping_cost = float64(sea_container_cost.Cost) * product.Length * product.Width * product.Height / float64(sea_container_cost.Volume)

		var update_au_inventory_sql = `UPDATE erp_inventory SET cost_price = %f WHERE product_id = %d AND (SELECT s.stock_type_id FROM erp_stock s WHERE s.id = erp_inventory.stock_id) = %d;`
		update_au_inventory_sql = fmt.Sprintf(update_au_inventory_sql, cn_cost_price+avg_sea_shipping_cost, cn_inventory.StockId, models.STOCK_TYPE_AU)
		_, err = db_session.Exec(update_au_inventory_sql)
		core.HandleError(err)
	}
}
