package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/core/collections"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/logistics_company_service"
	"matrix/modules/erp/service/logistics_mode_service"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/order_item_service"
	"matrix/modules/erp/service/shop_service"
	"matrix/modules/erp/service/stock_service"
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

// 查询条件的列表
func (c ErpOrderLogistics) get_good_query_list() []good_query {
	var good_query_list = make([]good_query, 0)

	good_query_list = append(good_query_list, good_query{Id: good_query_ss, Name: "单商品单数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_sm, Name: "单商品多数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_mm, Name: "多商品"})

	return good_query_list
}

func (c ErpOrderLogistics) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	var db_session = c.DbSession

	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	//err = db_session.Where(fmt.Sprintf("status_id = %d", models.LOGISTICS_COMPANY_ENABLED)).OrderBy("ord asc").Find(&logistics_company_list)
	err := db_session.OrderBy("ord asc").Find(&logistics_company_list)
	core.HandleError(err)

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(models.ORDER_STATUS_PHZ, db_session)

	return c.RenderJson(map[string]interface{}{
		"stock_id":               0,
		"stock_list":             stock_list,
		"logistics_company_id":   0,
		"logistics_company_list": logistics_company_list,
		"batch_code":             0,
		"batch_code_list":        batch_code_list,
		"good_query_id":          0,
		"good_query_list":        c.get_good_query_list(),
	})
}

