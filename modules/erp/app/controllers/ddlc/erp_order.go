package controllers

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/order_logistics_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/paypal_service"
	"matrix/modules/erp/service/shop_service"
	"matrix/modules/erp/service/stock_pick_crew_service"
	"matrix/modules/erp/setting"
)

type ErpOrder struct {
	*revel.Controller
}

var special_hold_type_id_list = []int64{models.ORDER_HOLD_TYPE_ENVELOPE, models.ORDER_HOLD_TYPE_EXPRESS, models.ORDER_HOLD_TYPE_INTERNATIONAL, models.ORDER_HOLD_TYPE_PICK_UP}

func (c ErpOrder) Index(status string) revel.Result {
	var db_session = c.DbSession

	var has_pass_order_new = false

	status = strings.ToLower(status) //大小写
	// 已创建 已支付 待配货 配货中 已发货 已完成
	active_menu_html := `<li class="active" data-menu="ddlc-order-%s-menu">%s</li>`
	status_chinese_name := ""
	var status_id int64 = 0
	if status == "all" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_ALL, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "所有订单")
		status_chinese_name = "所有订单"
		status_id = 0
	} else if status == "ycj" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YCJ, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "未支付")
		status_chinese_name = "未支付"
		status_id = models.ORDER_STATUS_YCJ
	} else if status == "yzf" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, c.Session)

		if pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_NEW, c.Session) {
			has_pass_order_new = true
		}

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已支付")
		status_chinese_name = "已支付"
		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "待配货")
		status_chinese_name = "待配货"
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "配货中")
		status_chinese_name = "配货中"
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已发货")
		status_chinese_name = "已发货"
		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YWC, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已完成")
		status_chinese_name = "已完成"
		status_id = models.ORDER_STATUS_YWC
	} else if status == "yqx" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YQX, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已取消")
		status_chinese_name = "已取消"
		status_id = models.ORDER_STATUS_YQX
	} else {
		panic("未知的订单状态！")
	}

	// 如果是待配货状态，除了显示待配货状态的订单外，还要显示特殊 hold 类型的 hold 状态的订单，但如果是只有仓库拣货人员的权限的话，就只能由墨尔本仓库的拣货人员来查看
	var can_view_special_hold_ = false
	if status_id == models.ORDER_STATUS_DPH {
		can_view_special_hold_ = can_view_special_hold(c.Session, db_session)
	}

	var has_pass_order_quxiao = pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_QUXIAO, c.Session)

	var is_button_quxiao_visible = false
	if has_pass_order_quxiao == true &&
		status_id != models.ORDER_STATUS_YFH &&
		status_id != models.ORDER_STATUS_YWC &&
		status_id != models.ORDER_STATUS_YQX {

		is_button_quxiao_visible = true
	}

	c.RenderArgs["active_menu_html"] = template.HTML(active_menu_html)
	c.RenderArgs["status"] = status
	c.RenderArgs["status_id"] = status_id
	c.RenderArgs["status_chinese_name"] = status_chinese_name
	c.RenderArgs["can_view_special_hold"] = can_view_special_hold_ // 这个后来在页面中没有使用，只是这里暂时没有去掉。更改来之，任何权限的人都可以看到 pick hold
	c.RenderArgs["icon_path"] = setting.ICON_PATH
	c.RenderArgs["has_pass_order_new"] = has_pass_order_new
	c.RenderArgs["is_button_quxiao_visible"] = is_button_quxiao_visible
	return c.RenderTemplate("erp/ddlc/order/order_index.html")
}

// 查询条件定义 开始
const (
	order_query_criteria_order_code          int64 = 1000
	order_query_criteria_order_number        int64 = 1100
	order_query_criteria_product_sku         int64 = 2000
	order_query_criteria_product_name        int64 = 3000
	order_query_criteria_buyer_user_id       int64 = 4000
	order_query_criteria_stock_name          int64 = 5000
	order_query_criteria_logistics_mode_name int64 = 6000
)

type order_query_criteria struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// 查询条件的列表
func (c ErpOrder) get_order_query_criteria_list() []order_query_criteria {
	var order_query_criteria_list = make([]order_query_criteria, 0)

	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_order_code, Name: "订单编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_order_number, Name: "订单代码"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_sku, Name: "商品编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_name, Name: "商品名称"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_user_id, Name: "买家账号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_stock_name, Name: "发货仓库"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_logistics_mode_name, Name: "物流方式"})

	return order_query_criteria_list
}

