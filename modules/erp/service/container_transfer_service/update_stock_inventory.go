package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/core"
	"fmt"
	"github.com/revel/revel"
	"matrix/modules/erp/service/inventory_history_service"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
)

func UpdateStockInventory(container_transfer hgbg_models.ContainerTransferInfo, web_session revel.Session, db_session *xorm.Session) {

	var entry_list = make([]hgbg_models.ContainerTransferEntryInfo, 0)
	core.HandleError(db_session.Where("transfer_id = ?", container_transfer.Id).Find(&entry_list))

	zt_stock_id, au_stock_id := get_zt_au_stock_id(container_transfer.Id, db_session)

	// 状态是 已出发的 更新在途仓的库存，扣减中国仓的库存
	if container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_DEPART {
		var cn_stock_id = get_cn_stock_id(db_session)

		for _, entry := range entry_list {
			update_stock_product_inventory(cn_stock_id, entry.ProductId, -entry.Qty, web_session, db_session)
			update_stock_product_inventory(zt_stock_id, entry.ProductId, entry.Qty, web_session, db_session)
		}
	}

	// 状态是 已到仓的 更新目的仓库的库存，扣减在途仓的库存
	if container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_STOCK {
		for _, entry := range entry_list {
			update_stock_product_inventory(zt_stock_id, entry.ProductId, -entry.Qty, web_session, db_session)
			update_stock_product_inventory(au_stock_id, entry.ProductId, entry.Qty, web_session, db_session)
		}
	}
}

func get_cn_stock_id(db_session *xorm.Session) int64 {
	var cn_stock_id_struct struct {
		StockId int64 `xorm:"'stock_id'"`
	}

	/*
		SELECT
		    s.id AS stock_id
		FROM erp_stock s
		WHERE s.stock_type_id = 1
		LIMIT 1;
	 */

	var sql = `
		SELECT
		    s.id AS stock_id
		FROM erp_stock s
		WHERE s.stock_type_id = %d
		LIMIT 1;`

	sql = fmt.Sprintf(sql, sppz_models.STOCK_TYPE_CN)

	_, err := db_session.Sql(sql).Get(&cn_stock_id_struct)
	core.HandleError(err)

	return cn_stock_id_struct.StockId
}

func get_zt_au_stock_id(transfer_id int64, db_session *xorm.Session) (zt_stock_id_ret, au_stock_id_ret int64) {

	var zt_au_stock_id_struct struct {
		ZtStockId int64 `xorm:"'zt_stock_id'"`
		AuStockId int64 `xorm:"'au_stock_id'"`
	}
	/*
	SELECT
		s.id              AS zt_stock_id,
		s.target_stock_id AS au_stock_id
	FROM erp_stock s
	WHERE s.stock_type_id = 1 AND s.dest_port_id IN (
		SELECT
			ct.dest_port_id
		FROM erp_container_transfer ct
		WHERE ct.id = 1
	)
	LIMIT 1;
 */
	var sql = `SELECT
		s.id              AS zt_stock_id,
		s.target_stock_id AS au_stock_id
	FROM erp_stock s
	WHERE s.stock_type_id = %d AND s.dest_port_id IN (
		SELECT
			ct.dest_port_id
		FROM erp_container_transfer ct
		WHERE ct.id = %d
	)
	LIMIT 1;`

	sql = fmt.Sprintf(sql, sppz_models.STOCK_TYPE_ZT, transfer_id)
	_, err := db_session.Sql(sql).Get(&zt_au_stock_id_struct)
	core.HandleError(err)

	return zt_au_stock_id_struct.ZtStockId, zt_au_stock_id_struct.AuStockId
}

func update_stock_product_inventory(stock_id int64, product_id int64, qty_diff int64, web_session revel.Session, db_session *xorm.Session) {
	var inventory kcgl_models.InventoryInfo
	_, err := db_session.Where("stock_id = ? AND product_id = ?", stock_id, product_id).Get(&inventory)
	core.HandleError(err)

	var login_user = core.GetLoginUser(web_session)

	if inventory.Id != 0 {
		var qty_before = inventory.Qty
		inventory.Qty = inventory.Qty + qty_diff
		var qty_after = inventory.Qty
		_, err = db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
		core.HandleError(err)

		inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("集装箱调拨更新库存 %d => %d, 工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
	} else {
		inventory.StockId = stock_id
		inventory.ProductId = product_id
		inventory.StorageLocId = 0 // 暂时设置为零吧
		inventory.Qty = qty_diff

		_, err := db_session.Insert(&inventory)
		core.HandleError(err)

		inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("集装箱调拨新增库存记录 %d => %d, 工号：%s", 0, qty_diff, login_user.LoginName), db_session)
	}
}
