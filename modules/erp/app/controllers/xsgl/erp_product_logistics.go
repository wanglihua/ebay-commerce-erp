package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/xsgl_models"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductLogistics struct {
	*revel.Controller
}

func (c ErpProductLogistics) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_LOGISTICS, c.Session)

	return c.RenderTemplate("erp/xsgl/product_logistics/product_logistics_index.html")
}

func (c ErpProductLogistics) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_LOGISTICS, c.Session)

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

	product_logistics_list := make([]xsgl_models.ProductLogisticsInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_logistics_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xsgl_models.ProductLogisticsInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_logistics_list,
		Total: count,
	})
}

type BrandDetailForm struct {
	ProductLogistics xsgl_models.ProductLogisticsInfo `json:"product_logistics"`
}

func (f *BrandDetailForm) IsCreate() bool {
	return f.ProductLogistics.Id == 0
}

func (f *BrandDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductLogistics.ProductCode).Message("SKU名称不能为空！")

	if f.ProductLogistics.Qty == 0 {
		validation.Error("数量不能为空！")
	}

	if f.ProductLogistics.Qty < 0 {
		validation.Error("数量不能为负！")
	}

	validation.Required(f.ProductLogistics.PostalCode).Message("地址邮编不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductLogistics) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_LOGISTICS, c.Session)

	// db_session := c.DbSession

	var detail_form BrandDetailForm
	detail_form.ProductLogistics = xsgl_models.ProductLogisticsInfo{}

	return c.RenderJson(detail_form)
}

func (c ErpProductLogistics) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_LOGISTICS, c.Session)

	db_session := c.DbSession

	var detail_form BrandDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_logistics_from_ui := detail_form.ProductLogistics
	var product_code = product_logistics_from_ui.PostalCode
	var qty = product_logistics_from_ui.Qty
	var postal_code = product_logistics_from_ui.PostalCode

	var product kcgl_models.ProductInfo
	_, err = db_session.Where("sku = ?", product_code).Get(&product)
	core.HandleError(err)

	_, err = db_session.Insert(&product_logistics_from_ui)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
}

func (c ErpProductLogistics) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_LOGISTICS, c.Session)

	db_session := c.DbSession

	brand_id_list := make([]int64, 0)
	c.Params.Bind(&brand_id_list, "id_list")

	affected, err := db_session.In("id", brand_id_list).Delete(new(sppz_models.BrandInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