// 查询条件定义  结束

// 优化查询定义 开始
const (
	good_query_ss int64 = 1000 // 单商品单数量
	good_query_sm int64 = 2000 // 单商品多数量
	good_query_mm int64 = 3000 // 多商品多数量
)

type good_query struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// 查询条件的列表
func (c ErpOrder) get_good_query_list() []good_query {
	var good_query_list = make([]good_query, 0)

	good_query_list = append(good_query_list, good_query{Id: good_query_ss, Name: "单商品单数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_sm, Name: "单商品多数量"})
	good_query_list = append(good_query_list, good_query{Id: good_query_mm, Name: "多商品"})

	return good_query_list
}

// 优化查询定义  结束

// index page data
func (c ErpOrder) PageData() revel.Result {
	//这个的权限就不用控制了

	var status_id int64
	c.Params.Bind(&status_id, "status_id")

	db_session := c.DbSession

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(status_id, db_session)

	shop_list := shop_service.Get_shop_select_item_list(db_session)

	var order_type_list = make([]models.OrderTypeInfo, 0)
	err := db_session.Find(&order_type_list)
	core.HandleError(err)

	var status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_ORDER)).Find(&status_list)
	core.HandleError(err)

	// pick up hold任何权限的人都可以看到
	var special_hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	var can_view_special_hold_ = can_view_special_hold(c.Session, db_session)
	if can_view_special_hold_ {
		err = db_session.In("id", special_hold_type_id_list).Find(&special_hold_type_list)
	} else {
		var special_hold_type_id_list_only_pick_up = []int64{models.ORDER_HOLD_TYPE_PICK_UP}
		err = db_session.In("id", special_hold_type_id_list_only_pick_up).Find(&special_hold_type_list)
	}
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"shop_id":             0,
		"shop_list":           shop_list,
		"order_type_id":       0,
		"order_type_list":     order_type_list,
		"status_id":           0,
		"status_list":         status_list,
		"query_criteria_id":   0,
		"query_criteria_list": c.get_order_query_criteria_list(),
		"good_query_id":       0,
		"good_query_list":     c.get_good_query_list(),
		"hold_type_list":      special_hold_type_list,
		"hold_type_id":        -1, // 表示不选择 hold
		"batch_code":          0,
		"batch_code_list":     batch_code_list,
	})
}

