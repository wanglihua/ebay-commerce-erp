package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
)

type ErpContainerTransferEntryExcel struct {
	*revel.Controller
}

func (c ErpContainerTransferEntryExcel) Import() revel.Result {
	db_session := c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type product_hs_code_import_t struct {
		Sku          string  `json:"sku"`
		CnHsCode     string  `json:"cn_hs_code"`
		AuHsCode     string  `json:"au_hs_code"`
		ErrorMessage string  `json:"error_message"`

		product models.ProductInfo `json:"-"`
	}

	product_hs_code_import_list := make([]product_hs_code_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	//revel.TRACE.Printf("sheet_count: %d\n", sheet_count)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//revel.TRACE.Printf("max_row: %d\n", max_row)
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			//row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_hs_code product_hs_code_import_t

			col_index := 0

			// SKU
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_hs_code.Sku = cell_val
			col_index++

			// cn_hs_code
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.CnHsCode = cell_val
			col_index++

			// au_hs_code
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.AuHsCode = cell_val
			col_index++

			product_hs_code_import_list = append(product_hs_code_import_list, product_hs_code)
		}
	}

	if len(product_hs_code_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	//验证导入的数据
	var product_hs_code_import_error_list = make([]product_hs_code_import_t, 0)

	for index, product_hs_code_import := range product_hs_code_import_list {
		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_hs_code_import.Sku)).Get(&product)
		core.HandleError(err)

		if has == false {
			product_hs_code_import.ErrorMessage = fmt.Sprintf("商品SKU %s 数据库中不存在！", product_hs_code_import.Sku)
			product_hs_code_import_error_list = append(product_hs_code_import_error_list, product_hs_code_import)
		} else {
			product_hs_code_import_list[index].product = product
		}
	}

	if len(product_hs_code_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "导入出现错误，详细请查看页面下部列表！",
			"error_data_list": product_hs_code_import_error_list,
		})
	}

	for _, product_hs_code_import := range product_hs_code_import_list {
		var product = product_hs_code_import.product

		product.Sku = product_hs_code_import.Sku
		product.CnHsCode = product_hs_code_import.CnHsCode
		product.AuHsCode = product_hs_code_import.AuHsCode

		update_col_list := []string{"cn_hs_code", "au_hs_code"}
		_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpContainerTransferEntryExcel) Export() revel.Result {
	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var containerTransfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	entry_list := make([]models.ContainerTransferEntryInfo, 0)
	err = db_session.Where(fmt.Sprintf("transfer_id = %d", transfer_id)).Find(&entry_list)
	core.HandleError(err)

	xlsx_file := c.export_entry_list_to_excel(db_session, entry_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("集装箱调拨明细导出_%s.xlsx", containerTransfer.BatchCode) //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpContainerTransferEntryExcel) export_entry_list_to_excel(db_session *xorm.Session, entry_list []models.ContainerTransferEntryInfo) *xlsx.File {

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00CCFFCC", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	//统一设置内容单元格格式的函数
	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	head_row_height := 0.8
	common_row_height := 0.6

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	//添加工作表
	sheet, err := xlsx_file.AddSheet("集装箱调拨明细")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		/*
            SKU，数量，箱数，价格
		 */

		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "箱数"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "价格"
		col_index++
	}

	for _, entry := range entry_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// SKU
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.ProductSku

		// 数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetInt64(entry.Qty)

		// 箱数
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetInt64(entry.BoxQty)

		// 价格
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloat(entry.Price)
	}

	return xlsx_file
}
