package tasks

import "matrix/modules/erp/service/inventory_service"

type UpdateInventoryPendingQtyTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t UpdateInventoryPendingQtyTask) Name() string {
	return "更新仓库库存待发货数量"
}

func (t UpdateInventoryPendingQtyTask) Run() {
	inventory_service.Calculate_inventory_pending_qty()
}
