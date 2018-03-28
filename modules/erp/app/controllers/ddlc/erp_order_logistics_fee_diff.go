package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/setting"
	"matrix/core/collections"
	"matrix/modules/erp/service/logistics_company_service"
	"matrix/modules/erp/service/logistics_mode_service"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"strings"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
)

type ErpOrderLogisticsFeeDiff struct {
	*revel.Controller
}

func (c ErpOrderLogisticsFeeDiff) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddlc/order_logistics_fee_diff/order_logistics_fee_diff_index.html")
}

const (
	FEE_DIFF_ACTUAL_FEE_GREATER int64 = 100
	FEE_DIFF_ACTUAL_FEE_SMALLER int64 = 200
	FEE_DIFF_ACTUAL_FEE_EQUAL   int64 = 300
	FEE_DIFF_ACTUAL_FEE_MISS    int64 = 400
	FEE_DIFF_ACTUAL_FEE_EXIST   int64 = 500
)

type fee_diff_type_t struct {
	Code int64  `json:"code"`
	Name string `json:"name"`
}

func (c ErpOrderLogisticsFeeDiff) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	var db_session = c.DbSession

	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var logistics_company_list = make([]ddpz_models.LogisticsCompanyInfo, 0)
	err := db_session.OrderBy("ord asc").Find(&logistics_company_list)
	core.HandleError(err)

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list_limit_100(0, db_session)

	var fee_diff_type_list = make([]fee_diff_type_t, 0)
	fee_diff_type_list = append(fee_diff_type_list, fee_diff_type_t{Code: FEE_DIFF_ACTUAL_FEE_GREATER, Name: "实际费用大于预估"})
	fee_diff_type_list = append(fee_diff_type_list, fee_diff_type_t{Code: FEE_DIFF_ACTUAL_FEE_SMALLER, Name: "实际费用小于预估"})
	fee_diff_type_list = append(fee_diff_type_list, fee_diff_type_t{Code: FEE_DIFF_ACTUAL_FEE_EQUAL, Name: "实际费用等于预估"})
	fee_diff_type_list = append(fee_diff_type_list, fee_diff_type_t{Code: FEE_DIFF_ACTUAL_FEE_MISS, Name: "实际费用未导入"})
	fee_diff_type_list = append(fee_diff_type_list, fee_diff_type_t{Code: FEE_DIFF_ACTUAL_FEE_EXIST, Name: "实际费用已导入"})

	return c.RenderJson(map[string]interface{}{
		"stock_id":               0,
		"stock_list":             stock_list,
		"logistics_company_id":   0,
		"logistics_company_list": logistics_company_list,
		"batch_code":             0,
		"batch_code_list":        batch_code_list,
		"fee_diff_type_code":     0,
		"fee_diff_type_list":     fee_diff_type_list,
	})
}

func (c ErpOrderLogisticsFeeDiff) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	list_data, count := c.get_list_data(false)

	return c.RenderJson(core.GridResult{
		Data:  list_data,
		Total: count,
	})
}

func (c ErpOrderLogisticsFeeDiff) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单物流费用偏差" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderLogisticsFeeDiff) get_export_excel(db_session *xorm.Session) *xlsx.File {
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
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单代码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "城市"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "邮编"

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
		cell.Value = "长(cm)"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "宽(cm)"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "高(cm)"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "体积重（Kg）"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际重（Kg）"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "拣货费用"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "预估物流费用"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际物流费用"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "运费偏差"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "支付时间"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "发货时间"
	}

	order_logistics_view_list, _ := c.get_list_data(true)

	for _, order_logistics_view := range order_logistics_view_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 订单编号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.Order.Code

		// 订单代码
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.Order.FullTranId

		// 城市
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.OrderShippingAddress.CityName

		// 邮编
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.OrderShippingAddress.PostalCode

		// 发货仓库
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.StockName

		// 物流公司
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.CompanyName

		// 物流方式
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics_view.ModeName

		// 长(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.Length)

		// 宽(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.Width)

		// 高(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.Height)

		// 体积重（Kg）
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if order_logistics_view.OrderLogistics.DimWeight != 0 {
			cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.DimWeight)
		} else {
			cell.Value = ""
		}

		// 实际重（Kg）
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.Weight)

		// 拣货费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.PickFee)

		// 预估物流费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.TotalFee)

		// 实际物流费用
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.ActualFee)

		// 运费偏差
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", order_logistics_view.OrderLogistics.ActualFee-order_logistics_view.OrderLogistics.TotalFee)

		// 支付时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if order_logistics_view.Order.PaidTime.Valid {
			cell.Value = order_logistics_view.Order.PaidTime.Time.Format("2006-01-02 15:04")
		} else {
			cell.Value = ""
		}

		// 发货时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if order_logistics_view.Order.ShippedTime.Valid {
			cell.Value = order_logistics_view.Order.ShippedTime.Time.Format("2006-01-02 15:04")
		} else {
			cell.Value = ""
		}
	}

	return xlsx_file
}

