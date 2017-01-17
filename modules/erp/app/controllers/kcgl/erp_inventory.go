package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpInventory struct {
	*revel.Controller
}

func (c ErpInventory) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/inventory/inventory_index.html")
}

type ErpInventoryViewItem struct {
	models.InventoryInfo  `xorm:"extends" json:"ivt"`
	models.StockInfo      `xorm:"extends" json:"st"`
	models.StorageLocInfo `xorm:"extends" json:"loc"`
	models.ProductInfo    `xorm:"extends" json:"prd"`
}

func (c ErpInventory) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT ivt.*, st.*, loc.*, prd.*
		FROM erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    INNER JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	query := db_session.
		Select("ivt.*, st.*, loc.*, prd.*").
		Table("erp_inventory").Alias("ivt").
		Join("INNER", []string{"erp_stock", "st"}, "ivt.stock_id = st.id").
		Join("INNER", []string{"erp_storage_loc", "loc"}, "ivt.storage_loc_id = loc.id").
		Join("INNER", []string{"erp_product", "prd"}, "ivt.product_id = prd.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("ivt.id")
	}

	inventory_list := make([]ErpInventoryViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpInventoryViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_list,
		Total: count,
	})
}

type InventoryDetailForm struct {
	Inventory models.InventoryInfo `json:"inventory"`
}

func (f *InventoryDetailForm) IsCreate() bool {
	return f.Inventory.Id == 0
}

func (f *InventoryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Inventory.StockId).Message("仓库不能为空！")

	validation.Required(f.Inventory.StorageLocId).Message("库位不能为空！")

	validation.Required(f.Inventory.ProductId).Message("商品不能为空！")

	validation.Required(f.Inventory.Qty).Message("数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventory) DetailData() revel.Result {
	db_session := c.DbSession

	var inventory_id int64
	c.Params.Bind(&inventory_id, "id")

	var inventory models.InventoryInfo
	if inventory_id != 0 {
		has, err := db_session.Id(inventory_id).Get(&inventory)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品库存不存在！"})
		}
	}

	var detail_form InventoryDetailForm
	detail_form.Inventory = inventory

	return c.RenderJson(detail_form)
}

func (c ErpInventory) Save() revel.Result {
	db_session := c.DbSession

	var detail_form InventoryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory := detail_form.Inventory

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&inventory)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(inventory.Id).Update(&inventory)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventory) Delete() revel.Result {
	db_session := c.DbSession

	inventory_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_id_list, "id_list")

	affected, err := db_session.In("id", inventory_id_list).Delete(new(models.InventoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
