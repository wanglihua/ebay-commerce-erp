package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductUser struct {
	*revel.Controller
}

func (c ErpProductUser) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_USER, c.Session)

	return c.RenderTemplate("erp/sppz/product_user/product_user_index.html")
}

func (c ErpProductUser) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_USER, c.Session)

	var db_session = c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var filter_type int64
	c.Params.Bind(&filter_type, "filter_type")

	type product_user_view_item_t struct {
		Product kcgl_models.ProductInfo `xorm:"extends" json:"p"`
		Status  models.StatusInfo       `xorm:"extends" json:"s"`
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

	if filter_type == 1 { // 未设置人员数据的记录
		filter_has_zero_data := core.SqlOr("sale_user_id = 0", "service_user_id = 0", "dev_user_id = 0")
		filter = core.SqlAnd(filter, filter_has_zero_data)
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	var product_list = make([]product_user_view_item_t, 0, limit)
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

type product_user_detail_form_t struct {
	Product kcgl_models.ProductInfo `json:"product"`
}

func (c ErpProductUser) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_USER, c.Session)

	var db_session = c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_user_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductUser) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_USER, c.Session)

	var db_session = c.DbSession

	var detail_form product_user_detail_form_t
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

	// var product_in_db_befor_modify = product_in_db

	product_in_db.SaleUserId = detail_form.Product.SaleUserId
	product_in_db.ServiceUserId = detail_form.Product.ServiceUserId
	product_in_db.DevUserId = detail_form.Product.DevUserId

	updata_col_List := []string{"sale_user_id", "service_user_id", "dev_user_id"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品人员数据未能更新！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}
