package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"encoding/json"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductHsCode struct {
	*revel.Controller
}

func (c ErpProductHsCode) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

	return c.RenderTemplate("erp/hgbg/product_hs_code/product_hs_code_index.html")
}

func (c ErpProductHsCode) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var filter_type int64
	c.Params.Bind(&filter_type, "filter_type")

	type product_weight_view_item_t struct {
		Product kcgl_models.ProductInfo `xorm:"extends" json:"p"`
		Status  models.StatusInfo  `xorm:"extends" json:"s"`
	}

	/*
		SELECT p.*, s.*
		FROM erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("p.*, s.*")
	sb.Count("count(p.id) AS product_count")
	sb.From(`erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id`)

	if filter_type == 1 { //有零值数据的记录
		filter_has_zero_data := core.SqlOr("cn_hs_code = ''", "au_hs_code = ''")
		filter = core.SqlAnd(filter, filter_has_zero_data)
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	var product_list = make([]product_weight_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_list)
	core.HandleError(err)

	type product_count_t struct {
		Count int64 `xorm:"'product_count'"`
	}
	var product_count product_count_t
	_, err = db_session.Sql(sb.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

type product_hs_code_detail_form_t struct {
	Product kcgl_models.ProductInfo `json:"product"`
}

func (c ErpProductHsCode) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_hs_code_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductHsCode) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

	db_session := c.DbSession

	var detail_form product_hs_code_detail_form_t
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var product_id = detail_form.Product.Id

	var product_in_db kcgl_models.ProductInfo
	has, err := db_session.Id(product_id).Get(&product_in_db)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品不存在！",
		})
	}

	var hs_code_cn hgbg_models.HsCodeCnReportElemInfo
	has, err = db_session.Where("cn_name = ?", detail_form.Product.CnDeclName).Get(&hs_code_cn)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "中国海关申报品名不存在！",
		})
	}

	var hs_code_au hgbg_models.HsCodeAuDutyInfo
	has, err = db_session.Where("en_name = ?", detail_form.Product.AuDeclName).Get(&hs_code_au)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "澳洲海关申报品名不存在！",
		})
	}

	product_in_db.CnDeclName = detail_form.Product.CnDeclName
	product_in_db.AuDeclName = detail_form.Product.AuDeclName
	// product_in_db.CnHsCode = detail_form.Product.CnHsCode
	// product_in_db.AuHsCode = detail_form.Product.AuHsCode

	// updata_col_List := []string{"cn_hs_code", "au_hs_code", "cn_decl_name", "au_decl_name"}
	updata_col_List := []string{"cn_decl_name", "au_decl_name"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品海关编码数据未能更新！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}
