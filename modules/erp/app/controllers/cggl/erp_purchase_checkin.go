package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpPurchaseCheckin struct {
	*revel.Controller
}

func (c ErpPurchaseCheckin) Index() revel.Result {
	return c.RenderTemplate("erp/cggl/purchase_checkin/purchase_checkin_index.html")
}

type PurchaseCheckinViewItem struct {
	kcgl_models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	models.StatusInfo                 `xorm:"extends" json:"s"`
}

func (c ErpPurchaseCheckin) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT pbe.*, p.*, s.*
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, p.*, s.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_purchase_bill_entry pbe
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id	`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pbe.id ASC")
	}

	sb.Limit(offset, limit)

	purchase_bill_entry_list := make([]PurchaseCheckinViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	var purchase_bill_entry_count struct {
		Count int64 `xorm:"'entry_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&purchase_bill_entry_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_entry_list,
		Total: purchase_bill_entry_count.Count,
	})
}
