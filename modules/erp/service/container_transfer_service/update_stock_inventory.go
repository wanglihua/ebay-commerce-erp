package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/core"
	"fmt"
)

func UpdateStockInventory(container_transfer models.ContainerTransferInfo, db_session *xorm.Session) {

	var entry_list = make([]models.ContainerTransferEntryInfo, 0)
	core.HandleError(db_session.Where("transfer_id = ?", container_transfer.Id).Find(&entry_list))

	zt_stock_id, au_stock_id := get_zt_au_stock_id(container_transfer.Id, db_session)

	_, _ = zt_stock_id, au_stock_id

	// 状态是 已出发的 更新在途仓的库存，扣减中国仓的库存
	if container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_DEPART {
		var cn_stock_id = get_cn_stock_id(db_session)
		_ = cn_stock_id

		for _, entry := range entry_list {
			_ = entry
		}
	}

	// 状态是 已到仓的 更新目的仓库的库存，扣减在途仓的库存
	if container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_STOCK {

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

	sql = fmt.Sprintf(sql, models.STOCK_TYPE_CN)

	_, err := db_session.Sql(sql).Get(&cn_stock_id_struct)
	core.HandleError(err)

	return cn_stock_id_struct.StockId
}

func get_zt_au_stock_id(transfer_id int64, db_session *xorm.Session) (zt_stock_id_ret, au_stock_id_ret int64) {

	var zt_au_stock_id_struct struct {
		ZtStockId int64  `xorm:"'zt_stock_id'"`
		AuStockId int64  `xorm:"'au_stock_id'"`
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

	sql = fmt.Sprintf(sql, models.STOCK_TYPE_ZT, transfer_id)
	_, err := db_session.Sql(sql).Get(&zt_au_stock_id_struct)
	core.HandleError(err)

	return zt_au_stock_id_struct.ZtStockId, zt_au_stock_id_struct.AuStockId
}
