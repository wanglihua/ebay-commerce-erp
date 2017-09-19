package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"fmt"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strings"
	"time"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/spzl_models"
	"matrix/modules/erp/models/sppz_models"
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

	type product_feedback_sku_view_item_t struct {
		Product            kcgl_models.ProductInfo            `xorm:"extends" json:"p"`
		ProductFeedbackSku spzl_models.ProductFeedbackSkuInfo `xorm:"extends" json:"pfu"`
	}

	/*
SELECT
    p.*,
    pfu.*
FROM erp_product_feedback_sku pfu
    INNER JOIN erp_product p ON pfu.product_id = p.id
WHERE p.sku = '12344'
ORDER BY p.id;

SELECT count(pfu.id) AS sku_count
FROM erp_product_feedback_sku pfu
    INNER JOIN erp_product p ON pfu.product_id = p.id
WHERE p.sku = '12344';
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`p.*, pfu.*`)
	sb.Count("count(pfu.id) AS sku_count")
	sb.From(`erp_product_feedback_sku pfu
    INNER JOIN erp_product p ON pfu.product_id = p.id`)

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

	product_list := make([]product_feedback_sku_view_item_t, 0, limit)
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

func (c ErpProductFeedback) Index(id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var db_session = c.DbSession

	var product_id = id
	c.RenderArgs["product_id"] = product_id

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)
	c.RenderArgs["product"] = product

	return c.RenderTemplate("erp/spzl/product_feedback/product_feedback_index.html")
}

func (c ErpProductFeedback) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var db_session = c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	/*
		"page": {
			"feedback_cate_id": 0,
			"feedback_cate_list": [],
			"feedback_sku": "",
			"feedback_supplier_id": 0,
			"feedback_supplier_list": []
		}
	 */

	var feedback_cate_list = make([]spzl_models.ProductFeedbackCateInfo, 0)
	var filter_cate = `id IN (SELECT DISTINCT cate_id
FROM erp_product_feedback pf
WHERE pf.product_id = %d)`
	filter_cate = fmt.Sprintf(filter_cate, product_id)
	core.HandleError(db_session.Where(filter_cate).Find(&feedback_cate_list))

	var feedback_supplier_list = make([]sppz_models.SupplierInfo, 0)
	var filter_supplier = `id IN (SELECT DISTINCT supplier_id
FROM erp_product_feedback pf
WHERE pf.product_id = %d)`
	filter_supplier = fmt.Sprintf(filter_supplier, product_id)
	core.HandleError(db_session.Where(filter_supplier).Find(&feedback_supplier_list))

	return c.RenderJson(map[string]interface{}{
		"success":                true,
		"feedback_cate_id":       0,
		"feedback_cate_list":     feedback_cate_list,
		"feedback_supplier_id":   0,
		"feedback_supplier_list": feedback_supplier_list,
	})
}

func (c ErpProductFeedback) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	var feedback_cate_id int64
	c.Params.Bind(&feedback_cate_id, "feedback_cate_id")

	var feedback_supplier_id int64
	c.Params.Bind(&feedback_supplier_id, "feedback_supplier_id")

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_feedback_view_item_t struct {
		ProductFeedback     spzl_models.ProductFeedbackInfo     `xorm:"extends" json:"pf"`
		ProductFeedbackCate spzl_models.ProductFeedbackCateInfo `xorm:"extends" json:"c"`
		Product             kcgl_models.ProductInfo             `xorm:"extends" json:"p"`
		Supplier            sppz_models.SupplierInfo            `xorm:"extends" json:"sup"`
	}

	/*
SELECT
    pf.*,
    c.*,
    p.*,
    sup.*
FROM erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`pf.*,
    c.*,
    p.*,
    sup.*`)
	sb.Count("count(pf.id) AS feedback_count")
	sb.From(`erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("p.id = %d", product_id))

	if feedback_cate_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.cate_id = %d", feedback_cate_id))
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

	type cate_sum_t struct {
		CateId   int64  `xorm:"'cate_id'" json:"cate_id"`
		CateName string `xorm:"'cate_name'" json:"cate_name"`
		Count    int64  `xorm:"'pf_count'" json:"count"`
	}
	/*
SELECT
    pf.cate_id AS cate_id,
    cate.name  AS cate_name,
    count(*)   AS pf_count
FROM erp_product_feedback pf
    INNER JOIN erp_product_feedback_cate cate ON pf.cate_id = cate.id
WHERE pf.product_id = 17629
GROUP BY pf.cate_id;
	 */
	var cate_sum_sql = `SELECT
    pf.cate_id AS cate_id,
    cate.name  AS cate_name,
    count(*)   AS pf_count
FROM erp_product_feedback pf
    INNER JOIN erp_product_feedback_cate cate ON pf.cate_id = cate.id
WHERE pf.product_id = %d
GROUP BY pf.cate_id;`
	cate_sum_sql = fmt.Sprintf(cate_sum_sql, product_id)

	var cate_sum_list = make([]cate_sum_t, 0)
	core.HandleError(db_session.Sql(cate_sum_sql).Find(&cate_sum_list))

	return c.RenderJson(map[string]interface{}{
		"grid_data": core.GridResult{
			Data:  product_feedback_list,
			Total: feedback_count.Count,
		},
		"cate_sum_list": cate_sum_list,
	})
}

