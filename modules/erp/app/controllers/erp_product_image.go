package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductImage struct {
	*revel.Controller
}

func (c ErpProductImage) Index() revel.Result {
	return c.RenderTemplate("erp/product_image/product_image_index.html")
}

func (c ErpProductImage) ListData() revel.Result {
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

	product_image_list := make([]models.ProductImageInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_image_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductImageInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_image_list,
		Total: count,
	})
}

type ProductImageDetailForm struct {
	ProductImage         models.ProductImageInfo        `json:"product_image"`
	ProductImageTypeList []models.ProductImageTypeInfo  `json:"product_image_type_list"`
	PlatformList         []models.ECommercePlatformInfo `json:"platform_list"`
}

func (f *ProductImageDetailForm) IsCreate() bool {
	return f.ProductImage.Id == 0
}

func (f *ProductImageDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductImage.ProductCode).Message("商品编码不能为空！")
	if f.ProductImage.ProductCode != "" {
		validation.MinSize(f.ProductImage.ProductCode, 2).Message("商品编码长度不能小于2！")
	}
	if f.ProductImage.ProductCode != "" {
		validation.MaxSize(f.ProductImage.ProductCode, 100).Message("商品编码长度不能大于100！")
	}

	validation.Required(f.ProductImage.TypeId).Message("图片属性不能为空！")

	validation.Required(f.ProductImage.PlatformId).Message("所属平台不能为空！")

	validation.Required(f.ProductImage.Url).Message("图片链接不能为空！")
	if f.ProductImage.Url != "" {
		validation.MinSize(f.ProductImage.Url, 2).Message("图片链接长度不能小于2！")
	}
	if f.ProductImage.Url != "" {
		validation.MaxSize(f.ProductImage.Url, 300).Message("图片链接长度不能大于300！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductImage) DetailData() revel.Result {
	db_session := c.DbSession

	var product_image_id int64
	c.Params.Bind(&product_image_id, "id")

	var product_image models.ProductImageInfo
	if product_image_id != 0 {
		has, err := db_session.Id(product_image_id).Get(&product_image)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品图片不存在！"})
		}
	}

	var product_image_type_list []models.ProductImageTypeInfo
	err := db_session.Find(&product_image_type_list)
	core.HandleError(err)

	var platform_list []models.ECommercePlatformInfo
	err = db_session.Find(&platform_list)
	core.HandleError(err)

	var detail_form ProductImageDetailForm
	detail_form.ProductImage = product_image
	detail_form.ProductImageTypeList = product_image_type_list
	detail_form.PlatformList = platform_list

	return c.RenderJson(detail_form)
}

func (c ErpProductImage) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductImageDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_image := detail_form.ProductImage

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_image)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_image.Id).Update(&product_image)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductImage) Delete() revel.Result {
	db_session := c.DbSession

	product_image_id_list := make([]int64, 0)
	c.Params.Bind(&product_image_id_list, "id_list")

	affected, err := db_session.In("id", product_image_id_list).Delete(new(models.ProductImageInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
