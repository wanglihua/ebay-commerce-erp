package wzt

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
)

type Exporter struct {
}

func (c Exporter) GetExcelFileName(p_transfer *models.ContainerTransferInfo) string {
	return ""
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *xlsx.File {
	return nil
}
