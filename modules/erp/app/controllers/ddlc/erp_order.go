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
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/order_logistics_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/paypal_service"
	"matrix/modules/erp/service/shop_service"
	"matrix/modules/erp/service/stock_pick_crew_service"
	"matrix/modules/erp/setting"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/order_history_service"
	"time"
	"encoding/json"
	"sync"
	"matrix/modules/erp/service/order_item_service"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
	"math"
	"strconv"
	"matrix/modules/erp/service/order_ebay_service"
	"matrix/modules/erp/service/order_feedback_service"
	"matrix/modules/erp/service/order_shipping_address_service"
	"matrix/modules/erp/models/ppgl_models"
)

type ErpOrder struct {
	*revel.Controller
}

func (c ErpOrder) Index(status string) revel.Result {
	var db_session = c.DbSession

	var has_pass_order_new = false

	if pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_NEW, c.Session) {
		has_pass_order_new = true
	}

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

	var has_pass_order_huifu = pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_HUIFU, c.Session)
	var is_button_huifu_visible = false
	if has_pass_order_huifu == true &&
		status_id == models.ORDER_STATUS_YQX {

		is_button_huifu_visible = true
	}

	// is_button_upload_ebay_ship_visible
	var has_pass_upload_ebay_ship = pass_service.Is_login_user_has_pass(pass_res.DDLC_ORDER_FHSC, c.Session)
	var is_button_upload_ebay_ship_visible = false
	if has_pass_upload_ebay_ship &&
		(status_id == models.ORDER_STATUS_DPH ||
			status_id == models.ORDER_STATUS_PHZ ||
			status_id == models.ORDER_STATUS_YFH ||
			status_id == models.ORDER_STATUS_YWC) {
		is_button_upload_ebay_ship_visible = true
	}

	c.RenderArgs["active_menu_html"] = template.HTML(active_menu_html)
	c.RenderArgs["status"] = status
	c.RenderArgs["status_id"] = status_id
	c.RenderArgs["status_chinese_name"] = status_chinese_name
	c.RenderArgs["can_view_special_hold"] = can_view_special_hold_ // 这个后来在页面中没有使用，只是这里暂时没有去掉。更改来之，任何权限的人都可以看到 pick hold
	c.RenderArgs["icon_path"] = setting.ICON_PATH
	c.RenderArgs["has_pass_order_new"] = has_pass_order_new
	c.RenderArgs["is_button_quxiao_visible"] = is_button_quxiao_visible
	c.RenderArgs["is_button_huifu_visible"] = is_button_huifu_visible
	c.RenderArgs["is_button_upload_ebay_ship_visible"] = is_button_upload_ebay_ship_visible

	return c.RenderTemplate("erp/ddlc/order/order_index.html")
}

// 优化查询定义  结束

// index page data
func (c ErpOrder) PageData() revel.Result {
	//这个的权限就不用控制了

	var status_id int64
	c.Params.Bind(&status_id, "status_id")

	db_session := c.DbSession

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list_limit_100(status_id, db_session)

	shop_list := shop_service.Get_shop_select_item_list(db_session)

	var order_type_list = make([]ddpz_models.OrderTypeInfo, 0)
	err := db_session.Find(&order_type_list)
	core.HandleError(err)

	var status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_ORDER)).Find(&status_list)
	core.HandleError(err)

	// pick up hold任何权限的人都可以看到
	var special_hold_type_list = make([]ddpz_models.OrderHoldTypeInfo, 0)
	var can_view_special_hold_ = can_view_special_hold(c.Session, db_session)
	if can_view_special_hold_ {
		err = db_session.In("id", special_hold_type_id_list).Find(&special_hold_type_list)
	} else {
		var special_hold_type_id_list_only_pick_up = []int64{ddpz_models.ORDER_HOLD_TYPE_PICK_UP}
		err = db_session.In("id", special_hold_type_id_list_only_pick_up).Find(&special_hold_type_list)
	}
	core.HandleError(err)

	var hold_type_list = make([]ddpz_models.OrderHoldTypeInfo, 0)
	core.HandleError(db_session.Find(&hold_type_list))

	var order_tag_list = make([]ddpz_models.OrderTagInfo, 0)
	err = db_session.OrderBy("name ASC").Find(&order_tag_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"shop_id":                0,
		"shop_list":              shop_list,
		"order_type_id":          0,
		"order_type_list":        order_type_list,
		"status_id":              0,
		"status_list":            status_list,
		"query_criteria_id1":     0,
		"query_criteria_id2":     0,
		"query_criteria_list":    get_order_query_criteria_list(),
		"good_query_id":          0,
		"good_query_list":        get_good_query_list(),
		"special_hold_type_list": special_hold_type_list,
		"special_hold_type_id":   -1, // 表示不选择 hold
		"hold_type_list":         hold_type_list,
		"hold_type_id":           0,
		"batch_code":             0,
		"batch_code_list":        batch_code_list,
		"order_tag_id":           0,
		"order_tag_list":         order_tag_list,
		"time_query_id":          time_query_paid_time,
		"time_query_list":        get_time_query_list(),
	})
}

