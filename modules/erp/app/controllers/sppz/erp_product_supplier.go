package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"fmt"
	"matrix/modules/erp/models/sppz_models"
	"strings"
)

type ErpProductSupplier struct {
	*revel.Controller
}

func (c ErpProductSupplier) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

	return c.RenderTemplate("erp/sppz/product_supplier/product_supplier_index.html")
}

func (c ErpProductSupplier) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

	var db_session = c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var supplier_name_query string
	c.Params.Bind(&supplier_name_query, "supplier_name_query")
	supplier_name_query = strings.TrimSpace(supplier_name_query)

	type product_supplier_view_item_t struct {
		Product  kcgl_models.ProductInfo  `xorm:"extends" json:"p"`
		Supplier sppz_models.SupplierInfo `xorm:"extends" json:"sup"`
		Status   models.StatusInfo        `xorm:"extends" json:"s"`
	}

	/*
SELECT p.*, sup.*, s.*
FROM erp_product p
    LEFT JOIN erp_supplier sup ON sup.id = p.supplier_id
    INNER JOIN erp_status s ON s.id = p.status_id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("p.*, sup.*, s.*")
	sb.Count("count(p.id) AS product_count")
	sb.From(`erp_product p
    LEFT JOIN erp_supplier sup ON sup.id = p.supplier_id
    INNER JOIN erp_status s ON s.id = p.status_id`)

	if supplier_name_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("sup.name LIKE '%%%s%%'", supplier_name_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	var product_supplier_list = make([]product_supplier_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_supplier_list)
	core.HandleError(err)

	type product_count_t struct {
		Count int64 `xorm:"'product_count'"`
	}
	var product_count product_count_t
	_, err = db_session.Sql(sb.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_supplier_list,
		Total: product_count.Count,
	})
}

type product_supplier_detail_form_t struct {
	Product      kcgl_models.ProductInfo `json:"product"`
	SupplierName string                  `json:"supplier_name"`
}

func (c ErpProductSupplier) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

	var db_session = c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var supplier sppz_models.SupplierInfo
	_, err = db_session.Id(product.SupplierId).Get(&supplier)
	core.HandleError(err)

	var detail_form product_supplier_detail_form_t
	detail_form.Product = product
	detail_form.SupplierName = supplier.Name

	var supplier_list = make([]sppz_models.SupplierInfo, 0)
	core.HandleError(db_session.OrderBy("convert(name using gbk) ASC").Find(&supplier_list))

	return c.RenderJson(map[string]interface{}{
		"detail":        detail_form,
		"supplier_list": supplier_list,
	})
}

func (c ErpProductSupplier) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

	var db_session = c.DbSession

	var detail_form product_supplier_detail_form_t
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

	product_in_db.SupplierId = detail_form.Product.SupplierId

	updata_col_List := []string{"supplier_id"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品供应商未能更新！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductSupplier) SupplierSelectListData() revel.Result {
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
