package wzt

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
	return fmt.Sprintf("五洲通出口报关_%s.xlsx", p_transfer.BatchCode)
}

func (c Exporter) GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File {
	var excel_template_hgbg_out_dir = filepath.Join(revel.BasePath, "/static/excel-template/hgbg/out")

	xlsx_file, err := excelize.OpenFile(filepath.Join(excel_template_hgbg_out_dir, "wzt.xlsx"))
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
		xlsx_file.SetCellValue(sheet1_name, "C132", p_transfer.BatchCode)

		var total_amount float64 = 0
		for index, entry := range entry_list {
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("B%d", 12+index), fmt.Sprintf("%s %s", entry.CnHsCnName, entry.CnHsEnName))
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("C%d", 12+index), entry.Qty)
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("E%d", 12+index), entry.Price)
			var amount = entry.Price * float64(entry.Qty)
			xlsx_file.SetCellValue(sheet1_name, fmt.Sprintf("F%d", 12+index), amount)

			total_amount += amount
		}

		xlsx_file.SetCellValue(sheet1_name, "F112", total_amount)
		xlsx_file.SetCellValue(sheet1_name, "F115", total_amount)
		xlsx_file.SetCellValue(sheet1_name, "F117", total_amount)
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

	{
		// --------------------------------------------------------------------------------------------------------
		var sheet5_name = "Sheet5"

		/*
	SELECT
		DISTINCT
		hcc.*
	FROM erp_container_transfer_entry cte
		INNER JOIN erp_hs_code_cn hcc ON hcc.id = cte.cn_hs_id
	WHERE cte.transfer_id = 1
	ORDER BY hcc.hs_code;
		 */

		var sql = `
SELECT
    DISTINCT
    hcc.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_hs_code_cn hcc ON hcc.id = cte.cn_hs_id
WHERE cte.transfer_id = %d
ORDER BY hcc.hs_code;`
		sql = fmt.Sprintf(sql, p_transfer.Id)

		var hs_code_cn_list = make([]models.HsCodeCnReportElemInfo, 0)
		core.HandleError(db_session.Sql(sql).Find(&hs_code_cn_list))

		// 92D050
		header_style1, _ := xlsx_file.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"","ident":1},"font":{"bold":true},"fill":{"type":"pattern","color":["#92D050"],"pattern":1},"border":[{"type":"right","color":"000000","style":1}]}`)
		header_style2, _ := xlsx_file.NewStyle(`{"alignment":{"horizontal":"center", "vertical":"","ident":1},"font":{"bold":true},"fill":{"type":"pattern","color":["#92D050"],"pattern":1},"border":[{"type":"right","color":"000000","style":1}]}`)

		// {"alignment":{"horizontal":"center"}}
		var row_num = 1
		xlsx_file.SetCellValue(sheet5_name, "A1", "编号")
		xlsx_file.SetCellStyle(sheet5_name, "A1", "A1", header_style2)
		row_num++

		for index, hs_code_cn := range hs_code_cn_list {
			// --------------------------------------------------------------------------------------
			var first_cell_axis = fmt.Sprintf("A%d", row_num)
			xlsx_file.SetCellValue(sheet5_name, first_cell_axis, index+1)

			var product_code_title_cell_axis = fmt.Sprintf("B%d", row_num)
			xlsx_file.SetCellValue(sheet5_name, product_code_title_cell_axis, "商品编码")
			xlsx_file.SetCellStyle(sheet5_name, product_code_title_cell_axis, product_code_title_cell_axis, header_style2)

			var product_code_cell_axis = fmt.Sprintf("C%d", row_num)
			xlsx_file.SetCellValue(sheet5_name, product_code_cell_axis, hs_code_cn.HsCode)
			xlsx_file.SetCellStyle(sheet5_name, product_code_cell_axis, product_code_title_cell_axis, header_style1)

			row_num++

			// --------------------------------------------------------------------------------------
			xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "品名")
			xlsx_file.SetCellStr(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.CnName)

			row_num++

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem01 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 1")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem01)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem02 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 2")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem02)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem03 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 3")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem03)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem04 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 4")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem04)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem05 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 5")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem05)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem06 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 6")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem06)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem07 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 7")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem07)

				row_num++
			}

			// --------------------------------------------------------------------------------------
			//if hs_code_cn.Elem08 != "" {
			{
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 8")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem08)

				row_num++
			}

			if hs_code_cn.Elem09 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 9")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem09)

				row_num++
			}

			if hs_code_cn.Elem10 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 10")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem10)

				row_num++
			}

			if hs_code_cn.Elem11 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 11")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem11)

				row_num++
			}

			if hs_code_cn.Elem12 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 12")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem12)

				row_num++
			}

			if hs_code_cn.Elem13 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 13")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem13)

				row_num++
			}

			if hs_code_cn.Elem14 != "" {
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("B%d", row_num), "申报元素 14")
				xlsx_file.SetCellValue(sheet5_name, fmt.Sprintf("C%d", row_num), hs_code_cn.Elem14)

				row_num++
			}

			xlsx_file.MergeCell(sheet5_name, first_cell_axis, fmt.Sprintf("A%d", row_num-1))

			xlsx_file.MergeCell(sheet5_name, fmt.Sprintf("A%d", row_num), fmt.Sprintf("C%d", row_num))
			row_num++
		}

	}

	return xlsx_file
}
