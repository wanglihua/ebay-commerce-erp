package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"fmt"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/spzl_models"
	"matrix/modules/erp/models/sppz_models"
	"strings"
)

type ErpProductFeedbackQuery struct {
	*revel.Controller
}

func (c ErpProductFeedbackQuery) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	return c.RenderTemplate("erp/spzl/product_feedback_query/product_feedback_query_index.html")
}

func (c ErpProductFeedbackQuery) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	var db_session = c.DbSession

	/*
		"page": {
			"feedback_cate_id": 0,
			"feedback_cate_list": [],
			"feedback_sku": "",
			"feedback_supplier_id": 0,
			"supplier_list": []
		}
	 */

	var feedback_cate_list = make([]spzl_models.ProductFeedbackCateInfo, 0)
	core.HandleError(db_session.OrderBy("convert(name using gbk) ASC").Find(&feedback_cate_list))

	var supplier_list = make([]sppz_models.SupplierInfo, 0)
	core.HandleError(db_session.OrderBy("convert(name using gbk) ASC").Find(&supplier_list))

	return c.RenderJson(map[string]interface{}{
		"success":            true,
		"sku":                "",
		"feedback_cate_id":   0,
		"feedback_cate_list": feedback_cate_list,
		"supplier_id":        0,
		"supplier_list":      supplier_list,
	})
}

func (c ErpProductFeedbackQuery) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	var sku string
	c.Params.Bind(&sku, "sku")
	sku = strings.TrimSpace(sku)

	var feedback_cate_id int64
	c.Params.Bind(&feedback_cate_id, "feedback_cate_id")

	var supplier_id int64
	c.Params.Bind(&supplier_id, "supplier_id")

	var db_session = c.DbSession

	var product kcgl_models.ProductInfo
	if sku != "" {
		_, err := db_session.Where("sku = ?", sku).Get(&product)
		core.HandleError(err)
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_feedback_view_item_t struct {
		ProductFeedback     spzl_models.ProductFeedbackInfo     `xorm:"extends" json:"pf"`
		ProductFeedbackCate spzl_models.ProductFeedbackCateInfo `xorm:"extends" json:"c"`
		Product             kcgl_models.ProductInfo             `xorm:"extends" json:"p"`
		Supplier            sppz_models.SupplierInfo            `xorm:"extends" json:"sup"`
	}

	/*
SELECT
    pf.*, c.*, p.*, sup.*
FROM erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select(`pf.*, c.*, p.*, sup.*`)
	sb.Count("count(pf.id) AS feedback_count")
	sb.From(`erp_product_feedback pf
    INNER JOIN erp_product p ON p.id = pf.product_id
    LEFT JOIN erp_product_feedback_cate c ON c.id = pf.cate_id
    LEFT JOIN erp_supplier sup ON sup.id = pf.supplier_id`)

	if sku != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.id = %d", product.Id))
	}

	if feedback_cate_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.cate_id = %d", feedback_cate_id))
	}

	if supplier_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pf.supplier_id = %d", supplier_id))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pf.id DESC")
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