func (c ErpOrder) ListData() revel.Result {
	// 权限控制是分流程状态的，在下面的状态检查里面做控制了
	db_session := c.DbSession

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	var status string //all，yzf, dph, phz, yfh, ywc
	c.Params.Bind(&status, "status")

	var query_criteria_id int64 = 0
	c.Params.Bind(&query_criteria_id, "query_criteria_id")
	var query_criteria_value string = ""
	c.Params.Bind(&query_criteria_value, "query_criteria_value")

	var good_query_id int64 = 0
	c.Params.Bind(&good_query_id, "good_query_id")

	var order_type_id int64 = 0
	c.Params.Bind(&order_type_id, "order_type_id")

	var hold_type_id int64 = 0
	c.Params.Bind(&hold_type_id, "hold_type_id")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var status_id int64 = 0 // 0： all status; other: 具体的 status_id

	// 先看 URL 里面的状态分类
	status = strings.ToLower(status) //大小写
	if status == "all" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_ALL, c.Session)

		status_id = 0
	} else if status == "ycj" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YCJ, c.Session)

		status_id = models.ORDER_STATUS_YCJ
	} else if status == "yzf" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, c.Session)

		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, c.Session)

		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, c.Session)

		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YWC, c.Session)

		status_id = models.ORDER_STATUS_YWC
	} else if status == "yqx" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YQX, c.Session)

		status_id = models.ORDER_STATUS_YQX
	} else {
		panic("未知的订单状态！")
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	// 分仓库查看
	// 只有待配货，配货中，已发货，已完成的订单可以分仓库查看
	if status_id == models.ORDER_STATUS_DPH || status_id == models.ORDER_STATUS_PHZ || status_id == models.ORDER_STATUS_YFH || status_id == models.ORDER_STATUS_YWC {
		// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
		var login_user = core.GetLoginUser(c.Session)
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
		c.Params.Bind(&status_id_query, "status_id") //查询条件的status_id，比url里面换算过来的status_id更为优先

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

	if query_criteria_id != 0 {
		if query_criteria_id == order_query_criteria_order_code {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code like '%%%s%%'", query_criteria_value))
		} else if query_criteria_id == order_query_criteria_order_number {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.full_tran_id like '%%%s%%'", query_criteria_value))
		} else if query_criteria_id == order_query_criteria_product_sku {
			order_id_list := order_service.Get_order_id_list_by_product_sku_query(db_session, query_criteria_value)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id == order_query_criteria_product_name {
			order_id_list := order_service.Get_order_id_list_by_product_name_query(db_session, query_criteria_value)

			var order_id_strs = "0"
			if len(order_id_list) != 0 {
				order_id_str_list := core.Int_list_to_str_list(order_id_list)
				order_id_strs = strings.Join(order_id_str_list, ",")
			}
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", order_id_strs))
		} else if query_criteria_id == order_query_criteria_buyer_user_id {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_user_id like '%%%s%%'", query_criteria_value))
		} else if query_criteria_id == order_query_criteria_stock_name {
			// SELECT * FROM erp_order o WHERE o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.stock_id IN (SELECT s.id FROM erp_stock s WHERE s.name LIKE '%%%s%%')", query_criteria_value))
		} else if query_criteria_id == order_query_criteria_logistics_mode_name {
			// SELECT * FROM erp_order o WHERE o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%name%')
			filter = core.SqlAnd(filter, fmt.Sprintf("o.mode_id IN (SELECT lm.id FROM erp_logistics_mode lm WHERE lm.name LIKE '%%%s%%')", query_criteria_value))
		}
	}

	// hold 订单 有了一个专门的 hold状态， 也有专门的 hold 状态的页面，在这里就不用多考虑了  ??
	/*
	if hold_type_id == 0 { //取不是hold的订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NULL")
	} else if hold_type_id == -1 { //取hold的所有订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NOT NULL")
	} else { //取得具体 hold 的订单类型
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}
	*/

	if status_id != 0 { // 不是所有订单的时候
		// 如果是待配货状态，除了显示待配货状态的订单外，还要显示特殊 hold 类型的 hold 状态的订单，但如果是只有仓库拣货人员的权限的话，就只能由墨尔本仓库的拣货人员来查看
		if status_id == models.ORDER_STATUS_DPH {

			if can_view_special_hold(c.Session, db_session) { // 是待配货状态，且也有权限查看 hold
				// -1 表示不选中任何 hold
				if hold_type_id == -1 { // 不选中任何 hold 时，只显示 DPH 的订单
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_DPH))
				} else if hold_type_id == 0 { // 0 表示全选全部的 hold
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d AND o.hold_type_id IN (%s)", models.ORDER_STATUS_HOLD, strings.Join(core.Int_list_to_str_list(special_hold_type_id_list), ",")))
				} else if hold_type_id > 0 { // 选中了具体的 hold类型
					filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d AND o.hold_type_id = %d", models.ORDER_STATUS_HOLD, hold_type_id))
				}
			} else { // 是待配货状态，但没有权限查看 hold
				filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", status_id))
			}
		} else {
			// 不是 待配货状态
			filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", status_id))
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
		}
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
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

func (c ErpOrder) Detail(order_id int64) revel.Result {
	var db_session = c.DbSession

	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var order models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	c.RenderArgs["order_id"] = order_id
	c.RenderArgs["is_hold"] = order.StatusId == models.ORDER_STATUS_HOLD
	c.RenderArgs["is_dph"] = order.StatusId == models.ORDER_STATUS_DPH
	return c.RenderTemplate("erp/ddlc/order/order_detail.html")
}

func (c ErpOrder) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "id")

	revel.TRACE.Println("order_id:", order_id)

	var order models.OrderInfo
	has_order, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)
	if has_order == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的订单不存在！",
		})
	}

	var status models.StatusInfo
	_, err = db_session.Id(order.StatusId).Get(&status)
	core.HandleError(err)

	var shop models.ShopAccountInfo
	_, err = db_session.Id(order.ShopId).Get(&shop)
	core.HandleError(err)

	var order_type models.OrderTypeInfo
	_, err = db_session.Id(order.TypeId).Get(&order_type)
	core.HandleError(err)

	var order_hold_type models.OrderHoldTypeInfo
	_, err = db_session.Id(order.HoldTypeId).Get(&order_hold_type)
	core.HandleError(err)

	var shipping_address models.OrderShippingAddressInfo // 取 ebay_shipping_address 中的一个， 优先取 paypal_shipping_address

	// 这里面取得是平台的地址  eBay
	var ebay_shipping_address models.OrderShippingAddressInfo
	has, err := db_session.Where(fmt.Sprintf("order_id = %d AND address_owner='eBay'", order_id)).OrderBy("update_time DESC").Limit(1, 0).Get(&ebay_shipping_address)
	core.HandleError(err)
	if has == true {
		shipping_address = ebay_shipping_address
	}

	// 这里面取得是PayPal的地址  PayPal
	var paypal_shipping_address models.OrderShippingAddressInfo
	has, err = db_session.Where(fmt.Sprintf("order_id = %d AND address_owner='PayPal'", order_id)).OrderBy("update_time DESC").Limit(1, 0).Get(&paypal_shipping_address)
	core.HandleError(err)
	if has == true {
		shipping_address = paypal_shipping_address
	}

	shipping_address.OrderId = order.Id //如果地址是空的，也先给一个 order_id

	type order_item_view_t struct {
		OrderItem         models.OrderItemInfo          `xorm:"extends" json:"oi"`
		Product           models.ProductInfo            `xorm:"extends" json:"p"`
		ProductSotoreItem models.ProductStoreItemInfo   `xorm:"extends" json:"psi"`
	}

	sql := `
SELECT oi.*, p.*, psi.*
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
    LEFT JOIN erp_product_store_item psi ON psi.id = oi.product_store_item_id
WHERE oi.order_id = %d;`
	sql = fmt.Sprintf(sql, order_id)

	order_item_list := make([]order_item_view_t, 0)
	err = db_session.Sql(sql).Find(&order_item_list)
	core.HandleError(err)

	var inventory_list = make([]inventory_service.Product_inventory_view_item_t, 0)
	for _, order_item := range order_item_list {
		var inventory_single_product_list = inventory_service.Get_product_inventory_list(db_session, order_item.Product.Id)
		inventory_list = append(inventory_list, inventory_single_product_list...)
	}

	var order_logistics models.OrderLogisticsInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order_id)).Get(&order_logistics)
	core.HandleError(err)

	var logistics_mode models.LogisticsModeInfo
	if !(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid == true &&
		(order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_ENVELOPE ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_INTERNATIONAL ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS)) {
		//如果是Hold状态且是信封Hold的就不显示匹配的物流，因为这些信息可以是前一次匹配时遗留的，后来指定了 Hold
		_, err = db_session.Id(order_logistics.ModeId).Get(&logistics_mode)
		core.HandleError(err)
	}

	var logistics_stock models.StockInfo
	if !(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid == true &&
		(order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_ENVELOPE ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_INTERNATIONAL ||
			order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS)) {
		//如果是Hold状态且是信封Hold的就不显示匹配的物流，因为这些信息可以是前一次匹配时遗留的，后来指定了 Hold
		_, err = db_session.Id(order_logistics.StockId).Get(&logistics_stock)
		core.HandleError(err)
	}

	var logistics_setting models.LogisticsSettingInfo
	_, err = db_session.Limit(1, 0).Get(&logistics_setting)
	core.HandleError(err)

	var order_cost_amount = order_logistics.CostAmount * order.ExchangeRate

	// 物流费用
	var logistics_fee float64 = (order_logistics.TotalFee + order_logistics.PickFee) * order.ExchangeRate

	// gst 消费税
	var order_gst float64 = order_service.Get_order_gst_fee(order.AmountPaid)

	// 毛利

	var order_logistics_total_fee = order_logistics.TotalFee
	if order_logistics.ActualFee > 0 {
		order_logistics_total_fee = order_logistics.ActualFee
	}

	// var profit = (total_amount + express_fee) - (order_gst + platform_fee + payment_fee + logistics_fee + pick_fee + cost_amount)
	var order_profit = ((order.AmountPaid + order.ExpressFee) - (order_gst + order.AmountPaid*(logistics_setting.Ptfl+logistics_setting.Jyfl) +
		order_logistics_total_fee + order_logistics.PickFee + order_logistics.CostAmount)) * order.ExchangeRate

	var au_stock_list = make([]models.StockInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_type_id = %d", models.STOCK_TYPE_AU)).Find(&au_stock_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"order":                   order,
		"status":                  status,
		"shop":                    shop,
		"order_type":              order_type,
		"ebay_shipping_address":   ebay_shipping_address,
		"paypal_shipping_address": paypal_shipping_address,
		"shipping_address":        shipping_address,
		"order_hold_type":         order_hold_type,
		"order_item_list":         order_item_list,
		"inventory_list":          inventory_list,
		"order_logistics":         order_logistics,
		"logistics_setting":       logistics_setting,
		"order_cost_amount":       order_cost_amount, // rmb
		"logistics_fee":           logistics_fee,     // rmb
		"order_gst":               order_gst,         // aud
		"order_profit":            order_profit,      //rmb
		"au_stock_list":           au_stock_list,
		"logistics_stock":         logistics_stock,
		"logistics_mode":          logistics_mode,
	})
}

