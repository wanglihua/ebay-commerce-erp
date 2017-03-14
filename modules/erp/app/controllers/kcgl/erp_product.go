package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/product_cate_service"
	"matrix/modules/erp/service/product_inventory_service"
	"matrix/modules/erp/service/product_log_service"
	"matrix/modules/erp/service/status_service"
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
	Product                   models.ProductInfo                                          `json:"product"`
	ProductImageList          []models.ProductImageInfo                                   `json:"product_image_list"`
	SupplierList              []models.SupplierInfo                                       `json:"supplier_list"`
	ProductCateSelectItemList []product_cate_service.ProductCateSelectItem                `json:"product_cate_list"`
	PackageMaterialList       []models.PackageMaterialInfo                                `json:"package_material_list"`
	StatusList                []models.StatusInfo                                         `json:"status_list"`
	InventoryList             []product_inventory_service.Product_inventory_view_item_t   `json:"inventory_list"`
	OrderItemList             []order_service.Order_item_view_item_t                      `json:"order_item_list"`
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

func (c ErpProduct) Detail(product_id int64) revel.Result {
	//return c.RenderText(fmt.Sprintf("product_id %d", product_id))
	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("erp/kcgl/product/product_detail.html")
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

	var product_image_list = make([]models.ProductImageInfo, 0)
	err := db_session.Where(fmt.Sprintf("product_id = %d", product_id)).Find(&product_image_list)
	core.HandleError(err)

	var supplier_list []models.SupplierInfo
	err = db_session.Asc(models.SupplierCols.Name).Find(&supplier_list)
	core.HandleError(err)

	var product_cate_list = product_cate_service.GetProductCateSelectItemList(db_session)

	var package_material_list [] models.PackageMaterialInfo
	err = db_session.Asc(models.PackageMaterialCols.Name).Find(&package_material_list)
	core.HandleError(err)

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_PRODUCT)

	inventory_list := product_inventory_service.Get_product_inventory_list(db_session, product_id)
	order_item_list := order_service.Get_order_item_data_list_by_product_id(db_session, product_id)

	var detail_form ProductDetailForm
	detail_form.Product = product
	detail_form.ProductImageList = product_image_list
	detail_form.SupplierList = supplier_list
	detail_form.ProductCateSelectItemList = product_cate_list
	detail_form.PackageMaterialList = package_material_list
	detail_form.StatusList = status_list
	detail_form.InventoryList = inventory_list
	detail_form.OrderItemList = order_item_list

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

	product_from_ui := detail_form.Product

	var affected int64
	if detail_form.IsCreate() {
		sku_count, err := db_session.Where("sku = ?", product_from_ui.Sku).Count(&product_from_ui)
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！" })
		}

		name_count, err := db_session.Where("name = ?", product_from_ui.Name).Count(&product_from_ui)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&product_from_ui)
		core.HandleError(err)

		product_log_service.Insert_product_log(nil, &product_from_ui, product_log_service.Product_operation_type.Insert, db_session, c.Session)
	} else {
		sku_count, err := db_session.Where("id <> ? and sku = ?", product_from_ui.Id, product_from_ui.Sku).Count(&product_from_ui)
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！" })
		}

		name_count, err := db_session.Where("id <> ? and name = ?", product_from_ui.Id, product_from_ui.Name).Count(&product_from_ui)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		var product_in_db models.ProductInfo
		_, err = db_session.Id(product_from_ui.Id).Get(&product_in_db)
		core.HandleError(err)

		affected, err = db_session.Id(product_from_ui.Id).Update(&product_from_ui)
		core.HandleError(err)

		product_log_service.Insert_product_log(&product_from_ui, &product_in_db, product_log_service.Product_operation_type.Update, db_session, c.Session)

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

	for _, product_id := range product_id_list {
		var product models.ProductInfo
		has, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)

		if has {
			product_log_service.Insert_product_log(nil, &product, product_log_service.Product_operation_type.Delete, db_session, c.Session)
		}
	}

	affected, err := db_session.In("id", product_id_list).Delete(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpProduct) GetProductInfo() revel.Result {
	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	var product models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"product": product,
	})
}

func (c ErpProduct) ProductOrderList(product_id int64) revel.Result {
	return c.RenderEmpty()
}

func (c ErpProduct) ProductOrderListData(product_id int64) revel.Result {
	return c.RenderEmpty()
}

func (c ErpProduct) ProductMessageList(product_id int64) revel.Result {
	return c.RenderEmpty()
}

func (c ErpProduct) ProductMessageListData(product_id int64) revel.Result {
	return c.RenderEmpty()
}

func (c ErpProduct) ProductLogList(product_id int64) revel.Result {
	return c.RenderEmpty()
}

func (c ErpProduct) ProductLogListData(product_id int64) revel.Result {
	return c.RenderEmpty()
}
