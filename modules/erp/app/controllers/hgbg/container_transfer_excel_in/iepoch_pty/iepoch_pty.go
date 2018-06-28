package iepoch_pty

import (
	"matrix/modules/erp/models/hgbg_models"
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

func (c Exporter) GetExcelFileName(p_transfer *hgbg_models.ContainerTransferInfo) string {
	return fmt.Sprintf("iepoch_pty_ltd_进口清关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *hgbg_models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_in_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/in")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_in_dir, "iepoch_pty.xlsx"))
	core.HandleError(err)

	// --------------------------------------------------------------------------------------------------------
	var sheet1_name = "Sheet1"

	xlsx_file.SetCellStr(sheet1_name, "B10", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "B11", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet1_name, "G12", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

	var start_port hgbg_models.ContainerTransferStartPortInfo
	_, err = db_session.Id(p_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "B12", start_port.PortName)

	var dest_port hgbg_models.ContainerTransferDestPortInfo
	_, err = db_session.Id(p_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)
	xlsx_file.SetCellStr(sheet1_name, "B13", dest_port.PortName)

	var item_list = container_transfer_service.GetContainerTransferEntryItem(p_transfer.Id, db_session)

	item_list = container_transfer_service.EntryCombineAndPriceRateForAu(p_transfer, item_list, db_session)

	var total_amount float64 = 0
	for index, item := range item_list {
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("A%d", 16+index), fmt.Sprintf("%s(%s)", item.ContainerTransferEntry.AuHsEnName, item.ContainerTransferEntry.AuHsCnName))
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("D%d", 16+index), item.ContainerTransferEntry.Qty)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("E%d", 16+index), item.ContainerTransferEntry.Price)
		var amount = item.ContainerTransferEntry.Price * float64(item.ContainerTransferEntry.Qty)
		xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("F%d", 16+index), amount)

		total_amount += amount
	}

	xlsx_file.SetCellValue(sheet1_name, "F216", total_amount)

	// --------------------------------------------------------------------------------------------------------
	var sheet2_name = "Sheet2"

	xlsx_file.SetCellStr(sheet2_name, "B10", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "B11", p_transfer.BatchCode)
	xlsx_file.SetCellStr(sheet2_name, "G12", time.Time(p_transfer.PackingTime).Format("2006-01-02"))
	xlsx_file.SetCellStr(sheet2_name, "E10", start_port.PortName)
	xlsx_file.SetCellStr(sheet2_name, "G10", dest_port.PortName)

	var total_qty int64 = 0
	var total_pkgs int64 = 0
	var total_net_weight float64 = 0
	var total_vol float64 = 0
	var total_gross_weight float64 = 0
	for index, item := range item_list {
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("B%d", 16+index), fmt.Sprintf("%s(%s)", item.ContainerTransferEntry.AuHsEnName, item.ContainerTransferEntry.AuHsCnName))
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("C%d", 16+index), item.ContainerTransferEntry.Qty)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("D%d", 16+index), item.ContainerTransferEntry.BoxQty)

		var net_weight = item.ContainerTransferEntry.WeightSum
		var gross_weight = net_weight * 1.01
		var vol = item.ContainerTransferEntry.VolumeSum

		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("E%d", 16+index), net_weight)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("F%d", 16+index), gross_weight)
		xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("G%d", 16+index), vol)

		total_qty += item.ContainerTransferEntry.Qty
		total_pkgs += item.ContainerTransferEntry.BoxQty
		total_net_weight += net_weight
		total_gross_weight += gross_weight
		total_vol += vol
	}

	xlsx_file.SetCellValue(sheet2_name, "H16", total_vol)

	xlsx_file.SetCellValue(sheet2_name, "C216", total_qty)
	xlsx_file.SetCellValue(sheet2_name, "D216", total_pkgs)
	xlsx_file.SetCellValue(sheet2_name, "E216", total_net_weight)
	xlsx_file.SetCellValue(sheet2_name, "F216", total_gross_weight)
	xlsx_file.SetCellValue(sheet2_name, "G216", total_vol)
	xlsx_file.SetCellValue(sheet2_name, "H216", total_vol)

	return xlsx_file
}
