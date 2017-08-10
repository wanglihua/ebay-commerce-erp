package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpProductStyle struct {
	*revel.Controller
}

func (c ErpProductStyle) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STYLE, c.Session)

	return c.RenderTemplate("erp/sppz/product_style/product_style_index.html")
}

func (c ErpProductStyle) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STYLE, c.Session)

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

	product_style_list := make([]models.ProductStyleInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_style_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductStyleInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_style_list,
		Total: count,
	})
}

type ProductStyleDetailForm struct {
	ProductStyle models.ProductStyleInfo `json:"product_style"`
}

func (f *ProductStyleDetailForm) IsCreate() bool {
	return f.ProductStyle.Id == 0
}

func (f *ProductStyleDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductStyle.Name).Message("名称不能为空！")
	if f.ProductStyle.Name != "" {
		validation.MinSize(f.ProductStyle.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductStyle.Name != "" {
		validation.MaxSize(f.ProductStyle.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ProductStyle.Description == "" {
		validation.MaxSize(f.ProductStyle.Description, 800).Message("描述长度不能大于800！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductStyle) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STYLE, c.Session)

	db_session := c.DbSession

	var product_style_id int64
	c.Params.Bind(&product_style_id, "id")

	var product_style models.ProductStyleInfo
	if product_style_id != 0 {
		has, err := db_session.Id(product_style_id).Get(&product_style)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品款式不存在！"})
		}
	}

	var detail_form ProductStyleDetailForm
	detail_form.ProductStyle = product_style

	return c.RenderJson(detail_form)
}

func (c ErpProductStyle) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STYLE, c.Session)

	db_session := c.DbSession

	var detail_form ProductStyleDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_style := detail_form.ProductStyle

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_style)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_style.Id).Update(&product_style)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductStyle) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STYLE, c.Session)
	
	db_session := c.DbSession

	product_style_id_list := make([]int64, 0)
	c.Params.Bind(&product_style_id_list, "id_list")

	affected, err := db_session.In("id", product_style_id_list).Delete(new(models.ProductStyleInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
