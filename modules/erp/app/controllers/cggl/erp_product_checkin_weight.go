package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductCheckinWeight struct {
	*revel.Controller
}

func (c ErpProductCheckinWeight) Index() revel.Result {
	return c.RenderTemplate("erp/product_checkin_weight/product_checkin_weight_index.html")
}

func (c ErpProductCheckinWeight) ListData() revel.Result {
	var db_session = c.DbSession

	type product_checkin_weight_item_t struct {
		PurchaseBill         cggl_models.PurchaseBillInfo         `xorm:"extends" json:"pb"`
		ProductCheckinWeight cggl_models.ProductCheckinWeightInfo `xorm:"extends" json:"cw"`
	}

	/*
SELECT
    pb.*, cw.*
FROM erp_product_checkin_weight cw
    INNER JOIN erp_purchase_bill_entry be ON cw.entry_id = be.id
    INNER JOIN erp_purchase_bill pb ON be.bill_id = pb.id;
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var sb = core.NewSqlBuilder()
	sb.Select("pb.*, cw.*")
	sb.Count("count(pb.id) AS count")
	sb.From(`erp_product_checkin_weight cw
    INNER JOIN erp_purchase_bill_entry be ON cw.entry_id = be.id
    INNER JOIN erp_purchase_bill pb ON be.bill_id = pb.id`)
	sb.Where(filter)
	sb.Limit(offset, limit)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("id ASC")
	}

	var product_checkin_weight_list = make([]product_checkin_weight_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_checkin_weight_list)
	core.HandleError(err)

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_checkin_weight_list,
		Total: count.Count,
	})
}

type ProductCheckinWeightDetailForm struct {
	ProductCheckinWeight cggl_models.ProductCheckinWeightInfo `json:"product_checkin_weight"`
}

func (f *ProductCheckinWeightDetailForm) IsCreate() bool {
	return f.ProductCheckinWeight.Id == 0
}

func (f *ProductCheckinWeightDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductCheckinWeight.BillId).Message("采购单不能为空！")
	validation.Required(f.ProductCheckinWeight.EntryId).Message("采购单明细不能为空！")
	validation.Required(f.ProductCheckinWeight.ProductId).Message("商品不能为空！")

	// validation.Required(f.ProductCheckinWeight.Length).Message("长度不能为空！")
	// validation.Required(f.ProductCheckinWeight.Width).Message("宽度不能为空！")
	// validation.Required(f.ProductCheckinWeight.Height).Message("高度不能为空！")
	// validation.Required(f.ProductCheckinWeight.Weight).Message("重量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductCheckinWeight) Detail(entry_id int64) revel.Result {
	var db_session = c.DbSession

	/*
SELECT p.*
FROM erp_purchase_bill_entry be
    INNER JOIN erp_product p ON be.product_id = p.id
WHERE be.id = 123;
	 */
	var sql_get_product = `
SELECT p.*
FROM erp_purchase_bill_entry be
    INNER JOIN erp_product p ON be.product_id = p.id
WHERE be.id = %d;`
	sql_get_product = fmt.Sprintf(sql_get_product, entry_id)

	var product kcgl_models.ProductInfo
	_, err := db_session.Sql(sql_get_product).Get(&product)
	core.HandleError(err)

	c.RenderArgs["entry_id"] = entry_id
	c.RenderArgs["product"] = product
	return c.RenderTemplate("erp/product_checkin_weight/product_checkin_weight_detail.html")
}

func (c ErpProductCheckinWeight) DetailData() revel.Result {
	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var product_checkin_weight cggl_models.ProductCheckinWeightInfo
	if entry_id != 0 {
		has, err := db_session.Where("entry_id = ?", entry_id).Get(&product_checkin_weight)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品入库长宽高重量不存在！"})
		}
	}

	var detail_form ProductCheckinWeightDetailForm
	detail_form.ProductCheckinWeight = product_checkin_weight

	return c.RenderJson(detail_form)
}

func (c ErpProductCheckinWeight) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductCheckinWeightDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var product_checkin_weight_from_ui = detail_form.ProductCheckinWeight

	var entry_id = product_checkin_weight_from_ui.EntryId
	var product_checkin_weight cggl_models.ProductCheckinWeightInfo
	has, err := db_session.Where("entry_id = ?", entry_id).Get(&product_checkin_weight)
	core.HandleError(err)

	var affected int64
	if !has {
		affected, err = db_session.Insert(&product_checkin_weight)
		core.HandleError(err)
	} else {
		affected, err = db_session.Where("entry_id = ?", product_checkin_weight_from_ui.EntryId).Update(&product_checkin_weight)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductCheckinWeight) Delete() revel.Result {
	db_session := c.DbSession

	product_checkin_weight_id_list := make([]int64, 0)
	c.Params.Bind(&product_checkin_weight_id_list, "id_list")

	affected, err := db_session.In("id", product_checkin_weight_id_list).Delete(new(cggl_models.ProductCheckinWeightInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
