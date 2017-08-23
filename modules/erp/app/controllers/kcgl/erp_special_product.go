package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpSpecialProduct struct {
	*revel.Controller
}

func (c ErpSpecialProduct) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	return c.RenderTemplate("erp/kcgl/special_product/special_product_index.html")
}

func (c ErpSpecialProduct) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type special_product_view_item_t struct {
		SpecialProduct kcgl_models.SpecialProductInfo `xorm:"extends" json:"sp"`
		Product        kcgl_models.ProductInfo        `xorm:"extends" json:"p"`
	}

	/*
SELECT
    sp.*,
    p.*
FROM erp_special_product sp
    INNER JOIN erp_product p ON p.id = sp.product_id
	 */

	query := db_session.
	Select("sp.*,p.*").
		Table("erp_special_product").Alias("sp").
		Join("INNER", []string{"erp_product", "p"}, "p.id = sp.product_id").
		Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("sp.id")
	}

	special_product_list := make([]special_product_view_item_t, 0, limit)
	err := data_query.Limit(limit, offset).Find(&special_product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(special_product_view_item_t))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  special_product_list,
		Total: count,
	})
}

type SpecialProductDetailForm struct {
	SkuList     string         `json:"sku_list"`
	SpecialId   int64          `json:"special_id"`
	Sku         string         `json:"sku"`
	Description string         `json:"description"`
	IsCreate    bool           `json:"is_create"`
}

func (f *SpecialProductDetailForm) Valid(validation *revel.Validation) bool {
	return validation.HasErrors() == false
}

func (c ErpSpecialProduct) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	var db_session = c.DbSession

	var id int64
	c.Params.Bind(&id, "id")

	var detail_form SpecialProductDetailForm

	if id != 0 {
		var special_product kcgl_models.SpecialProductInfo
		_, err := db_session.Id(id).Get(&special_product)
		core.HandleError(err)

		var product kcgl_models.ProductInfo
		_, err = db_session.Id(special_product.ProductId).Get(&product)
		core.HandleError(err)

		detail_form.SpecialId = special_product.Id
		detail_form.Sku = product.Sku
		detail_form.Description = special_product.Description
		detail_form.IsCreate = false
	} else {
		detail_form.IsCreate = true
	}

	return c.RenderJson(detail_form)
}

func (c ErpSpecialProduct) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	var db_session = c.DbSession

	var detail_form SpecialProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	if detail_form.SpecialId == 0 {
		var sku = detail_form.Sku

		var product kcgl_models.ProductInfo
		has, err := db_session.Where("sku = ?", sku).Get(&product)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("商品SKU %s 不存在！", sku)})
		}

		count, err := db_session.Where("product_id = ?", product.Id).Count(&kcgl_models.SpecialProductInfo{})
		core.HandleError(err)

		if count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("商品SKU %s 之前已经被设置为特殊SKU！", sku)})
		}

		var special_product kcgl_models.SpecialProductInfo
		special_product.ProductId = product.Id
		special_product.Description = detail_form.Description

		_, err = db_session.Insert(&special_product)
		core.HandleError(err)
	} else {
		var special_product kcgl_models.SpecialProductInfo
		_, err := db_session.Where("id = ?", detail_form.SpecialId).Get(&special_product)
		core.HandleError(err)

		special_product.Description = detail_form.Description

		_, err = db_session.Id(detail_form.SpecialId).Cols("description").Update(&special_product)
		core.HandleError(err)
	}

	return c.RenderJson(detail_form)
}

func (c ErpSpecialProduct) BatchAdd() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	db_session := c.DbSession

	var detail_form SpecialProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var sku_list_str = strings.TrimSpace(detail_form.SkuList)
	if sku_list_str == "" {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "SKU列表不能为空！",
		})
	}

	var sku_list = strings.Split(sku_list_str, ";")
	if len(sku_list) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "SKU列表为空！"})
	}

	for _, sku := range sku_list {
		if strings.TrimSpace(sku) != "" {
			var product kcgl_models.ProductInfo
			has, err := db_session.Where("sku = ?", sku).Get(&product)
			core.HandleError(err)

			if has == false {
				return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("商品SKU %s 不存在！", sku)})
			}

			var special_product kcgl_models.SpecialProductInfo
			count, err := db_session.Where("product_id = ?", product.Id).Count(&kcgl_models.SpecialProductInfo{})
			core.HandleError(err)

			if count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("商品SKU %s 之前已经被设置为特殊SKU！", sku)})
			}

			special_product.ProductId = product.Id
			_, err = db_session.Insert(&special_product)
			core.HandleError(err)
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}

func (c ErpSpecialProduct) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_SPECIAL_SKU, c.Session)

	db_session := c.DbSession

	special_product_id_list := make([]int64, 0)
	c.Params.Bind(&special_product_id_list, "id_list")

	affected, err := db_session.In("id", special_product_id_list).Delete(new(kcgl_models.SpecialProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
