package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderItem struct {
	*revel.Controller
}

func (c ErpOrderItem) DetailData() revel.Result {
	var item_id int64
	c.Params.Bind(&item_id, "item_id")

	var db_session = c.DbSession

	var order_item models.OrderItemInfo
	_, err := db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order_item_form OrderItemForm
	order_item_form.ItemId = order_item.Id
	order_item_form.OrderId = order_item.OrderId
	order_item_form.Sku = order_item.ItemSku
	order_item_form.Qty = order_item.Quantity

	return c.RenderJson(order_item_form)
}

type OrderItemForm struct {
	ItemId  int64  `json:"item_id"`
	OrderId int64  `json:"order_id"`
	Sku     string `json:"sku"`
	Qty     int    `json:"qty"`
}

func (c ErpOrderItem) Save() revel.Result {

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细保存成功！",
	})
}

func (c ErpOrderItem) Delete() revel.Result {
	var item_id int64
	c.Params.Bind(&item_id, "item_id")

	var db_session = c.DbSession

	_, err := db_session.Id(item_id).Delete(new(models.OrderItemInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细删除成功！",
	})
}
