package controllers

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
)

type ErpOrderHoldExcel struct {
	*revel.Controller
}

// 最佳仓缺货 Hold 导出
func (c ErpOrderHoldExcel) ExportBestStockOos() revel.Result {
	return nil
}

func (c ErpOrderExport) get_export_best_stock_oos_excel(db_session *xorm.Session) *xlsx.File {
	return nil
}
