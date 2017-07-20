package elife_pty

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
	return fmt.Sprintf("elife_pty_ltd_进口清关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_in_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/in")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_in_dir, "elife_pty.xlsx"))
	core.HandleError(err)

	var sheet1_name = "Sheet1"

	xlsx_file.SetCellStr(sheet1_name, "F11", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "F13", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "F16", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	var start_port models.ContainerTransferStartPortInfo
	_, err = db_session.Id(p_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "B11", start_port.PortName)

	var dest_port models.ContainerTransferDestPortInfo
	_, err = db_session.Id(p_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "B12", dest_port.PortName)

	var sheet2_name = "Sheet2"

	xlsx_file.SetCellStr(sheet2_name, "B12", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "B14", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G14", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	return xlsx_file
}