func (c ErpOrder) ListData() revel.Result {
	// 权限控制是分流程状态的，在下面的状态检查里面做控制了
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if strings.Contains(order, "stock_name") {
		order = strings.Replace(order, "`stock_name`", "o.stock_id", -1) // replace `stock_name`，不能是 stock_name
	}

	filter, _ = Get_order_filter(db_session, c.Session, c.Params, filter)

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}

func (c ErpOrder) Detail(order_id int64) revel.Result {
	var db_session = c.DbSession

	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var order ddlc_models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	c.RenderArgs["order_id"] = order_id
	c.RenderArgs["is_hold"] = order.StatusId == models.ORDER_STATUS_HOLD
	c.RenderArgs["is_dph"] = order.StatusId == models.ORDER_STATUS_DPH

	// 不是 已发货 已完成 已取消 状态的都能修改
	var can_modify = order_service.Can_order_modify(&order)
	c.RenderArgs["can_modify"] = can_modify

	var can_modify_disabled = template.HTML("")
	if can_modify == false {
		can_modify_disabled = template.HTML(` disabled="disabled" `)
	}
	c.RenderArgs["can_modify_disabled"] = can_modify_disabled

	c.RenderArgs["can_edit_waybill_number"] = order_service.Can_order_modify_waybill_number(&order)

	var is_pickup = false
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
		is_pickup = true
	}

	c.RenderArgs["is_pickup"] = is_pickup

	/* 暂时先设置为总是可以解锁快递
	var is_express = false
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_EXPRESS {
		is_express = true
	}
	c.RenderArgs["is_express"] = is_express
	*/
	c.RenderArgs["is_express"] = true

	var is_pickup_paid int64 = 0
	var pickup_order ddlc_models.PickupOrderInfo
	_, err = db_session.Where("order_id = ?", order_id).Get(&pickup_order)
	core.HandleError(err)

	is_pickup_paid = pickup_order.IsPaid

	c.RenderArgs["is_pickup_paid"] = is_pickup_paid

	var can_reship = order.StatusId == models.ORDER_STATUS_YFH || order.StatusId == models.ORDER_STATUS_YWC
	c.RenderArgs["can_reship"] = can_reship

	var order_feedback_list = order_feedback_service.Get_order_feedback_list(order_id, db_session)
	c.RenderArgs["order_feedback_list"] = order_feedback_list
	c.RenderArgs["order_feedback_list_len"] = len(order_feedback_list)

	// todo 具体什么状态下可以退款，什么状态下不能退款，需要判断一下
	c.RenderArgs["can_refund"] = true

	return c.RenderTemplate("erp/ddlc/order/order_detail.tpl") //使用 tpl 扩展名是为了这个页面不被IDE识别为HTML，这个页面在HTML编辑模式下格式化会打乱产生问题
}

