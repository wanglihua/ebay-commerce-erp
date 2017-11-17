package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"fmt"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
	"matrix/modules/erp/models/sppz_models"
)

type ErpSupplierExcel struct {
	*revel.Controller
}

func (c ErpSupplierExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

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

	type supplier_import_t struct {
		Supplier sppz_models.SupplierInfo
		CityName string
	}

	supplier_import_list := make([]supplier_import_t, 0)

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
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var supplier_import supplier_import_t

			col_index := 0

			if len(row.Cells) < 5 {
				break
			}

			// *名称 联系人 电话 联系方式 地区 地址 描述

			// 名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			supplier_import.Supplier.Name = cell_val
			col_index++

			// 名称 *联系人 电话 联系方式 地区 地址 描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.Supplier.Contact = core.MakeNullString(cell_val, true)
			col_index++

			// 名称 联系人 *电话 联系方式 地区 地址 描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.Supplier.Phone = core.MakeNullString(cell_val, true)
			col_index++

			// 名称 联系人 电话 *联系方式 地区 地址 描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.Supplier.Contact = core.MakeNullString(cell_val, true)
			col_index++

			// 名称 联系人 电话 联系方式 *地区 地址 描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.CityName = cell_val
			col_index++

			// 名称 联系人 电话 联系方式 地区 *地址 描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.Supplier.Address = core.MakeNullString(cell_val, true)
			col_index++

			// 名称 联系人 电话 联系方式 地区 地址 *描述
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			supplier_import.Supplier.Desc = core.MakeNullString(cell_val, true)
			col_index++

			supplier_import_list = append(supplier_import_list, supplier_import)
		}
	}

	if len(supplier_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入供应商列表为空！",
		})
	}

	for index, supplier_import := range supplier_import_list {
		// 设置 region_id, 实际是 city_id
		if supplier_import.CityName != "" {
			var region sppz_models.RegionInfo
			has, err := db_session.Where("city = ? OR district = ?", supplier_import.CityName, supplier_import.CityName).Get(&region)
			core.HandleError(err)
			if has == false {
				panic(fmt.Sprintf("指定的地区 %s 不存在！", supplier_import.CityName))
			}

			supplier_import_list[index].Supplier.RegionId = region.Id
		} else {
			supplier_import.Supplier.RegionId = 0
		}
	}

	for _, supplier_import := range supplier_import_list {
		var supplier sppz_models.SupplierInfo
		has, err := db_session.Where(fmt.Sprintf("name = '%s'", supplier_import.Supplier.Name)).Get(&supplier)
		core.HandleError(err)

		supplier.Name = supplier_import.Supplier.Name
		supplier.Contact = supplier_import.Supplier.Contact
		supplier.ContactInfo = supplier_import.Supplier.ContactInfo
		supplier.Address = supplier_import.Supplier.Address
		supplier.Phone = supplier_import.Supplier.Phone
		supplier.Desc = supplier_import.Supplier.Desc
		supplier.RegionId = supplier_import.Supplier.RegionId

		if has == true {
			update_col_list := []string{"name", "contact", "contact_info", "address", "phone", "desc", "region_id"}
			_, err = db_session.Id(supplier.Id).Cols(update_col_list...).Update(&supplier)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&supplier)
			core.HandleError(err)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "供应商数据导入成功！",
	})
}

func (c ErpSupplierExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_EXPORT, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("s.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	/*
SELECT s.*, r.*
FROM erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("s.*, r.*")
	sb.From(`erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id`)
	sb.Where(filter)
	sb.OrderBy("s.id ASC")

	supplier_list := make([]supplier_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&supplier_list)
	core.HandleError(err)

	xlsx_file := c.export_supplier_list_to_excel(db_session, supplier_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "供应商导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpSupplierExcel) export_supplier_list_to_excel(db_session *xorm.Session, supplier_list []supplier_view_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品供应商")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	// 名称 联系人 电话 联系方式 地区 地址 描述

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "联系人"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "电话"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "联系方式"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "地区"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "地址"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 45)
		set_head_cell_style(cell)
		cell.Value = "描述"
		col_index++
	}

	for _, supplier := range supplier_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 名称 联系人 电话 联系方式 地区 地址 描述

		// 名称
		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.Name

		// 联系人
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.Contact.String

		// 电话
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.Phone.String

		// 联系方式
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.ContactInfo.String

		// 地区
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Region.City + supplier.Region.District

		// 地址
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.Address.String

		// 描述
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = supplier.Supplier.Desc.String
	}

	return xlsx_file
}
