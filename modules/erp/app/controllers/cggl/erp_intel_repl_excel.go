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
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/sppz_models"
	auth_models "matrix/modules/auth/models"
)

type ErpIntelReplExcel struct {
	*revel.Controller
}

type purchase_bill_entry_item_t struct {
	PurchaseBillEntry cggl_models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	PurchaseBill      cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	Supplier          sppz_models.SupplierInfo          `xorm:"extends" json:"sup"`
	ApplyUser         auth_models.UserInfo              `xorm:"extends" json:"ayu"`
	Product           kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	Status            models.StatusInfo                 `xorm:"extends" json:"s"`
}

func (c ErpIntelReplExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var db_session = c.DbSession

	var bill_id_list = make([]int64, 0)
	c.Params.Bind(&bill_id_list, "bill_id_list")

	if len(bill_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "未指定导出采购单！",
		})
	}

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

	/*
		SELECT pbe.*, pb.*, sup.*, ayu.*, p.*, s.*
		FROM erp_purchase_bill_entry pbe
			INNER JOIN erp_purchase_bill pb ON pb.id = pbe.bill_id
			INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
			INNER JOIN auth_user ayu ON ayu.id = pb.apply_user_id
			INNER JOIN erp_product p ON p.id = pbe.product_id
			INNER JOIN erp_status s ON s.id = pbe.status_id;
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("pbe.*, pb.*, sup.*, ayu.*, p.*, s.*")
	sb.From(`erp_purchase_bill_entry pbe
			INNER JOIN erp_purchase_bill pb ON pb.id = pbe.bill_id
			INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
			INNER JOIN auth_user ayu ON ayu.id = pb.apply_user_id
			INNER JOIN erp_product p ON p.id = pbe.product_id
			INNER JOIN erp_status s ON s.id = pbe.status_id`)

	var filter = fmt.Sprintf("pbe.bill_id IN (%s)", strings.Join(bill_id_str_list, ","))

	sb.Where(filter)

	sb.OrderBy("pbe.bill_id ASC, pbe.id ASC")

	entry_list := make([]purchase_bill_entry_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&entry_list)
	core.HandleError(err)

	xlsx_file := c.get_export_excel_content(db_session, entry_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("采购单明细导出_%s.xlsx", strings.Join(bill_code_list, "_"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpIntelReplExcel) get_export_excel_content(db_session *xorm.Session, entry_list []purchase_bill_entry_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("采购单明细")
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
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "供应商"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 35)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "价格"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "金额"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "运费"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "是否开票"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "收税点"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退税点"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退税额"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "配件费用"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "其它费用"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "申请人"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "采购时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "下单时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "入库时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "描述"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++
	}

	for _, entry := range entry_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 批次号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.PurchaseBill.Code

		// 供应商
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Supplier.Name

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.Product.Sku

		// 名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Product.CnName

		// 价格
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.Price, "0.00")

		// 数量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetInt64(entry.PurchaseBillEntry.Qty)

		// 金额
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.Amount, "0.00")

		// 运费
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBill.AdvanceExpenses+entry.PurchaseBill.FinalExpenses, "0.00")

		// 是否开票
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if entry.PurchaseBillEntry.NeedInvoice == 0 {
			cell.SetString("否")
		} else {
			cell.SetString("是")
		}

		// 收税点
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.TaxPoint, "0.00")

		// 退税点
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.DrawbackPoint, "0.00")

		// 退税额
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.DrawbackAmount, "0.00")

		// 配件费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.FittingCost, "0.00")

		// 其它费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetFloatWithFormat(entry.PurchaseBillEntry.OtherCost, "0.00")

		// 申请人
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetString(entry.ApplyUser.LoginName)

		// 采购时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if entry.PurchaseBill.ApplyTime.Valid == false {
			cell.SetString("")
		} else {
			cell.SetString(entry.PurchaseBill.ApplyTime.Time.Format("2006-01-02 15:03"))
		}

		// 下单时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if entry.PurchaseBill.PayTime.Valid == false {
			cell.SetString("")
		} else {
			cell.SetString(entry.PurchaseBill.PayTime.Time.Format("2006-01-02 15:03"))
		}

		// 入库时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if entry.PurchaseBill.CheckinTime.Valid == false {
			cell.SetString("")
		} else {
			cell.SetString(entry.PurchaseBillEntry.CheckinTime.Time.Format("2006-01-02 15:03"))
		}

		// 描述
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.PurchaseBillEntry.Desc.String

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.Status.Name
	}

	return xlsx_file
}
