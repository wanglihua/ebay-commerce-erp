package ozo

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/Luxurioust/excelize"
	"path/filepath"
	"github.com/revel/revel"
	"matrix/core"
	"time"
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

	var start_port models.ContainerTransferStartPortInfo
	_, err = db_session.Id(p_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)

	var dest_port models.ContainerTransferDestPortInfo
	_, err = db_session.Id(p_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)

	var entry_list = make([]models.ContainerTransferEntryInfo, 0)
	core.HandleError(db_session.Where("transfer_id = ?", p_transfer.Id).Find(&entry_list))

	{
		// --------------------------------------------------------------------------------------------------------
		var sheet1_name = "Sheet1"

		xlsx_file.SetCellStr(sheet1_name, "E1", p_transfer.BatchCode)
		xlsx_file.SetCellStr(sheet1_name, "E2", time.Time(p_transfer.PackingTime).Format("2006-01-02"))

		xlsx_file.SetCellStr(sheet1_name, "B125", start_port.TradeItem)

		xlsx_file.SetCellStr(sheet1_name, "E125", dest_port.PortName)

		var total_amount float64 = 0
		for index, entry := range entry_list {
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("B%d", 12+index), entry.CnHsCnName)
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("C%d", 12+index), entry.Qty)
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("E%d", 12+index), entry.Price)
			var amount = entry.Price * float64(entry.Qty)
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("F%d", 12+index), amount)

			total_amount += amount
		}

		xlsx_file.SetCellValue(sheet1_name, "F112", total_amount)
		xlsx_file.SetCellValue(sheet1_name, "F115", total_amount)
		xlsx_file.SetCellValue(sheet1_name, "F117", total_amount)
		xlsx_file.SetCellValue(sheet1_name, "C132", p_transfer.BatchCode)
	}

	{
		// --------------------------------------------------------------------------------------------------------
		var sheet2_name = "Sheet2"

		xlsx_file.SetCellStr(sheet2_name, "F4", p_transfer.BatchCode)
		xlsx_file.SetCellStr(sheet2_name, "F5", p_transfer.BatchCode)
		xlsx_file.SetCellStr(sheet2_name, "F7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))
		xlsx_file.SetCellStr(sheet2_name, "C10", start_port.PortName)
		xlsx_file.SetCellStr(sheet2_name, "C11", dest_port.PortName)
		xlsx_file.SetCellStr(sheet2_name, "F10", start_port.TradeItem)

		var total_qty int64 = 0
		var total_amount float64 = 0
		for index, entry := range entry_list {
			xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("B%d", 14+index), entry.CnHsCnName)
			xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("C%d", 14+index), entry.CnHsEnName)
			xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("D%d", 14+index), entry.Qty)
			xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("E%d", 14+index), entry.Price)
			var amount = entry.Price * float64(entry.Qty)
			xlsx_file.SetCellValue(sheet2_name, fmt.Sprintf("F%d", 14+index), amount)

			total_qty += entry.Qty
			total_amount += amount
		}

		xlsx_file.SetCellValue(sheet2_name, "D114", total_qty)
		xlsx_file.SetCellValue(sheet2_name, "F114", total_amount)
	}

	{
		// --------------------------------------------------------------------------------------------------------
		var sheet3_name = "Sheet3"

		xlsx_file.SetCellStr(sheet3_name, "G4", p_transfer.BatchCode)
		xlsx_file.SetCellStr(sheet3_name, "G5", p_transfer.BatchCode)
		xlsx_file.SetCellStr(sheet3_name, "G7", time.Time(p_transfer.PackingTime).Format("2006-01-02"))
		xlsx_file.SetCellStr(sheet3_name, "C10", start_port.PortName)
		xlsx_file.SetCellStr(sheet3_name, "G10", start_port.TradeItem)
		xlsx_file.SetCellStr(sheet3_name, "D11", dest_port.PortName)

		var total_qty int64 = 0
		var total_pkgs int64 = 0
		var total_net_weight float64 = 0
		var total_vol float64 = 0
		var total_gross_weight float64 = 0
		for index, entry := range entry_list {
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("B%d", 14+index), entry.CnHsCnName)
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("C%d", 14+index), entry.CnHsEnName)
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("D%d", 14+index), entry.Qty)
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("E%d", 14+index), entry.BoxQty)

			var product models.ProductInfo
			_, err := db_session.Id(entry.ProductId).Get(&product)
			core.HandleError(err)

			var net_weight = product.Weight * float64(entry.Qty)
			var gross_weight = net_weight * 1.01
			var vol = (product.Length * product.Width * product.Height / 1000000) * float64(entry.Qty)

			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("F%d", 14+index), net_weight)
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("G%d", 14+index), gross_weight)
			xlsx_file.SetCellValue(sheet3_name, fmt.Sprintf("H%d", 14+index), vol)

			total_qty += entry.Qty
			total_pkgs += entry.BoxQty
			total_net_weight += net_weight
			total_gross_weight += gross_weight
			total_vol += vol
		}

		xlsx_file.SetCellValue(sheet3_name, "I14", total_vol)

		xlsx_file.SetCellValue(sheet3_name, "D114", total_qty)
		xlsx_file.SetCellValue(sheet3_name, "E114", total_pkgs)
		xlsx_file.SetCellValue(sheet3_name, "F114", total_net_weight)
		xlsx_file.SetCellValue(sheet3_name, "G114", total_gross_weight)
		xlsx_file.SetCellValue(sheet3_name, "H114", total_vol)
	}

	{
		// --------------------------------------------------------------------------------------------------------
		var sheet4_name = "Sheet4"

		for index, entry := range entry_list {
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("B%d", 2+index), entry.CnHsCode)
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("C%d", 2+index), entry.CnHsCnName)
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("D%d", 2+index), entry.Qty)
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("E%d", 2+index), "个")
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("F%d", 2+index), entry.Price)
			var amount = entry.Price * float64(entry.Qty)
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("G%d", 2+index), amount)
			xlsx_file.SetCellValue(sheet4_name, fmt.Sprintf("H%d", 2+index), "澳币")

		}
	}

	return xlsx_file
}
