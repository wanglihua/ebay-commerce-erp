package tasks

import (
	"matrix/modules/erp/service/sale_report_service"
)

type SaleReportTask struct {
}

func (t SaleReportTask) Name() string {
	return "销售统计分析"
}

func (t SaleReportTask) Run() {
	sale_report_service.CalculateSaleReport()
}
