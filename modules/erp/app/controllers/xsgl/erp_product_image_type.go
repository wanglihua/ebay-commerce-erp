package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductImageType struct {
	*revel.Controller
}

func (c ErpProductImageType) Index() revel.Result {
	return c.RenderTemplate("erp/xsgl/product_image_type/product_image_type_index.html")
}

func (c ErpProductImageType) ListData() revel.Result {
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

	product_image_type_list := make([]models.ProductImageTypeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_image_type_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductImageTypeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_image_type_list,
		Total: count,
	})
}

type ProductImageTypeDetailForm struct {
	ProductImageType models.ProductImageTypeInfo `json:"product_image_type"`
}

func (f *ProductImageTypeDetailForm) IsCreate() bool {
	return f.ProductImageType.Id == 0
}

func (f *ProductImageTypeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductImageType.Name).Message("名称不能为空！")
	if f.ProductImageType.Name != "" {
		validation.MinSize(f.ProductImageType.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductImageType.Name != "" {
		validation.MaxSize(f.ProductImageType.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ProductImageType.Desc.Valid == true {
		validation.MinSize(f.ProductImageType.Desc.String, 2).Message("描述长度不能小于2！")
	}
	if f.ProductImageType.Desc.Valid == true {
		validation.MaxSize(f.ProductImageType.Desc.String, 400).Message("描述长度不能大于400！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductImageType) DetailData() revel.Result {
	db_session := c.DbSession

	var product_image_type_id int64
	c.Params.Bind(&product_image_type_id, "id")

	var product_image_type models.ProductImageTypeInfo
	if product_image_type_id != 0 {
		has, err := db_session.Id(product_image_type_id).Get(&product_image_type)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品图片类型不存在！"})
		}
	}

	var detail_form ProductImageTypeDetailForm
	detail_form.ProductImageType = product_image_type

	return c.RenderJson(detail_form)
}

func (c ErpProductImageType) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductImageTypeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_image_type := detail_form.ProductImageType

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_image_type)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_image_type.Id).Update(&product_image_type)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductImageType) Delete() revel.Result {
	db_session := c.DbSession

	product_image_type_id_list := make([]int64, 0)
	c.Params.Bind(&product_image_type_id_list, "id_list")

	for _, product_image_type_id := range product_image_type_id_list {
		if product_image_type_id == models.PRODUCT_IMAGE_TYPE_CCTP ||
			product_image_type_id == models.PRODUCT_IMAGE_TYPE_NRTP ||
			product_image_type_id == models.PRODUCT_IMAGE_TYPE_KCTP ||
			product_image_type_id == models.PRODUCT_IMAGE_TYPE_ZJTP {
			return c.RenderJson(core.JsonResult{Success: false, Message: "内置图片属性：橱窗图片、内容图片、库存图片、质检图片，不能删除!"})
		}
	}

	affected, err := db_session.In("id", product_image_type_id_list).Delete(new(models.ProductImageTypeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
