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
	"matrix/modules/erp/service/order_history_service"
	"matrix/modules/erp/models/sppz_models"
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

func (c ErpOrderHold) get_time_query_list() []time_query {
	var time_query_list = make([]time_query, 0)

	time_query_list = append(time_query_list, time_query{Id: time_query_paid_time, Name: "支付时间"})
	time_query_list = append(time_query_list, time_query{Id: time_query_ship_time, Name: "发货时间"})

	return time_query_list
}

func (c ErpOrderHold) PageData() revel.Result {
	var db_session = c.DbSession

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(models.ORDER_STATUS_HOLD, db_session)

	// hold类型列表
	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	err := db_session.OrderBy("ord ASC").Find(&hold_type_list)
	core.HandleError(err)

	var au_stock_list = make([]sppz_models.StockInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&au_stock_list)
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
		"order_tag_id":    0,
		"order_tag_list":  order_tag_list,
		"time_query_id":   time_query_paid_time,
		"time_query_list": c.get_time_query_list(),
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

	var order_tag_id int64
	c.Params.Bind(&order_tag_id, "order_tag_id")

	var time_query_id int64
	c.Params.Bind(&time_query_id, "time_query_id")

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "query_begin_time")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	c.Params.Bind(&query_end_time, "query_end_time")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if time_query_id == time_query_paid_time {
		if query_begin_date_err == nil {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
		}

		if query_end_date_err == nil {
			query_end_date = query_end_date.Add(time.Hour * 24)
			filter = core.SqlAnd(filter, fmt.Sprintf("o.paid_time < '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
		}
	}

	if time_query_id == time_query_ship_time {
		if query_begin_date_err == nil {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
		}

		if query_end_date_err == nil {
			query_end_date = query_end_date.Add(time.Hour * 24)
			filter = core.SqlAnd(filter, fmt.Sprintf("o.shipped_time < '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
		}
	}

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_HOLD))

	if hold_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	} else { // 如果没选Hold类型，就不显示那四个特殊个Hold类型
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id <> %d AND o.hold_type_id <> %d AND o.hold_type_id <> %d AND o.hold_type_id <> %d", models.ORDER_HOLD_TYPE_INTERNATIONAL, models.ORDER_HOLD_TYPE_EXPRESS, models.ORDER_HOLD_TYPE_PICK_UP, models.ORDER_HOLD_TYPE_ENVELOPE))
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
	}

	if order_code != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.code = '%s'", order_code))
	}

	if order_tag_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.tag_id = %d", order_tag_id))
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}

// 订单从 Hold 状态 搬至 已处理Hold 状态
// 但 信封、快递、国际件、pick up 四个 Hold 类型除外
func (c ErpOrderHold) Processed() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HOLD, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_HOLD ||
			( order.StatusId == models.ORDER_STATUS_HOLD &&
				(order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_ENVELOPE ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_EXPRESS ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_INTERNATIONAL ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态或Hold类型不能搬至已处理Hold状态！",
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.PreHoldTypeId = order.HoldTypeId.Int64
		order.StatusId = models.ORDER_STATUS_PROCESSED_HOLD
		_, err = db_session.Id(order.Id).Cols("pre_hold_type_id", "status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单 Hold 转 已处理Hold， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单搬至已处理Hold成功！",
	})
}
