package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strings"
	"matrix/modules/erp/models/kcgl_models"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"strconv"
)

type ErpProductPurchaseSettingExcel struct {
	*revel.Controller
}

func (c ErpProductPurchaseSettingExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

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

	type product_weight_t struct {
		sku    string
		length float64
		width  float64
		height float64
		weight float64
	}

	product_weight_list := make([]product_weight_t, 0)

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

			var product_weight product_weight_t

			col_index := 1 //从第二列开始导入，略去第一列的编号

			if len(row.Cells) < 5 {
				break
			}

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_weight.sku = cell_val
			col_index++

			// length
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_length, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.length = product_length
			col_index++

			// width
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_width, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.width = product_width
			col_index++

			// height
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_height, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.height = product_height
			col_index++

			// weight
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			weight, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.weight = weight
			col_index++

			product_weight_list = append(product_weight_list, product_weight)
		}
	}

	if len(product_weight_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	has_zero_data := false
	for _, product_weight := range product_weight_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_weight.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_weight.sku
		product.Length = product_weight.length
		product.Width = product_weight.width
		product.Height = product_weight.height
		product.Weight = product_weight.weight

		if product_weight.length == 0 || product_weight.width == 0 || product_weight.height == 0 || product_weight.weight == 0 {
			has_zero_data = true
		}

		if has == true {
			update_col_list := []string{"length", "width", "height", "weight"}
			_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
			core.HandleError(err)
		} else {
			/*
			// 这里只用于维护和更新长宽高数据，不进行数据的插入操作
			return c.RenderJson(map[string]interface{}{
				"success": true,
				"message": "商品SKU数据库中不存在！",
			})
			*/

			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_weight.sku))

			/*
			product.CnName = product.Sku
			product.Name = product.Sku
			product.StatusId = models.PRODUCT_STATUS_ZC
			_, err = db_session.Insert(&product)
			core.HandleError(err)
			*/
		}
	}

	message := "数据导入成功！"
	if has_zero_data == true {
		message = "数据导入成功, 部分长宽高数据为零！"
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func (c ErpProductPurchaseSettingExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	product_list := make([]kcgl_models.ProductInfo, 0)
	err := db_session.Where(filter).Find(&product_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品长宽高重量导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductPurchaseSettingExcel) export_product_list_to_excel(db_session *xorm.Session, product_list []kcgl_models.ProductInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品长宽高重量列表")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	var find_product_status_name = func(status_id int64) string {
		for _, status := range product_status_list {
			if status.Id == status_id {
				return status.Name
			}
		}

		return ""
	}

	//序号	SKU	中文名称	产品长度（cm）	产品宽度（cm）	产品高度（cm）	实际重量(kg)	状态

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品长度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品宽度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品高度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际重量(kg)"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++
	}

	for index, product := range product_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product.CnName

		// 	产品长度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Length)

		// 产品宽度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Width)

		// 	产品高度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Height)

		// 	实际重量(kg)
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.3f", product.Weight)

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = find_product_status_name(product.StatusId)
	}

	return xlsx_file
}
