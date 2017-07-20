package giant_star

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"path/filepath"
	"github.com/revel/revel"
	"matrix/core"
	"github.com/Luxurioust/excelize"
	"time"
)

type Exporter struct {
}

func (c Exporter) GetExcelFileName(p_transfer *models.ContainerTransferInfo) string {
	return fmt.Sprintf("Giant_Star_进口清关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_in_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/in")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_in_dir, "giant_star.xlsx"))
	core.HandleError(err)

	var sheet1_name = "Sheet1"

	xlsx_file.SetCellStr(sheet1_name, "F4", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "F5", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "F7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	var start_port models.ContainerTransferStartPortInfo
	_, err = db_session.Id(p_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "C10", start_port.PortName)

	var dest_port models.ContainerTransferDestPortInfo
	_, err = db_session.Id(p_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "C11", dest_port.PortName)

	var sheet2_name = "Sheet2"

	xlsx_file.SetCellStr(sheet2_name, "G4", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G5", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	xlsx_file.SetCellStr(sheet1_name, "C10", start_port.PortName)
	xlsx_file.SetCellStr(sheet1_name, "D11", dest_port.PortName)

	return xlsx_file
}
