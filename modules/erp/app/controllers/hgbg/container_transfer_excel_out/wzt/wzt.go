package wzt

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"fmt"
)

type Exporter struct {
}

func (c Exporter) GetExcelFileName(p_transfer *models.ContainerTransferInfo) string {
	return fmt.Sprintf("五洲通出口报关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *xlsx.File {
	return nil
}
