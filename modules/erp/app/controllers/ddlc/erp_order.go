package controllers

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
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
		active_menu_html = fmt.Sprintf(active_menu_html, status, "所有订单")
		status_chinese_name = "所有订单"
	} else if status == "yzf" {
		active_menu_html = fmt.Sprintf(active_menu_html, status, "已支付")
		status_chinese_name = "已支付"
	} else if status == "dph" {
		active_menu_html = fmt.Sprintf(active_menu_html, status, "待配货")
		status_chinese_name = "待配货"
	} else if status == "phz" {
		active_menu_html = fmt.Sprintf(active_menu_html, status, "配货中")
		status_chinese_name = "配货中"
	} else if status == "yfh" {
		active_menu_html = fmt.Sprintf(active_menu_html, status, "已发货")
		status_chinese_name = "已发货"
	} else if status == "ywc" {
		active_menu_html = fmt.Sprintf(active_menu_html, status, "已完成")
		status_chinese_name = "已完成"
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
	db_session := c.DbSession

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id") // -1: all hold; 0: none hold; other: 具体的 hold_type_id

	var status string //all，yzf, dph, phz, yfh, ywc
	c.Params.Bind(&status, "status")

	var status_id int64 = 0 // 0： all status; other: 具体的 status_id

	// 先看 URL 里面的状态分类
	status = strings.ToLower(status) //大小写
	if status == "all" {
		status_id = 0
	} else if status == "yzf" {
		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		status_id = models.ORDER_STATUS_YWC
	}

	// 再看查询条件里面的状态分类
	c.Params.Bind(&status_id, "status_id") //查询条件的status_id，比url里面换算过来的status_id更为优先

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

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

	// todo: hold 订单是不是也应该可以看成是一个状态，状态中再加一个hold状态
	if hold_type_id == 0 { //取不是hold的订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NULL")
	} else if hold_type_id == -1 { //取hold的所有订单
		filter = core.SqlAnd(filter, "o.hold_type_id IS NOT NULL")
	} else { //取得具体 hold 的订单类型
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}

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
	c.RenderArgs["static_version"] = core.GetStaticVersion()
	c.RenderArgs["order_id"] = order_id
	return c.RenderTemplate("erp/ddlc/order/order_detail.html")
}

func (c ErpOrder) DetailData() revel.Result {
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
    INNER JOIN erp_product_store_item psi ON psi.id = oi.product_store_item_id
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
