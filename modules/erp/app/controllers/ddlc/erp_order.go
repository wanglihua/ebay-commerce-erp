package controllers

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/paypal_service"
	"matrix/modules/erp/service/shop_service"
)

type ErpOrder struct {
	*revel.Controller
}

func (c ErpOrder) Index(status string) revel.Result {
	status = strings.ToLower(status) //大小写
	// 已创建 已支付 待配货 配货中 已发货 已完成
	active_menu_html := `<li class="active" data-menu="ddlc-order-%s-menu">%s</li>`
	status_chinese_name := ""
	if status == "all" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_ALL, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "所有订单")
		status_chinese_name = "所有订单"
	} else if status == "ycj" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YCJ, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "未支付")
		status_chinese_name = "未支付"
	} else if status == "yzf" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已支付")
		status_chinese_name = "已支付"
	} else if status == "dph" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "待配货")
		status_chinese_name = "待配货"
	} else if status == "phz" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "配货中")
		status_chinese_name = "配货中"
	} else if status == "yfh" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已发货")
		status_chinese_name = "已发货"
	} else if status == "ywc" {
		pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YWC, c.Session)

		active_menu_html = fmt.Sprintf(active_menu_html, status, "已完成")
		status_chinese_name = "已完成"
	} else {
		panic("未知的订单状态！")
	}

	c.RenderArgs["active_menu_html"] = template.HTML(active_menu_html)
	c.RenderArgs["status"] = status
	c.RenderArgs["status_chinese_name"] = status_chinese_name
	return c.RenderTemplate("erp/ddlc/order/order_index.html")
}

const (
	order_query_criteria_order_code    int64 = 1000
	order_query_criteria_product_sku   int64 = 2000
	order_query_criteria_product_name  int64 = 3000
	order_query_criteria_buyer_user_id int64 = 4000
)

type order_query_criteria struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// 查询条件的列表
func (c ErpOrder) get_order_query_criteria_list() []order_query_criteria {
	var order_query_criteria_list = make([]order_query_criteria, 0)

	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_order_code, Name: "订单编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_sku, Name: "商品编号"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_product_name, Name: "商品名称"})
	order_query_criteria_list = append(order_query_criteria_list, order_query_criteria{Id: order_query_criteria_buyer_user_id, Name: "买家账号"})

	return order_query_criteria_list
}

// index page data
func (c ErpOrder) PageData() revel.Result {
	//这个的权限就不用控制了

	db_session := c.DbSession

	shop_list := shop_service.Get_shop_select_item_list(db_session)

	var order_type_list = make([]models.OrderTypeInfo, 0)
	err := db_session.Find(&order_type_list)
	core.HandleError(err)

	var status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_ORDER)).Find(&status_list)
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
	})
}

func (c ErpOrder) ListData() revel.Result {
	// 权限控制是分流程状态的，在下面的状态检查里面做控制了
	db_session := c.DbSession

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	/* Hold 订单有单独的界面，不用在这里考虑 ??
	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id") // -1: all hold; 0: none hold; other: 具体的 hold_type_id
	*/

	var status string //all，yzf, dph, phz, yfh, ywc
	c.Params.Bind(&status, "status")

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
	} else {
		panic("未知的订单状态！")
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	// 分仓库查看
	// 只有待配货，配货中，已发货，已完成的订单可以分仓库查看
	if status_id == models.ORDER_STATUS_DPH || status_id == models.ORDER_STATUS_PHZ || status_id == models.ORDER_STATUS_YFH || status_id == models.ORDER_STATUS_YWC {
		// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
		/*
		SELECT *
		FROM erp_order o
		WHERE o.id IN (
		SELECT ol.order_id
		FROM erp_order_logistics ol
		WHERE ol.stock_id IN (
			SELECT spc.stock_id
			FROM erp_stock_pick_crew spc
			WHERE spc.user_id = 0))
		 */
		var login_user = core.GetLoginUser(c.Session)
		var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
		var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
		var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服
		if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
			filter = core.SqlAnd(filter, fmt.Sprintf(`o.id IN (
    SELECT ol.order_id
    FROM erp_order_logistics ol
    WHERE ol.stock_id IN (
        SELECT spc.stock_id
        FROM erp_stock_pick_crew spc
        WHERE spc.user_id = %d))`, login_user.UserId))
		}
	}

	//todo: 界面上显示配货的仓库信息， left join

	if status_id == 0 { // 后来改成只有在所有订单的的界面下，才有状态的下拉框
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

	var order_type_id int64 = 0
	c.Params.Bind(&order_type_id, "order_type_id")
	if order_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("type_id = %d", order_type_id))
	}

	var query_criteria_id int64 = 0
	c.Params.Bind(&query_criteria_id, "query_criteria_id")
	var query_criteria_value string = ""
	c.Params.Bind(&query_criteria_value, "query_criteria_value")
	if query_criteria_id != 0 {
		if query_criteria_id == order_query_criteria_order_code {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.code like '%%%s%%'", query_criteria_value))
		} else if query_criteria_id == order_query_criteria_product_sku {
			order_id_list := order_service.Get_order_id_list_by_product_sku_query(db_session, query_criteria_value)
			order_id_str_list := core.Int_list_to_str_list(order_id_list)
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", strings.Join(order_id_str_list, ",")))
		} else if query_criteria_id == order_query_criteria_product_name {
			order_id_list := order_service.Get_order_id_list_by_product_name_query(db_session, query_criteria_value)
			order_id_str_list := core.Int_list_to_str_list(order_id_list)
			filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (%s)", strings.Join(order_id_str_list, ",")))
		} else if query_criteria_id == order_query_criteria_buyer_user_id {
			filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_user_id like '%%%s%%'", query_criteria_value))
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

	if status_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", status_id))
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}

func (c ErpOrder) Detail(order_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	c.RenderArgs["order_id"] = order_id
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

	// 这里面取得是平台的地址
	var shipping_address models.OrderShippingAddressInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d AND address_owner = 'eBay'", order_id)).Limit(1, 0).Get(&shipping_address)
	core.HandleError(err)

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

	return c.RenderJson(map[string]interface{}{
		"order":            order,
		"status":           status,
		"shop":             shop,
		"order_type":       order_type,
		"shipping_address": shipping_address,
		"order_hold_type":  order_hold_type,
		"order_item_list":  order_item_list,
	})
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
