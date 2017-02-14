package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/service/product_cate_service"
)

type ErpProduct struct {
	*revel.Controller
}

func (c ErpProduct) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product/product_index.html")
}

type ErpProductViewItem struct {
	models.ProductInfo          `xorm:"extends" json:"prd"`
	models.ProductCateInfo      `xorm:"extends" json:"cate"`
	models.SupplierInfo         `xorm:"extends" json:"sup"`
	models.PackageMaterialInfo  `xorm:"extends" json:"pm"`
	models.StatusInfo           `xorm:"extends" json:"st"`
}

func (c ErpProduct) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT prd.*, cate.*, sup.*, pm.*, st.*
FROM erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("prd.*, cate.*, sup.*, pm.*, st.*")
	sql_builder.Count("count(prd.id) as product_count")
	sql_builder.From(`erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id`)

	sql_builder.Where(filter)

	//query extra filter here

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]ErpProductViewItem, 0, limit)
	err := db_session.Sql(sql_builder.SqlList()).Find(&product_list)
	core.HandleError(err)

	var product_count struct {
		Count int64 `xorm:"'product_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

type ProductDetailForm struct {
	Product                   models.ProductInfo                           `json:"product"`
	SupplierList              []models.SupplierInfo                        `json:"supplier_list"`
	ProductCateSelectItemList []product_cate_service.ProductCateSelectItem `json:"product_cate_list"`
	PackageMaterialList       [] models.PackageMaterialInfo                `json:"package_material_list"`
	StatusList                []models.StatusInfo                          `json:"status_list"`
}

func (f *ProductDetailForm) IsCreate() bool {
	return f.Product.Id == 0
}

func (f *ProductDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Product.Sku).Message("SKU不能为空！")
	if f.Product.Sku != "" {
		validation.MinSize(f.Product.Sku, 2).Message("SKU长度不能小于2！")
	}
	if f.Product.Sku != "" {
		validation.MaxSize(f.Product.Sku, 100).Message("SKU长度不能大于100！")
	}

	validation.Required(f.Product.Name).Message("名称不能为空！")
	if f.Product.Name != "" {
		validation.MinSize(f.Product.Name, 2).Message("名称长度不能小于2！")
	}
	if f.Product.Name != "" {
		validation.MaxSize(f.Product.Name, 100).Message("名称长度不能大于100！")
	}

	validation.Required(f.Product.CostPrice).Message("成本价不能为空！")

	validation.Required(f.Product.Weight).Message("重量不能为空！")

	validation.Required(f.Product.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProduct) DetailData() revel.Result {
	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product models.ProductInfo
	if product_id != 0 {
		has, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品不存在！"})
		}
	}

	if product.Id == 0 {
		product.StatusId = models.PRODUCT_STATUS_ZC
	}

	var supplier_list []models.SupplierInfo
	err := db_session.Asc(models.SupplierCols.Name).Find(&supplier_list)
	core.HandleError(err)

	var product_cate_list = product_cate_service.GetProductCateSelectItemList(db_session)

	var package_material_list [] models.PackageMaterialInfo
	err = db_session.Asc(models.PackageMaterialCols.Name).Find(&package_material_list)
	core.HandleError(err)

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_PRODUCT)

	var detail_form ProductDetailForm
	detail_form.Product = product
	detail_form.SupplierList = supplier_list
	detail_form.ProductCateSelectItemList = product_cate_list
	detail_form.PackageMaterialList = package_material_list
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpProduct) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product := detail_form.Product

	var affected int64
	if detail_form.IsCreate() {
		sku_count, err := db_session.Where("sku = ?", product.Sku).Count(&product)
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！" })
		}

		name_count, err := db_session.Where("name = ?", product.Name).Count(&product)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&product)
		core.HandleError(err)
	} else {
		sku_count, err := db_session.Where("id <> ? and sku = ?", product.Id, product.Sku).Count(&product)
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！" })
		}

		name_count, err := db_session.Where("id <> ? and name = ?", product.Id, product.Name).Count(&product)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(product.Id).Update(&product)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProduct) Delete() revel.Result {
	db_session := c.DbSession

	product_id_list := make([]int64, 0)
	c.Params.Bind(&product_id_list, "id_list")

	affected, err := db_session.In("id", product_id_list).Delete(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
