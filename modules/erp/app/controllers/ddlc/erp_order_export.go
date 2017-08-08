package controllers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
	"matrix/modules/auth/service/pass_service"
)

type ErpOrderExport struct {
	*revel.Controller
}

func (c ErpOrderExport) Export() revel.Result {
	var status string
	c.Params.Bind(&status, "status")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id")

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var db_session = c.DbSession

	var status_id int64 = 0
	if status == "dph" {
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "hold" {
		status_id = models.ORDER_STATUS_HOLD
	}

	// 先看 URL 里面的状态分类
	status = strings.ToLower(status) //大小写
	if status == "all" {
		status_id = 0
	} else if status == "ycj" {
		status_id = models.ORDER_STATUS_YCJ
	} else if status == "yzf" {
		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		status_id = models.ORDER_STATUS_YWC
	} else if status == "yqx" {
		status_id = models.ORDER_STATUS_YQX
	}

	var filter = "1=1"

	if len(id_list) == 0 {
		if status_id != 0 {
			filter = fmt.Sprintf("o.status_id = %d", status_id)
		}

		if batch_code != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
		}

		if hold_type_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	xlsx_file := c.get_export_excel(db_session, filter, status_id)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单编码订单代码" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

type order_code_num_export_view_t struct {
	Order    models.OrderInfo          `xorm:"extends" json:"o"`
	HoldType models.OrderHoldTypeInfo  `xorm:"extends" json:"ht"`
}

func (c ErpOrderExport) get_export_excel(db_session *xorm.Session, filter string, status_id int64) *xlsx.File {

	/*
SELECT o.*, ht.*
FROM erp_order o
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
	 */

	var sql = `SELECT o.*, ht.*
FROM erp_order o
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
WHERE `

	sql = sql + filter

	var order_list = make([]order_code_num_export_view_t, 0)
	err := db_session.Sql(sql).Find(&order_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("订单编码订单代码")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
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
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyFont = true

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

	row_height := 0.6

	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0

		var cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "订单编码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单代码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "hold原因"

		if status_id == models.ORDER_STATUS_HOLD {
			col_index++
			cell = header_row.AddCell()
			sheet.SetColWidth(col_index, col_index, 20)
			set_head_cell_style(cell)
			cell.Value = "最佳仓"

			col_index++
			cell = header_row.AddCell()
			sheet.SetColWidth(col_index, col_index, 20)
			set_head_cell_style(cell)
			cell.Value = "最佳物流"
		}
	}

	for _, order := range order_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 订单编码
		var cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.Order.Code

		// 订单代码
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.Order.FullTranId

		// hold原因
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.HoldType.Name

		if status_id == models.ORDER_STATUS_HOLD {

			var order_logistics_logic models.OrderLogisticsLogicInfo
			_, err := db_session.Where(fmt.Sprintf("order_id = %d", order.Order.Id)).Get(&order_logistics_logic)
			core.HandleError(err)

			// 最佳仓

			var stock models.StockInfo
			_, err = db_session.Id(order_logistics_logic.StockId).Get(&stock)
			core.HandleError(err)

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.Value = stock.Name

			// 最佳物流
			var mode models.LogisticsModeInfo
			_, err = db_session.Id(order_logistics_logic.ModeId).Get(&mode)
			core.HandleError(err)

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.Value = mode.Name
		}
	}

	return xlsx_file
}

func (c ErpOrderExport) ExportLogistics() revel.Result {
	var db_session = c.DbSession

	xlsx_file := c.get_export_logisics_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单物流费用" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderExport) get_export_logisics_excel(db_session *xorm.Session) *xlsx.File {

	type order_logistics_view_item_t struct {
		OrderLogistics   models.OrderLogisticsInfo    `xorm:"extends" json:"ol"`
		Order            models.OrderInfo             `xorm:"extends" json:"o"`
		Stock            models.StockInfo             `xorm:"extends" json:"s"`
		LogisticsCompany models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
		LogisticsMode    models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
	}

	/*
SELECT
    ol.*, o.*, s.*, lc.*, lm.*
FROM
    erp_order o
    LEFT JOIN erp_order_logistics ol ON o.id = ol.order_id
    LEFT JOIN erp_stock s ON s.id = ol.stock_id
    LEFT JOIN erp_logistics_company lc ON lc.id = ol.company_id
    LEFT JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.batch_code = 0
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
      AND (ol.actual_fee - ol.total_fee) >= 0.01
      AND (ol.actual_fee - ol.total_fee) <= 100.00
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*, s.*, lc.*, lm.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order o
    LEFT JOIN erp_order_logistics ol ON o.id = ol.order_id
    LEFT JOIN erp_stock s ON s.id = ol.stock_id
    LEFT JOIN erp_logistics_company lc ON lc.id = ol.company_id
    LEFT JOIN erp_logistics_mode lm ON lm.id = ol.mode_id`)

	var filter = "(1=1)"

	filter = c.apply_query_filter(filter)

	sb.Where(filter)

	sb.OrderBy("ol.id ASC")

	order_logistics_list := make([]order_logistics_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_list)
	core.HandleError(err)

	for index, order_logistics := range order_logistics_list {
		if order_logistics.Order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PICK_UP {
			order_logistics_list[index].LogisticsMode.Name = "PICK UP"
		}

		if order_logistics.Order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_ENVELOPE {
			order_logistics_list[index].LogisticsMode.Name = "信封"
		}

		if order_logistics.Order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_EXPRESS {
			order_logistics_list[index].LogisticsMode.Name = "快递"
		}

		if order_logistics.Order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_INTERNATIONAL {
			order_logistics_list[index].LogisticsMode.Name = "国际件"
		}
	}

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("订单物流")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
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
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
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

	row_height := 0.6

	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "订单"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "订单代码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "发货仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "物流公司"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "物流方式"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "预估物流费用"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "运单号"
	}

	for index, order_logistics := range order_logistics_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// 订单
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.Order.Code

		// 订单代码
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.Order.FullTranId

		// 发货仓库
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.Stock.Name

		// 物流公司
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.LogisticsCompany.Name

		// 物流方式
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.LogisticsMode.Name

		// 预估物流费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.TotalFee)

		// 运单号
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order_logistics.Order.WaybillNumber
	}

	return xlsx_file
}

func (c ErpOrderExport) apply_query_filter(filter string) string {
	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "logistics_company_id")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var order_code_query string
	c.Params.Bind(&order_code_query, "order_code_query")
	order_code_query = strings.TrimSpace(order_code_query)

	var waybill_number_query string
	c.Params.Bind(&waybill_number_query, "waybill_number_query")
	waybill_number_query = strings.TrimSpace(waybill_number_query)

	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	var good_query_id int64 = 0
	c.Params.Bind(&good_query_id, "good_query_id")

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	if len(id_list) == 0 {
		if stock_id != 0 {
			// 这里是查询上的对仓库的限制，如果是拣货人员的话，下面还有是不是这个仓库的拣货人员的限制
			filter = core.SqlAnd(filter, fmt.Sprintf("ol.stock_id = %d", stock_id))
		}

		if logistics_company_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("ol.company_id = %d", logistics_company_id))
		}

		if batch_code != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
		}

		if order_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code = '%s'", order_code_query))
		}

		if waybill_number_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.waybill_number = '%s'", waybill_number_query))
		}

		if begin_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
		}

		if end_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	// #5 订单物流中只显示和只导出配货中的订单
	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_YWC)) //配货中

	/*
	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id <> %d", models.ORDER_STATUS_YQC))  // 已取消
	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id <> %d", models.ORDER_STATUS_HOLD)) // hold
	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id <> %d", models.ORDER_STATUS_YWC))  // 已完成
	*/

	/*
	if fee_diff_min_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("(ol.actual_fee - ol.total_fee) >= %f", fee_diff_min_query))
	}

	if fee_diff_max_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("(ol.actual_fee - ol.total_fee) <= %f", fee_diff_max_query))
	}
	*/

	// 分仓库查看
	// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
	var login_user = core.GetLoginUser(c.Session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
	var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
	var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服
	if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
		filter = core.SqlAnd(filter, fmt.Sprintf(`ol.stock_id IN (SELECT spc.stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = %d)`, login_user.UserId))
	}

	// 优化查询
	if good_query_id != 0 {
		if good_query_id == good_query_ss {
			/*
# 单产品单数量
SELECT * FROM erp_order o WHERE o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND 1 = 1 # order filter
    WHERE oi.quantity = 1
    GROUP BY oi.order_id
    HAVING count(oi.id) = 1
);
			 */

			var good_query = `
o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND %s
    WHERE oi.quantity = 1
    GROUP BY oi.order_id
    HAVING count(oi.id) = 1
)
`
			filter = fmt.Sprintf(good_query, filter)

		} else if good_query_id == good_query_sm {
			/*
# 单产品多数量
SELECT * FROM erp_order o WHERE o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND 1 = 1 # order filter
    WHERE oi.quantity > 1
    GROUP BY oi.order_id
    HAVING count(oi.id) = 1
);
			 */

			var good_query = `
o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND %s
    WHERE oi.quantity > 1
    GROUP BY oi.order_id
    HAVING count(oi.id) = 1
)
`
			filter = fmt.Sprintf(good_query, filter)

		} else if good_query_id == good_query_mm {
			/*
# 多产品
SELECT * FROM erp_order o WHERE o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND 1 = 1 # order filter
    GROUP BY oi.order_id
    HAVING count(oi.id) > 1
);
			 */

			var good_query = `
o.id IN (
    SELECT oi.order_id AS order_id
    FROM erp_order_item oi
        INNER JOIN erp_order o ON oi.order_id = o.id AND %s
    GROUP BY oi.order_id
    HAVING count(oi.id) > 1
)
`
			filter = fmt.Sprintf(good_query, filter)
		}
	}

	return filter
}
