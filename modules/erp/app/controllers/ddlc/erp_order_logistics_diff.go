package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
	"fmt"
	"strings"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_logistics_service"
	"runtime/debug"
	xorm_core "github.com/go-xorm/core"
)

type ErpOrderLogisticsDiff struct {
	*revel.Controller
}

func (c ErpOrderLogisticsDiff) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	var db_session = c.DbSession

	var stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where("stock_type_id = ?", sppz_models.STOCK_TYPE_AU).Find(&stock_list))

	c.RenderArgs["stock_list"] = stock_list
	return c.RenderTemplate("erp/ddlc/order_logistics_diff/order_logistics_diff_index.html")
}

/*
func (c ErpOrderLogisticsDiff) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
	})
}
*/

func (c ErpOrderLogisticsDiff) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	var db_session = c.DbSession

	var stock_query int64
	c.Params.Bind(&stock_query, "stock_query")

	var stock_second_query int64
	c.Params.Bind(&stock_second_query, "stock_second_query")

	var shipped_time_begin_query_str string
	c.Params.Bind(&shipped_time_begin_query_str, "shipped_time_begin_query")
	shipped_time_begin_query, shipped_time_begin_query_err := core.Parse_time(shipped_time_begin_query_str)

	var shipped_time_end_query_str string
	c.Params.Bind(&shipped_time_end_query_str, "shipped_time_end_query")
	shipped_time_end_query, shipped_time_end_query_err := core.Parse_time(shipped_time_end_query_str)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type logistics_diff_view_item_t struct {
		OrderLogisticsLogic  ddlc_models.OrderLogisticsLogicInfo  `xorm:"extends" json:"oll"`
		Order                ddlc_models.OrderInfo                `xorm:"extends" json:"o"`
		OrderShippingAddress ddlc_models.OrderShippingAddressInfo `xorm:"extends" json:"osa"`
		OrderHoldType        ddpz_models.OrderHoldTypeInfo        `xorm:"extends" json:"ht"`
		LogisticsMode        ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm"`
		LogisticsModeSecond  ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm_sec"`
	}

	/*
SELECT oll.*, o.*, osa.*, ht.*, lm.*, lm_sec.*
FROM erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    LEFT JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    LEFT JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second;
	 */

	sb := core.NewSqlBuilder()
	sb.Select("oll.*, o.*, osa.*, ht.*, lm.*, lm_sec.*")
	sb.Count("count(oll.id) as count")
	sb.From(`erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    LEFT JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    LEFT JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second`)

	if stock_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("oll.stock_id = %d", stock_query))
	}

	if stock_second_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("lm_sec.stock_id = %d", stock_second_query))
	}

	if shipped_time_begin_query_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time >= '%04d-%02d-%02d'", shipped_time_begin_query.Year(), shipped_time_begin_query.Month(), shipped_time_begin_query.Day()))
	}

	if shipped_time_end_query_err == nil {
		shipped_time_end_query = shipped_time_end_query.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time < '%04d-%02d-%02d'", shipped_time_end_query.Year(), shipped_time_end_query.Month(), shipped_time_end_query.Day()))
	}

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d OR o.status_id = %d ", models.ORDER_STATUS_YFH, models.ORDER_STATUS_YWC)) // 只显示已发货和已完成的

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("oll.id DESC")
	}

	sb.Limit(offset, limit)

	logistics_mode_list := make([]logistics_diff_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_mode_list)
	core.HandleError(err)

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_mode_list,
		Total: count.Count,
	})
}

func (c ErpOrderLogisticsDiff) Match() revel.Result {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				revel.ERROR.Printf("匹配订单物流理论值出现错误： %s\n", err)
				revel.ERROR.Println(string(debug.Stack()))
			}
		}()

		// get db driver and db url
		db_driver, _ := revel.Config.String("db.driver")
		db_url, _ := revel.Config.String("db.url")

		db_engine, err := xorm.NewEngine(db_driver, db_url)
		core.HandleError(err)

		db_engine.TZLocation = time.Local
		db_engine.SetMaxIdleConns(3)
		db_engine.SetMaxOpenConns(100)

		if revel.DevMode == true {
			db_engine.ShowSQL()
			db_engine.ShowExecTime()
			xorm_core.LogQuerySql = true
		}

		db_session := db_engine.NewSession()
		db_session = db_session.NoAutoCondition(true)

		defer func() {
			db_session.Close()
		}()

		var order_list = make([]ddlc_models.OrderInfo, 0)

		var sql = `SELECT *
FROM erp_order o
WHERE o.id NOT IN (
    SELECT oll.order_id
    FROM erp_order_logistics_logic oll
)`
		core.HandleError(db_session.Sql(sql).Find(&order_list))

		if len(order_list) > 0 {
			var p_order_logistics = new(order_logistics_service.OrderLogisticsService)
			p_order_logistics.CalculateListLogic(order_list, db_session)
		}
	}()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动订单理论物流匹配成功！",
	})
}