func (c ErpOrder) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "id")

	revel.TRACE.Println("order_id:", order_id)

	var order ddlc_models.OrderInfo
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

	var shop xtpz_models.ShopAccountInfo
	_, err = db_session.Id(order.ShopId).Get(&shop)
	core.HandleError(err)

	var order_type ddpz_models.OrderTypeInfo
	_, err = db_session.Id(order.TypeId).Get(&order_type)
	core.HandleError(err)

	var order_hold_type ddpz_models.OrderHoldTypeInfo
	_, err = db_session.Id(order.HoldTypeId).Get(&order_hold_type)
	core.HandleError(err)

	ebay_shipping_address, paypal_shipping_address, shipping_address := order_shipping_address_service.Get_order_shipping_address(order_id, db_session)

	shipping_address.OrderId = order.Id //如果地址是空的，也先给一个 order_id

	type order_item_view_t struct {
		OrderItem         ddlc_models.OrderItemInfo        `xorm:"extends" json:"oi"`
		Product           kcgl_models.ProductInfo          `xorm:"extends" json:"p"`
		ProductSotoreItem sppz_models.ProductStoreItemInfo `xorm:"extends" json:"psi"`
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

	for index, inventory := range inventory_list {
		if inventory.Stock.StockTypeId == sppz_models.STOCK_TYPE_ZT {
			// 在途仓获取到港时间
			//inventory_list[index].DeptTime = inventory_service.Get_zt_stock_product_dept_time(inventory.Stock.Id, inventory.Product.Id, db_session)

			// 在途仓获取到仓时间
			var zt_product_info = inventory_service.Get_zt_stock_product_stock_time(inventory.Stock.Id, inventory.Product.Id, db_session)
			inventory_list[index].TransferId = zt_product_info.TransferId
			inventory_list[index].BatchCode = zt_product_info.BatchCode
			inventory_list[index].StockTime = zt_product_info.StockTime.Format("2006-01-02")

			// inventory_list[index].BatchCode, inventory_list[index].StockTime = inventory_service.Get_zt_stock_product_stock_time(inventory.Stock.Id, inventory.Product.Id, db_session)
		}
	}

	var order_logistics ddlc_models.OrderLogisticsInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order_id)).Get(&order_logistics)
	core.HandleError(err)

	var logistics_mode ddpz_models.LogisticsModeInfo
	if !(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid == true &&
		(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS)) {
		//如果是Hold状态且是信封Hold的就不显示匹配的物流，因为这些信息可以是前一次匹配时遗留的，后来指定了 Hold
		_, err = db_session.Id(order_logistics.ModeId).Get(&logistics_mode)
		core.HandleError(err)
	}

	var logistics_stock sppz_models.StockInfo
	if !(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid == true &&
		(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS)) {
		//如果是Hold状态且是信封Hold的就不显示匹配的物流，因为这些信息可以是前一次匹配时遗留的，后来指定了 Hold
		_, err = db_session.Id(order_logistics.StockId).Get(&logistics_stock)
		core.HandleError(err)
	}

	var logistics_setting ddpz_models.LogisticsSettingInfo
	_, err = db_session.Limit(1, 0).Get(&logistics_setting)
	core.HandleError(err)

	var order_cost_amount = order_logistics.CostAmount * order.ExchangeRate

	// 物流费用

	var logistics_fee float64 = (order_logistics.TotalFee + order_logistics.PickFee) * order.ExchangeRate

	if order_logistics.ActualFee != 0 {
		logistics_fee = (order_logistics.ActualFee + order_logistics.PickFee) * order.ExchangeRate
	}

	// gst 消费税
	var order_gst float64 = order_service.Get_order_gst_fee(&order)

	// 重新按实际运费计算毛利
	var order_profit_au = order_service.CalculateOrderProfit(&order, &order_logistics)
	var order_profit = order_profit_au * order.ExchangeRate

	if math.Abs(order.ProfitAu-order_profit_au) > core.FLOAT_MIN_DIFF { // 之前许多订单没有按实际运费计算毛利的，这个机会计算一下
		order.ProfitAu = order_profit_au
		order.ProfitRmb = order_profit

		_, err := db_session.Id(order.Id).Cols("profit_au", "profit_rmb").Update(&order)
		core.HandleError(err)

		// 平均分配利润到订单项
		var order_item_list = make([]ddlc_models.OrderItemInfo, 0)
		core.HandleError(db_session.Where("order_id = ?", order.Id).Find(&order_item_list))

		var order_cost_amount = order_service.GetOrderCostAmountSum(order.Id, db_session)

		for _, oi := range order_item_list {
			// 利润分配
			oi.GrossProfit = oi.CostPrice * float64(oi.Quantity) * order.ProfitAu / order_cost_amount
			_, err := db_session.Id(oi.Id).Cols("gross_profit").Update(&oi)
			core.HandleError(err)
		}
	}

	var au_stock_list = make([]sppz_models.StockInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&au_stock_list)
	core.HandleError(err)

	var pickup_pay_method_list = make([]ddpz_models.PickupPayMethodInfo, 0)
	core.HandleError(db_session.Find(&pickup_pay_method_list))

	var pickup_actual_pay_method_list = make([]ddpz_models.PickupPayMethodInfo, 0)
	core.HandleError(db_session.In("id", []int64{ddpz_models.PICKUP_PAY_METHOD_CARD, ddpz_models.PICKUP_PAY_METHOD_CASH}).Find(&pickup_actual_pay_method_list))

	var pickup_order ddlc_models.PickupOrderInfo
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
		_, err := db_session.Where("order_id = ?", order.Id).Get(&pickup_order)
		core.HandleError(err)
		pickup_order.OrderId = order.Id
	}

	pickup_order.ActualAmount = order.AmountPaid

	var refund_list = make([]ppgl_models.OrderRefundInfo, 0)
	err = db_session.Where("order_id = ?", order.Id).Find(&refund_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"order":                         order,
		"status":                        status,
		"shop":                          shop,
		"order_type":                    order_type,
		"ebay_shipping_address":         ebay_shipping_address,
		"paypal_shipping_address":       paypal_shipping_address,
		"shipping_address":              shipping_address,
		"order_hold_type":               order_hold_type,
		"order_item_list":               order_item_list,
		"inventory_list":                inventory_list,
		"order_logistics":               order_logistics,
		"logistics_setting":             logistics_setting,
		"order_cost_amount":             order_cost_amount, // rmb
		"logistics_fee":                 logistics_fee,     // rmb
		"order_gst":                     order_gst,         // aud
		"order_profit":                  order.ProfitRmb,   // rmb
		"au_stock_list":                 au_stock_list,
		"logistics_stock":               logistics_stock,
		"logistics_mode":                logistics_mode,
		"pickup":                        pickup_order,
		"pickup_pay_method_list":        pickup_pay_method_list,
		"pickup_actual_pay_method_list": pickup_actual_pay_method_list,
		"refund_list":                   refund_list,
	})
}

func (c ErpOrder) ReMatchHoldOrderLogistics() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.DDLC_ORDER_VIEW_DETAIL, pass_res.DDLC_ORDER_HOLD}, c.Session)
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_MATCH, c.Session)

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

	var order_list = make([]ddlc_models.OrderInfo, 0)
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

	var p_order_logistics = new(order_logistics_service.OrderLogisticsService)
	var login_user = core.GetLoginUser(c.Session)
	p_order_logistics.LoginUser = login_user
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

