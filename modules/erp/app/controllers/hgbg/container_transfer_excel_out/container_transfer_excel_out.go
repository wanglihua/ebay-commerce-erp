package container_transfer_excel_in

import (
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
)

type ContainerTransferExcelOutExportInterface interface {
	GetExcelFileName() string
	GetExcelFile(transfer_id int64, db_session *xorm.Session) *xlsx.File
}
