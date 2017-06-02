package tasks

import "matrix/modules/erp/service/inventory_service"

type UpdateAuInventoryCostTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t UpdateAuInventoryCostTask) Name() string {
	return "更新澳仓成本"
}

func (t UpdateAuInventoryCostTask) Run() {

	inventory_service.Update_au_stock_cost()
}
