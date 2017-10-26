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
	"math"
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

	type product_purchase_setting_t struct {
		sku              string
		purchase_days    int64
		stock_alert_days int64
	}

	var product_purchase_setting_list = make([]product_purchase_setting_t, 0)

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
			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_weight product_purchase_setting_t

			col_index := 0 // 第一列是编号

			// 编号，只是判断一下是不是空的
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			col_index++

			//SKU
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_weight.sku = cell_val
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			purchase_days, err := strconv.ParseInt(cell_val, 10, 64)
			core.HandleError(err)
			product_weight.purchase_days = purchase_days
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			stock_alert_days, err := strconv.ParseInt(cell_val, 10, 64)
			core.HandleError(err)
			product_weight.stock_alert_days = stock_alert_days
			col_index++

			product_purchase_setting_list = append(product_purchase_setting_list, product_weight)
		}
	}

	if len(product_purchase_setting_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	has_zero_data := false
	for _, product_purchase_setting := range product_purchase_setting_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_purchase_setting.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_purchase_setting.sku
		product.PurchaseDays = product_purchase_setting.purchase_days
		product.StockAlertDays = product_purchase_setting.stock_alert_days
		product.StockAlertQty = int64(math.Floor(float64(product.StockAlertDays)*product.DailyMeanSales + 0.5))

		if product_purchase_setting.purchase_days == 0 || product_purchase_setting.stock_alert_days == 0 {
			has_zero_data = true
		}

		if has == true {
			update_col_list := []string{"purchase_days", "stock_alert_days", "stock_alert_qty"}
			_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
			core.HandleError(err)
		} else {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_purchase_setting.sku))
		}
	}

	message := "数据导入成功！"
	if has_zero_data == true {
		message = "数据导入成功, 部分数据为零！"
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

	file_name := "商品采购设置导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

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
	sheet, err := xlsx_file.AddSheet("商品采购设置")
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

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 50)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "采购天数"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "库存警戒天数"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "日均销量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "警戒库存"
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
		set_text_cell_style(cell, "left")
		cell.Value = product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.CnName

		// 采购天数
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%d", product.PurchaseDays)

		// 库存警戒天数
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%d", product.StockAlertDays)

		// 日均销量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", product.DailyMeanSales)

		// 警戒库存
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%d", product.StockAlertQty)

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = find_product_status_name(product.StatusId)
	}

	return xlsx_file
}
