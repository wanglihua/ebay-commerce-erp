package ozo_estore

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
	return fmt.Sprintf("ozoestore_进口清关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_in_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/in")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_in_dir, "ozo_estore.xlsx"))
	core.HandleError(err)

	// --------------------------------------------------------------------------------------------------------
	var sheet1_name = "Sheet1"

	xlsx_file.SetCellStr(sheet1_name, "G4", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "G5", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "G7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	var start_port models.ContainerTransferStartPortInfo
	_, err = db_session.Id(p_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "G13", start_port.PortName)

	var dest_port models.ContainerTransferDestPortInfo
	_, err = db_session.Id(p_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "G14", dest_port.PortName)

	var entry_list = make([]models.ContainerTransferEntryInfo, 0)
	core.HandleError(db_session.Where("transfer_id = ?", p_transfer.Id).Find(&entry_list))

	entry_list = container_transfer_service.EntryCombineAndPriceRateForAu(p_transfer, entry_list, db_session)

	var total_amount float64 = 0
	for index, entry := range entry_list {
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("A%d", 22+index), fmt.Sprintf("%s(%s)", entry.AuHsEnName, entry.AuHsCnName))
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("D%d", 22+index), entry.Price)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("F%d", 22+index), entry.Qty)
		var amount = entry.Price * float64(entry.Qty)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("H%d", 22+index), amount)

		total_amount += amount
	}

	xlsx_file.SetCellValue(sheet1_name, "H122", total_amount)

	// --------------------------------------------------------------------------------------------------------
	var sheet2_name = "Sheet2"

	xlsx_file.SetCellStr(sheet2_name, "G13", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G14", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G16", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	var total_qty int64 = 0
	var total_pkgs int64 = 0
	var total_net_weight float64 = 0
	var total_vol float64 = 0
	var total_gross_weight float64 = 0
	for index, entry := range entry_list {
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("A%d", 19+index), fmt.Sprintf("%s(%s)", entry.AuHsEnName, entry.AuHsCnName))
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("D%d", 19+index), entry.Qty)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("E%d", 19+index), entry.BoxQty)

		var net_weight = entry.WeightSum
		var gross_weight = net_weight * 1.01
		var vol = entry.VolumeSum

		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("F%d", 19+index), net_weight)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("G%d", 19+index), vol)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("H%d", 19+index), gross_weight)

		total_qty += entry.Qty
		total_pkgs += entry.BoxQty
		total_net_weight += net_weight
		total_gross_weight += gross_weight
		total_vol += vol
	}

	xlsx_file.SetCellValue(sheet2_name, "I19", total_vol)

	xlsx_file.SetCellValue(sheet2_name, "D119", total_qty)
	xlsx_file.SetCellValue(sheet2_name, "E119", total_pkgs)
	xlsx_file.SetCellValue(sheet2_name, "F119", total_net_weight)
	xlsx_file.SetCellValue(sheet2_name, "G119", total_vol)
	xlsx_file.SetCellValue(sheet2_name, "H119", total_gross_weight)
	xlsx_file.SetCellValue(sheet2_name, "I119", total_vol)

	return xlsx_file
}
