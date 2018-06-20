package controllers

import (
	"fmt"
	"matrix/modules/erp/service/order_feedback_service"
	"strings"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"time"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/core"
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/service/stock_pick_crew_service"
	"matrix/modules/erp/service/order_service"
)

var special_hold_type_id_list = []int64{ddpz_models.ORDER_HOLD_TYPE_ENVELOPE, ddpz_models.ORDER_HOLD_TYPE_EXPRESS, ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL, ddpz_models.ORDER_HOLD_TYPE_PICK_UP}

// 查询条件定义 开始
const (
	order_query_criteria_order_code                 int64 = 1000
	order_query_criteria_order_number               int64 = 1100
	order_query_criteria_product_sku                int64 = 1200
	order_query_criteria_product_name               int64 = 1300
	order_query_criteria_buyer_user_id              int64 = 1400
	order_query_criteria_stock_name                 int64 = 1500
	order_query_criteria_logistics_mode_name        int64 = 1600
	order_query_criteria_waybill_number             int64 = 1700
	order_query_criteria_buyer_email                int64 = 1800
	order_query_criteria_buyer_email_accurate       int64 = 1900
	order_query_criteria_buyer_phone                int64 = 2000
	order_query_criteria_buyer_phone_accurate       int64 = 2100
	order_query_criteria_buyer_name                 int64 = 2200
	order_query_criteria_buyer_name_accurate        int64 = 2300
	order_query_criteria_buyer_postal_code          int64 = 2400
	order_query_criteria_buyer_postal_code_accurate int64 = 2500

	// accurate
)

type order_query_criteria struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// 查询条件的列表
func get_order_query_criteria_list() []order_query_criteria {
	var order_query_criteria_list = make([]order_query_criteria, 0)

	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_order_code, Name: "订单编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_order_number, Name: "订单代码"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_sku, Name: "商品编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_name, Name: "商品名称"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_user_id, Name: "买家账号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_stock_name, Name: "发货仓库"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_logistics_mode_name, Name: "物流方式"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_waybill_number, Name: "运输单号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_email, Name: "买家邮箱"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_email_accurate, Name: "买家邮箱（精确）"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_phone, Name: "买家电话"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_phone_accurate, Name: "买家电话（精确）"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_name, Name: "买家姓名"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_name_accurate, Name: "买家姓名（精确）"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_postal_code, Name: "买家邮编"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_postal_code_accurate, Name: "买家邮编（精确）"})

	return order_query_criteria_list
}

// 查询条件定义  结束

// 优化查询定义 开始
const (
	good_query_ss                       int64 = 1000 // 单商品单数量
	good_query_sm                       int64 = 2000 // 单商品多数量
	good_query_mm                       int64 = 3000 // 多商品多数量
	good_query_empty_waybill_number     int64 = 4000 // 运单号为空
	good_query_not_empty_waybill_number int64 = 5000 // 运单号不为空
	good_query_hold_order               int64 = 6000 // hold订单
	good_query_feedback_negative        int64 = 7000 // 差评订单
)

type good_query struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// 查询条件的列表
func get_good_query_list() []good_query {
	var good_query_list = make([]good_query, 0)

	good_query_list = append(good_query_list, good_query{Id: good_query_ss, Name: "单商品单数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_sm, Name: "单商品多数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_mm, Name: "多商品"})
	good_query_list = append(good_query_list, good_query{Id: good_query_empty_waybill_number, Name: "运单号为空"})
	good_query_list = append(good_query_list, good_query{Id: good_query_not_empty_waybill_number, Name: "运单号不为空"})
	good_query_list = append(good_query_list, good_query{Id: good_query_hold_order, Name: "hold订单"})
	good_query_list = append(good_query_list, good_query{Id: good_query_feedback_negative, Name: "差评订单"})

	return good_query_list
}

const (
	time_query_paid_time int64 = 100 // 支付时间
	time_query_ship_time int64 = 200 // 发货时间
)

type time_query struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func get_time_query_list() []time_query {
	var time_query_list = make([]time_query, 0)

	time_query_list = append(time_query_list, time_query{Id: time_query_paid_time, Name: "支付时间"})
	time_query_list = append(time_query_list, time_query{Id: time_query_ship_time, Name: "发货时间"})

	return time_query_list
}