func (c ErpOrder) ReMatchHoldOrderLogistics() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.DDLC_ORDER_VIEW_DETAIL, pass_res.DDLC_ORDER_HOLD}, c.Session)

	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	var prefer_stock_id int64
	c.Params.Bind(&prefer_stock_id, "prefer_stock_id")

	var prefer_mode_id int64
	c.Params.Bind(&prefer_mode_id, "prefer_mode_id")

	var unlock_express int64
	c.Params.Bind(&unlock_express, "unlock_express")

	var is_unlock_express bool
	if unlock_express <= 0 {
		is_unlock_express = false
	} else {
		is_unlock_express = true
	}

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_HOLD && order.StatusId != models.ORDER_STATUS_DPH {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "当前状态不是Hold状态或待配货状态，不能重新匹配物流!",
			})
		}
	}

	var p_order_logistics = new(order_logistics_service.OrderLogistics)
	var match_option order_logistics_service.Match_option_t
	match_option.Prefer_stock_id = prefer_stock_id
	match_option.Prefer_mode_id = prefer_mode_id
	match_option.Unlock_express = is_unlock_express
	err = p_order_logistics.CalculateListWithOption(order_list, match_option, db_session)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "重新匹配物流完成！",
	})
}

func (c ErpOrder) Save() revel.Result {
	// 如果是新增保存，就检查有没有订单新增的权限
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_NEW, c.Session)

	return c.RenderEmpty()
}

