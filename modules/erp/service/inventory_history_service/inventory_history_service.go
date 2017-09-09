package inventory_history_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
)

func Insert_inventory_history(inventory kcgl_models.InventoryInfo, description string, db_session *xorm.Session) {

	var inventory_history kcgl_models.InventoryHistoryInfo

	inventory_history.StockId = inventory.StockId
	inventory_history.StorageLocId = inventory.StorageLocId
	inventory_history.ProductId = inventory.ProductId
	inventory_history.Qty = inventory.Qty
	inventory_history.Price = inventory.CostPrice
	inventory_history.Description = description

	_, err := db_session.Insert(&inventory_history)
	core.HandleError(err)
}
