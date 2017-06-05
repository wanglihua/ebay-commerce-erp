package inventory_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

// 扣减订单库存
func Take_order_inventory(order_id int64, db_session *xorm.Session) {

	/*
	SELECT oi.item_product_id AS product_id, oi.quantity AS quantity FROM erp_order_item oi WHERE oi.order_id = 0;
	*/
	type OrderProductQuantityItem struct {
		ProductId int64 `xorm:"'product_id'" json:"product_id"`
		Quantity  int64 `xorm:"'quantity'" json:"quantity"`
	}

	var order_product_quantity_item_list = make([]OrderProductQuantityItem, 0)
	sql_get_order_product_item := `SELECT oi.item_product_id AS product_id, oi.quantity AS quantity FROM erp_order_item oi WHERE oi.order_id = %d;`
	sql_get_order_product_item = fmt.Sprintf(sql_get_order_product_item, order_id)
	err := db_session.Sql(sql_get_order_product_item).Find(&order_product_quantity_item_list)
	core.HandleError(err)

	// 获取订单配货仓库
	/*
		SELECT ol.stock_id AS stock_id FROM erp_order_logistics ol WHERE ol.order_id = 0;
	 */
	type OrderStockId struct {
		StockId int64 `xorm:"'stock_id'" json:"stock_id"`
	}

	var order_stock_id OrderStockId
	sql_get_order_stock := `SELECT ol.stock_id AS stock_id FROM erp_order_logistics ol WHERE ol.order_id = %d;`
	sql_get_order_stock = fmt.Sprintf(sql_get_order_stock, order_id)
	_, err = db_session.Sql(sql_get_order_stock).Get(&order_stock_id)
	core.HandleError(err)

	var stock_id = order_stock_id.StockId
	if stock_id == 0 {
		//订单没有明确的发货仓库，不能扣减库存
		panic("订单没有明确的发货仓库，不能扣减库存！")
	}

	// 分别扣减库存
	for _, order_product_quantity_item := range order_product_quantity_item_list {
		var product_id = order_product_quantity_item.ProductId
		var quantity = order_product_quantity_item.Quantity

		/*
			SELECT * FROM erp_inventory WHERE stock_id = 0 AND product_id = 0;
		 */
		var sql_get_inventory_record = `SELECT * FROM erp_inventory WHERE stock_id = %d AND product_id = %d;`
		sql_get_inventory_record = fmt.Sprintf(sql_get_inventory_record, stock_id, product_id)

		var inventory models.InventoryInfo
		has, err := db_session.Sql(sql_get_inventory_record).Get(&inventory)
		core.HandleError(err)

		if has && inventory.Qty >= quantity {
			inventory.Qty = inventory.Qty - quantity
			_, err = db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
			core.HandleError(err)
		} else {
			panic("仓库中不存在相应库存，无法扣减相应库存！")
		}
	}
}
