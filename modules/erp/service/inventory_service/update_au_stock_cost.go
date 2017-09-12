package inventory_service

import (
	"github.com/go-xorm/xorm"
	"time"
	"fmt"
	"matrix/modules/erp/service/sea_container_cost_service"
	"github.com/revel/revel"

	"runtime/debug"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/service/inventory_history_service"
	"math"
)

func Update_au_stock_cost() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("更新澳仓成本出现错误： %s\n", err)
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
	db_session.IsAutoCommit = true

	defer func() {
		db_session.Close()
	}()

	// db_session.Begin()

	/*
SELECT nvt.* FROM erp_inventory nvt
INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE s.stock_type_id = 0;

UPDATE erp_inventory SET cost_price = 0 WHERE product_id = 0 AND (SELECT s.stock_type_id FROM erp_stock s WHERE s.id = erp_inventory.stock_id) = 0;
*/

	var inventory_cn_list = make([]kcgl_models.InventoryInfo, 0)

	var get_cn_inventory_list = `SELECT nvt.* FROM erp_inventory nvt
    INNER JOIN erp_stock s ON s.id = nvt.stock_id
WHERE s.stock_type_id = %d;`
	get_cn_inventory_list = fmt.Sprintf(get_cn_inventory_list, sppz_models.STOCK_TYPE_CN)

	err = db_session.Sql(get_cn_inventory_list).Find(&inventory_cn_list)
	core.HandleError(err)

	// db_session.Commit()

	// 澳仓自动计算库存成本
	{
		db_session.Begin()

		// 加载海运成本设置
		// 不用汇率换算了，仓库的成本都是RMB的
		var sea_container_cost = sea_container_cost_service.Get_sea_container_cost(db_session)

		var au_stock_list = make([]sppz_models.StockInfo, 0)
		core.HandleError(db_session.Where("stock_type_id = ?", sppz_models.STOCK_TYPE_AU).Find(&au_stock_list))

		// 汇率换算
		//var aud_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)

		// db_session.Commit()

		for _, cn_inventory := range inventory_cn_list {
			// 获取中国仓商品成本
			var cn_cost_price = cn_inventory.CostPrice

			var product kcgl_models.ProductInfo
			_, err = db_session.Id(cn_inventory.ProductId).Get(&product)
			core.HandleError(err)

			db_session.Begin()
			// 如果澳仓库存不存的的，需要自动添加零库存的记录，方便记录澳仓成本。
			for _, au_stock := range au_stock_list {
				var au_stock_inventory kcgl_models.InventoryInfo
				has, err := db_session.Where("product_id = ? AND stock_id = ?", cn_inventory.ProductId, au_stock.Id).Get(&au_stock_inventory)
				core.HandleError(err)

				if has == false {
					au_stock_inventory.ProductId = cn_inventory.ProductId
					au_stock_inventory.StockId = au_stock.Id
					au_stock_inventory.StorageLocId = 0
					au_stock_inventory.Qty = 0
					au_stock_inventory.PendingQty = 0
					au_stock_inventory.CostPrice = 0

					_, err = db_session.Insert(&au_stock_inventory)
					core.HandleError(err)
				}
			}
			// db_session.Commit()

			db_session.Begin()
			// 计算商品体积，计算平摊海运成本
			var avg_sea_shipping_cost = float64(sea_container_cost.Cost) * product.Length * product.Width * product.Height / float64(sea_container_cost.Volume)

			/*
			var update_au_inventory_sql = `UPDATE erp_inventory SET cost_price = %f WHERE product_id = %d AND (SELECT s.stock_type_id FROM erp_stock s WHERE s.id = erp_inventory.stock_id) = %d;`
			update_au_inventory_sql = fmt.Sprintf(update_au_inventory_sql, cn_cost_price+avg_sea_shipping_cost, cn_inventory.ProductId, sppz_models.STOCK_TYPE_AU)
			_, err = db_session.Exec(update_au_inventory_sql)
			core.HandleError(err)
			*/

			/*
SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = 300 AND nvt.product_id = 2;
			 */
			var sql_get_product_au_inventory = `SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = %d AND nvt.product_id = %d;`

			sql_get_product_au_inventory = fmt.Sprintf(sql_get_product_au_inventory, sppz_models.STOCK_TYPE_AU, product.Id)
			var au_inventory_list = make([]kcgl_models.InventoryInfo, 0)
			core.HandleError(db_session.Sql(sql_get_product_au_inventory).Find(&au_inventory_list))

			// db_session.Commit()

			for _, au_inventory := range au_inventory_list {
				var price_before = au_inventory.CostPrice
				au_inventory.CostPrice = cn_cost_price + avg_sea_shipping_cost
				var price_after = au_inventory.CostPrice

				if math.Abs(price_after-price_before) > core.FLOAT_MIN_DIFF {
					db_session.Begin()
					_, err := db_session.Id(au_inventory.Id).Cols("cost_price").Update(&au_inventory)
					core.HandleError(err) // todo: Error 1213: Deadlock found when trying to get lock; try restarting transaction

					/*
					if err != nil {
						revel.ERROR.Printf("更新澳仓成本出现错误： %s\n", err)
						revel.ERROR.Println(string(debug.Stack()))
					}
					*/

					if err == nil {
						inventory_history_service.Insert_inventory_history(au_inventory, fmt.Sprintf("自动更新澳仓成本。 %0.2f => %0.2f", price_before, price_after), db_session)
					}
					// db_session.Commit()
				}
			}
		}
	}

	db_session.Begin()

	// 自动添加不存在的对应的在途仓记录 {
	var zt_stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where("stock_type_id = ?", sppz_models.STOCK_TYPE_ZT).Find(&zt_stock_list))

	// db_session.Commit()

	// 汇率换算
	//var aud_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)

	db_session.Begin()
	for _, cn_inventory := range inventory_cn_list {

		// 如果在途仓库存不存的的，需要自动添加零库存的记录。
		for _, zt_stock := range zt_stock_list {
			var zt_stock_inventory kcgl_models.InventoryInfo
			has, err := db_session.Where("product_id = ? AND stock_id = ?", cn_inventory.ProductId, zt_stock.Id).Get(&zt_stock_inventory)
			core.HandleError(err)

			if has == false {
				zt_stock_inventory.ProductId = cn_inventory.ProductId
				zt_stock_inventory.StockId = zt_stock.Id
				zt_stock_inventory.StorageLocId = 0
				zt_stock_inventory.Qty = 0
				zt_stock_inventory.PendingQty = 0
				zt_stock_inventory.CostPrice = 0

				_, err = db_session.Insert(&zt_stock_inventory)
				core.HandleError(err)
			}
		}
	}
	// db_session.Commit()
}
