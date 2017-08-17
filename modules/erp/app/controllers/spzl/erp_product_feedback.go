package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductFeedback struct {
	*revel.Controller
}

func (c ErpProductFeedback) Index() revel.Result {
	return c.RenderTemplate("erp/spzl/product_feedback/product_feedback_index.html")
}

func (c ErpProductFeedback) ListData() revel.Result {
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

	product_feedback_list := make([]models.ProductFeedbackInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_feedback_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductFeedbackInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_feedback_list,
		Total: count,
	})
}

type ProductFeedbackDetailForm struct {
	ProductFeedback models.ProductFeedbackInfo      `json:"product_feedback"`
}

func (f *ProductFeedbackDetailForm) IsCreate() bool {
	return f.ProductFeedback.Id == 0
}

func (f *ProductFeedbackDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductFeedback.CateId).Message("反馈类别不能为空！")
	validation.Required(f.ProductFeedback.ProductId).Message("商品不能为空！")
	validation.Required(f.ProductFeedback.SupplierId).Message("供应商不能为空！")
	validation.Required(f.ProductFeedback.Code).Message("反馈编码不能为空！")

	if f.ProductFeedback.Description != "" {
		validation.MaxSize(f.ProductFeedback.Description, 800).Message("描述长度不能大于800！")
	}

	validation.Required(f.ProductFeedback.UserId).Message("反馈人id不能为空！")

	validation.Required(f.ProductFeedback.UserLoginName).Message("反馈人工号不能为空！")
	if f.ProductFeedback.UserLoginName != "" {
		validation.MinSize(f.ProductFeedback.UserLoginName, 2).Message("反馈人工号长度不能小于2！")
	}
	if f.ProductFeedback.UserLoginName != "" {
		validation.MaxSize(f.ProductFeedback.UserLoginName, 100).Message("反馈人工号长度不能大于100！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductFeedback) DetailData() revel.Result {
	db_session := c.DbSession

	var product_feedback_id int64
	c.Params.Bind(&product_feedback_id, "id")

	var product_feedback models.ProductFeedbackInfo
	if product_feedback_id != 0 {
		has, err := db_session.Id(product_feedback_id).Get(&product_feedback)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品质量反馈不存在！"})
		}
	}

	var detail_form ProductFeedbackDetailForm
	detail_form.ProductFeedback = product_feedback

	var product_feedback_cate_list = make([]models.ProductFeedbackCateInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&product_feedback_cate_list))

	var supplier_list = make([]models.SupplierInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&supplier_list))

	return c.RenderJson(map[string]interface{}{
		"detail":                     detail_form,
		"product_feedback_cate_list": product_feedback_cate_list,
		"supplier_list":              supplier_list,
	})
}

func (c ErpProductFeedback) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductFeedbackDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_feedback := detail_form.ProductFeedback

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_feedback)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_feedback.Id).Update(&product_feedback)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductFeedback) Delete() revel.Result {
	db_session := c.DbSession

	product_feedback_id_list := make([]int64, 0)
	c.Params.Bind(&product_feedback_id_list, "id_list")

	affected, err := db_session.In("id", product_feedback_id_list).Delete(new(models.ProductFeedbackInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
