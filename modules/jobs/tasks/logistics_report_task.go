package tasks

import (
	"matrix/modules/erp/service/logistics_report_service"
)

type LogisticsReportTask struct {
}

func (t LogisticsReportTask) Name() string {
	return "物流报表分析"
}

func (t LogisticsReportTask) Run() {
	logistics_report_service.Logistics_report()
}
