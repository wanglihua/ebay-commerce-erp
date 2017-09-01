package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/service/product_log_service"
)

type ErpProductPurchaseSetting struct {
	*revel.Controller
}

func (c ErpProductPurchaseSetting) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

	return c.RenderTemplate("erp/cggl/product_purchase_setting/product_purchase_setting_index.html")
}

func (c ErpProductPurchaseSetting) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

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
		filter_has_zero_data := core.SqlOr("purchase_days = 0", "stock_alert_days = 0", "stock_alert_qty = 0", "daily_mean_sales = 0")
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

type product_purchase_setting_detail_form_t struct {
	Product kcgl_models.ProductInfo `json:"product"`
}

func (c ErpProductPurchaseSetting) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_purchase_setting_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductPurchaseSetting) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PRODUCT_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form product_purchase_setting_detail_form_t
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

	var product_in_db_before_modify = product_in_db

	product_in_db.PurchaseDays = detail_form.Product.PurchaseDays
	product_in_db.StockAlertDays = detail_form.Product.StockAlertDays

	updata_col_List := []string{"stock_alert_days", "purchase_days"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品采购设置数据未能更新！",
		})
	}

	product_log_service.Insert_product_log(&detail_form.Product, &product_in_db_before_modify, core.Db_operation_types.Update, db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}