func (c ErpOrderLogistics) apply_query_filter(filter string) string {
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
		filter = core.SqlAnd(filter, fmt.Sprintf("ol.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	// #5 订单物流中只显示和只导出配货中的订单
	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_PHZ)) //配货中

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

func (c ErpOrderLogistics) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	db_session := c.DbSession

	/*
	var fee_diff_min_query float64
	var fee_diff_max_query float64
	c.Params.Bind(&fee_diff_min_query, "fee_diff_min_query")
	c.Params.Bind(&fee_diff_max_query, "fee_diff_max_query")
	*/

	type order_logistics_view_item_t struct {
		OrderLogistics models.OrderLogisticsInfo  `xorm:"extends" json:"ol"`
		Order          models.OrderInfo           `xorm:"extends" json:"o"`
		StockName      string                     `json:"stock_name"`
		CompanyName    string                     `json:"company_name"`
		ModeName       string                     `json:"mode_name"`
		Products       []string                   `json:"products"`
		OrderCode      string                     `json:"order_code"` //交易代码 店铺号 + 交易号
		Sfkltjz        string                     `json:"sfkltjz"`    // 是否考虑体积重
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
	var order_id_products_map = order_item_service.Get_order_id_products_map(order_id_set.List(), db_session)
	var shop_id_name_map = shop_service.Get_shop_id_name_map(shop_id_set.List(), db_session)

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

		var order_sku_qty_list = order_id_products_map[order_logistics.OrderLogistics.OrderId]
		if len(order_sku_qty_list) != 0 {
			var products = make([]string, 0)
			for _, order_id_sku_qty := range order_sku_qty_list {
				products = append(products, fmt.Sprintf("%s × %d", order_id_sku_qty.Sku, order_id_sku_qty.Qty))
			}
			order_logistics_list[index].Products = products
		}

		if order_logistics.Order.ShopId != 0 {
			order_logistics_list[index].OrderCode = shop_id_name_map[order_logistics.Order.ShopId].ShopCode + order_logistics.Order.TranId
		}
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
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
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

func (c ErpOrderLogistics) Export1() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel1(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单物流费用" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderLogistics) get_export_excel1(db_session *xorm.Session) *xlsx.File {

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

func (c ErpOrderLogistics) Export2() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel2(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单物流" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderLogistics) get_export_excel2(db_session *xorm.Session) *xlsx.File {

	type order_logistics_view_item_t struct {
		OrderLogistics       models.OrderLogisticsInfo          `xorm:"extends" json:"ol"`
		Order                models.OrderInfo                   `xorm:"extends" json:"o"`
		Stock                models.StockInfo                   `xorm:"extends" json:"s"`
		LogisticsCompany     models.LogisticsCompanyInfo        `xorm:"extends" json:"lc"`
		LogisticsMode        models.LogisticsModeInfo           `xorm:"extends" json:"lm"`
		OrderShippingAddress models.OrderShippingAddressInfo    `xorm:"extends" json:"osa"`
		ShopAccount          models.ShopAccountInfo             `xorm:"extends" json:"sa"`
	}

	/*
SELECT ol.*, o.*, s.*, lc.*, lm.*, osa.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = ol.address_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE ol.stock_id = 0 AND ol.company_id = 0
      AND o.batch_code = 0
      AND o.paid_time > '2017-01-01' AND o.paid_time < '2017-04-30'
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ol.*, o.*, s.*, lc.*, lm.*, osa.*, sa.*")
	sb.Count("count(ol.id) AS order_logistics_count")
	sb.From(`erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_stock s ON s.id = ol.stock_id
    INNER JOIN erp_logistics_company lc ON lc.id = ol.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = ol.address_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id`)

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

	var col_count_without_skus = 0
	// 仓库 物流公司 物流方式 长(cm) 宽(cm) 高(cm) 实际重量（Kg） 邮政编码 客户姓名 电话 客户email 邮寄地址 所属城市 所属地区 国家 订单代码 仓位 商品编号 商品数量 仓位 商品编号 商品数量
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "仓库"

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
		cell.Value = "实际重量（Kg）"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "邮政编码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "客户姓名"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "电话"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "客户email"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "邮寄地址"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "所属城市"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "所属地区"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "国家"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "订单代码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "店铺名"

		col_count_without_skus = col_index + 1
	}

	var max_product_sku_count = 0

	for _, order_logistics := range order_logistics_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 仓库
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics.Stock.Name

		// 物流公司
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics.LogisticsCompany.Name

		// 物流方式
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = order_logistics.LogisticsMode.Name

		// 长(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Length)

		// 宽(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Width)

		// 高(cm)
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Height)

		// 实际重量（Kg）
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)

		// 邮政编码
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.PostalCode)

		// 客户姓名
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.Name)

		// 电话
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.Phone)

		// 客户email
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.Order.BuyerEmail)

		// 邮寄地址
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.Street1 + " " + order_logistics.OrderShippingAddress.Street2)

		// 所属城市
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.CityName)

		// 所属地区
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.Province)

		// 国家
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.OrderShippingAddress.CountryName)

		// 订单代码
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(order_logistics.ShopAccount.Code + order_logistics.Order.TranId)

		// 店铺名
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("FROM:%s", order_logistics.ShopAccount.Name)

		var product_detail_list = c.get_order_logistics_product_detail_list(order_logistics.Order.Id, order_logistics.OrderLogistics.StockId)
		var product_detail_list_len = len(product_detail_list)

		var set_even_product_detail_background = func(cell *xlsx.Cell) {
			style := cell.GetStyle()

			fill := *xlsx.NewFill("solid", "00FFFF00", "00000000") //CCFFCC
			style.Fill = fill
			style.ApplyFill = true

			cell.SetStyle(style)
		}

		//设置标题
		if product_detail_list_len > max_product_sku_count {
			var header_row = sheet.Rows[0]
			for i := 0; i < product_detail_list_len-max_product_sku_count; i++ {
				var col_index = col_count_without_skus + max_product_sku_count*3
				cell = header_row.AddCell()
				sheet.SetColWidth(col_index, col_index, 10)
				set_head_cell_style(cell)
				cell.Value = "仓位"
				if product_detail_list_len%2 == 0 {
					set_even_product_detail_background(cell)
				}

				col_index++
				cell = header_row.AddCell()
				sheet.SetColWidth(col_index, col_index, 20)
				set_head_cell_style(cell)
				cell.Value = "商品编号"
				if product_detail_list_len%2 == 0 {
					set_even_product_detail_background(cell)
				}

				col_index++
				cell = header_row.AddCell()
				sheet.SetColWidth(col_index, col_index, 10)
				set_head_cell_style(cell)
				cell.Value = "商品数量"
				if product_detail_list_len%2 == 0 {
					set_even_product_detail_background(cell)
				}
			}
			max_product_sku_count = product_detail_list_len
		}

		// 设置文字
		for index, product_detail := range product_detail_list {
			// 仓位
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = product_detail.StorageLocName
			if (index+1)%2 == 0 {
				set_even_product_detail_background(cell)
			}

			// 商品编号
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = product_detail.ProductSku
			if (index+1)%2 == 0 {
				set_even_product_detail_background(cell)
			}

			// 商品数量
			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.Value = fmt.Sprint(product_detail.Quantity)
			if (index+1)%2 == 0 {
				set_even_product_detail_background(cell)
			}
		}

	}

	return xlsx_file
}

type order_logistics_product_detail_t struct {
	StorageLocName string `xorm:"'storage_loc_name'" json:"storage_loc_name"`
	ProductSku     string `xorm:"'product_sku'" json:"product_sku"`
	Quantity       int64  `xorm:"'quantity'" json:"quantity"`
}

func (c ErpOrderLogistics) get_order_logistics_product_detail_list(order_id int64, stock_id int64) []order_logistics_product_detail_t {
	var db_session = c.DbSession

	/*
SELECT sl.name AS storage_loc_name, p.sku AS product_sku, oi.quantity
FROM
    erp_order_item oi
    INNER JOIN erp_inventory nvt ON oi.item_product_id = nvt.product_id AND nvt.stock_id = 0
    LEFT JOIN erp_storage_loc sl ON nvt.stock_id = sl.stock_id AND nvt.storage_loc_id = sl.id
    INNER JOIN erp_product p ON oi.item_product_id = p.id
WHERE oi.order_id = 0
ORDER BY oi.transaction_price * oi.quantity DESC
	 */

	var order_logistics_product_detail_list = make([]order_logistics_product_detail_t, 0)
	sql := `SELECT sl.name AS storage_loc_name, p.sku AS product_sku, oi.quantity
FROM
    erp_order_item oi
    INNER JOIN erp_inventory nvt ON oi.item_product_id = nvt.product_id AND nvt.stock_id = %d
    LEFT JOIN erp_storage_loc sl ON nvt.stock_id = sl.stock_id AND nvt.storage_loc_id = sl.id
    INNER JOIN erp_product p ON oi.item_product_id = p.id
WHERE oi.order_id = %d
ORDER BY oi.transaction_price * oi.quantity DESC`

	sql = fmt.Sprintf(sql, stock_id, order_id)

	err := db_session.Sql(sql).Find(&order_logistics_product_detail_list)
	core.HandleError(err)

	return order_logistics_product_detail_list
}