func Get_order_filter(db_session *xorm.Session, web_session revel.Session, params *revel.Params, filter string) string {

	var shop_id int64
	params.Bind(&shop_id, "shop_id")

	var status string //all，yzf, dph, phz, yfh, ywc
	params.Bind(&status, "status")

	var query_criteria_id1 int64 = 0
	params.Bind(&query_criteria_id1, "query_criteria_id1")
	var query_criteria_value1 string = ""
	params.Bind(&query_criteria_value1, "query_criteria_value1")
	query_criteria_value1 = strings.TrimSpace(query_criteria_value1)

	var query_criteria_id2 int64 = 0
	params.Bind(&query_criteria_id2, "query_criteria_id2")
	var query_criteria_value2 string = ""
	params.Bind(&query_criteria_value2, "query_criteria_value2")
	query_criteria_value2 = strings.TrimSpace(query_criteria_value2)

	var good_query_id int64 = 0
	params.Bind(&good_query_id, "good_query_id")

	var order_type_id int64 = 0
	params.Bind(&order_type_id, "order_type_id")

	var special_hold_type_id int64 = 0
	params.Bind(&special_hold_type_id, "special_hold_type_id")

	var hold_type_id int64 = 0
	params.Bind(&hold_type_id, "hold_type_id")

	var batch_code int64
	params.Bind(&batch_code, "batch_code")

	var time_query_id int64
	params.Bind(&time_query_id, "time_query_id")

	var query_begin_time string
	params.Bind(&query_begin_time, "query_begin_time")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	params.Bind(&query_end_time, "query_end_time")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	var order_tag_id int64
	params.Bind(&order_tag_id, "order_tag_id")

	var status_id int64 = 0 // 0： all status; other: 具体的 status_id

	// 先看 URL 里面的状态分类
	status = strings.ToLower(status) //大小写
	if status == "all" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_ALL, web_session)

		status_id = 0
	} else if status == "ycj" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YCJ, web_session)

		status_id = models.ORDER_STATUS_YCJ
	} else if status == "yzf" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, web_session)

		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, web_session)

		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, web_session)

		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, web_session)

		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YWC, web_session)

		status_id = models.ORDER_STATUS_YWC
	} else if status == "yqx" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YQX, web_session)

		status_id = models.ORDER_STATUS_YQX
	} else {
		panic("未知的订单状态！")
	}

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

	// 分仓库查看
	// 只有待配货，配货中，已发货，已完成的订单可以分仓库查看
	if status_id == models.ORDER_STATUS_DPH ||
		status_id == models.ORDER_STATUS_PHZ ||
		status_id == models.ORDER_STATUS_YFH ||
		status_id == models.ORDER_STATUS_YWC ||
		( // pick up 的 hold状态，也是分仓库查看的
			status_id == models.ORDER_STATUS_HOLD &&
				special_hold_type_id == ddpz_models.ORDER_HOLD_TYPE_PICK_UP) {
		// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
		var login_user = core.GetLoginUser(web_session)
		var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
		var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
		var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服
		if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
			filter = core.SqlAnd(filter, fmt.Sprintf(`o.stock_id IN (SELECT spc.stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = %d)`, login_user.UserId))
		}
	}

	if status_id == 0 { //表示所有订单的页面。 后来改成只有在所有订单的的界面下，才有状态的下拉框
		var status_id_query int64 = 0 // 下拉框里面的状态，不是通过 Url 里面来的
		// 再看查询条件里面的状态分类
		params.Bind(&status_id_query, "status_id") //查询条件的status_id，比url里面换算过来的status_id更为优先

		if status_id_query != 0 { // 只有在下拉框的状态为不零的时候，下拉框的状态才起作用
			status_id = status_id_query
		}
	}

	if shop_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shop_id = %d", shop_id))
	}

	if order_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.type_id = %d", order_type_id))
	}

	if query_criteria_id1 != 0 {
		if query_criteria_id1 == order_query_criteria_order_code {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code like '%%%s%%'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_order_number {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.full_tran_id like '%%%s%%'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_product_sku {
			order_id_list := order_service.Get_order_id_list_by_product_sku_query(db_session, query_criteria_value1)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id1 == order_query_criteria_product_name {
			order_id_list := order_service.Get_order_id_list_by_product_name_query(db_session, query_criteria_value1)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id1 == order_query_criteria_buyer_user_id {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_user_id like '%%%s%%'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_stock_name {
			// SELECT * FROM erp_order o WHERE o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_logistics_mode_name {
			// SELECT * FROM erp_order o WHERE o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_waybill_number {
			//waybill_number
			filter = core.SqlAnd(filter, fmt.Sprintf("o.waybill_number like '%%%s%%'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_email {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_email like '%%%s%%'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_email_accurate {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_email = '%s'", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_phone {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_phone_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone = '%s')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_name {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_name_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name = '%s')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_postal_code {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id1 == order_query_criteria_buyer_postal_code_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code = '%s')", query_criteria_value1))
		}
	}

	if query_criteria_id2 != 0 {
		if query_criteria_id2 == order_query_criteria_order_code {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code like '%%%s%%'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_order_number {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.full_tran_id like '%%%s%%'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_product_sku {
			order_id_list := order_service.Get_order_id_list_by_product_sku_query(db_session, query_criteria_value2)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id2 == order_query_criteria_product_name {
			order_id_list := order_service.Get_order_id_list_by_product_name_query(db_session, query_criteria_value2)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id2 == order_query_criteria_buyer_user_id {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_user_id like '%%%s%%'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_stock_name {
			// SELECT * FROM erp_order o WHERE o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%%%s%%')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_logistics_mode_name {
			// SELECT * FROM erp_order o WHERE o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%%%s%%')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_waybill_number {
			//waybill_number
			filter = core.SqlAnd(filter, fmt.Sprintf("o.waybill_number like '%%%s%%'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_email {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_email like '%%%s%%'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_email_accurate {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_email = '%s'", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_phone {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone LIKE '%%%s%%')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_phone_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.phone = '%s')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_name {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name LIKE '%%%s%%')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_name_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.name = '%s')", query_criteria_value2))
		} else if query_criteria_id2 == order_query_criteria_buyer_postal_code {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code LIKE '%123%';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code LIKE '%%%s%%')", query_criteria_value1))
		} else if query_criteria_id2 == order_query_criteria_buyer_postal_code_accurate {
			// SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code = '123';
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT osa.order_id FROM erp_order_shipping_address osa WHERE osa.postal_code = '%s')", query_criteria_value1))
		}
	}

	// hold 订单 有了一个专门的 hold状态， 也有专门的 hold 状态的页面，在这里就不用多考虑了  ??
	/*
	if special_hold_type_id == 0 { //取不是hold的订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NULL")
	} else if special_hold_type_id == -1 { //取hold的所有订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NOT NULL")
	} else { //取得具体 hold 的订单类型
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", special_hold_type_id))
	}
	*/

	if status_id != 0 { // 不是所有订单的时候
		// 如果是待配货状态，除了显示待配货状态的订单外，还要显示特殊 hold 类型的 hold 状态的订单，但如果是只有仓库拣货人员的权限的话，就只能由墨尔本仓库的拣货人员来查看
		if status_id == models.ORDER_STATUS_DPH {

			if can_view_special_hold(web_session, db_session) || special_hold_type_id == ddpz_models.ORDER_HOLD_TYPE_PICK_UP { // 是待配货状态，且也有权限查看 hold, // pickup 的任何人员都可以看到
				// -1 表示不选中任何 hold
				if special_hold_type_id == -1 { // 不选中任何 hold 时，只显示 DPH 的订单
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_DPH))
				} else if special_hold_type_id == 0 { // 0 表示全选全部的 hold
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d AND o.hold_type_id IN (%s)", models.ORDER_STATUS_HOLD, strings.Join(core.Int_list_to_str_list(special_hold_type_id_list), ",")))
				} else if special_hold_type_id > 0 { // 选中了具体的 hold类型
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d AND o.hold_type_id = %d", models.ORDER_STATUS_HOLD, special_hold_type_id))
				}
			} else { // 是待配货状态，但没有权限查看 hold
				filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", status_id))
			}
		} else {
			// 不是 待配货状态
			filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", status_id))
		}
	}

	revel.TRACE.Println("hold_type_id", hold_type_id)
	if status_id == models.ORDER_STATUS_YWC {
		revel.TRACE.Println("hold_type_id", hold_type_id)
		if hold_type_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
		}
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
	}

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
		} else if good_query_id == good_query_empty_waybill_number {
			filter = core.SqlAnd(filter, "(o.waybill_number IS NULL OR o.waybill_number = '')")
		} else if good_query_id == good_query_not_empty_waybill_number {
			filter = core.SqlAnd(filter, "(o.waybill_number IS NOT NULL AND o.waybill_number <> '')")
		} else if good_query_id == good_query_hold_order {
			filter = core.SqlAnd(filter, "(o.hold_type_id IS NOT NULL AND o.hold_type_id <> 0)")
		} else if good_query_id == good_query_feedback_negative {
			/*
SELECT of.order_id
FROM erp_order_feedback of
WHERE of.comment_type = 1;
			 */

			var good_query = `
o.id IN (
    SELECT of.order_id
    FROM erp_order_feedback of
    WHERE of.comment_type = %d  
)`
			good_query = fmt.Sprintf(good_query, order_feedback_service.ORDER_FEEDBACK_NEGATIVE)

			filter = core.SqlAnd(filter, good_query)
		}
	}

	if order_tag_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.tag_id = %d", order_tag_id))
	}

	return filter
}

func can_view_special_hold(web_session revel.Session, db_session *xorm.Session) bool {
	var login_user = core.GetLoginUser(web_session)
	var is_pick_crew_has_mel_stock = stock_pick_crew_service.Is_pick_crew_has_mel_stock(login_user.UserId, db_session)

	var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
	var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
	var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服

	// 如果是待配货状态，除了显示待配货状态的订单外，还要显示特殊 hold 类型的 hold 状态的订单，但如果是只有仓库拣货人员的权限的话，就只能由墨尔本仓库的拣货人员来查看

	if login_user.IsSuperAdmin || is_admin || is_custom_service || is_sales || is_pick_crew_has_mel_stock {
		return true
	}

	return false
}
