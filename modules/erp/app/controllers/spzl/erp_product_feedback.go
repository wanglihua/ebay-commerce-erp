package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strings"
)

type ErpProductFeedback struct {
	*revel.Controller
}

func (c ErpProductFeedback) SkuList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	return c.RenderTemplate("erp/spzl/product_feedback/product_feedback_sku_list.html")
}

func (c ErpProductFeedback) SkuListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	db_session := c.DbSession

	var feedback_sku string
	c.Params.Bind(&feedback_sku, "feedback_sku")
	feedback_sku = strings.TrimSpace(feedback_sku)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT DISTINCT p.*
FROM erp_product_feedback pf
    INNER JOIN erp_product p ON pf.product_id = p.id
WHERE p.sku = '12344';

SELECT count(DISTINCT p.id) AS sku_count
FROM erp_product_feedback pf
    INNER JOIN erp_product p ON pf.product_id = p.id
WHERE p.sku = '12344';
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`DISTINCT p.*`)
	sb.Count("count(DISTINCT p.id) AS sku_count")
	sb.From(`erp_product_feedback pf
    INNER JOIN erp_product p ON pf.product_id = p.id`)

	if feedback_sku != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sku = '%s'", feedback_sku))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	product_list := make([]models.ProductInfo, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_list)
	core.HandleError(err)

	var sku_count struct {
		Count int64 `xorm:"'sku_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&sku_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: sku_count.Count,
	})
}

func (c ErpProductFeedback) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	return c.RenderTemplate("erp/spzl/product_feedback/product_feedback_index.html")
}

func (c ErpProductFeedback) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var db_session = c.DbSession

	/*
		"page": {
			"feedback_cate_id": 0,
			"feedback_cate_list": [],
			"feedback_sku": "",
			"feedback_supplier_id": 0,
			"feedback_supplier_list": []
		}
	 */

	var feedback_cate_list = make([]models.ProductFeedbackCateInfo, 0)
	core.HandleError(db_session.Find(&feedback_cate_list))

	var feedback_supplier_list = make([]models.SupplierInfo, 0)
	core.HandleError(db_session.Find(&feedback_supplier_list))

	return c.RenderJson(map[string]interface{}{
		"success":                true,
		"feedback_cate_id":       0,
		"feedback_cate_list":     feedback_cate_list,
		"feedback_sku":           "",
		"feedback_supplier_id":   0,
		"feedback_supplier_list": feedback_supplier_list,
	})
}

func (c ErpProductFeedback) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var feedback_code string
	c.Params.Bind(&feedback_code, "feedback_code")
	feedback_code = strings.TrimSpace(feedback_code)

	var feedback_cate_id int64
	c.Params.Bind(&feedback_cate_id, "feedback_cate_id")

	var feedback_sku string
	c.Params.Bind(&feedback_sku, "feedback_sku")
	feedback_sku = strings.TrimSpace(feedback_sku)

	var feedback_supplier_id int64
	c.Params.Bind(&feedback_supplier_id, "feedback_supplier_id")

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_feedback_view_item_t struct {
		ProductFeedback     models.ProductFeedbackInfo     `xorm:"extends" json:"pf"`
		ProductFeedbackCate models.ProductFeedbackCateInfo `xorm:"extends" json:"c"`
		Product             models.ProductInfo             `xorm:"extends" json:"p"`
		Supplier            models.SupplierInfo            `xorm:"extends" json:"sup"`
	}

	/*
SELECT
    pf.*,
    c.*,
    p.*,
    sup.*
FROM erp_product_feedback pf
    INNER JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    INNER JOIN erp_product p ON p.id = pf.product_id
    INNER JOIN erp_supplier sup ON sup.id = pf.supplier_id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`pf.*,
    c.*,
    p.*,
    sup.*`)
	sb.Count("count(pf.id) AS feedback_count")
	sb.From(`erp_product_feedback pf
    INNER JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    INNER JOIN erp_product p ON p.id = pf.product_id
    INNER JOIN erp_supplier sup ON sup.id = pf.supplier_id`)

	if feedback_cate_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.cate_id = %d", feedback_cate_id))
	}

	if feedback_sku != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sku = '%s'", feedback_sku))
	}

	if feedback_supplier_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.supplier_id = %d", feedback_supplier_id))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pf.id ASC")
	}

	sb.Limit(offset, limit)

	product_feedback_list := make([]product_feedback_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_feedback_list)
	core.HandleError(err)

	var feedback_count struct {
		Count int64 `xorm:"'feedback_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&feedback_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_feedback_list,
		Total: feedback_count.Count,
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
	validation.Required(f.ProductFeedback.SupplierId).Message("供应商不能为空！")

	if f.ProductFeedback.Description != "" {
		validation.MaxSize(f.ProductFeedback.Description, 800).Message("描述长度不能大于800！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductFeedback) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

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

	// 获取商品的SKU
	var product models.ProductInfo
	_, err := db_session.Id(product_feedback.ProductId).Get(&product)
	core.HandleError(err)

	product_feedback.Sku = product.Sku

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
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	db_session := c.DbSession

	var detail_form ProductFeedbackDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_feedback := detail_form.ProductFeedback

	var product models.ProductInfo
	has, err := db_session.Where("sku = ?", product_feedback.Sku).Get(&product)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "输入的SKU不存在！",
		})
	}

	product_feedback.ProductId = product.Id

	var affected int64
	if detail_form.IsCreate() {
		var login_user = core.GetLoginUser(c.Session)
		product_feedback.UserId = login_user.UserId
		product_feedback.UserLoginName = login_user.LoginName

		affected, err = db_session.Insert(&product_feedback)
		core.HandleError(err)

		//product_feedback.Code = fmt.Sprintf("FB%04d", product_feedback.Id)
		//_, err = db_session.Id(product_feedback.Id).Cols("code").Update(&product_feedback)
		//core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_feedback.Id).Cols("cate_id", "product_id", "supplier_id", "description").Update(&product_feedback)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductFeedback) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	db_session := c.DbSession

	product_feedback_id_list := make([]int64, 0)
	c.Params.Bind(&product_feedback_id_list, "id_list")

	affected, err := db_session.In("id", product_feedback_id_list).Delete(new(models.ProductFeedbackInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
