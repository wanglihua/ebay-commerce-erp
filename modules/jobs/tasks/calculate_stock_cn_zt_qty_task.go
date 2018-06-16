package tasks

import (
	"matrix/modules/erp/service/purchase_bill_service"
)

type CalculateStockCnZtQtyTask struct {
}

func (t CalculateStockCnZtQtyTask) Name() string {
	return "更新中国在途仓库存数量"
}

func (t CalculateStockCnZtQtyTask) Run() {
	purchase_bill_service.Calculate_stock_cn_zt_qty()
}
