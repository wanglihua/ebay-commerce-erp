package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
)

type ErpProductCheckinWeight struct {
	*revel.Controller
}

func (c ErpProductCheckinWeight) Index() revel.Result {
	return c.RenderTemplate("erp/product_checkin_weight/product_checkin_weight_index.html")
}

func (c ErpProductCheckinWeight) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_checkin_weight_list := make([]cggl_models.ProductCheckinWeightInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_checkin_weight_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(cggl_models.ProductCheckinWeightInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_checkin_weight_list,
		Total: count,
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

	validation.Required(f.ProductCheckinWeight.Length).Message("长度不能为空！")

	validation.Required(f.ProductCheckinWeight.Width).Message("宽度不能为空！")

	validation.Required(f.ProductCheckinWeight.Height).Message("高度不能为空！")

	validation.Required(f.ProductCheckinWeight.Weight).Message("重量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductCheckinWeight) Detail(entry_id int64) revel.Result {
	c.RenderArgs["entry_id"] = entry_id
	return c.RenderTemplate("erp/product_checkin_weight/product_checkin_weight_detail.html")
}

func (c ErpProductCheckinWeight) DetailData() revel.Result {
	db_session := c.DbSession

	var product_checkin_weight_id int64
	c.Params.Bind(&product_checkin_weight_id, "id")

	var product_checkin_weight cggl_models.ProductCheckinWeightInfo
	if product_checkin_weight_id != 0 {
		has, err := db_session.Id(product_checkin_weight_id).Get(&product_checkin_weight)
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

	product_checkin_weight := detail_form.ProductCheckinWeight

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_checkin_weight)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_checkin_weight.Id).Update(&product_checkin_weight)
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
