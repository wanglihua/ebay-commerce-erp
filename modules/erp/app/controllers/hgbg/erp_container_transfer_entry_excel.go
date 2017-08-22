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
	"strconv"
	"matrix/modules/erp/service/container_transfer_service"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/service/container_transfer_log_service"
)

type ErpContainerTransferEntryExcel struct {
	*revel.Controller
}

func (c ErpContainerTransferEntryExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var containerTransfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type transfer_entry_import_t struct {
		ProductSku   string  `json:"product_sku"`
		Qty          string  `json:"qty"`
		BoxQty       string  `json:"box_qty"`
		NeedInvoice  string  `json:"need_invoice"`
		ErrorMessage string  `json:"error_message"`

		product    models.ProductInfo            `json:"-"`
		hs_code_au models.HsCodeAuDutyInfo       `json:"-"`
		hs_code_cn models.HsCodeCnReportElemInfo `json:"-"`
		qty        int64                         `json:"-"`
		box_qty    int64                         `json:"-"`
		price      float64                       `json:"-"`
	}

	transfer_entry_import_list := make([]transfer_entry_import_t, 0)

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

			var transfer_entry_import transfer_entry_import_t

			col_index := 0

			// SKU
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			transfer_entry_import.ProductSku = cell_val
			col_index++

			// qty
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			transfer_entry_import.Qty = cell_val
			col_index++

			// box_qty
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			transfer_entry_import.BoxQty = cell_val
			col_index++

			// 是否开票
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			transfer_entry_import.NeedInvoice = cell_val
			col_index++

			transfer_entry_import_list = append(transfer_entry_import_list, transfer_entry_import)
		}
	}

	if len(transfer_entry_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入列表为空！",
		})
	}

	//验证导入的数据
	var container_transfer_entry_import_error_list = make([]transfer_entry_import_t, 0)

	for index, transfer_entry_import := range transfer_entry_import_list {
		var buffer bytes.Buffer

		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", transfer_entry_import.ProductSku)).Get(&product)
		core.HandleError(err)

		if has == false {
			buffer.WriteString("商品SKU不存在！")

		} else {
			transfer_entry_import_list[index].product = product
		}

		if product.AuHsCode == "" {
			buffer.WriteString("商品澳洲海关编码未设置")
		} else {
			var hsCodeAu models.HsCodeAuDutyInfo
			has, err := db_session.Where("hs_code = ?", product.AuHsCode).Get(&hsCodeAu)
			core.HandleError(err)

			if has == false {
				buffer.WriteString("商品澳洲海关编码已设置，但该海关编码不存在！")
			} else {
				transfer_entry_import_list[index].hs_code_au = hsCodeAu
			}
		}

		if product.CnHsCode == "" {
			buffer.WriteString("商品中国海关编码未设置！")
		} else {
			var hsCodeCn models.HsCodeCnReportElemInfo
			has, err := db_session.Where("hs_code = ?", product.CnHsCode).Get(&hsCodeCn)
			core.HandleError(err)

			if has == false {
				buffer.WriteString("商品中国海关编码已设置，但该海关编码不存在！")
			} else {
				transfer_entry_import_list[index].hs_code_cn = hsCodeCn
			}
		}

		qty, err := strconv.ParseInt(transfer_entry_import.Qty, 10, 64)
		if err != nil {
			buffer.WriteString("数量格式输入不对！")
		} else {
			transfer_entry_import_list[index].qty = qty
		}

		box_qty, err := strconv.ParseInt(transfer_entry_import.BoxQty, 10, 64)
		if err != nil {
			buffer.WriteString("箱数格式输入不对！")
		} else {
			transfer_entry_import_list[index].box_qty = box_qty
		}

		if transfer_entry_import.NeedInvoice != "是" && transfer_entry_import.NeedInvoice != "否" {
			buffer.WriteString("是否开票输入不对！")
		}

		var au_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)
		//transfer_entry_import_list[index].price = inventory_service.GetProductCnStockCostPrice(product.Id, db_session) / au_exchange_rate
		//transfer_entry_import_list[index].price = product_service.GetProductPrice(product.Id, db_session) / au_exchange_rate
		transfer_entry_import_list[index].price = product.CostPrice / au_exchange_rate

		if product.CostPrice == 0 {
			buffer.WriteString("商品的采购价为零！")
		}

		if product.Length == 0 || product.Width == 0 || product.Height == 0 || product.Weight == 0 {
			buffer.WriteString("商品长宽高重量数据有为零！")
		}

		transfer_entry_import.ErrorMessage = buffer.String()

		if transfer_entry_import.ErrorMessage != "" {
			container_transfer_entry_import_error_list = append(container_transfer_entry_import_error_list, transfer_entry_import)
		}
	}

	if len(container_transfer_entry_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "导入出现错误，详细请查看页面下部列表！",
			"error_data_list": container_transfer_entry_import_error_list,
		})
	}

	for _, transfer_entry_import := range transfer_entry_import_list {
		var transfer_entry models.ContainerTransferEntryInfo
		transfer_entry.TransferId = transfer_id
		transfer_entry.BatchCode = containerTransfer.BatchCode
		transfer_entry.ProductId = transfer_entry_import.product.Id
		transfer_entry.ProductSku = transfer_entry_import.ProductSku
		transfer_entry.Qty = transfer_entry_import.qty
		transfer_entry.BoxQty = transfer_entry_import.box_qty
		transfer_entry.Price = transfer_entry_import.price

		transfer_entry.AuHsId = transfer_entry_import.hs_code_au.Id
		transfer_entry.AuHsCode = transfer_entry_import.hs_code_au.HsCode
		transfer_entry.AuHsDutyRate = transfer_entry_import.hs_code_au.DutyRate
		transfer_entry.AuHsCnName = transfer_entry_import.hs_code_au.CnName
		transfer_entry.AuHsEnName = transfer_entry_import.hs_code_au.EnName

		transfer_entry.CnHsId = transfer_entry_import.hs_code_cn.Id
		transfer_entry.CnHsCode = transfer_entry_import.hs_code_cn.HsCode
		transfer_entry.CnHsCnName = transfer_entry_import.hs_code_cn.CnName
		transfer_entry.CnHsEnName = transfer_entry_import.hs_code_cn.EnName

		_, err := db_session.Insert(&transfer_entry)
		core.HandleError(err)
	}

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(transfer_id, db_session)

	container_transfer_log_service.Insert_container_transfer_log(containerTransfer, "导入调拨明细", db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

type cte_item_t struct {
	ContainerTransferEntry models.ContainerTransferEntryInfo `xorm:"extends" json:"cte"`
	Product                models.ProductInfo                `xorm:"extends" json:"p"`
}

func (c ErpContainerTransferEntryExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var containerTransfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	/*
		SELECT cte.*, p.*
			FROM erp_container_transfer_entry cte
		INNER JOIN erp_product p ON cte.product_id = p.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("cte.*, p.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id`)

	var filter = fmt.Sprintf("cte.transfer_id = %d", transfer_id)
	sb.Where(filter)

	sb.OrderBy("cte.id asc")

	entry_list := make([]cte_item_t, 0)
	err = db_session.Sql(sb.SqlList()).Find(&entry_list)
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

func (c ErpContainerTransferEntryExcel) export_entry_list_to_excel(db_session *xorm.Session, entry_list []cte_item_t) *xlsx.File {

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
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "数量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "箱数"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "开票"
		col_index++
		/*
		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "价格"
		col_index++
		*/
	}

	for _, entry := range entry_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// SKU
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = entry.Product.Sku

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = entry.Product.CnName

		// 数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetInt64(entry.ContainerTransferEntry.Qty)

		// 箱数
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetInt64(entry.ContainerTransferEntry.BoxQty)

		/*
		// 开票
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = func() string {
			if entry.ContainerTransferEntry.NeedInvoice == 0 {
				return "否"
			} else {
				return "是"
			}
		}()
		*/

		/*
		// 价格
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloat(entry.Price)
		*/
	}

	return xlsx_file
}
