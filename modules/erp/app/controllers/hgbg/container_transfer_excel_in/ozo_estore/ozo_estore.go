package ozo_estore

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"path/filepath"
	"github.com/revel/revel"
	"matrix/core"
	"github.com/Luxurioust/excelize"
)

type Exporter struct {
}

func (c Exporter) GetExcelFileName(p_transfer *models.ContainerTransferInfo) string {
	return fmt.Sprintf("ozoestore_进口清关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_in_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/in")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_in_dir, "ozo_estore.xlsx"))
	core.HandleError(err)

	var sheet1_name = "Sheet1"

	xlsx_file.SetCellStr(sheet1_name, "G4", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "G5", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "L1", sheet1_name)

	xlsx_file.SetSheetName(sheet1_name, "COMMERCIAL_INVOICE")

	return xlsx_file
}