func (c ErpOrder) SaveDesc() revel.Result {
	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var order_desc string
	c.Params.Bind(&order_desc, "order_desc")

	var order ddlc_models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	var desc_before = order.Desc.String
	order.Desc = core.MakeNullString(order_desc, true)

	_, err = db_session.Id(order_id).Cols("desc").Update(&order)
	core.HandleError(err)

	var desc_after = order_desc

	var login_user = core.GetLoginUser(c.Session)
	order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单修改描述, %s => %s， 工号： %s", desc_before, desc_after, login_user.LoginName), db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单描述保存成功！",
	})
}

func (c ErpOrder) SavePickupActual() revel.Result {
	var db_session = c.DbSession

	var pickup_order ddlc_models.PickupOrderInfo
	err := json.Unmarshal(c.GetRequestBody(), &pickup_order)
	core.HandleError(err)

	if pickup_order.ActualPayMethodId == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "现场付款方式不能为空！",
		})
	}

	var pickup_order_in_db ddlc_models.PickupOrderInfo
	has, err := db_session.Where("order_id = ?", pickup_order.OrderId).Get(&pickup_order_in_db)
	core.HandleError(err)

	if has == true {
		pickup_order_in_db.ActualIsPaid = pickup_order.ActualIsPaid
		pickup_order_in_db.ActualPayMethodId = pickup_order.ActualPayMethodId
		pickup_order_in_db.Description = pickup_order.Description
		_, err = db_session.Id(pickup_order.Id).Cols("actual_is_paid", "actual_pay_method_id", "description").Update(&pickup_order_in_db)
		core.HandleError(err)
	} else {
		_, err = db_session.Insert(&pickup_order)
		core.HandleError(err)
	}

	var order ddlc_models.OrderInfo
	_, err = db_session.Id(pickup_order.OrderId).Get(&order)
	core.HandleError(err)

	order.AdjustmentAmount = pickup_order.ActualAmount
	order.AdjustmentAmountCurrency = "AUD"
	order.AdjustmentAmountRmb = pickup_order.ActualAmount * order.ExchangeRate
	order.AmountPaid = pickup_order.ActualAmount
	order.AmountPaidCurrency = "AUD"
	order.AmountPaidRmb = pickup_order.ActualAmount * order.ExchangeRate
	order.TotalAmount = pickup_order.ActualAmount
	order.TotalAmountCurrency = "AUD"
	order.TotalAmountRmb = pickup_order.ActualAmount * order.ExchangeRate

	_, err = db_session.Id(order.Id).
		Cols("adjustment_amount", "adjustment_amount_currency", "adjustment_amount_rmb", "amount_paid", "amount_paid_currency", "amount_paid_rmb", "total_amount", "total_amount_currency", "total_amount_rmb").
		Update(&order)
	core.HandleError(err)

	// 获得当前订单的物流信息
	var order_logistics ddlc_models.OrderLogisticsInfo
	_, err = db_session.Where("order_id = ?", order.Id).Get(&order_logistics)
	core.HandleError(err)

	order_service.Update_order_logistics_cost_amount_and_order_profit_and_order_item_price_profit(order, order_logistics, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存PICK现场收款信息成功！",
	})
}

// SavePickupActual

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
	pass_service.Check_login_user_pass_and([]int64{pass_res.DDLC_ORDER_YZF, pass_res.DDLC_ORDER_LOGISTICS_MATCH}, c.Session)

	var db_session = c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var p_order_logistics = new(order_logistics_service.OrderLogisticsService)
	var login_user = core.GetLoginUser(c.Session)
	p_order_logistics.LoginUser = login_user

	var err error
	if len(id_list) == 0 {
		err = p_order_logistics.CalculateListPaidOrders(login_user, db_session)
	} else {
		err = p_order_logistics.CalculateListWithUser(id_list, login_user, db_session)
	}

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

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_HOLD ||
			( order.StatusId == models.ORDER_STATUS_HOLD &&
				(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_EXPRESS ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态或Hold类型不能搬至待配货状态！",
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_DPH
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单Hold转 待配货， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单搬至待配货成功！",
	})
}

func (c ErpOrder) PickupDaiPeiHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YZF, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	if len(order_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单列表不能为空！",
		})
	}

	var order_list = make([]ddlc_models.OrderInfo, 0)
	core.HandleError(db_session.In("id", order_id_list).Find(&order_list))

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_YZF || order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "当前选择订单不全是已支付状态的PICKUP订单！",
			})
		}
	}

	/*
	var order_item_zero_count = order_item_service.Get_order_item_zero_count_in_order_list(order_list, db_session)

	if order_item_zero_count > 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单明细不能为空！",
		})
	}
	*/

	var order_code_list = order_item_service.Get_order_item_zero_code_list_in_order_list(order_list, db_session)
	if len(order_code_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单明细不能为空！订单号： " + strings.Join(order_code_list, ","),
		})
	}

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		_, err := db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "PICKUP订单待配货成功！",
	})
}

