package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"strings"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/cggl_models"
)

type ErpPurchaseAbnormalExcel struct {
	*revel.Controller
}

func (c ErpPurchaseAbnormalExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	var db_session = c.DbSession

	var bill_id_list = make([]int64, 0)
	c.Params.Bind(&bill_id_list, "bill_id_list")

	var filter = "1=1"

	/*
SELECT pbe.*, pb.*, p.*, at.*, s.*
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_purchase_abnormal_type at ON pbe.abnormal_type_id = at.id
    INNER JOIN erp_status s ON pbe.status_id = s.id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, pb.*, p.*, at.*, s.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_purchase_abnormal_type at ON pbe.abnormal_type_id = at.id
    INNER JOIN erp_status s ON pbe.status_id = s.id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id = %d OR pbe.status_id = %d", models.PURCHASE_BILL_STATUS_YC, models.PURCHASE_BILL_STATUS_YCL)) // 异常 和 异常已处理

	if len(bill_id_list) != 0 {
		var bill_id_str_list = make([]string, 0)
		for _, bill_id := range bill_id_list {
			bill_id_str_list = append(bill_id_str_list, fmt.Sprint(bill_id))
		}

		var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
		core.HandleError(db_session.In("id", bill_id_list).Find(&purchase_bill_list))

		var bill_code_list = make([]string, 0)
		for _, bill := range purchase_bill_list {
			bill_code_list = append(bill_code_list, bill.Code)
		}

		filter = core.SqlAnd(filter, fmt.Sprintf("pbe.id IN (%s)", strings.Join(bill_id_str_list, ",")))
	} else {
		var bill_code_query string
		c.Params.Bind(&bill_code_query, "bill_code_query")
		bill_code_query = strings.TrimSpace(bill_code_query)

		var status_id_query int64
		c.Params.Bind(&status_id_query, "status_id_query")

		if bill_code_query != "" && bill_code_query != "0" {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.code = '%s'", bill_code_query))
		}

		if status_id_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id = %d", status_id_query))
		}
	}

	sb.Where(filter)

	sb.OrderBy("pbe.id ASC")

	purchase_bill_entry_list := make([]PurchaseAbnormalViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	xlsx_file := c.export_excel(db_session, purchase_bill_entry_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "异常采购单导出.xlsx"

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpPurchaseAbnormalExcel) export_excel(db_session *xorm.Session, entry_list []PurchaseAbnormalViewItem) *xlsx.File {
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
	sheet, err := xlsx_file.AddSheet("异常采购单")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "批次号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "异常类型"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 35)
		set_head_cell_style(cell)
		cell.Value = "名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "入库数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "异常数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "采购时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "入库时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "偏差金额"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "处理结果"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 38)
		set_head_cell_style(cell)
		cell.Value = "备注"
		col_index++
	}

	for _, entry := range entry_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 采购批次号
		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.PurchaseBill.Code

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Status.Name

		// 异常类型
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.PurchaseAbnormalType.Name

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Product.Sku

		// 名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Product.CnName

		// 数量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetInt64(entry.PurchaseBillEntry.Qty)

		// 入库数量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetInt64(entry.PurchaseBillEntry.CheckinQty)

		// 异常数量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetInt64(entry.PurchaseBillEntry.Qty - entry.PurchaseBillEntry.CheckinQty)

		// 采购时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.PurchaseBill.ApplyTime.Time.Format("2006-01-02 15:03")

		// 入库时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.PurchaseBill.CheckinTime.Time.Format("2006-01-02 15:03")

		// 偏差金额
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.DiffAmount, "0.00")

		// 处理结果
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = ""

		if entry.PurchaseBillEntry.AbnormalProcessResultId != 0 {
			var purchase_abnormal_process_result cggl_models.PurchaseAbnormalProcessResultInfo
			_, err = db_session.Id(entry.PurchaseBillEntry.AbnormalProcessResultId).Get(&purchase_abnormal_process_result)
			core.HandleError(err)
			cell.Value = purchase_abnormal_process_result.Name
		}

		// 备注
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.PurchaseBillEntry.Desc.String
	}

	return xlsx_file
}
