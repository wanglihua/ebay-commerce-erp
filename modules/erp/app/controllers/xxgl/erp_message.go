package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {
	return c.RenderTemplate("erp/xxgl/message/message_index.html")
}

func (c ErpMessage) ListData() revel.Result {
	db_session := c.DbSession

	type message_view_item_t struct {
		models.MessageInfo           `xorm:"extends" json:"msg"`
		models.ShopAccountInfo       `xorm:"extends" json:"sa"`
		models.ECommercePlatformInfo `xorm:"extends" json:"pf"`
		models.ProductInfo           `xorm:"extends" json:"p"`
		models.StatusInfo            `xorm:"extends" json:"s"`
	}

	/*
SELECT msg.*, sa.*, pf.*, p.*, s.*
FROM erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    INNER JOIN erp_status s ON s.id = msg.status_id
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	sb := core.NewSqlBuilder()
	sb.Select("msg.*, sa.*, pf.*, p.*, s.*")
	sb.Count("count(msg.id) as msg_count")
	sb.From(`
erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    INNER JOIN erp_status s ON s.id = msg.status_id`)
	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("msg.id DESC")
	}

	sb.Limit(offset, limit)

	message_list := make([]message_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	var msg_count struct {
		Count int64 `xorm:"msg_count"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&msg_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  message_list,
		Total: msg_count.Count,
	})
}
