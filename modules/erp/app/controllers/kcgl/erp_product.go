
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProduct struct {
	*revel.Controller
}

func (c ErpProduct) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product/product_index.html")
}

func (c ErpProduct) ListData() revel.Result {
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

	product_list := make([]models.ProductInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: count,
	})
}

type ProductDetailForm struct {
	Product models.ProductInfo `json:"product"`
}

func (f *ProductDetailForm) IsCreate() bool {
	return f.Product.Id == 0
}

func (f *ProductDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.Product.Sku).Message("SKU不能为空！")
	if f.Product.Sku != "" {
		validation.MinSize(f.Product.Sku, 2).Message("SKU长度不能小于2！")
	}
	if f.Product.Sku != "" {
		validation.MaxSize(f.Product.Sku, 100).Message("SKU长度不能大于100！")
	}

	validation.Required(f.Product.Name).Message("名称不能为空！")
	if f.Product.Name != "" {
		validation.MinSize(f.Product.Name, 2).Message("名称长度不能小于2！")
	}
	if f.Product.Name != "" {
		validation.MaxSize(f.Product.Name, 100).Message("名称长度不能大于100！")
	}

	validation.Required(f.Product.EnglishName).Message("英文名称不能为空！")
	if f.Product.EnglishName != "" {
		validation.MinSize(f.Product.EnglishName, 2).Message("英文名称长度不能小于2！")
	}
	if f.Product.EnglishName != "" {
		validation.MaxSize(f.Product.EnglishName, 100).Message("英文名称长度不能大于100！")
	}

	validation.Required(f.Product.CostPrice).Message("成本价不能为空！")

	validation.Required(f.Product.Weight).Message("重量不能为空！")

	validation.Required(f.Product.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProduct) DetailData() revel.Result {
	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product models.ProductInfo
	if product_id != 0 {
		has, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品不存在！"})
		}
	}

	var detail_form ProductDetailForm
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProduct) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product := detail_form.Product

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&product)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(product.Id).Update(&product)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProduct) Delete() revel.Result {
	db_session := c.DbSession

	product_id_list := make([]int64, 0)
	c.Params.Bind(&product_id_list, "id_list")

	affected, err := db_session.In("id", product_id_list).Delete(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

