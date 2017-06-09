package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpOrderItem struct {
	*revel.Controller
}

func (c ErpOrderItem) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var item_id int64
	c.Params.Bind(&item_id, "item_id")

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var db_session = c.DbSession

	var order_item models.OrderItemInfo
	_, err := db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order_item_form OrderItemForm
	order_item_form.ItemId = order_item.Id
	order_item_form.OrderId = order_id
	order_item_form.Sku = order_item.ItemSku
	order_item_form.Qty = order_item.Quantity

	return c.RenderJson(map[string]interface{}{
		"success":           true,
		"order_item_detail": order_item_form,
	})
}

type OrderItemForm struct {
	ItemId  int64  `json:"item_id"` //这是 order_item 表的主键
	OrderId int64  `json:"order_id"`
	Sku     string `json:"sku"`
	Qty     int    `json:"qty"`
}

func (c ErpOrderItem) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var order_item_form OrderItemForm
	err := json.Unmarshal(c.GetRequestBody(), &order_item_form)
	core.HandleError(err)

	c.Validation.Required(order_item_form.Sku).Message("SKU 不能为空！")
	c.Validation.Required(order_item_form.Qty).Message("数量 不能为空！")

	if c.Validation.HasErrors() {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	//获取 product_id
	var product models.ProductInfo
	_, err = db_session.Where(fmt.Sprintf("sku = '%s'", order_item_form.Sku)).Get(&product)
	core.HandleError(err)

	if product.Id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的SKU不存在!",
		})
	}

	var order_item models.OrderItemInfo
	if order_item_form.ItemId != 0 {
		_, err = db_session.Id(order_item_form.ItemId).Get(&order_item)
		core.HandleError(err)
	}

	order_item.OrderId = order_item_form.OrderId
	order_item.TransactionPriceCurrency = "AUD"
	order_item.ItemProductId = product.Id
	order_item.ItemSku = order_item_form.Sku
	order_item.Quantity = order_item_form.Qty

	if order_item_form.ItemId == 0 {
		_, err = db_session.Insert(&order_item)
		core.HandleError(err)
	} else {
		var update_cols = []string{"item_product_id", "item_sku", "quantity"}
		_, err = db_session.Id(order_item_form.ItemId).Cols(update_cols...).Update(&order_item)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细保存成功！",
	})
}

func (c ErpOrderItem) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

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
