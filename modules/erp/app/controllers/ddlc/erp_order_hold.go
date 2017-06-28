package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/setting"
	"strings"
	"time"
)

type ErpOrderHold struct {
	*revel.Controller
}

func (c ErpOrderHold) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HOLD, c.Session)

	var has_pass_order_quxiao = pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_QUXIAO, c.Session)

	var is_button_quxiao_visible = false
	if has_pass_order_quxiao == true {
		is_button_quxiao_visible = true
	}

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	c.RenderArgs["is_button_quxiao_visible"] = is_button_quxiao_visible
	return c.RenderTemplate("erp/ddlc/order_hold/order_hold_index.html")
}

func (c ErpOrderHold) PageData() revel.Result {
	var db_session = c.DbSession

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(models.ORDER_STATUS_HOLD, db_session)

	// hold类型列表
	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	err := db_session.OrderBy("ord ASC").Find(&hold_type_list)
	core.HandleError(err)

	var au_stock_list = make([]models.StockInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_type_id = %d", models.STOCK_TYPE_AU)).Find(&au_stock_list)
	core.HandleError(err)

	var order_tag_list = make([]models.OrderTagInfo, 0)
	err = db_session.OrderBy("name ASC").Find(&order_tag_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"batch_code":      0,
		"batch_code_list": batch_code_list,
		"hold_type_id":    0,
		"hold_type_list":  hold_type_list,
		"au_stock_list":   au_stock_list,
		"order_tag_list":  order_tag_list,
	})
}

func (c ErpOrderHold) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HOLD, c.Session)

	db_session := c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var order_code string
	c.Params.Bind(&order_code, "order_code")
	order_code = strings.TrimSpace(order_code)

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id")

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "query_begin_time")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	c.Params.Bind(&query_end_time, "query_end_time")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if query_begin_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
	}

	if query_end_date_err == nil {
		query_end_date = query_end_date.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time <= '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
	}

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_HOLD))

	if hold_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
	}

	if order_code != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.code = '%s'", order_code))
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