type order_logistics_view_item_t struct {
	OrderLogistics       ddlc_models.OrderLogisticsInfo       `xorm:"extends" json:"ol"`
	Order                ddlc_models.OrderInfo                `xorm:"extends" json:"o"`
	OrderShippingAddress ddlc_models.OrderShippingAddressInfo `xorm:"extends" json:"osa"`
	StockName            string                               `json:"stock_name"`
	CompanyName          string                               `json:"company_name"`
	ModeName             string                               `json:"mode_name"`
}

func (c ErpOrderLogisticsFeeDiff) get_list_data(is_for_export bool) (list_data_return []order_logistics_view_item_t, count_return int64) {
	db_session := c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "logistics_company_id")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	var fee_diff_min_query float64
	var fee_diff_max_query float64
	c.Params.Bind(&fee_diff_min_query, "fee_diff_min_query")
	c.Params.Bind(&fee_diff_max_query, "fee_diff_max_query")

	var fee_diff_type_code int64
	c.Params.Bind(&fee_diff_type_code, "fee_diff_type_code")

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	/*
SELECT ol.*, o.*, osa.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = ol.address_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.batch_code = 0
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
      AND (ol.actual_fee - ol.total_fee) >= 0.01
      AND (ol.actual_fee - ol.total_fee) <= 100.00;
	 */

	var filter, order string
	var offset, limit int64

	if is_for_export == false {
		filter, order, offset, limit = core.GetGridRequestParam(c.Request)
	} else {
		filter, order, offset, limit = "1=1", "", 0, 0
	}

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*, osa.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = ol.address_id`)

	if len(id_list) == 0 {
		if stock_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("ol.stock_id = %d", stock_id))
		}

		if logistics_company_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("ol.company_id = %d", logistics_company_id))
		}

		if batch_code != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
		}

		if begin_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
		}

		if end_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
		}

		if fee_diff_min_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("(ol.actual_fee - ol.total_fee) >= %f", fee_diff_min_query))
		}

		if fee_diff_max_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("(ol.actual_fee - ol.total_fee) <= %f", fee_diff_max_query))
		}

		if fee_diff_type_code != 0 {
			if fee_diff_type_code == FEE_DIFF_ACTUAL_FEE_GREATER {
				filter = core.SqlAnd(filter, "ol.actual_fee > ol.total_fee")
			} else if fee_diff_type_code == FEE_DIFF_ACTUAL_FEE_SMALLER {
				filter = core.SqlAnd(filter, "ol.actual_fee < ol.total_fee AND ol.actual_fee <> 0")
			} else if fee_diff_type_code == FEE_DIFF_ACTUAL_FEE_EQUAL {
				filter = core.SqlAnd(filter, "ol.actual_fee = ol.total_fee AND ol.actual_fee <> 0")
			} else if fee_diff_type_code == FEE_DIFF_ACTUAL_FEE_MISS {
				filter = core.SqlAnd(filter, "ol.actual_fee = 0")
			} else if fee_diff_type_code == FEE_DIFF_ACTUAL_FEE_EXIST {
				filter = core.SqlAnd(filter, "ol.actual_fee <> 0")
			}
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	// 分仓库查看
	// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
	var login_user = core.GetLoginUser(c.Session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
	var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
	var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服
	if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
		filter = core.SqlAnd(filter, fmt.Sprintf(`ol.stock_id IN (SELECT spc.stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = %d)`, login_user.UserId))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("o.id DESC")
	}

	if is_for_export == false {
		sb.Limit(offset, limit)
	}

	order_logistics_list := make([]order_logistics_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_list)
	core.HandleError(err)

	var stock_id_set = collections.NewIntSet()
	var company_id_set = collections.NewIntSet()
	var mode_id_set = collections.NewIntSet()

	for _, order_logistics := range order_logistics_list {
		stock_id_set.Add(order_logistics.OrderLogistics.StockId)
		company_id_set.Add(order_logistics.OrderLogistics.CompanyId)
		mode_id_set.Add(order_logistics.OrderLogistics.ModeId)
	}

	var stock_id_name_map = stock_service.Get_stock_id_name_map(stock_id_set.List(), db_session)
	var company_id_name_map = logistics_company_service.Get_company_id_name_map(company_id_set.List(), db_session)
	var mode_id_name_map = logistics_mode_service.Get_logistics_mode_id_name_map(mode_id_set.List(), db_session)

	for index, order_logistics := range order_logistics_list {
		if order_logistics.OrderLogistics.StockId != 0 {
			order_logistics_list[index].StockName = stock_id_name_map[order_logistics.OrderLogistics.StockId]
		}

		if order_logistics.OrderLogistics.CompanyId != 0 {
			order_logistics_list[index].CompanyName = company_id_name_map[order_logistics.OrderLogistics.CompanyId]
		}

		if order_logistics.OrderLogistics.ModeId != 0 {
			order_logistics_list[index].ModeName = mode_id_name_map[order_logistics.OrderLogistics.ModeId]
		}
	}

	if is_for_export == false {
		var order_logistics_count struct {
			Count int64 `xorm:"'order_logistics_count'"`
		}
		_, err = db_session.Sql(sb.SqlCount()).Get(&order_logistics_count)
		core.HandleError(err)

		return order_logistics_list, order_logistics_count.Count
	} else {
		return order_logistics_list, 0
	}
}
