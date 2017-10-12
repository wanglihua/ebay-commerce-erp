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
	"strconv"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/purchase_bill_log_service"
	"matrix/modules/erp/service/purchase_bill_service"
	"matrix/modules/erp/models/sppz_models"
	auth_models "matrix/modules/auth/models"
)

type ErpPurchaseBillEntryExcel struct {
	*revel.Controller
}

func (c ErpPurchaseBillEntryExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	revel.TRACE.Println("bill_id", bill_id)

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err := db_session.Id(bill_id).Get(&purchase_bill)
	core.HandleError(err)

	if purchase_bill_service.CanPurchaseBillModify(purchase_bill, c.Session) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单当前状态不能导入明细！",
		})
	}

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type purchase_bill_entry_import_t struct {
		ProductSku    string `json:"product_sku"`
		Price         string `json:"price"`
		MelQty        string `json:"mel_qty"`
		BneQty        string `json:"bne_qty"`
		NeedInvoice   string `json:"need_invoice"`
		TaxPoint      string `json:"tax_point"`
		DrawbackPoint string `json:"drawback_point"`
		FittingCost   string `json:"fitting_cost"`
		OtherCost     string `json:"other_cost"`
		Desc          string `json:"desc"`

		product        kcgl_models.ProductInfo `json:"-"`
		price          float64                 `json:"-"`
		mel_qty        int64                   `json:"-"`
		bne_qty        int64                   `json:"-"`
		need_invoice   int64                   `json:"-"`
		tax_point      float64                 `json:"-"`
		drawback_point float64                 `json:"-"`
		fitting_cost   float64                 `json:"-"`
		other_cost     float64                 `json:"-"`
	}

	purchase_bill_entry_import_list := make([]purchase_bill_entry_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			var cell *xlsx.Cell
			var cell_val string
			var err error

			var purchase_bill_entry_import purchase_bill_entry_import_t

			col_index := 0

			// SKU
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			purchase_bill_entry_import.ProductSku = cell_val
			col_index++

			// Price
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.Price = cell_val
			col_index++

			// MelQty
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.MelQty = cell_val
			col_index++

			// BneQty
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.BneQty = cell_val
			col_index++

			// NeedInvoice
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.NeedInvoice = cell_val
			col_index++

			// TaxPoint
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.TaxPoint = cell_val
			col_index++

			// DrawbackPoint
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.DrawbackPoint = cell_val
			col_index++

			// FittingCost
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.FittingCost = cell_val
			col_index++

			// OtherCost
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.OtherCost = cell_val
			col_index++

			// 备注
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			purchase_bill_entry_import.Desc = cell_val
			col_index++

			purchase_bill_entry_import_list = append(purchase_bill_entry_import_list, purchase_bill_entry_import)
		}
	}

	if len(purchase_bill_entry_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入列表为空！",
		})
	}

	for index, purchase_bill_entry_import := range purchase_bill_entry_import_list {
		var buffer bytes.Buffer

		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", purchase_bill_entry_import.ProductSku)).Get(&product)
		core.HandleError(err)

		if has == false {
			buffer.WriteString("商品SKU不存在！")

		} else {
			purchase_bill_entry_import_list[index].product = product
		}

		var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
		has, err = db_session.Where("product_id = ?", product.Id).Get(&purchase_bill_entry)
		core.HandleError(err)

		if has == true {
			buffer.WriteString(fmt.Sprintf("当前采购单的 商品 %s 的明细已存在！", product.Sku))
		}

		price, err := strconv.ParseFloat(purchase_bill_entry_import.Price, 64)
		if err != nil {
			buffer.WriteString("价格格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].price = price
		}

		mel_qty, err := strconv.ParseInt(purchase_bill_entry_import.MelQty, 10, 64)
		if err != nil {
			buffer.WriteString("MEL数量格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].mel_qty = mel_qty
		}

		bne_qty, err := strconv.ParseInt(purchase_bill_entry_import.BneQty, 10, 64)
		if err != nil {
			buffer.WriteString("BNE数量格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].bne_qty = bne_qty
		}

		if strings.TrimSpace(purchase_bill_entry_import.NeedInvoice) == "是" {
			purchase_bill_entry_import.need_invoice = 1
		} else {
			purchase_bill_entry_import.need_invoice = 0
		}

		tax_point, err := strconv.ParseFloat(purchase_bill_entry_import.TaxPoint, 64)
		if err != nil {
			buffer.WriteString("收税点格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].tax_point = tax_point
		}

		drawback_point, err := strconv.ParseFloat(purchase_bill_entry_import.DrawbackPoint, 64)
		if err != nil {
			buffer.WriteString("退税点格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].drawback_point = drawback_point
		}

		fitting_cost, err := strconv.ParseFloat(purchase_bill_entry_import.FittingCost, 64)
		if err != nil {
			buffer.WriteString("配件费用格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].fitting_cost = fitting_cost
		}

		other_cost, err := strconv.ParseFloat(purchase_bill_entry_import.OtherCost, 64)
		if err != nil {
			buffer.WriteString("其它费用格式输入不对！")
		} else {
			purchase_bill_entry_import_list[index].other_cost = other_cost
		}

		var error_str = buffer.String()

		if error_str != "" {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("导入出现错误，详细：%s", error_str),
			})
		}
	}

	for _, purchase_bill_entry_import := range purchase_bill_entry_import_list {
		var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
		purchase_bill_entry.BillId = bill_id
		purchase_bill_entry.ProductId = purchase_bill_entry_import.product.Id
		purchase_bill_entry.Price = purchase_bill_entry_import.price
		purchase_bill_entry.MelQty = purchase_bill_entry_import.mel_qty
		purchase_bill_entry.BneQty = purchase_bill_entry_import.bne_qty
		purchase_bill_entry.Qty = purchase_bill_entry_import.mel_qty + purchase_bill_entry_import.bne_qty
		purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(purchase_bill_entry.Qty)
		purchase_bill_entry.NeedInvoice = purchase_bill_entry_import.need_invoice
		purchase_bill_entry.TaxPoint = purchase_bill_entry_import.tax_point
		purchase_bill_entry.DrawbackPoint = purchase_bill_entry_import.drawback_point
		purchase_bill_entry.FittingCost = purchase_bill_entry_import.fitting_cost
		purchase_bill_entry.OtherCost = purchase_bill_entry_import.other_cost

		if purchase_bill_entry_import.Desc == "" {
			purchase_bill_entry.Desc = core.MakeNullString("", false)
		} else {
			purchase_bill_entry.Desc = core.MakeNullString(purchase_bill_entry_import.Desc, true)
		}

		purchase_bill_entry.StatusId = purchase_bill.StatusId

		purchase_bill_entry.DrawbackAmount = purchase_bill_service.CalculateBillEntryDrawbackAmount(purchase_bill_entry)

		_, err := db_session.Insert(&purchase_bill_entry)
		core.HandleError(err)
	}

	purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单明细导入", db_session, c.Session)

	// 更新采购单总金额
	purchase_bill_service.UpdatePurchaseBillTotalAmount(purchase_bill.Id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

type purchase_bill_entry_item_t struct {
	PurchaseBillEntry cggl_models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	PurchaseBill      cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	Supplier          sppz_models.SupplierInfo          `xorm:"extends" json:"sup"`
	ApplyUser         auth_models.UserInfo              `xorm:"extends" json:"ayu"`
	Product           kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	Status            models.StatusInfo                 `xorm:"extends" json:"s"`
}

func (c ErpPurchaseBillEntryExcel) Export() revel.Result {
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

	xlsx_file := c.export_entry_list_to_excel(db_session, entry_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("采购单明细导出_%s.xlsx", strings.Join(bill_code_list, "_"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpPurchaseBillEntryExcel) export_entry_list_to_excel(db_session *xorm.Session, entry_list []purchase_bill_entry_item_t) *xlsx.File {

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
