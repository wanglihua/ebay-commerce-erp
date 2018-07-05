package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/core"
	"fmt"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/service/excel_service"
	"matrix/core/collections"
	"matrix/modules/erp/service/logistics_company_service"
	"matrix/modules/erp/service/logistics_mode_service"
	"strings"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
)

type ErpLogisticsReport struct {
	*revel.Controller
}

func (c ErpLogisticsReport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	return c.RenderTemplate("erp/ddbb/logistics_report/logistics_report_index.html")
}

func (c ErpLogisticsReport) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	var db_session = c.DbSession

	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var logistics_company_list = make([]ddpz_models.LogisticsCompanyInfo, 0)
	err := db_session.OrderBy("ord asc").Find(&logistics_company_list)
	core.HandleError(err)

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err = db_session.OrderBy("name asc").Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_id":               0,
		"stock_list":             stock_list,
		"logistics_company_id":   0,
		"logistics_company_list": logistics_company_list,
		"logistics_mode_id":      0,
		"logistics_mode_list":    logistics_mode_list,
	})
}

func (c ErpLogisticsReport) GetLogisticsModeList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "logistics_company_id")

	var filter = "1=1"

	if stock_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("stock_id = %d", stock_id))
	}

	if logistics_company_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("company_id = %d", logistics_company_id))
	}

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Where(filter).OrderBy("name asc").Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"logistics_mode_list": logistics_mode_list,
	})
}

func (c ErpLogisticsReport) apply_query_filter(filter string) string {
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

		if order_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code = '%s'", order_code_query))
		}

		if begin_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
		}

		if end_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	// #5 订单物流中只显示和只导出配货中的订单
	filter = core.SqlAnd(filter, fmt.Sprintf("(o.status_id = %d OR o.status_id = %d)", models.ORDER_STATUS_YFH, models.ORDER_STATUS_YWC)) // 已发货，已完成

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

	return filter
}

func (c ErpLogisticsReport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	db_session := c.DbSession

	/*
	var fee_diff_min_query float64
	var fee_diff_max_query float64
	c.Params.Bind(&fee_diff_min_query, "fee_diff_min_query")
	c.Params.Bind(&fee_diff_max_query, "fee_diff_max_query")
	*/

	type order_logistics_view_item_t struct {
		OrderLogistics ddlc_models.OrderLogisticsInfo `xorm:"extends" json:"ol"`
		Order          ddlc_models.OrderInfo          `xorm:"extends" json:"o"`
		StockName      string                         `json:"stock_name"`
		CompanyName    string                         `json:"company_name"`
		ModeName       string                         `json:"mode_name"`
		Sfkltjz        string                         `json:"sfkltjz"` // 是否考虑体积重
	}

	/*
SELECT ol.*, o.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.batch_code = 0
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
      AND (ol.actual_fee - ol.total_fee) >= 0.01
      AND (ol.actual_fee - ol.total_fee) <= 100.00;
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id`)

	filter = c.apply_query_filter(filter)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ol.id ASC")
	}

	sb.Limit(offset, limit)

	order_logistics_list := make([]order_logistics_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_list)
	core.HandleError(err)

	var stock_id_set = collections.NewIntSet()
	var company_id_set = collections.NewIntSet()
	var mode_id_set = collections.NewIntSet()
	var order_id_set = collections.NewIntSet()
	var shop_id_set = collections.NewIntSet()

	for _, order_logistics := range order_logistics_list {
		stock_id_set.Add(order_logistics.OrderLogistics.StockId)
		company_id_set.Add(order_logistics.OrderLogistics.CompanyId)
		mode_id_set.Add(order_logistics.OrderLogistics.ModeId)
		order_id_set.Add(order_logistics.OrderLogistics.OrderId)
		shop_id_set.Add(order_logistics.Order.ShopId)
	}

	var stock_id_name_map = stock_service.Get_stock_id_name_map(stock_id_set.List(), db_session)
	var company_id_name_map = logistics_company_service.Get_company_id_name_map(company_id_set.List(), db_session)
	var mode_map = logistics_mode_service.Get_logistics_mode_map(mode_id_set.List(), db_session)

	for index, order_logistics := range order_logistics_list {
		if order_logistics.OrderLogistics.StockId != 0 {
			order_logistics_list[index].StockName = stock_id_name_map[order_logistics.OrderLogistics.StockId]
		}

		if order_logistics.OrderLogistics.CompanyId != 0 {
			order_logistics_list[index].CompanyName = company_id_name_map[order_logistics.OrderLogistics.CompanyId]
		}

		var sfkltjz = "--"
		if order_logistics.OrderLogistics.ModeId != 0 {
			order_logistics_list[index].ModeName = mode_map[order_logistics.OrderLogistics.ModeId].Name

			if mode_map[order_logistics.OrderLogistics.ModeId].Sfkltjz == 0 {
				sfkltjz = "否"
			} else {
				sfkltjz = "是"
			}

		}
		order_logistics_list[index].Sfkltjz = sfkltjz
	}

	var order_logistics_count struct {
		Count int64 `xorm:"'order_logistics_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&order_logistics_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_logistics_list,
		Total: order_logistics_count.Count,
	})
}

func (c ErpLogisticsReport) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)
	
	var db_session = c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "物流报表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpLogisticsReport) get_export_excel(db_session *xorm.Session) *xlsx.File {

	type order_logistics_view_item_t struct {
		OrderLogistics   ddlc_models.OrderLogisticsInfo   `xorm:"extends" json:"ol"`
		Order            ddlc_models.OrderInfo            `xorm:"extends" json:"o"`
		Stock            sppz_models.StockInfo            `xorm:"extends" json:"s"`
		LogisticsCompany ddpz_models.LogisticsCompanyInfo `xorm:"extends" json:"lc"`
		LogisticsMode    ddpz_models.LogisticsModeInfo    `xorm:"extends" json:"lm"`
	}

	/*
SELECT ol.*, o.*, s.*, lc.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.batch_code = 0
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
      AND (ol.actual_fee - ol.total_fee) >= 0.01
      AND (ol.actual_fee - ol.total_fee) <= 100.00
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*, s.*, lc.*, lm.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id`)

	var filter = "(1=1)"

	filter = c.apply_query_filter(filter)

	sb.Where(filter)

	sb.OrderBy("ol.id ASC")

	order_logistics_list := make([]order_logistics_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("订单物流")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		excel_service.Set_head_cell_style(cell)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		excel_service.Set_text_cell_style(cell, horizontal_align)
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
