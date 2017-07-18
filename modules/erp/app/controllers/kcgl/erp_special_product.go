package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
	"fmt"
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

	type special_product_view_item_t struct {
		SpecialProduct models.SpecialProductInfo `xorm:"extends" json:"sp"`
		Product        models.ProductInfo        `xorm:"extends" json:"p"`
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
	SkuList string `json:"sku_list"`
}

func (f *SpecialProductDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.SkuList).Message("SKU列表不能为空！")

	return validation.HasErrors() == false
}

func (c ErpSpecialProduct) DetailData() revel.Result {
	var detail_form SpecialProductDetailForm
	detail_form.SkuList = ""

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

	var sku_list_str = detail_form.SkuList
	var sku_list = strings.Split(sku_list_str, ";")
	if len(sku_list) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "SKU列表为空！"})
	}

	for _, sku := range sku_list {
		if strings.TrimSpace(sku) != "" {
			var product models.ProductInfo
			has, err := db_session.Where("sku = ?", sku).Get(&product)
			core.HandleError(err)

			if has == false {
				return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("商品SKU %s 不存在！", sku)})
			}

			var special_product models.SpecialProductInfo
			count, err := db_session.Where("product_id = ?", product.Id).Count(&models.SpecialProductInfo{})
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
	db_session := c.DbSession

	special_product_id_list := make([]int64, 0)
	c.Params.Bind(&special_product_id_list, "id_list")

	affected, err := db_session.In("id", special_product_id_list).Delete(new(models.SpecialProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
