package tasks

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/sea_container_cost_service"
)

type UpdateAuInventoryCostTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t UpdateAuInventoryCostTask) Name() string {
	return "更新澳仓成本"
}

func (t UpdateAuInventoryCostTask) Run() {

	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("更新澳仓库存出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	// get db driver and db url

	if t.db_driver == "" {
		t.db_driver, _ = revel.Config.String("db.driver")
	}

	if t.db_url == "" {
		t.db_url, _ = revel.Config.String("db.url")
	}

	//只有确实要操作数据库的时候，才去初始化数据库引擎

	db_engine, err := xorm.NewEngine(t.db_driver, t.db_url)
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
