package controllers

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
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

func (c ErpOrder) PageData() revel.Result {
	db_session := c.DbSession

	shop_list := get_shop_select_item_list(db_session)
	shop_list = append([]shop_select_item_t{shop_select_item_t{Id: 0, Name: ""}}, shop_list...)

	return c.RenderJson(map[string]interface{}{
		"shop_list": shop_list,
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

	var status_id int64 = 0          // 0： all status; other: 具体的 status_id
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

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if shop_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shop_id = %d", shop_id))
	}

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

	order_list, count := get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}

func (c ErpOrder) DetailData() revel.Result {
	db_session := c.DbSession

	var id int64
	c.Params.Bind(&id, "id")

	var order models.OrderInfo
	has_order, err := db_session.Id(id).Get(&order)
	core.HandleError(err)
	if has_order == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的订单不存在！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"order": order,
	})
}
