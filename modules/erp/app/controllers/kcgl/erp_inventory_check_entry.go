package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"strings"
)

type ErpInventoryCheckEntry struct {
	*revel.Controller
}

func (c ErpInventoryCheckEntry) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/inventory_check_entry/inventory_check_entry_index.html")
}

func (c ErpInventoryCheckEntry) ListData() revel.Result {
	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	var inventory_check_bill models.InventoryCheckInfo
	if bill_id != 0 {
		has, err := db_session.Id(bill_id).Get(&inventory_check_bill)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的盘点单不存在！"})
		}
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpInventoryCheckEntryViewItem struct {
		models.InventoryCheckEntryInfo `xorm:"extends" json:"ice"`
		models.InventoryCheckInfo      `xorm:"extends" json:"ic"`
		models.StorageLocInfo          `xorm:"extends" json:"sl"`
		models.ProductInfo             `xorm:"extends" json:"p"`
	}

	/*
	SELECT ice.*, ic.*, sl.*, p.*
	FROM erp_inventory_check_entry ice
		INNER JOIN erp_inventory_check ic ON ice.check_id = ic.id
		INNER JOIN erp_storage_loc sl ON ice.storage_loc_id = sl.id
		INNER JOIN erp_product p ON ice.product_id = p.id
	*/

	sb := core.NewSqlBuilder()
	sb.Select("ice.*, ic.*, sl.*, p.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_inventory_check_entry ice
		INNER JOIN erp_inventory_check ic ON ice.check_id = ic.id
		INNER JOIN erp_storage_loc sl ON ice.storage_loc_id = sl.id
		INNER JOIN erp_product p ON ice.product_id = p.id`)

	sb.Where(core.SqlAnd(filter, fmt.Sprintf("ice.check_id = %d", bill_id)))

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ice.id ASC")
	}

	sb.Limit(offset, limit)

	inventory_check_entry_list := make([]ErpInventoryCheckEntryViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_check_entry_list)
	core.HandleError(err)

	var entry_count struct {
		Count int64 `xorm:"'entry_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&entry_count)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"bill_code": func() string {
			if inventory_check_bill.Code == "" {
				return "请先单击盘点单行"
			}
			return inventory_check_bill.Code
		}(),
		"list_data": core.GridResult{
			Data:  inventory_check_entry_list,
			Total: entry_count.Count,
		},
	})
}

type InventoryCheckEntryDetailForm struct {
	InventoryCheck      models.InventoryCheckInfo      `json:"inventory_check"`
	InventoryCheckEntry models.InventoryCheckEntryInfo `json:"check_entry"`
	ProductName         string                         `json:"product_name"`
	StockName           string                         `json:"stock_name"`
	StorageLocList      []models.StorageLocInfo        `json:"storage_loc_list"`
}

func (f *InventoryCheckEntryDetailForm) IsCreate() bool {
	return f.InventoryCheckEntry.Id == 0
}

func (f *InventoryCheckEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.InventoryCheckEntry.CheckId).Message("盘点批次不能为空！")

	validation.Required(f.InventoryCheckEntry.StorageLocId).Message("库位不能为空！")

	validation.Required(f.InventoryCheckEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.InventoryCheckEntry.PreQty).Message("盘点前库存不能为空！")

	validation.Required(f.InventoryCheckEntry.CheckQty).Message("盘点库存不能为空！")

	validation.Required(f.InventoryCheckEntry.QtyDiff).Message("库存差值不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventoryCheckEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	if bill_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请先选择盘点单！"})
	}

	var inventory_check models.InventoryCheckInfo
	has, err := db_session.Id(bill_id).Get(&inventory_check)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的盘点单不存在！"})
	}

	var stock models.StockInfo
	_, err = db_session.Id(inventory_check.StockId).Get(&stock)

	storage_loc_list := make([]models.StorageLocInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_id = %d", inventory_check.StockId)).OrderBy("name asc").Find(&storage_loc_list)
	core.HandleError(err)

	var inventory_check_entry_id int64
	c.Params.Bind(&inventory_check_entry_id, "id")

	var inventory_check_entry models.InventoryCheckEntryInfo
	var product_name string
	if inventory_check_entry_id != 0 {
		has, err := db_session.Id(inventory_check_entry_id).Get(&inventory_check_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的库存盘点明细不存在！"})
		}

		var product models.ProductInfo
		product_id := inventory_check_entry.ProductId
		_, err = db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		product_name = product.Name
	}

	var detail_form InventoryCheckEntryDetailForm
	detail_form.InventoryCheck = inventory_check
	detail_form.InventoryCheckEntry = inventory_check_entry
	detail_form.ProductName = product_name
	detail_form.StockName = stock.Name
	detail_form.StorageLocList = storage_loc_list

	return c.RenderJson(detail_form)
}

func (c ErpInventoryCheckEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form InventoryCheckEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory_check_entry := detail_form.InventoryCheckEntry
	inventory_check_entry.CheckId = detail_form.InventoryCheck.Id

	check_id := detail_form.InventoryCheck.Id
	var inventory_check models.InventoryCheckInfo
	_, err = db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_YCJ && inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_WTG {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单不是已创建或未通过状态，信息不能修改！"})
	}

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&inventory_check_entry)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(inventory_check_entry.Id).Update(&inventory_check_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventoryCheckEntry) Delete() revel.Result {
	db_session := c.DbSession

	inventory_check_entry_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_check_entry_id_list, "id_list")

	var check_id_list = make([]struct {
		CheckId int64 `xorm:"'check_id'" json:"check_id"`
	}, 0)

	entry_id_str_list := make([]string, 0)
	for _, entry_id_int := range entry_id_str_list {
		entry_id_str_list = append(entry_id_str_list, entry_id_int)
	}

	err := db_session.Sql(fmt.Sprintf(`SELECT check_id FROM erp_inventory_check_entry WHERE id IN (%s)`, strings.Join(entry_id_str_list, ","))).Find(&check_id_list)
	core.HandleError(err)

	if len(check_id_list) != 1 {
		return c.RenderJson(core.JsonResult{Success: true, Message: "请选择一个盘点单下的数据进行删除!"})
	}

	check_id := check_id_list[0]
	var inventory_check models.InventoryCheckInfo
	_, err = db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_YCJ && inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_WTG {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单不是已创建或未通过状态，明细信息不能删除！"})
	}

	affected, err := db_session.In("id", inventory_check_entry_id_list).Delete(new(models.InventoryCheckEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpInventoryCheckEntry) ProductSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpProductViewItem struct {
		models.ProductInfo          `xorm:"extends" json:"prd"`
		models.ProductCateInfo      `xorm:"extends" json:"cate"`
		models.SupplierInfo         `xorm:"extends" json:"sup"`
		models.PackageMaterialInfo  `xorm:"extends" json:"pm"`
		models.StatusInfo           `xorm:"extends" json:"st"`
	}

	/*
SELECT prd.*, cate.*, sup.*, pm.*, st.*
FROM erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("prd.*, cate.*, sup.*, pm.*, st.*")
	sql_builder.Count("count(prd.id) as product_count")
	sql_builder.From(`erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id`)

	sql_builder.Where(filter)

	//query extra filter here

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]ErpProductViewItem, 0, limit)
	err := db_session.Sql(sql_builder.SqlList()).Find(&product_list)
	core.HandleError(err)

	var product_count struct {
		Count int64 `xorm:"'product_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}
