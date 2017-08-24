package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpSupplierPayMethod struct {
	*revel.Controller
}

func (c ErpSupplierPayMethod) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_PAY_METHOD, c.Session)

	return c.RenderTemplate("erp/sppz/supplier_pay_method/supplier_pay_method_index.html")
}

func (c ErpSupplierPayMethod) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_PAY_METHOD, c.Session)

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

	supplier_pay_method_list := make([]sppz_models.SupplierPayMethodInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&supplier_pay_method_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(sppz_models.SupplierPayMethodInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  supplier_pay_method_list,
		Total: count,
	})
}

type SupplierPayMethodDetailForm struct {
	SupplierPayMethod sppz_models.SupplierPayMethodInfo `json:"pay_method"`
}

func (f *SupplierPayMethodDetailForm) IsCreate() bool {
	return f.SupplierPayMethod.Id == 0
}

func (f *SupplierPayMethodDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.SupplierPayMethod.Name).Message("名称不能为空！")
	if f.SupplierPayMethod.Name != "" {
		validation.MinSize(f.SupplierPayMethod.Name, 2).Message("名称长度不能小于2！")
	}
	if f.SupplierPayMethod.Name != "" {
		validation.MaxSize(f.SupplierPayMethod.Name, 100).Message("名称长度不能大于100！")
	}

	if f.SupplierPayMethod.Description == "" {
		validation.MaxSize(f.SupplierPayMethod.Description, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpSupplierPayMethod) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_PAY_METHOD, c.Session)

	db_session := c.DbSession

	var supplier_pay_method_id int64
	c.Params.Bind(&supplier_pay_method_id, "id")

	var supplier_pay_method sppz_models.SupplierPayMethodInfo
	if supplier_pay_method_id != 0 {
		has, err := db_session.Id(supplier_pay_method_id).Get(&supplier_pay_method)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的供应商付款方式不存在！"})
		}
	}

	var detail_form SupplierPayMethodDetailForm
	detail_form.SupplierPayMethod = supplier_pay_method

	return c.RenderJson(detail_form)
}

func (c ErpSupplierPayMethod) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_PAY_METHOD, c.Session)

	db_session := c.DbSession

	var detail_form SupplierPayMethodDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	supplier_pay_method := detail_form.SupplierPayMethod

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", supplier_pay_method.Name).Count(new(sppz_models.SupplierPayMethodInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Insert(&supplier_pay_method)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", supplier_pay_method.Id, supplier_pay_method.Name).Count(new(sppz_models.SupplierPayMethodInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Id(supplier_pay_method.Id).Update(&supplier_pay_method)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpSupplierPayMethod) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER_PAY_METHOD, c.Session)

	db_session := c.DbSession

	supplier_pay_method_id_list := make([]int64, 0)
	c.Params.Bind(&supplier_pay_method_id_list, "id_list")

	affected, err := db_session.In("id", supplier_pay_method_id_list).Delete(new(sppz_models.SupplierPayMethodInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