func (c ErpOrder) PeiHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_DPH, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)

	if len(order_id_list) != 0 {
		err := db_session.In("id", order_id_list).Find(&order_list)
		core.HandleError(err)
	} else {
		err := db_session.Where(fmt.Sprintf("status_id = %d", models.ORDER_STATUS_DPH)).Find(&order_list)
		core.HandleError(err)
	}

	var order_id_str_list = make([]string, 0)
	for _, order := range order_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order.Id))
	}

	for _, order := range order_list {
		if order.StockId == 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "有的订单的发货仓库未设置！",
			})
		}
	}

	var login_user_stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	// 检查用户对发货仓库是否有查看权限
	for _, order := range order_list {
		var has_stock = false
		for _, stock_id := range login_user_stock_id_list {
			if order.StockId == stock_id {
				has_stock = true
				break
			}
		}

		if has_stock == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "有的订单的发货仓库，当前用户没权限查看！",
			})
		}
	}

	/*
	// 检查有权限的仓库当前是不是有配货中的订单

	if len(login_user_stock_id_list) != 0 {
		count, err := db_session.
		In("stock_id", login_user_stock_id_list).
			Where(fmt.Sprintf("status_id = %d", models.ORDER_STATUS_PHZ)).
			Count(&ddlc_models.OrderInfo{})
		core.HandleError(err)

		if count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "有订单还在配货中，不能搬至已配货！",
			})
		}
	} else {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "当前用户没权限查看任何仓库，不能搬动订单至配货中！",
		})
	}
	*/

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_DPH &&
			(!(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid &&
				(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_EXPRESS ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP))) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能配货！",
			})
		}
	}

	type stock_id_t struct {
		StockId int64 `xorm:"'stock_id'"`
	}

	/*
SELECT DISTINCT stock_id
FROM erp_order o
WHERE o.status_id = 0;


SELECT DISTINCT stock_id
FROM erp_order o
WHERE o.id IN (-1, -2, -3);
	*/

	var phz_stock_id_list = make([]stock_id_t, 0)
	var sql_phz_stock_id_list = `SELECT DISTINCT stock_id
FROM erp_order o
WHERE o.status_id = %d;`
	sql_phz_stock_id_list = fmt.Sprintf(sql_phz_stock_id_list, models.ORDER_STATUS_PHZ)
	core.HandleError(db_session.Sql(sql_phz_stock_id_list).Find(&phz_stock_id_list))

	var ph_stock_id_list = make([]stock_id_t, 0)
	var sql_ph_stock_id_list = `SELECT DISTINCT stock_id
FROM erp_order o
WHERE o.id IN (%s);`
	sql_ph_stock_id_list = fmt.Sprintf(sql_ph_stock_id_list, strings.Join(order_id_str_list, ","))
	core.HandleError(db_session.Sql(sql_ph_stock_id_list).Find(&ph_stock_id_list))

	for _, phz_stock_id := range phz_stock_id_list {
		for _, ph_stock_id := range ph_stock_id_list {
			if phz_stock_id == ph_stock_id {
				return c.RenderJson(map[string]interface{}{
					"success": false,
					"message": "有订单还在配货中，不能搬至已配货！",
				})
			}
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_PHZ
		_, err := db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单配货， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单配货成功！",
	})
}

func (c ErpOrder) ForcePeiHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_FORCE_PEIHUO, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)

	if len(order_id_list) != 0 {
		err := db_session.In("id", order_id_list).Find(&order_list)
		core.HandleError(err)
	} else {
		err := db_session.Where(fmt.Sprintf("status_id = %d", models.ORDER_STATUS_DPH)).Find(&order_list)
		core.HandleError(err)
	}

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_DPH &&
			(!(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid &&
				(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_EXPRESS ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
					order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP))) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不能配货！",
			})
		}
	}
	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_PHZ
		_, err := db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单配货， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单取消波次限制配货成功！",
	})
}

var fahuo_mutex sync.Mutex

func (c ErpOrder) FaHuo() revel.Result {
	// 通过互斥锁来解决并发发货多次同时扣减库存的问题
	fahuo_mutex.Lock()
	defer func() {
		fahuo_mutex.Unlock()
	}()

	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或特殊hold状态，不能发货！",
			})
		}
	}

	for _, order := range order_list {
		if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
			var pickup_order_info ddlc_models.PickupOrderInfo
			_, err := db_session.Where("order_id = ?", order.Id).Get(&pickup_order_info)
			core.HandleError(err)

			if pickup_order_info.IsPaid == 0 && pickup_order_info.ActualIsPaid == 0 {
				return c.RenderJson(map[string]interface{}{
					"success": false,
					"message": "PICKUP订单未收款不能搬至已发货！",
				})
			}
		}
	}

	var mel_stock_id = stock_service.Get_mel_stock_id(db_session)

	for _, order := range order_list {
		if order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
			// 是Hold直接发货的，把 Hold订单的仓库直接设置为墨尔本的仓库，但 pickup 的有自己的仓库
			order.StockId = mel_stock_id
			_, err = db_session.Id(order.Id).Cols("stock_id").Update(&order)
			core.HandleError(err)

			// order_list[index].Version = order_list[index].Version + 1
		}
	}

	/*
	// 前面更新过之后，version 就会改变所以再取一下
	order_list = make([]ddlc_models.OrderInfo, 0)
	err = db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)
	*/

	var login_user = core.GetLoginUser(c.Session)

	var take_inventory_result = true
	var time_now = time.Now()
	for _, order := range order_list {

		if inventory_service.Take_order_inventory(order, db_session) == false {
			take_inventory_result = false
		}

		order.StatusId = models.ORDER_STATUS_YFH
		order.ShippedTime = core.MakeNullLongTime(time_now, true)
		_, err = db_session.Id(order.Id).Cols("status_id", "shipped_time").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单发货， 工号： %s", login_user.LoginName), db_session)
	}

	// 上传 ebay 发货状态
	/* 临时取消这个操作，通过手工完成，手工成熟后，再启用
	go order_ebay_service.Upload_order_ebay_status(order_id_list)
	*/

	/*
	var upload_order_ebay_status_result = order_ebay_service.Upload_order_ebay_status(order_id_list)

	if upload_order_ebay_status_result.Success == false {
		return c.RenderJson(upload_order_ebay_status_result)
	}
	*/

	// 更新库存待配货数量
	go inventory_service.Calculate_inventory_pending_qty()

	if take_inventory_result == true {
		return c.RenderJson(map[string]interface{}{
			"success": true,
			"message": "订单发货成功！",
		})
	} else {
		return c.RenderJson(map[string]interface{}{
			"success": true,
			"message": "订单发货成功，但部分库存扣减失败！",
		})
	}
}

