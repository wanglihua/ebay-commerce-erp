package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpInventoryCheckApprove struct {
	*revel.Controller
}

func (c ErpInventoryCheckApprove) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_CHECK_APPROVE, c.Session)

	return c.RenderTemplate("erp/kcgl/inventory_check_approve/inventory_check_approve_index.html")
}

type ErpInventoryCheckApproveViewItem struct {
	models.InventoryCheckInfo `xorm:"extends" json:"ic"`
	models.StockInfo          `xorm:"extends" json:"st"`
	CheckUser   auth_models.UserInfo `xorm:"extends" json:"cu"`
	ApproveUser auth_models.UserInfo `xorm:"extends" json:"au"`
	models.StatusInfo         `xorm:"extends" json:"ss"`
}

func (c ErpInventoryCheckApprove) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_CHECK_APPROVE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT ic.*, st.*, cu.*, au.*, ss.*
		FROM erp_inventory_check ic.
		    INNER JOIN erp_stock st ON ic.stock_id = st.id
		    INNER JOIN auth_user cu ON ic.check_user_id = cu.id
		    LEFT JOIN auth_user au ON ic.approve_user_id = au.id
		    INNER JOIN erp_status ss ON ic.status_id = ss.id
	*/

	query := db_session.
	Select("ic.*, st.*, cu.*, au.*, ss.*").
		Table("erp_inventory_check").Alias("ic").
		Join("INNER", []string{"erp_stock", "st"}, "ic.stock_id = st.id").
		Join("INNER", []string{"auth_user", "cu"}, "ic.check_user_id = cu.id").
		Join("LEFT", []string{"auth_user", "au"}, "ic.approve_user_id = au.id").
		Join("INNER", []string{"erp_status", "ss"}, "ic.status_id = ss.id").
		Where(filter)

	//query extra filter here
	query = query.Where("ic.status_id = ?", models.INVENTORY_CHECK_STATUS_SPZ)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("ic.id")
	}

	inventory_check_list := make([]ErpInventoryCheckApproveViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_check_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpInventoryCheckApproveViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_check_list,
		Total: count,
	})
}

// 这里是不是要加个信号量，让检查状态和更新库存的操作互斥
func (c ErpInventoryCheckApprove) Approve() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_CHECK_APPROVE, c.Session)

	db_session := c.DbSession

	var check_id int64
	c.Params.Bind(&check_id, "check_id")

	var inventory_check models.InventoryCheckInfo
	_, err := db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_SPZ {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单状态不是审核中！"})
	}

	//修改状态
	inventory_check.StatusId = models.INVENTORY_CHECK_STATUS_YWC
	_, err = db_session.Id(check_id).Cols("status_id").Update(&inventory_check)
	core.HandleError(err)

	//更新库存
	inventory_check_entry_list := make([]models.InventoryCheckEntryInfo, 0)
	err = db_session.Sql(fmt.Sprintf("SELECT * FROM erp_inventory_check_entry WHERE check_id = %d", check_id)).Find(&inventory_check_entry_list)
	core.HandleError(err)

	for _, inventory_check_entry := range inventory_check_entry_list {
		inventory := new(models.InventoryInfo)
		// inventory.Qty = inventory_check_entry.CheckQty
		inventory.Qty = inventory.Qty + inventory_check_entry.QtyDiff //利用库存差异来做应该更合理一些

		affected, err := db_session.
		Where("storage_loc_id = ? AND product_id = ?", inventory_check_entry.StorageLocId, inventory_check_entry.ProductId).
			Cols("qty").
			Update(inventory)
		core.HandleError(err)

		// 如果不存在相对应的库存记录
		if affected == 0 {
			var inventory models.InventoryInfo
			inventory.StockId = inventory_check.StockId
			inventory.StorageLocId = inventory_check_entry.StorageLocId
			inventory.ProductId = inventory_check_entry.ProductId
			inventory.Qty = inventory_check_entry.CheckQty

			_, err = db_session.Insert(&inventory)
			core.HandleError(err)
		}
	}

	//记录库存变更

	return c.RenderJson(core.JsonResult{
		Success: true,
		Message: "审批成功！",
	})
}

func (c ErpInventoryCheckApprove) Back() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_CHECK_APPROVE, c.Session)

	db_session := c.DbSession

	var check_id int64
	c.Params.Bind(&check_id, "check_id")

	var inventory_check models.InventoryCheckInfo
	_, err := db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_SPZ {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单状态不是审核中！"})
	}

	//修改状态
	inventory_check.StatusId = models.INVENTORY_CHECK_STATUS_WTG
	_, err = db_session.Id(check_id).Cols("status_id").Update(&inventory_check)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{
		Success: true,
		Message: "回退成功！",
	})
}
