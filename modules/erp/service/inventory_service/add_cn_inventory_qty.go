package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/core"
	"matrix/modules/erp/service/stock_service"
	"github.com/revel/revel"
	"matrix/modules/erp/service/inventory_history_service"
)

// 添加中国仓库的库存
func AddCnInventoryQty(product_id int64, qty_added int64, db_session *xorm.Session, web_session revel.Session) {

	/*
SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = 100 AND nvt.product_id = 1;
	 */

	var sql_get_product_cn_inventory = `SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = %d AND nvt.product_id = %d;`

	sql_get_product_cn_inventory = fmt.Sprintf(sql_get_product_cn_inventory, sppz_models.STOCK_TYPE_CN, product_id)
	var inventory kcgl_models.InventoryInfo
	has, err := db_session.Sql(sql_get_product_cn_inventory).Get(&inventory)
	core.HandleError(err)

	var qty_before = inventory.Qty
	if has == true {
		inventory.Qty = inventory.Qty + qty_added
		_, err = db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
		core.HandleError(err)
	} else {
		var cn_stock_id = stock_service.Get_cn_stock_id(db_session)
		inventory.StockId = cn_stock_id
		inventory.ProductId = product_id
		inventory.Qty = qty_added

		_, err = db_session.Insert(&inventory)
		core.HandleError(err)
	}
	var qty_after = inventory.Qty

	var login_user = core.GetLoginUser(web_session)
	inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("采购单入库更新库存 %d => %d，工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
}

// 添加样品间仓库的库存
func AddSampleInventoryQty(product_id int64, qty_added int64, db_session *xorm.Session, web_session revel.Session) {

	/*
SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = 400 AND nvt.product_id = 1;
	 */

	var sql_get_product_cn_inventory = `SELECT nvt.*
FROM erp_inventory nvt
    INNER JOIN erp_stock s ON nvt.stock_id = s.id
WHERE s.stock_type_id = %d AND nvt.product_id = %d;`

	sql_get_product_cn_inventory = fmt.Sprintf(sql_get_product_cn_inventory, sppz_models.STOCK_TYPE_SAMPLE, product_id)
	var inventory kcgl_models.InventoryInfo
	has, err := db_session.Sql(sql_get_product_cn_inventory).Get(&inventory)
	core.HandleError(err)

	var qty_before = inventory.Qty
	if has == true {
		inventory.Qty = inventory.Qty + qty_added
		_, err = db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
		core.HandleError(err)
	} else {
		var sample_stock_id = stock_service.Get_sample_stock_id(db_session)
		inventory.StockId = sample_stock_id
		inventory.ProductId = product_id
		inventory.Qty = qty_added

		_, err = db_session.Insert(&inventory)
		core.HandleError(err)
	}
	var qty_after = inventory.Qty

	var login_user = core.GetLoginUser(web_session)
	inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("采购单入库更新库存 %d => %d，工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
}