func (c ErpOrder) UploadEbayStatus() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_FHSC, c.Session)

	var db_session = c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var upload_order_ebay_status_result = order_ebay_service.Upload_order_ebay_status(order_id_list)

	if upload_order_ebay_status_result.Success == false {
		return c.RenderJson(upload_order_ebay_status_result)
	}

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		var logistics_company ddpz_models.LogisticsCompanyInfo
		_, err := db_session.Id(order.CompanyId).Get(&logistics_company)
		core.HandleError(err)
		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单发货信息上传ebay, 物流公司：%s 运单号：%s，工号： %s", logistics_company.EbayName, order.WaybillNumber, login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单发货信息上传ebay成功！",
	})
}

func (c ErpOrder) ShiJiWuHuo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或hold，不能设置实际无货！",
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		order.HoldTypeId = core.MakeNullInt(ddpz_models.ORDER_HOLD_TYPE_OOS_ACTUAL_STOCK, true)
		_, err = db_session.Id(order.Id).Cols("status_id", "hold_type_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单实际无货，工号： %s", login_user.LoginName), db_session)
	}

	go inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单设置实际无货成功！",
	})
}

func (c ErpOrder) SpecialNotShip() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var desc string
	c.Params.Bind(&desc, "desc")

	if len(order_id_list) != 1 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择一条订单",
		})
	}

	if strings.TrimSpace(desc) == "" {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请输入描述!",
		})
	}

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_PHZ && order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是配货中或hold，不能设置 特殊情况实际未发货！",
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		order.HoldTypeId = core.MakeNullInt(ddpz_models.ORDER_HOLD_TYPE_SPECIAL_NOT_SHIP, true)

		if order.Desc.Valid == false {
			order.Desc = core.MakeNullString(desc, true)
		} else {
			order.Desc = core.MakeNullString(order.Desc.String+desc, true)
		}

		_, err = db_session.Id(order.Id).Cols("status_id", "hold_type_id", "desc").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("特殊情况实际未发货，工号： %s", login_user.LoginName), db_session)
	}

	go inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单设置 特殊情况实际未发货 成功！",
	})
}

func (c ErpOrder) SpecialHoldSpecialNotShip() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_PHZ, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var desc string
	c.Params.Bind(&desc, "desc")

	if len(order_id_list) != 1 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择一条订单",
		})
	}

	if strings.TrimSpace(desc) == "" {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请输入描述!",
		})
	}

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	for _, order := range order_list {
		if (order.StatusId != models.ORDER_STATUS_HOLD) ||
			(order.StatusId == models.ORDER_STATUS_HOLD && (
				order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_ENVELOPE &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_EXPRESS &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL &&
					order.HoldTypeId.Int64 != ddpz_models.ORDER_HOLD_TYPE_PICK_UP)) {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单当前状态不是特殊hold，不能设置 特殊Hold特殊情况实际未发货！",
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		order.HoldTypeId = core.MakeNullInt(ddpz_models.ORDER_HOLD_TYPE_SPECIAL_HOLD_SPECIAL_NOT_SHIP, true)

		if order.Desc.Valid == false {
			order.Desc = core.MakeNullString(desc, true)
		} else {
			order.Desc = core.MakeNullString(order.Desc.String+desc, true)
		}

		_, err = db_session.Id(order.Id).Cols("status_id", "hold_type_id", "desc").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("特殊Hold特殊情况实际未发货，工号： %s", login_user.LoginName), db_session)
	}

	go inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单设置 特殊Hold特殊情况实际未发货 成功！",
	})
}

func (c ErpOrder) WanCheng() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_YFH, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)
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

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_YWC
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单完成，工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单完成成功！",
	})
}

func (c ErpOrder) PickupOptionData() revel.Result {
	var db_session = c.DbSession

	var stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&stock_list))

	return c.RenderJson(map[string]interface{}{
		"success":    true,
		"stock_list": stock_list,
	})
}

