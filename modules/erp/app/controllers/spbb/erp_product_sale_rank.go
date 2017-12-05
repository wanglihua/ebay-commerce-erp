package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/service/product_sale_analysis_service"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/spbb_models"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
)

type ErpProductSaleRank struct {
	*revel.Controller
}

func (c ErpProductSaleRank) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	return c.RenderTemplate("erp/spbb/product_sale_rank/product_sale_rank_index.html")
}

// 这个函数在这里没有被用到，因为查询条件比较简单
func (c ErpProductSaleRank) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	var db_session = c.DbSession

	var stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&stock_list))

	return c.RenderJson(map[string]interface{}{
		"analysis_type_id": 0,
		"analysis_type_list": []map[string]interface{}{
			map[string]interface{}{
				"id":   spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL,
				"name": "实际仓库发货",
			},
			map[string]interface{}{
				"id":   spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC,
				"name": "理论仓库发货",
			},
			map[string]interface{}{
				"id":   spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_SKU,
				"name": "不分仓统计",
			},
		},
		"stock_list": stock_list,
	})
}

type product_sale_analysis_rank_t struct {
	ProductSaleAnalysis spbb_models.ProductSaleAnalysisInfo `xorm:"extends" json:"sa"`
	Stock               sppz_models.StockInfo               `xorm:"extends" json:"s"`
}

