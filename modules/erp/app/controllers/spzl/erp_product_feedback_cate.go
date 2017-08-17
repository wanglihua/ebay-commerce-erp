package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"strings"
)

type ErpProductFeedbackCate struct {
	*revel.Controller
}

func (c ErpProductFeedbackCate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_CATE, c.Session)

	return c.RenderTemplate("erp/spzl/product_feedback_cate/product_feedback_cate_index.html")
}

func (c ErpProductFeedbackCate) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_CATE, c.Session)

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

	product_feedback_cate_list := make([]models.ProductFeedbackCateInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_feedback_cate_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductFeedbackCateInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_feedback_cate_list,
		Total: count,
	})
}

type ProductFeedbackCateDetailForm struct {
	ProductFeedbackCate models.ProductFeedbackCateInfo `json:"product_feedback_cate"`
}

func (f *ProductFeedbackCateDetailForm) IsCreate() bool {
	return f.ProductFeedbackCate.Id == 0
}

func (f *ProductFeedbackCateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductFeedbackCate.Name).Message("名称不能为空！")
	if f.ProductFeedbackCate.Name != "" {
		validation.MinSize(f.ProductFeedbackCate.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductFeedbackCate.Name != "" {
		validation.MaxSize(f.ProductFeedbackCate.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ProductFeedbackCate.Description != "" {
		validation.MaxSize(f.ProductFeedbackCate.Description, 800).Message("描述长度不能大于800！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductFeedbackCate) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_CATE, c.Session)

	db_session := c.DbSession

	var product_feedback_cate_id int64
	c.Params.Bind(&product_feedback_cate_id, "id")

	var product_feedback_cate models.ProductFeedbackCateInfo
	if product_feedback_cate_id != 0 {
		has, err := db_session.Id(product_feedback_cate_id).Get(&product_feedback_cate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品质量反馈类别不存在！"})
		}
	}

	var detail_form ProductFeedbackCateDetailForm
	detail_form.ProductFeedbackCate = product_feedback_cate

	return c.RenderJson(detail_form)
}

func (c ErpProductFeedbackCate) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_CATE, c.Session)

	db_session := c.DbSession

	var detail_form ProductFeedbackCateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_feedback_cate := detail_form.ProductFeedbackCate

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_feedback_cate)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_feedback_cate.Id).Update(&product_feedback_cate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductFeedbackCate) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_CATE, c.Session)

	db_session := c.DbSession

	product_feedback_cate_id_list := make([]int64, 0)
	c.Params.Bind(&product_feedback_cate_id_list, "id_list")

	var feedback_count struct {
		Count int64 `xorm:"'feedback_count'"`
	}

	/*
SELECT count(pf.id) AS feedback_count
FROM erp_product_feedback pf
WHERE pf.cate_id IN (123, 234, 345);
	 */

	var cate_feedback_sql = `SELECT count(pf.id) AS feedback_count
FROM erp_product_feedback pf
WHERE pf.cate_id IN (%s);`
	cate_feedback_sql = fmt.Sprintf(cate_feedback_sql, strings.Join(core.Int_list_to_str_list(product_feedback_cate_id_list), ","))
	_, err := db_session.Sql(cate_feedback_sql).Get(&feedback_count)
	core.HandleError(err)

	if feedback_count.Count > 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "有类别在商品质量反馈信息中被引用，不能被删除！",
		})
	}

	affected, err := db_session.In("id", product_feedback_cate_id_list).Delete(new(models.ProductFeedbackCateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
