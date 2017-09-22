package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
)

type ErpInventoryHistory struct {
	*revel.Controller
}

func (c ErpInventoryHistory) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	var inventory_id int64
	c.Params.Bind(&inventory_id, "inventory_id")

	c.RenderArgs["inventory_id"] = inventory_id

	return c.RenderTemplate("erp/kcgl/inventory_history/inventory_history_index.html")
}

type InventoryHistoryViewItem struct {
	InventoryHistory kcgl_models.InventoryHistoryInfo  `xorm:"extends" json:"his"`
	Stock            sppz_models.StockInfo             `xorm:"extends" json:"st"`
	StorageLoc       sppz_models.StorageLocInfo        `xorm:"extends" json:"loc"`
	Product          kcgl_models.ProductInfo           `xorm:"extends" json:"prd"`
}

func (c ErpInventoryHistory) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	var inventory_id int64
	c.Params.Bind(&inventory_id, "inventory_id")

	db_session := c.DbSession

	var inventory kcgl_models.InventoryInfo
	_, err := db_session.Id(inventory_id).Get(&inventory)
	core.HandleError(err)

	filter, _, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT his.*, st.*, loc.*, prd.*
		FROM erp_inventory_history his
		    INNER JOIN erp_stock st ON his.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON his.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON his.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select(`his.*, st.*, loc.*, prd.*`)
	sb.Count(`count(his.id) AS his_count`)
	sb.From(`erp_inventory_history his
		    INNER JOIN erp_stock st ON his.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON his.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON his.product_id = prd.id`)

	if inventory_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("his.stock_id = %d", inventory.StockId))
		filter = core.SqlAnd(filter, fmt.Sprintf("his.product_id = %d", inventory.ProductId))
	}

	sb.Where(filter)

	sb.OrderBy("his.id DESC")

	sb.Limit(offset, limit)

	inventory_history_list := make([]InventoryHistoryViewItem, 0, limit)
	err = db_session.Sql(sb.SqlList()).Find(&inventory_history_list)
	core.HandleError(err)

	var his_count struct {
		Count int64 `xorm:"his_count"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&his_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_history_list,
		Total: his_count.Count,
	})
}
