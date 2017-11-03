package controllers

import (
	"github.com/revel/revel"

	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"strings"
	"matrix/modules/erp/models/kcgl_models"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/spzl_models"
	"matrix/modules/erp/models/sppz_models"
)

type ErpProductFeedbackExcel struct {
	*revel.Controller
}

func (c ErpProductFeedbackExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	return c.RenderEmpty()
}

type product_feedback_export_item_t struct {
	ProductFeedback     spzl_models.ProductFeedbackInfo     `xorm:"extends" json:"pf"`
	ProductFeedbackCate spzl_models.ProductFeedbackCateInfo `xorm:"extends" json:"c"`
	Product             kcgl_models.ProductInfo             `xorm:"extends" json:"p"`
	Supplier            sppz_models.SupplierInfo            `xorm:"extends" json:"sup"`
}

func (c ErpProductFeedbackExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	var db_session = c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	var sku string
	c.Params.Bind(&sku, "sku")
	sku = strings.TrimSpace(sku)

	var feedback_cate_id int64
	c.Params.Bind(&feedback_cate_id, "feedback_cate_id")

	var supplier_id int64
	c.Params.Bind(&supplier_id, "supplier_id")

	var product kcgl_models.ProductInfo
	if sku != "" {
		_, err := db_session.Where("sku = ?", sku).Get(&product)
		core.HandleError(err)
	}

	/*
SELECT
    pf.*, c.*, p.*, sup.*
FROM erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`pf.*, c.*, p.*, sup.*`)
	sb.Count("count(pf.id) AS feedback_count")
	sb.From(`erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id`)

	if sku != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.id = %d", product.Id))
	}

	if feedback_cate_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.cate_id = %d", feedback_cate_id))
	}

	if supplier_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.supplier_id = %d", supplier_id))
	}

	sb.Where(filter)

	sb.OrderBy("pf.id DESC")

	product_feedback_list := make([]product_feedback_export_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&product_feedback_list)
	core.HandleError(err)

	xlsx_file := c.export_product_feedback_list_to_excel(db_session, product_feedback_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品质量反馈导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductFeedbackExcel) export_product_feedback_list_to_excel(db_session *xorm.Session, product_feedback_list []product_feedback_export_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品质量反馈")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "反馈类别"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 40)
		set_head_cell_style(cell)
		cell.Value = "商品名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 40)
		set_head_cell_style(cell)
		cell.Value = "供应商"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 50)
		set_head_cell_style(cell)
		cell.Value = "描述"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 50)
		set_head_cell_style(cell)
		cell.Value = "解决方案"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "反馈人"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "创建时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "修改时间"
		col_index++
	}

	for index, product_feedback := range product_feedback_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// 反馈类别
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.ProductFeedbackCate.Name

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.Product.Sku

		// 商品名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.Product.CnName

		// 供应商
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.Supplier.Name

		//  描述
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.ProductFeedback.Description

		// 解决方案
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_feedback.ProductFeedback.Solution

		// 反馈人
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product_feedback.ProductFeedback.UserLoginName

		// 创建时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = time.Time(product_feedback.ProductFeedback.CreateTime).Format("2006-01-02 15:04:05")

		// 修改时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = time.Time(product_feedback.ProductFeedback.UpdateTime).Format("2006-01-02 15:04:05")
	}

	return xlsx_file
}
