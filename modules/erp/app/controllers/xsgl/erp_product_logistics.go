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
	"matrix/modules/erp/service/product_logistics_service"
	"matrix/modules/erp/models/ddpz_models"
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
	var product_code = product_logistics_from_ui.ProductCode
	var qty = product_logistics_from_ui.Qty
	var postal_code = product_logistics_from_ui.PostalCode

	revel.TRACE.Println("product_code:", product_code)

	var product kcgl_models.ProductInfo
	has, err := db_session.Where("sku = ?", product_code).Get(&product)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的SKU商品不存在！"})
	}

	err, result := product_logistics_service.Match_product_qty_postage(product.Id, qty, postal_code, db_session)
	core.HandleError(err)

	var product_logistics xsgl_models.ProductLogisticsInfo

	has, err = db_session.Where("product_code = ? AND qty = ?", product_code, qty).Get(&product_logistics)
	core.HandleError(err)

	product_logistics.ProductId = product.Id
	product_logistics.ProductCode = product.Sku
	product_logistics.Qty = qty
	product_logistics.PostalCode = postal_code

	var login_user = core.GetLoginUser(c.Session)
	product_logistics.UserId = login_user.UserId
	product_logistics.LoginName = login_user.LoginName

	if result.Final_amount < 0 {
		product_logistics.ModeId = 0
		product_logistics.ModeName = ""
		product_logistics.FeePostalCode = ""
		product_logistics.Fee = 0
		product_logistics.Cost = 0
		product_logistics.Success = 0

		product_logistics.StockId = 0
		product_logistics.StockName = ""
		product_logistics.CompanyId = 0
		product_logistics.CompanyName = ""
	} else {
		product_logistics.ModeId = result.Mode_id
		product_logistics.ModeName = result.Mode.Name

		if result.Fee_rate_id != 0 {
			product_logistics.PostalCode = result.Fee_rate.Postcode
		} else if result.Fee_rate2_id != 0 {
			product_logistics.PostalCode = result.Fee_rate2.Postcode
		}

		product_logistics.Fee = result.Total_fee
		product_logistics.Cost = result.Pickup_fee
		product_logistics.Success = 1

		// stock
		product_logistics.StockId = result.Mode.StockId

		var stock sppz_models.StockInfo
		_, err = db_session.Id(result.Mode.StockId).Get(&stock)
		core.HandleError(err)
		product_logistics.StockName = stock.Name

		// company
		product_logistics.CompanyId = result.Mode.CompanyId

		var company ddpz_models.LogisticsCompanyInfo
		_, err = db_session.Id(result.Mode.CompanyId).Get(&company)
		core.HandleError(err)
		product_logistics.CompanyName = company.Name
	}

	if has == false {
		_, err = db_session.Insert(&product_logistics)
		core.HandleError(err)
	} else {
		_, err = db_session.Id(product_logistics.Id).Update(&product_logistics)
		core.HandleError(err)
	}

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
