package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
)

func UpdateStockInventory(container_transfer models.ContainerTransferInfo, db_session *xorm.Session) {

	// 状态是 已出发的 更新在途仓的库存，扣减中国仓的库存

	// 状态是 已到仓的 更新目的仓库的库存，扣减在途仓的库存

}