type ProductFeedbackDetailForm struct {
	ProductFeedback spzl_models.ProductFeedbackInfo `json:"product_feedback"`
	SupplierName    string                          `json:"supplier_name"`
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

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	var product_feedback spzl_models.ProductFeedbackInfo
	if product_feedback_id != 0 {
		has, err := db_session.Id(product_feedback_id).Get(&product_feedback)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品质量反馈不存在！"})
		}
	} else {
		product_feedback.ProductId = product_id
	}

	// 获取商品的SKU
	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_feedback.ProductId).Get(&product)
	core.HandleError(err)

	product_feedback.Sku = product.Sku

	var supplier sppz_models.SupplierInfo
	_, err = db_session.Id(product_feedback.SupplierId).Get(&supplier)
	core.HandleError(err)

	var detail_form ProductFeedbackDetailForm
	detail_form.ProductFeedback = product_feedback
	detail_form.SupplierName = supplier.Name

	var product_feedback_cate_list = make([]spzl_models.ProductFeedbackCateInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&product_feedback_cate_list))

	return c.RenderJson(map[string]interface{}{
		"detail":                     detail_form,
		"product_feedback_cate_list": product_feedback_cate_list,
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

	var product kcgl_models.ProductInfo
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
		affected, err = db_session.Id(product_feedback.Id).Cols("cate_id", "product_id", "supplier_id", "description", "solution").Update(&product_feedback)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	var feedback_sku spzl_models.ProductFeedbackSkuInfo
	has, err = db_session.Where("product_id = ?", product.Id).Get(&feedback_sku)
	core.HandleError(err)

	feedback_sku.ProductId = product.Id

	var time_now = time.Now()
	feedback_sku.LastCreate = core.Time(time_now)
	feedback_sku.LastUpdate = core.Time(time_now)

	if has == false {
		_, err = db_session.Insert(&feedback_sku)
		core.HandleError(err)
	} else {
		_, err = db_session.Id(feedback_sku.Id).Cols("last_update").Update(&feedback_sku) // 只更新修改时间
		core.HandleError(err)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductFeedback) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK, c.Session)

	var has_delete_pass = pass_service.Is_login_user_has_pass(pass_res.SPZL_PRODUCT_FEEDBACK_DELETE, c.Session)

	db_session := c.DbSession

	product_feedback_id_list := make([]int64, 0)
	c.Params.Bind(&product_feedback_id_list, "id_list")

	if len(product_feedback_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择要删除的商品质量反馈！",
		})
	}

	var product_feedback_id = product_feedback_id_list[0]

	var product_feedback spzl_models.ProductFeedbackInfo
	_, err := db_session.Id(product_feedback_id).Get(&product_feedback)
	core.HandleError(err)

	// 检查删除权限
	if has_delete_pass == false {
		// 如果没有删除权限的，就只能删除自己当日创建的
		var sql = `SELECT count(*) AS count
FROM erp_product_feedback pf
WHERE (pf.user_id <> %d) OR (pf.user_id = %d AND pf.create_time < '%s');`
		var login_user = core.GetLoginUser(c.Session)
		var time_now = time.Now()
		sql = fmt.Sprintf(sql, login_user.UserId, login_user.UserId, fmt.Sprintf("%04d-%02d-%02d", time_now.Year(), int(time_now.Month()), time_now.Day()))
		var count core.Count
		_, err := db_session.Sql(sql).Get(&count)
		core.HandleError(err)

		if count.Count > 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "只能删除本人今日创建的反馈！",
			})
		}
	}

	// 删除反馈
	affected, err := db_session.In("id", product_feedback_id_list).Delete(new(spzl_models.ProductFeedbackInfo))
	core.HandleError(err)

	// 更改反馈SKU
	// 一般删除的是同一个SKU下面的
	var product_id = product_feedback.ProductId
	count, err := db_session.Where("product_id = ?", product_id).Count(new(spzl_models.ProductFeedbackInfo))
	core.HandleError(err)

	if count == 0 {
		_, err := db_session.Where("product_id = ?", product_id).Delete(new(spzl_models.ProductFeedbackSkuInfo))
		core.HandleError(err)
	} else {
		var feedback_sku spzl_models.ProductFeedbackSkuInfo
		_, err = db_session.Where("product_id = ?", product_id).Get(&feedback_sku)
		core.HandleError(err)

		var time_now = time.Now()
		feedback_sku.LastUpdate = core.Time(time_now)

		_, err = db_session.Id(feedback_sku.Id).Cols("last_update").Update(&feedback_sku) // 只更新修改时间
		core.HandleError(err)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpProductFeedback) SupplierSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type supplier_view_item_t struct {
		Supplier sppz_models.SupplierInfo `xorm:"extends" json:"s"`
		Region   sppz_models.RegionInfo   `xorm:"extends" json:"r"`
	}

	/*
SELECT s.*, r.*
FROM erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("s.*, r.*")
	sb.Count("count(s.id) AS supplier_count")
	sb.From(`erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("o.id ASC")
	}

	sb.Limit(offset, limit)

	supplier_list := make([]supplier_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&supplier_list)
	core.HandleError(err)

	var supplier_count struct {
		Count int64 `xorm:"'supplier_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&supplier_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  supplier_list,
		Total: supplier_count.Count,
	})
}
