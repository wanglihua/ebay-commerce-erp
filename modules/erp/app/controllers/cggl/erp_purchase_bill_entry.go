package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpPurchaseBillEntry struct {
	*revel.Controller
}

func (c ErpPurchaseBillEntry) Index() revel.Result {
	return c.RenderTemplate("erp/cggl/purchase_bill_entry/purchase_bill_entry_index.html")
}

type ErpPurchaseBillEntryViewItem struct {
	models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	models.ProductInfo           `xorm:"extends" json:"p"`
	models.StatusInfo            `xorm:"extends" json:"s"`
}

func (c ErpPurchaseBillEntry) ListData() revel.Result {
	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	var purchase_bill models.PurchaseBillInfo
	if purchase_bill_id != 0 {
		has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
		}
	}

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

	sb.Where(core.SqlAnd(filter, fmt.Sprintf("bill_id = %d", purchase_bill_id)))

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pbe.id ASC")
	}

	sb.Limit(offset, limit)

	purchase_bill_entry_list := make([]ErpPurchaseBillEntryViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	var purchase_bill_entry_count struct {
		Count int64 `xorm:"'entry_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&purchase_bill_entry_count)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"purchase_bill_code": func() string {
			if purchase_bill.Code == "" {
				return "请先单击采购单行"
			}
			return purchase_bill.Code
		}(),
		"list_data": core.GridResult{
			Data:  purchase_bill_entry_list,
			Total: purchase_bill_entry_count.Count,
		},
	})
}

type PurchaseBillEntryDetailForm struct {
	PurchaseBill      models.PurchaseBillInfo      `json:"purchase_bill"`
	PurchaseBillEntry models.PurchaseBillEntryInfo `json:"purchase_bill_entry"`
	ProductName       string                       `json:"product_name"`
}

func (f *PurchaseBillEntryDetailForm) IsCreate() bool {
	return f.PurchaseBillEntry.Id == 0
}

func (f *PurchaseBillEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PurchaseBillEntry.BillId).Message("批次不能为空！")

	validation.Required(f.PurchaseBillEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.PurchaseBillEntry.Qty).Message("数量不能为空！")

	validation.Required(f.PurchaseBillEntry.Price).Message("单价不能为空！")

	validation.Required(f.PurchaseBillEntry.Amount).Message("金额不能为空！")

	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MinSize(f.PurchaseBillEntry.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MaxSize(f.PurchaseBillEntry.Desc.String, 400).Message("备注长度不能大于400！")
	}

	return validation.HasErrors() == false
}

func (c ErpPurchaseBillEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	if purchase_bill_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请先选择采购单！"})
	}

	var purchase_bill models.PurchaseBillInfo
	has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
	}

	var purchase_bill_entry_id int64
	c.Params.Bind(&purchase_bill_entry_id, "id")

	var purchase_bill_entry models.PurchaseBillEntryInfo
	var product_name string
	if purchase_bill_entry_id != 0 {
		has, err := db_session.Id(purchase_bill_entry_id).Get(&purchase_bill_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单明细不存在！"})
		}

		var product models.ProductInfo
		product_id := purchase_bill_entry.ProductId
		_, err = db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		product_name = product.Name
	}

	var detail_form PurchaseBillEntryDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.PurchaseBillEntry = purchase_bill_entry
	detail_form.ProductName = product_name

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseBillEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PurchaseBillEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	purchase_bill_entry := detail_form.PurchaseBillEntry
	purchase_bill_entry.BillId = detail_form.PurchaseBill.Id

	purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(purchase_bill_entry.Qty)

	var affected int64
	if detail_form.IsCreate() {
		billId_count, err := db_session.Where("bill_id = ? and product_id = ?", purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(&purchase_bill_entry)
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次、商品已存在！" })
		}

		purchase_bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YCJ

		affected, err = db_session.Insert(&purchase_bill_entry)
		core.HandleError(err)
	} else {
		billId_count, err := db_session.Where("id <> ? and bill_id = ? and product_id = ?", purchase_bill_entry.Id, purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(&purchase_bill_entry)
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次、商品已存在！" })
		}

		affected, err = db_session.Id(purchase_bill_entry.Id).Update(&purchase_bill_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBillEntry) Delete() revel.Result {
	db_session := c.DbSession

	purchase_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_entry_id_list, "id_list")

	affected, err := db_session.In("id", purchase_bill_entry_id_list).Delete(new(models.PurchaseBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpPurchaseBillEntry) ProductSelectListData() revel.Result {
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
