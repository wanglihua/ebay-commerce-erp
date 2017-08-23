package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/sppz_models"
)

type ErpProductCert struct {
	*revel.Controller
}

func (c ErpProductCert) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CERT, c.Session)

	return c.RenderTemplate("erp/sppz/product_cert/product_cert_index.html")
}

func (c ErpProductCert) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CERT, c.Session)

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

	product_cert_list := make([]sppz_models.ProductCertInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_cert_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(sppz_models.ProductCertInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_cert_list,
		Total: count,
	})
}

type ProductCertDetailForm struct {
	ProductCert sppz_models.ProductCertInfo `json:"product_cert"`
}

func (f *ProductCertDetailForm) IsCreate() bool {
	return f.ProductCert.Id == 0
}

func (f *ProductCertDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductCert.Name).Message("名称不能为空！")
	if f.ProductCert.Name != "" {
		validation.MinSize(f.ProductCert.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductCert.Name != "" {
		validation.MaxSize(f.ProductCert.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ProductCert.Description == "" {
		validation.MaxSize(f.ProductCert.Description, 800).Message("描述长度不能大于800！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductCert) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CERT, c.Session)

	db_session := c.DbSession

	var product_cert_id int64
	c.Params.Bind(&product_cert_id, "id")

	var product_cert sppz_models.ProductCertInfo
	if product_cert_id != 0 {
		has, err := db_session.Id(product_cert_id).Get(&product_cert)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品认证不存在！"})
		}
	}

	var detail_form ProductCertDetailForm
	detail_form.ProductCert = product_cert

	return c.RenderJson(detail_form)
}

func (c ErpProductCert) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CERT, c.Session)

	db_session := c.DbSession

	var detail_form ProductCertDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_cert := detail_form.ProductCert

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_cert)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_cert.Id).Update(&product_cert)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductCert) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CERT, c.Session)

	db_session := c.DbSession

	product_cert_id_list := make([]int64, 0)
	c.Params.Bind(&product_cert_id_list, "id_list")

	affected, err := db_session.In("id", product_cert_id_list).Delete(new(sppz_models.ProductCertInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
