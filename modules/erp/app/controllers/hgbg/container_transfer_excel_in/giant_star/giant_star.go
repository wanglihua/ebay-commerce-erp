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
	"matrix/modules/erp/service/container_transfer_service"
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

	// --------------------------------------------------------------------------------------------------------
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

	var entry_list = make([]models.ContainerTransferEntryInfo, 0)
	core.HandleError(db_session.Where("transfer_id = ?", p_transfer.Id).Find(&entry_list))

	entry_list = container_transfer_service.EntryCombineAndPriceRateForAu(p_transfer, entry_list, db_session)

	var total_qty int64 = 0
	var total_amount float64 = 0
	for index, entry := range entry_list {
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("B%d", 15+index), entry.AuHsCnName)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("C%d", 15+index), entry.AuHsEnName)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("D%d", 15+index), entry.Qty)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("E%d", 15+index), entry.Price)
		var amount = entry.Price * float64(entry.Qty)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("F%d", 15+index), amount)

		total_qty += entry.Qty
		total_amount += amount
	}

	xlsx_file.SetCellValue(sheet1_name, "D115", total_qty)
	xlsx_file.SetCellValue(sheet1_name, "F115", total_amount)

	// --------------------------------------------------------------------------------------------------------
	var sheet2_name = "Sheet2"

	xlsx_file.SetCellStr(sheet2_name, "G4", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G5", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	xlsx_file.SetCellStr(sheet2_name, "C10", start_port.PortName)
	xlsx_file.SetCellStr(sheet2_name, "D11", dest_port.PortName)

	total_qty = 0
	var total_pkgs int64 = 0
	var total_net_weight float64 = 0
	var total_vol float64 = 0
	var total_gross_weight float64 = 0
	for index, entry := range entry_list {
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("B%d", 14+index), entry.AuHsCnName)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("C%d", 14+index), entry.AuHsEnName)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("D%d", 14+index), entry.Qty)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("E%d", 14+index), entry.BoxQty)

		var net_weight = entry.WeightSum
		var gross_weight = net_weight * 1.01
		var vol = entry.VolumeSum

		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("F%d", 14+index), net_weight)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("G%d", 14+index), gross_weight)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("H%d", 14+index), vol)

		total_qty += entry.Qty
		total_pkgs += entry.BoxQty
		total_net_weight += net_weight
		total_gross_weight += gross_weight
		total_vol += vol
	}

	xlsx_file.SetCellValue(sheet2_name, "I14", total_vol)

	xlsx_file.SetCellValue(sheet2_name, "D114", total_qty)
	xlsx_file.SetCellValue(sheet2_name, "E114", total_pkgs)
	xlsx_file.SetCellValue(sheet2_name, "F114", total_net_weight)
	xlsx_file.SetCellValue(sheet2_name, "G114", total_gross_weight)
	xlsx_file.SetCellValue(sheet2_name, "H114", total_vol)

	return xlsx_file
}
