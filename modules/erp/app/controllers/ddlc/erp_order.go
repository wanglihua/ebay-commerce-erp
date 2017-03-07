package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrder struct {
	*revel.Controller
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

	var status string //空，yzf, dph, phz, yfh, ywc
	c.Params.Bind(&status, "status")

	var status_id int64 = 0 // 0： all status; other: 具体的 status_id
	if status == "" {
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
