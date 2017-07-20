package ozo

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/Luxurioust/excelize"
	"path/filepath"
	"github.com/revel/revel"
	"matrix/core"
)

type Exporter struct {
}

func (c Exporter) GetExcelFileName(p_transfer *models.ContainerTransferInfo) string {
	return fmt.Sprintf("澳中欧出口报关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_out_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/out")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_out_dir, "ozo.xlsx"))
	core.HandleError(err)

	return xlsx_file
}
