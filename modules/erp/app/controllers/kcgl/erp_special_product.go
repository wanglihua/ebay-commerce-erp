package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpSpecialProduct struct {
	*revel.Controller
}

func (c ErpSpecialProduct) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/special_product/special_product_index.html")
}

func (c ErpSpecialProduct) ListData() revel.Result {
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

	special_product_list := make([]models.SpecialProductInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&special_product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.SpecialProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  special_product_list,
		Total: count,
	})
}

type SpecialProductDetailForm struct {
	SpecialProduct models.SpecialProductInfo `json:"special_product"`
}

func (f *SpecialProductDetailForm) IsCreate() bool {
	return f.SpecialProduct.Id == 0
}

func (f *SpecialProductDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.SpecialProduct.ProductId).Message("ProductId不能为空！")

	return validation.HasErrors() == false
}

func (c ErpSpecialProduct) DetailData() revel.Result {
	db_session := c.DbSession

	var special_product_id int64
	c.Params.Bind(&special_product_id, "id")

	var special_product models.SpecialProductInfo
	if special_product_id != 0 {
		has, err := db_session.Id(special_product_id).Get(&special_product)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的特殊SKU不存在！"})
		}
	}

	var detail_form SpecialProductDetailForm
	detail_form.SpecialProduct = special_product

	return c.RenderJson(detail_form)
}

func (c ErpSpecialProduct) Save() revel.Result {
	db_session := c.DbSession

	var detail_form SpecialProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	special_product := detail_form.SpecialProduct

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&special_product)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(special_product.Id).Update(&special_product)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpSpecialProduct) Delete() revel.Result {
	db_session := c.DbSession

	special_product_id_list := make([]int64, 0)
	c.Params.Bind(&special_product_id_list, "id_list")

	affected, err := db_session.In("id", special_product_id_list).Delete(new(models.SpecialProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