func (c ErpOrder) Pickup() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.DDLC_ORDER_DPH, pass_res.DDLC_ORDER_HOLD}, c.Session)

	var db_session = c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")
	revel.TRACE.Println("order_id_list", order_id_list)

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")
	revel.TRACE.Println("stock_id", stock_id)

	var order_list = make([]ddlc_models.OrderInfo, 0)
	core.HandleError(db_session.In("id", order_id_list).Find(&order_list))

	for _, order := range order_list {
		if order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("订单 %s 当前状态 已经是 PICK UP", order.Code),
			})
		}
	}

	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_DPH && order.StatusId != models.ORDER_STATUS_HOLD {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("订单 %s 当前状态不能搬至 PICK UP", order.Code),
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.StatusId = models.ORDER_STATUS_HOLD
		order.HoldTypeId = core.MakeNullInt(ddpz_models.ORDER_HOLD_TYPE_PICK_UP, true)
		order.StockId = stock_id
		order.CompanyId = 0
		order.ModeId = 0

		_, err := db_session.Id(order.Id).Cols("status_id", "hold_type_id", "stock_id", "company_id", "mode_id").Update(&order)
		core.HandleError(err)

		// order logistics
		var order_logistics ddlc_models.OrderLogisticsInfo
		has, err := db_session.Where(fmt.Sprintf("order_id = %d", order.Id)).Get(&order_logistics)
		core.HandleError(err)

		order_logistics.OrderId = order.Id
		order_logistics.StockId = stock_id
		order_logistics.CompanyId = 0
		order_logistics.ModeId = 0
		order_logistics.TotalFee = 0
		order_logistics.PickFee = 0

		if has {
			_, err = db_session.Id(order_logistics.Id).Cols("stock_id", "company_id", "mode_id").Update(&order_logistics)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&order_logistics)
			core.HandleError(err)
		}

		// pickup order
		var pickup_order ddlc_models.PickupOrderInfo
		has, err = db_session.Where(fmt.Sprintf("order_id = %d", order.Id)).Get(&pickup_order)
		core.HandleError(err)

		pickup_order.OrderId = order.Id
		pickup_order.IsPaid = 1

		if has {
			_, err = db_session.Id(order_logistics.Id).Cols("is_paid").Update(&order_logistics)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&pickup_order)
			core.HandleError(err)
		}

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单PICK UP， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "设置为 PICK UP 成功！",
	})
}

// #10 订单作废的功能
func (c ErpOrder) QuXiao() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_QUXIAO, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {
		order.CancelStatusId = order.StatusId
		order.StatusId = models.ORDER_STATUS_YQX // 已取消
		_, err = db_session.Id(order.Id).Cols("status_id", "cancel_status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单取消， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单取消成功！",
	})
}

func (c ErpOrder) HuiFu() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_HUIFU, c.Session)

	db_session := c.DbSession

	order_id_list := make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	var order_list = make([]ddlc_models.OrderInfo, 0)
	err := db_session.In("id", order_id_list).Find(&order_list)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	for _, order := range order_list {

		/*
	ORDER_STATUS_YCJ            = 3850 // 订单状态 已创建
	ORDER_STATUS_YZF            = 3900 // 订单状态 已支付
	ORDER_STATUS_HOLD           = 3950 // 订单状态 Hold
	ORDER_STATUS_PROCESSED_HOLD = 3975 // 订单状态 已处理Hold
	ORDER_STATUS_DPH            = 4000 // 订单状态 待配货
	ORDER_STATUS_PHZ            = 4100 // 订单状态 配货中
	ORDER_STATUS_YFH            = 4200 // 订单状态 已发货
	ORDER_STATUS_YWC            = 4300 // 订单状态 已完成
	ORDER_STATUS_YQX            = 4350 // 订单状态 已取消
		 */

		if order.CancelStatusId != models.ORDER_STATUS_YCJ &&
			order.CancelStatusId != models.ORDER_STATUS_YZF &&
			order.CancelStatusId != models.ORDER_STATUS_HOLD &&
			order.CancelStatusId != models.ORDER_STATUS_PROCESSED_HOLD &&
			order.CancelStatusId != models.ORDER_STATUS_DPH &&
			order.CancelStatusId != models.ORDER_STATUS_PHZ &&
			order.CancelStatusId != models.ORDER_STATUS_YFH &&
			order.CancelStatusId != models.ORDER_STATUS_YWC &&
			order.CancelStatusId != models.ORDER_STATUS_YQX {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "订单取消前的记录的订单状态不合法！",
			})
		}

		order.StatusId = order.CancelStatusId
		_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
		core.HandleError(err)

		order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单取消， 工号： %s", login_user.LoginName), db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单恢复成功！",
	})
}