type logistics_diff_export_item_t struct {
	OrderLogisticsLogic  ddlc_models.OrderLogisticsLogicInfo  `xorm:"extends" json:"oll"`
	Order                ddlc_models.OrderInfo                `xorm:"extends" json:"o"`
	OrderShippingAddress ddlc_models.OrderShippingAddressInfo `xorm:"extends" json:"osa"`
	Stock                sppz_models.StockInfo                `xorm:"extends" json:"s"`
	LogisticsMode        ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm"`
	StockSecond          sppz_models.StockInfo                `xorm:"extends" json:"s_sec"`
	LogisticsModeSecond  ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm_sec"`
}

func (c ErpOrderLogisticsDiff) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("oll.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	var stock_query int64
	c.Params.Bind(&stock_query, "stock_query")

	var stock_second_query int64
	c.Params.Bind(&stock_second_query, "stock_second_query")

	var shipped_time_begin_query_str string
	c.Params.Bind(&shipped_time_begin_query_str, "shipped_time_begin_query")
	shipped_time_begin_query, shipped_time_begin_query_err := core.Parse_time(shipped_time_begin_query_str)

	var shipped_time_end_query_str string
	c.Params.Bind(&shipped_time_end_query_str, "shipped_time_end_query")
	shipped_time_end_query, shipped_time_end_query_err := core.Parse_time(shipped_time_end_query_str)

	if stock_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("lm.stock_id = %d", stock_query))
	}

	if stock_second_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("lm_sec.stock_id = %d", stock_second_query))
	}

	if shipped_time_begin_query_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time >= '%04d-%02d-%02d'", shipped_time_begin_query.Year(), shipped_time_begin_query.Month(), shipped_time_begin_query.Day()))
	}

	if shipped_time_end_query_err == nil {
		shipped_time_end_query = shipped_time_end_query.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time < '%04d-%02d-%02d'", shipped_time_end_query.Year(), shipped_time_end_query.Month(), shipped_time_end_query.Day()))
	}

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d OR o.status_id = %d ", models.ORDER_STATUS_YFH, models.ORDER_STATUS_YWC)) // 只导出已发货和已完成的

	/*
SELECT oll.*, o.*, osa.*, s.*, lm.*, s_sec.*, lm_sec.*
FROM erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    INNER JOIN erp_stock s ON s.id = oll.stock_id
    INNER JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_stock s_sec ON s_sec.id = oll.stock_id_second
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("oll.*, o.*, osa.*, s.*, lm.*, s_sec.*, lm_sec.*")
	sb.From(`erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    INNER JOIN erp_stock s ON s.id = oll.stock_id
    INNER JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_stock s_sec ON s_sec.id = oll.stock_id_second
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second`)
	sb.Where(filter)

	var logistics_diff_export_item_list = make([]logistics_diff_export_item_t, 0)
	core.HandleError(db_session.Sql(sb.SqlList()).Find(&logistics_diff_export_item_list))

	xlsx_file := c.export_product_list_to_excel(db_session, logistics_diff_export_item_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "最佳仓次佳仓偏差导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderLogisticsDiff) export_product_list_to_excel(db_session *xorm.Session, logistics_diff_list []logistics_diff_export_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("最佳仓次佳仓偏差列表")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单编号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单代码"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "邮编"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "计费重量"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "最佳仓"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "最佳仓物流"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "最佳仓费用"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "次佳仓"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "次佳仓物流"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "次佳仓费用"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "费用偏差"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "匹配时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "发货时间"
		col_index++
	}

	for index, logistics_diff := range logistics_diff_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.Order.Code

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.Order.FullTranId

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.OrderShippingAddress.PostalCode

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = fmt.Sprintf("%0.2f", logistics_diff.OrderLogisticsLogic.Weight)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.LogisticsMode.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", logistics_diff.OrderLogisticsLogic.TotalFee)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.StockSecond.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = logistics_diff.LogisticsModeSecond.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = fmt.Sprintf("%0.2f", logistics_diff.OrderLogisticsLogic.TotalFeeSecond)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", logistics_diff.OrderLogisticsLogic.FeeDiff)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%s", logistics_diff.OrderLogisticsLogic.UpdateTime.String())

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if logistics_diff.Order.ShippedTime.Valid {
			cell.Value = fmt.Sprintf("%s", logistics_diff.Order.ShippedTime.Time.Format("2006-01-02 15:04:05"))
		} else {
			cell.Value = ""
		}
	}

	return xlsx_file
}
