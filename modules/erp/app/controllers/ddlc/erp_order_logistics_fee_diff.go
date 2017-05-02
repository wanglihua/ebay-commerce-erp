package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/setting"
	"matrix/core/collections"
	"matrix/modules/erp/service/logistics_company_service"
	"matrix/modules/erp/service/logistics_mode_service"
)

type ErpOrderLogisticsFeeDiff struct {
	*revel.Controller
}

func (c ErpOrderLogisticsFeeDiff) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddlc/order_logistics_fee_diff/order_logistics_fee_diff_index.html")
}

func (c ErpOrderLogisticsFeeDiff) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

	var db_session = c.DbSession

	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	err := db_session.OrderBy("ord asc").Find(&logistics_company_list)
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

func (c ErpOrderLogisticsFeeDiff) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF, c.Session)

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

	type order_logistics_view_item_t struct {
		OrderLogistics models.OrderLogisticsInfo  `xorm:"extends" json:"ol"`
		Order          models.OrderInfo           `xorm:"extends" json:"o"`
		StockName      string                     `json:"stock_name"`
		CompanyName    string                     `json:"company_name"`
		ModeName       string                     `json:"mode_name"`
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
		sb.OrderBy("id ASC")
	}

	sb.Limit(offset, limit)

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