func (c ErpOrder) ReShip() revel.Result {
	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	// 1 复制订单
	// 1.1 获取现在的订单
	var order ddlc_models.OrderInfo
	has, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的订单id不存在！",
		})
	}

	// 1.2 生成新的订单的编号
	var code_prefix string = order.Code
	if strings.Contains(order.Code, "-") {
		var order_code_splits = strings.Split(order.Code, "-")
		code_prefix = order_code_splits[0]
	}

	// 订单代码
	var tran_id_prefix string = order.TranId
	if strings.Contains(order.TranId, "-") {
		var tran_id_splits = strings.Split(order.TranId, "-")
		tran_id_prefix = tran_id_splits[0]
	}

	var full_tran_id_prefix string = order.FullTranId
	if strings.Contains(order.FullTranId, "-") {
		var full_tran_id_splits = strings.Split(order.FullTranId, "-")
		full_tran_id_prefix = full_tran_id_splits[0]
	}

	var max_code_order ddlc_models.OrderInfo
	_, err = db_session.Where(fmt.Sprintf("code LIKE '%s%%'", code_prefix)).OrderBy("code DESC").Limit(1).Get(&max_code_order)
	core.HandleError(err)

	var code_num int64 = 0
	if strings.Contains(max_code_order.Code, "-") {
		var order_code_splits = strings.Split(max_code_order.Code, "-")
		var err error
		code_num, err = strconv.ParseInt(order_code_splits[len(order_code_splits)-1], 10, 64)
		core.HandleError(err)
	}

	var new_order_code = fmt.Sprintf("%s-%d", code_prefix, code_num+1)
	var new_tran_id = fmt.Sprintf("%s-%d", tran_id_prefix, code_num+1)
	var new_full_tran_id = fmt.Sprintf("%s-%d", full_tran_id_prefix, code_num+1)

	// 1.3 保存新的订单信息
	var new_order ddlc_models.OrderInfo
	new_order = order               // 复制原订单的信息
	new_order.Id = 0                // 置空 id 让数据库重新生成id
	new_order.Code = new_order_code // 使用新的订单编号
	new_order.TranId = new_tran_id
	new_order.FullTranId = new_full_tran_id
	var time_now = core.NullLongTimeFrom(time.Now())
	new_order.CreatedTime = time_now
	new_order.PaidTime = time_now
	new_order.CreateTime = core.Time(time.Now())
	new_order.UpdateTime = core.Time(time.Now())

	// 重发复制订单时， 收款，平台交易费，付款交易费都没有，但有拣货费， 利润先置零
	new_order.AmountPaid = 0
	new_order.AmountPaidRmb = 0
	new_order.ExpressFee = 0
	new_order.TotalAmount = 0
	new_order.TotalAmountRmb = 0
	new_order.AdjustmentAmount = 0
	new_order.AdjustmentAmountRmb = 0
	new_order.SubTotalAmount = 0
	new_order.SubTotalAmountRmb = 0
	new_order.Ptfy = 0
	new_order.Jyfy = 0
	new_order.ProfitAu = 0
	new_order.ProfitRmb = 0
	new_order.TagId = 0

	new_order.AddressId = 0 // 先把地址置空，等会儿再填入新的地址id

	new_order.WaybillNumber = "" // 不复制货运单号

	new_order.StatusId = models.ORDER_STATUS_YZF   // 设置重发订单的初始状态为已支付
	new_order.TypeId = ddpz_models.ORDER_TYPE_RSDD // 设置重发订单的初始状态为已支付

	_, err = db_session.Insert(&new_order)
	core.HandleError(err)

	var new_order_id int64 = new_order.Id

	// 2 复制明细

	// 2.1 获取原来订单的明细列表
	var order_item_list = make([]ddlc_models.OrderItemInfo, 0)
	err = db_session.Where("order_id = ?", order_id).Find(&order_item_list)
	core.HandleError(err)

	// 2.2 保存为新的订单的明细
	for _, order_item := range order_item_list {
		order_item.Id = 0                 // 置空 id 让数据库重新生成id
		order_item.OrderId = new_order_id //设置为新的订单的订单 id

		_, err = db_session.Insert(&order_item)
		core.HandleError(err)
	}

	// 3 复制地址

	// 3.1 获取原来订单的地址
	var order_shipping_address_list = make([]ddlc_models.OrderShippingAddressInfo, 0)
	err = db_session.Where("order_id = ?", order_id).Find(&order_shipping_address_list)
	core.HandleError(err)

	// 3.2 保存为新的订单地址
	var order_shipping_address_id int64 = 0
	for _, order_shipping_address := range order_shipping_address_list {
		order_shipping_address.Id = 0                 // 置空 id 让数据库重新生成id
		order_shipping_address.OrderId = new_order_id //设置为新的订单的订单 id

		_, err = db_session.Insert(&order_shipping_address)
		core.HandleError(err)

		if order_shipping_address_id == 0 {
			order_shipping_address_id = order_shipping_address.Id
		}
	}

	// 3.3 更新订单地址id
	new_order.AddressId = order_shipping_address_id
	_, err = db_session.Id(new_order_id).Cols("address_id").Update(&new_order)
	core.HandleError(err)

	// 4 复制订单物流
	var order_logistics_list = make([]ddlc_models.OrderLogisticsInfo, 0) // 用列表的方式，其实总是只有一个数据
	err = db_session.Where("order_id = ?", order_id).Find(&order_logistics_list)
	core.HandleError(err)

	for _, order_logistics := range order_logistics_list {
		order_logistics.Id = 0                 // 置空 id 让数据库重新生成id
		order_logistics.OrderId = new_order_id //设置为新的订单的订单 id
		order_logistics.ActualFee = 0

		_, err = db_session.Insert(&order_logistics)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success":      true,
		"message":      "订单重发复制成功！",
		"new_order_id": new_order_id, // 重发复制后的 “新” 的订单id
	})
}
