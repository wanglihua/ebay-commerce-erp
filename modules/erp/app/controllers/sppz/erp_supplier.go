package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/region_service"
	"fmt"
	"strings"
	"matrix/modules/erp/models/sppz_models"
)

type ErpSupplier struct {
	*revel.Controller
}

func (c ErpSupplier) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	return c.RenderTemplate("erp/sppz/supplier/supplier_index.html")
}

type supplier_view_item_t struct {
	Supplier sppz_models.SupplierInfo `xorm:"extends" json:"s"`
	Region   sppz_models.RegionInfo   `xorm:"extends" json:"r"`
}

func (c ErpSupplier) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

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

type SupplierDetailForm struct {
	Supplier     sppz_models.SupplierInfo       `json:"supplier"`
	ProvinceId   int64                     `json:"province_id"`
	ProvinceList []region_service.Region_t `json:"province_list"`
	CityList     []region_service.Region_t `json:"city_list"`
	// city_id 就是 supplier 中的 region_id
}

func (f *SupplierDetailForm) IsCreate() bool {
	return f.Supplier.Id == 0
}

func (f *SupplierDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Supplier.Name).Message("名称不能为空！")
	if f.Supplier.Name != "" {
		validation.MinSize(f.Supplier.Name, 2).Message("名称长度不能小于2！")
	}
	if f.Supplier.Name != "" {
		validation.MaxSize(f.Supplier.Name, 100).Message("名称长度不能大于100！")
	}

	if f.Supplier.Contact.Valid == true {
		validation.MinSize(f.Supplier.Contact.String, 2).Message("联系人长度不能小于2！")
	}
	if f.Supplier.Contact.Valid == true {
		validation.MaxSize(f.Supplier.Contact.String, 100).Message("联系人长度不能大于100！")
	}

	if f.Supplier.ContactInfo.Valid == true {
		validation.MaxSize(f.Supplier.ContactInfo.String, 500).Message("联系方式长度不能大于500！")
	}

	if f.Supplier.Address.Valid == true {
		validation.MinSize(f.Supplier.Address.String, 2).Message("地址长度不能小于2！")
	}
	if f.Supplier.Address.Valid == true {
		validation.MaxSize(f.Supplier.Address.String, 200).Message("地址长度不能大于200！")
	}

	if f.Supplier.Phone.Valid == true {
		validation.MinSize(f.Supplier.Phone.String, 2).Message("电话长度不能小于2！")
	}
	if f.Supplier.Phone.Valid == true {
		validation.MaxSize(f.Supplier.Phone.String, 100).Message("电话长度不能大于100！")
	}

	if f.Supplier.Desc.Valid == true {
		validation.MaxSize(f.Supplier.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpSupplier) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	db_session := c.DbSession

	var supplier_id int64
	c.Params.Bind(&supplier_id, "id")

	var detail_form SupplierDetailForm
	var province_list = region_service.Get_province_list(db_session)
	detail_form.ProvinceList = province_list

	var supplier sppz_models.SupplierInfo
	if supplier_id != 0 {
		has, err := db_session.Id(supplier_id).Get(&supplier)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的供应商不存在！"})
		}

		var province_id = region_service.Get_region_parent_id(supplier.RegionId, db_session)
		detail_form.ProvinceId = province_id

		var city_list = region_service.Get_city_list(province_id, db_session)
		detail_form.CityList = city_list
	}

	detail_form.Supplier = supplier

	return c.RenderJson(detail_form)
}

func (c ErpSupplier) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	db_session := c.DbSession

	var detail_form SupplierDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	supplier := detail_form.Supplier

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", supplier.Name).Count(new(sppz_models.SupplierInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Insert(&supplier)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", supplier.Id, supplier.Name).Count(new(sppz_models.SupplierInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Id(supplier.Id).Update(&supplier)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpSupplier) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	db_session := c.DbSession

	supplier_id_list := make([]int64, 0)
	c.Params.Bind(&supplier_id_list, "id_list")

	// 检查供应商是否在使用中

	// 检查在商品反馈中的使用

	var feedback_count struct {
		Count int64 `xorm:"'feedback_count'"`
	}
	/*
SELECT count(pf.id) AS feedback_count
FROM erp_product_feedback pf
WHERE pf.supplier_id IN (123, 234, 345);
	 */

	var supplier_feedback_sql = `SELECT count(pf.id) AS feedback_count
FROM erp_product_feedback pf
WHERE pf.supplier_id IN (%s);`
	supplier_feedback_sql = fmt.Sprintf(supplier_feedback_sql, strings.Join(core.Int_list_to_str_list(supplier_id_list), ","))
	_, err := db_session.Sql(supplier_feedback_sql).Get(&feedback_count)
	core.HandleError(err)

	if feedback_count.Count > 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "有供应商在商品质量反馈信息中被引用，不能被删除！",
		})
	}

	// todo: 检查在商品详情中的使用

	affected, err := db_session.In("id", supplier_id_list).Delete(new(sppz_models.SupplierInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