func (c ErpOrder) GetPaypalShippingAddress() revel.Result {
	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	shipping_address := paypal_service.Get_order_paypal_shipping_address(db_session, order_id)

	return c.RenderJson(map[string]interface{}{
		"shipping_address": shipping_address,
	})
}

func (c ErpOrder) MatchLogistics() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, c.Session)

	var db_session = c.DbSession

	var p_order_logistics = new(order_logistics_service.OrderLogistics)
	err := p_order_logistics.CalculateListPaidOrders(db_session)

	// core.HandleError(err)
	if err != nil {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "匹配物流出现错误！" + err.Error(),
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "已支付订单匹配物流完成！",
	})
}

// 订单从 Hold 状态 搬至 待配货 状态
// 但 信封、快递、国际件、pick up 四个 Hold 类型除外
func (c ErpOrder) DaiPeiHuo() revel.Result {
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
				"message": "订单当前状态或Hold类型不能搬至待配货状态！",
			})
		}
	}

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_DPH
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单搬至待配货成功！",
	})
}

func (c ErpOrder) PeiHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_DPH &&
			(!(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid &&
				(order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_ENVELOPE ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_EXPRESS ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_INTERNATIONAL ||
					order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_PICK_UP))) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能配货！",
			})
		}
	}

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_PHZ
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单配货成功！",
	})
}

func (c ErpOrder) FaHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或hold，不能发货！",
			})
		}
	}

	for _, order := range order_list {

		inventory_service.Take_order_inventory(order.Id, db_session)

		order.StatusId = models.ORDER_STATUS_YFH
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	// 更新库存待配货数量
	go inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单发货成功！",
	})
}

func (c ErpOrder) ShiJiWuHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或hold，不能设置实际无货！",
			})
		}
	}

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		order.HoldTypeId = core.MakeNullInt(models.ORDER_HOLD_TYPE_OOS_ACTUAL_STOCK, true)
		_, err = db_session.Id(order.Id).Cols("status_id", "hold_type_id").Update(&order)
		core.HandleError(err)
	}

	go inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单设置实际无货成功！",
	})
}

func (c ErpOrder) WanCheng() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_YFH {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能设置完成！",
			})
		}
	}

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_YWC
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单完成成功！",
	})
}

// #10 订单作废的功能
func (c ErpOrder) QuXiao() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_QUXIAO, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_YQX // 已取消
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单取消成功！",
	})
}
