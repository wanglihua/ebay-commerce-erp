package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
)

type ErpOrderAll struct {
	*revel.Controller
}

func (c ErpOrderAll) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_all/order_all_index.html")
}

func (c ErpOrderAll) PageData() revel.Result {
	db_session := c.DbSession

	shop_list := get_shop_select_item_list(db_session)
	shop_list = append([]shop_select_item_t{shop_select_item_t{Id: 0, Name: ""}}, shop_list...)

	return c.RenderJson(map[string]interface{}{
		"shop_list": shop_list,
	})
}

func (c ErpOrderAll) ListData() revel.Result {
	db_session := c.DbSession

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if shop_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shop_id = %d", shop_id))
	}

	order_list, count := get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