func (c ErpProductSaleRank) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	db_session := c.DbSession

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT sa.*, s.*
FROM erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id
	 */

	var sb = core.NewSqlBuilder()

	sb.Select("sa.*, s.*")
	sb.Count("count(1) AS sa_count")
	sb.From(`erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id`)

	var sql_args = make([]interface{}, 0)

	if sku_query != "" {
		filter = core.SqlAnd(filter, "sa.sku = ?")
		sql_args = append(sql_args, sku_query)
	}

	// 只取不分仓统计的数据
	filter = core.SqlAnd(filter, "sa.type = ?")
	sql_args = append(sql_args, spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_SKU)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sa.id ASC")
	}

	sb.Limit(offset, limit)

	product_sale_analysis_list := make([]product_sale_analysis_rank_t, 0, limit)
	err := db_session.Sql(sb.SqlList(), sql_args...).Find(&product_sale_analysis_list)
	core.HandleError(err)

	for index, product_sale_day_analysis := range product_sale_analysis_list {
		var product kcgl_models.ProductInfo
		_, err := db_session.Id(product_sale_day_analysis.ProductSaleAnalysis.ProductId).Get(&product)
		core.HandleError(err)

		product_sale_analysis_list[index].ProductSaleAnalysis.ProductName = product.CnName
	}

	var sa_count struct {
		Count int64 `xorm:"'sa_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount(), sql_args...).Get(&sa_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_sale_analysis_list,
		Total: sa_count.Count,
	})
}

func (c ErpProductSaleRank) Analyze() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	go product_sale_analysis_service.ProductSaleAnalyze()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "立即启动销售分析后台任务成功！",
	})
}

func (c ErpProductSaleRank) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_RANK, c.Session)

	var db_session = c.DbSession

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")

	var analysis_type_id int64
	c.Params.Bind(&analysis_type_id, "analysis_type_id")

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var product_sale_analysis_setting = product_sale_analysis_setting_service.Get_product_sale_analysis_setting(db_session)

	var filter = "1=1";

	/*
SELECT sa.*, s.*
FROM erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id
	 */

	var sb = core.NewSqlBuilder()

	sb.Select("sa.*, s.*")
	sb.From(`erp_product_sale_analysis sa
    LEFT JOIN erp_stock s ON sa.stock_id = s.id`)

	var sql_args = make([]interface{}, 0)

	if sku_query != "" {
		filter = core.SqlAnd(filter, "sa.sku = ?")
		sql_args = append(sql_args, sku_query)
	}

	if analysis_type_id != 0 {
		filter = core.SqlAnd(filter, "sa.type = ?")
		sql_args = append(sql_args, analysis_type_id)
	}

	if stock_id != 0 {
		filter = core.SqlAnd(filter, "sa.stock_id = ?")
		sql_args = append(sql_args, stock_id)
	}

	sb.Where(filter)

	sb.OrderBy("sa.daily_mean_sales ASC")

	product_sale_analysis_list := make([]product_sale_analysis_rank_t, 0)
	err := db_session.Sql(sb.SqlList(), sql_args...).Find(&product_sale_analysis_list)
	core.HandleError(err)

	for index, product_sale_day_analysis := range product_sale_analysis_list {
		if product_sale_day_analysis.ProductSaleAnalysis.Type == spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "实际仓库发货"
		}

		if product_sale_day_analysis.ProductSaleAnalysis.Type == spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_LOGIC {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "理论仓库发货"
		}

		if product_sale_day_analysis.ProductSaleAnalysis.Type == spbb_models.PRODUCT_SALE_ANALYSIS_TYPE_SKU {
			product_sale_analysis_list[index].ProductSaleAnalysis.TypeName = "不分仓统计"
		}

		var product kcgl_models.ProductInfo
		_, err := db_session.Id(product_sale_day_analysis.ProductSaleAnalysis.ProductId).Get(&product)
		core.HandleError(err)

		product_sale_analysis_list[index].ProductSaleAnalysis.ProductName = product.CnName
	}

	xlsx_file := c.export_analysis_to_excel(db_session, product_sale_analysis_list, product_sale_analysis_setting)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品销售分析导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductSaleRank) export_analysis_to_excel(db_session *xorm.Session, analysis_list []product_sale_analysis_rank_t, product_sale_analysis_setting spbb_models.ProductSaleAnalysisSettingInfo) *xlsx.File {

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

	// 创建工作薄
	xlsx_file := xlsx.NewFile()

	// 添加工作表
	sheet, err := xlsx_file.AddSheet("商品列表")
	core.HandleError(err)

	// 添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "统计类型"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "仓库"
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

		// 90天销量 60天销量 30天销量 15天销量

		if product_sale_analysis_setting.Level1Days != 0 {
			cell = row.AddCell()
			sheet.SetColWidth(col_index, col_index, 10)
			set_head_cell_style(cell)
			cell.Value = fmt.Sprintf("%d天销量", product_sale_analysis_setting.Level1Days)
			col_index++
		}

		if product_sale_analysis_setting.Level2Days != 0 {
			cell = row.AddCell()
			sheet.SetColWidth(col_index, col_index, 10)
			set_head_cell_style(cell)
			cell.Value = fmt.Sprintf("%d天销量", product_sale_analysis_setting.Level2Days)
			col_index++
		}

		if product_sale_analysis_setting.Level3Days != 0 {
			cell = row.AddCell()
			sheet.SetColWidth(col_index, col_index, 10)
			set_head_cell_style(cell)
			cell.Value = fmt.Sprintf("%d天销量", product_sale_analysis_setting.Level3Days)
			col_index++
		}

		if product_sale_analysis_setting.Level4Days != 0 {
			cell = row.AddCell()
			sheet.SetColWidth(col_index, col_index, 10)
			set_head_cell_style(cell)
			cell.Value = fmt.Sprintf("%d天销量", product_sale_analysis_setting.Level4Days)
			col_index++
		}

		if product_sale_analysis_setting.Level5Days != 0 {
			cell = row.AddCell()
			sheet.SetColWidth(col_index, col_index, 10)
			set_head_cell_style(cell)
			cell.Value = fmt.Sprintf("%d天销量", product_sale_analysis_setting.Level5Days)
			col_index++
		}

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "日均销量"
		col_index++
	}

	for _, analysis := range analysis_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 统计类型
		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = analysis.ProductSaleAnalysis.TypeName

		// 仓库
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = analysis.Stock.Name

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = analysis.ProductSaleAnalysis.Sku

		// 商品名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = analysis.ProductSaleAnalysis.ProductName

		if product_sale_analysis_setting.Level1Days != 0 {
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(analysis.ProductSaleAnalysis.Level1Qty)
		}

		if product_sale_analysis_setting.Level2Days != 0 {
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(analysis.ProductSaleAnalysis.Level2Qty)
		}

		if product_sale_analysis_setting.Level3Days != 0 {
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(analysis.ProductSaleAnalysis.Level3Qty)
		}

		if product_sale_analysis_setting.Level4Days != 0 {
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(analysis.ProductSaleAnalysis.Level4Qty)
		}

		if product_sale_analysis_setting.Level5Days != 0 {
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(analysis.ProductSaleAnalysis.Level5Qty)
		}

		// 日均销量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", analysis.ProductSaleAnalysis.DailyMeanSales)
	}

	return xlsx_file
}
