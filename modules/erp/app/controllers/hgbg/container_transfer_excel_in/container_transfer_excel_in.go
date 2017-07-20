package container_transfer_excel_in

import (
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/models"
)

type ContainerTransferExcelInExportInterface interface {
	GetExcelFileName(p_transfer *models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *xlsx.File
}
