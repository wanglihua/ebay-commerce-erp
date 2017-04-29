package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/setting"
)

type ErpOrderLogistics struct {
	*revel.Controller
}

func (c ErpOrderLogistics) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddlc/order_logistics/order_logistics_index.html")
}

func (c ErpOrderLogistics) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	var db_session = c.DbSession

	var stock_list = make([]models.StockInfo, 0)
	//err := db_session.Where(fmt.Sprintf("stock_type_id = %d AND status_id = %d", models.STOCK_TYPE_AU, models.STOCK_STATUS_ENABLED)).OrderBy("ord asc").Find(&stock_list) //只取澳洲的仓库
	err := db_session.Where(fmt.Sprintf("stock_type_id = %d", models.STOCK_TYPE_AU)).OrderBy("ord asc").Find(&stock_list) //只取澳洲的仓库
	core.HandleError(err)

	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	//err = db_session.Where(fmt.Sprintf("status_id = %d", models.LOGISTICS_COMPANY_ENABLED)).OrderBy("ord asc").Find(&logistics_company_list)
	err = db_session.OrderBy("ord asc").Find(&logistics_company_list)
	core.HandleError(err)

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(db_session)

	return c.RenderJson(map[string]interface{}{
		"stock_id":               0,
		"stock_list":             stock_list,
		"logistics_company_id":   0,
		"logistics_company_list": logistics_company_list,
		"batch_code":             0,
		"batch_code_list":        batch_code_list,
	})
}

func (c ErpOrderLogistics) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

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

	/*
	var fee_diff_min_query float64
	var fee_diff_max_query float64
	c.Params.Bind(&fee_diff_min_query, "fee_diff_min_query")
	c.Params.Bind(&fee_diff_max_query, "fee_diff_max_query")
	*/

	type order_logistics_view_item_t struct {
		models.OrderLogisticsInfo    `xorm:"extends" json:"ol"`
		models.OrderInfo             `xorm:"extends" json:"o"`
		models.StockInfo             `xorm:"extends" json:"s"`
		models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
		models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
	}

	/*
SELECT ol.*, o.*, s.*, lc.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.id IN (SELECT ib.order_id FROM erp_order_import_batch ib WHERE ib.batch_code = 0)
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
      AND (ol.actual_fee - ol.total_fee) >= 0.01
      AND (ol.actual_fee - ol.total_fee) <= 100.00
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*, s.*, lc.*, lm.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id`)

	if stock_id != 0 {
		// 这里是查询上的对仓库的限制，如果是拣货人员的话，下面还有是不是这个仓库的拣货人员的限制
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.stock_id = %d", stock_id))
	}

	if logistics_company_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.company_id = %d", logistics_company_id))
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT ib.order_id FROM erp_order_import_batch ib WHERE ib.batch_code = %d)", batch_code))
	}

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
	}

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
	/*
	SELECT *
	FROM erp_order_logistics ol
	WHERE ol.stock_id IN (
		SELECT spc.stock_id
		FROM erp_stock_pick_crew spc
		WHERE spc.user_id = 0)
	 */
	var login_user = core.GetLoginUser(c.Session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
	var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
	var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服
	if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
		filter = core.SqlAnd(filter, fmt.Sprintf(`ol.stock_id IN (
        SELECT spc.stock_id
        FROM erp_stock_pick_crew spc
        WHERE spc.user_id = %d)`, login_user.UserId))
	}

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

type OrderLogisticsDetailForm struct {
	OrderLogistics models.OrderLogisticsInfo `json:"order_logistics"`
}

func (f *OrderLogisticsDetailForm) IsCreate() bool {
	return f.OrderLogistics.Id == 0
}

func (f *OrderLogisticsDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.OrderLogistics.OrderId).Message("订单不能为空！")

	validation.Required(f.OrderLogistics.AddressId).Message("物流地址不能为空！")

	validation.Required(f.OrderLogistics.ModeId).Message("物流方式不能为空！")

	validation.Required(f.OrderLogistics.FeeRateId).Message("物流费率不能为空！")

	validation.Required(f.OrderLogistics.Length).Message("长度不能为空！")

	validation.Required(f.OrderLogistics.Width).Message("宽度不能为空！")

	validation.Required(f.OrderLogistics.Height).Message("高度不能为空！")

	validation.Required(f.OrderLogistics.Weight).Message("重量不能为空！")

	validation.Required(f.OrderLogistics.TotalFee).Message("物流费用不能为空！")

	validation.Required(f.OrderLogistics.ActualFee).Message("实际费用不能为空！")

	return validation.HasErrors() == false
}

func (c ErpOrderLogistics) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	db_session := c.DbSession

	var order_logistics_id int64
	c.Params.Bind(&order_logistics_id, "id")

	var order_logistics models.OrderLogisticsInfo
	if order_logistics_id != 0 {
		has, err := db_session.Id(order_logistics_id).Get(&order_logistics)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单物流不存在！"})
		}
	}

	var detail_form OrderLogisticsDetailForm
	detail_form.OrderLogistics = order_logistics

	return c.RenderJson(detail_form)
}

func (c ErpOrderLogistics) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	db_session := c.DbSession

	var detail_form OrderLogisticsDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	order_logistics := detail_form.OrderLogistics

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&order_logistics)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(order_logistics.Id).Update(&order_logistics)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpOrderLogistics) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	db_session := c.DbSession

	order_logistics_id_list := make([]int64, 0)
	c.Params.Bind(&order_logistics_id_list, "id_list")

	affected, err := db_session.In("id", order_logistics_id_list).Delete(new(models.OrderLogisticsInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpOrderLogistics) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单物流费用" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderLogistics) get_export_excel(db_session *xorm.Session) *xlsx.File {
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

	type order_logistics_view_item_t struct {
		OrderLogistics   models.OrderLogisticsInfo    `xorm:"extends" json:"ol"`
		Order            models.OrderInfo             `xorm:"extends" json:"o"`
		Stock            models.StockInfo             `xorm:"extends" json:"s"`
		LogisticsCompany models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
		LogisticsMode    models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
	}

	/*
SELECT ol.*, o.*, s.*, lc.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.id IN (SELECT ib.order_id FROM erp_order_import_batch ib WHERE ib.batch_code = 0)
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

	if stock_id != 0 {
		// 这里是查询上的对仓库的限制，如果是拣货人员的话，下面还有是不是这个仓库的拣货人员的限制
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.stock_id = %d", stock_id))
	}

	if logistics_company_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.company_id = %d", logistics_company_id))
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT ib.order_id FROM erp_order_import_batch ib WHERE ib.batch_code = %d)", batch_code))
	}

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	sb.Where(filter)

	sb.OrderBy("ol.id ASC")

	order_logistics_list := make([]order_logistics_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
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

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
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
		cell.Value = fmt.Sprint(order_logistics.OrderLogistics.TotalFee)
	}

	return xlsx_file
}
