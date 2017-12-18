package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strconv"
)

type ErpIntelReplExcel struct {
	*revel.Controller
}

func (c ErpIntelReplExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	var db_session = c.DbSession

	/*
sku_query:
min_profit_query:
max_profit_query:
status_id_query
	 */

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")
	sku_query = strings.TrimSpace(sku_query)

	var min_profit_query_str string
	c.Params.Bind(&min_profit_query_str, "min_profit_query")
	min_profit_query_str = strings.TrimSpace(min_profit_query_str)
	var min_profit_query float64
	var err error
	if min_profit_query_str != "" {
		min_profit_query, err = strconv.ParseFloat(min_profit_query_str, 64)
		if err != nil {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "最小利润率输入有误！",
			})
		}
	}

	var max_profit_query_str string
	c.Params.Bind(&max_profit_query_str, "max_profit_query")
	max_profit_query_str = strings.TrimSpace(max_profit_query_str)
	var max_profit_query float64

	if max_profit_query_str != "" {
		max_profit_query, err = strconv.ParseFloat(max_profit_query_str, 64)
		if err != nil {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "最大利润率输入有误！",
			})
		}
	}

	var product_status_id_query int64
	c.Params.Bind(&product_status_id_query, "product_status_id_query")

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	var sale_user_id_query int64
	c.Params.Bind(&sale_user_id_query, "sale_user_id_query")

	var dev_user_id_query int64
	c.Params.Bind(&dev_user_id_query, "dev_user_id_query")

	// aoData.push({name: 'query_begin_time', value: $("#query_begin_time").val()});
	var query_begin_time_str string
	c.Params.Bind(&query_begin_time_str, "query_begin_time")
	query_begin_time_str = strings.TrimSpace(query_begin_time_str)

	query_begin_time, query_begin_time_err := core.Parse_time(query_begin_time_str)

	var query_end_time_str string
	c.Params.Bind(&query_end_time_str, "query_end_time")
	query_end_time_str = strings.TrimSpace(query_end_time_str)

	query_end_time, query_end_time_err := core.Parse_time(query_end_time_str)

	filter, order, _, _ := core.GetGridRequestParam(c.Request)

	/*

SELECT ir.*, p.*, ps.*, s.*, sup.*
FROM erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_status ps ON p.status_id = ps.id
    INNER JOIN erp_status s ON ir.status_id = s.id
    LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id;

	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ir.*, p.*, ps.*, s.*, sup.*")
	sb.From(`erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_status ps ON p.status_id = ps.id
    INNER JOIN erp_status s ON ir.status_id = s.id
	LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id`)

	if sku_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sku = '%s'", sku_query))
	}

	if min_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit >= %f", min_profit_query))
	}

	if max_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit <= %f", max_profit_query))
	}

	if product_status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.status_id = %d", product_status_id_query))
	}

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.status_id = %d", status_id_query))
	}

	if sale_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sale_user_id = %d", sale_user_id_query))
	}

	if dev_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.dev_user_id = %d", dev_user_id_query))
	}

	if query_begin_time_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.update_time >= '%s'", query_begin_time.Format("2006-01-02 15:04")))
	}

	if query_end_time_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.update_time <= '%s'", query_end_time.Format("2006-01-02 15:04")))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ir.id DESC")
	}

	var intel_repl_list = make([]intel_repl_view_item_t, 0)
	core.HandleError(db_session.Sql(sb.SqlList()).Find(&intel_repl_list))

	xlsx_file := c.get_export_excel_content(db_session, intel_repl_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("智能补货导出_%s.xlsx", time.Now().Format("2016-01-02 15:04"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpIntelReplExcel) get_export_excel_content(db_session *xorm.Session, item_list []intel_repl_view_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("智能补货")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "商品"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "供应商"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 35)
		set_head_cell_style(cell)
		cell.Value = "商品状态"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "毛利率"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "生成时间"
		col_index++
	}

	for _, item := range item_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)
		//
		//// 批次号
		//cell := row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.Value = entry.PurchaseBill.Code
		//
		//// 供应商
		//cell = row.AddCell()
		//set_text_cell_style(cell, "left")
		//cell.Value = entry.Supplier.Name
		//
		//// SKU
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.Value = entry.Product.Sku
		//
		//// 名称
		//cell = row.AddCell()
		//set_text_cell_style(cell, "left")
		//cell.Value = entry.Product.CnName
		//
		//// 价格
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBillEntry.Price, "0.00")
		//
		//// 数量
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetInt64(entry.PurchaseBillEntry.Qty)
		//
		//// 金额
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBillEntry.Amount, "0.00")
		//
		//// 运费
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBill.AdvanceExpenses+entry.PurchaseBill.FinalExpenses, "0.00")
		//
		//// 是否开票
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//if entry.PurchaseBillEntry.NeedInvoice == 0 {
		//	cell.SetString("否")
		//} else {
		//	cell.SetString("是")
		//}
		//
		//// 收税点
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBillEntry.TaxPoint, "0.00")
		//
		//// 退税点
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBillEntry.DrawbackPoint, "0.00")
		//
		//// 退税额
		//cell = row.AddCell()
		//set_text_cell_style(cell, "center")
		//cell.SetFloatWithFormat(entry.PurchaseBillEntry.DrawbackAmount, "0.00")
	}

	return xlsx_file
}
