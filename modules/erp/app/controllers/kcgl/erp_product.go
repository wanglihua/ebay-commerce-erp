package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/message_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/product_cate_service"
	"matrix/modules/erp/service/product_log_service"
	"matrix/modules/erp/service/product_service"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/cggl_models"
	"github.com/go-xorm/xorm"
	auth_models "matrix/modules/auth/models"
)

type ErpProduct struct {
	*revel.Controller
}

func (c ErpProduct) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession
	var product_cate_list = product_cate_service.GetProductCateSelectItemList(db_session)

	c.RenderArgs["can_save"] = pass_service.Is_login_user_has_pass(pass_res.KCGL_PRODUCT_DETAIL_SAVE, c.Session)

	c.RenderArgs["product_cate_list"] = product_cate_list

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	var product_status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.Where("type = ?", models.STATUS_TYPE_PRODUCT).OrderBy("id").Find(&product_status_list))

	c.RenderArgs["sale_user_list"] = user_list
	c.RenderArgs["service_user_list"] = user_list
	c.RenderArgs["dev_user_list"] = user_list

	c.RenderArgs["product_status_list"] = product_status_list

	return c.RenderTemplate("erp/kcgl/product/product_index.html")
}

type ErpProductViewItem struct {
	Product         kcgl_models.ProductInfo                           `xorm:"extends" json:"prd"`
	ProductCate     kcgl_models.ProductCateInfo                       `xorm:"extends" json:"cate"`
	Supplier        sppz_models.SupplierInfo                          `xorm:"extends" json:"sup"`
	PackageMaterial sppz_models.PackageMaterialInfo                   `xorm:"extends" json:"pm"`
	Status          models.StatusInfo                                 `xorm:"extends" json:"st"`
	InventoryList   []inventory_service.Product_inventory_view_item_t `xorm:"-" json:"inventory_list"`
}

func (c ErpProduct) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_cate_query int64
	c.Params.Bind(&product_cate_query, "product_cate_query")

	var product_sku_query string
	c.Params.Bind(&product_sku_query, "product_sku_query")

	var product_name_query string
	c.Params.Bind(&product_name_query, "product_name_query")

	var product_status_query int64
	c.Params.Bind(&product_status_query, "product_status_query")

	var sale_user_id_query int64
	c.Params.Bind(&sale_user_id_query, "sale_user_id_query")

	var service_user_id_query int64
	c.Params.Bind(&service_user_id_query, "service_user_id_query")

	var dev_user_id_query int64
	c.Params.Bind(&dev_user_id_query, "dev_user_id_query")

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

	//query extra filter here
	if product_cate_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.cate_id = %d", product_cate_query))
	}

	if product_status_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.status_id = %d", product_status_query))
	}

	if strings.TrimSpace(product_sku_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_sku_query))
	}

	if strings.TrimSpace(product_name_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.cn_name LIKE '%%%s%%'", product_name_query))
	}

	if sale_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sale_user_id = %d", sale_user_id_query))
	}

	if service_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.service_user_id = %d", service_user_id_query))
	}

	if dev_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.dev_user_id = %d", dev_user_id_query))
	}

	sql_builder.Where(filter)

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]ErpProductViewItem, 0, limit)
	err := db_session.Sql(sql_builder.SqlList()).Find(&product_list)
	core.HandleError(err)

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))

	var get_user_desc_in_user_list = func(user_id int64) string {
		for _, user := range user_list {
			if user.Id == user_id {
				return fmt.Sprintf("%s（%s）", user.LoginName, user.NickName)
			}
		}

		return ""
	}

	for index, product := range product_list {
		product_list[index].Product.SaleUser = get_user_desc_in_user_list(product.Product.SaleUserId)
		product_list[index].Product.ServiceUser = get_user_desc_in_user_list(product.Product.ServiceUserId)
		product_list[index].Product.DevUser = get_user_desc_in_user_list(product.Product.DevUserId)
	}

	var product_id_list = make([]int64, 0)
	for _, product := range product_list {
		product_id_list = append(product_id_list, product.Product.Id)
	}

	inventory_list_product_list := inventory_service.Get_inventory_by_product_list(db_session, product_id_list)

	for index, _ := range product_list {
		inventory_list_product := make([]inventory_service.Product_inventory_view_item_t, 0)
		for _, inventory := range inventory_list_product_list {
			if inventory.ProductId == product_list[index].Product.Id {
				inventory_list_product = append(inventory_list_product, inventory)
			}
		}

		product_list[index].InventoryList = inventory_list_product
	}

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
	Product                   kcgl_models.ProductInfo                           `json:"product"`
	ProductImageList          []kcgl_models.ProductImageInfo                    `json:"product_image_list"`
	SupplierList              []sppz_models.SupplierInfo                        `json:"supplier_list"`
	ProductCateSelectItemList []product_cate_service.ProductCateSelectItem      `json:"product_cate_list"`
	PackageMaterialList       []sppz_models.PackageMaterialInfo                 `json:"package_material_list"`
	StatusList                []models.StatusInfo                               `json:"status_list"`
	ProductBrandList          []sppz_models.BrandInfo                           `json:"product_brand_list"`
	ProductStyleList          []sppz_models.ProductStyleInfo                    `json:"product_style_list"`
	ProductCertList           []sppz_models.ProductCertInfo                     `json:"product_cert_list"`
	InventoryList             []inventory_service.Product_inventory_view_item_t `json:"inventory_list"`
	OrderItemList             []order_service.Order_item_view_item_t            `json:"order_item_list"`
	ProductLogList            []kcgl_models.ProductLogInfo                      `json:"product_log_list"`
	ProductPurchaseList       []ProductPurchaseListItem                         `json:"product_purchase_list"`
	MessageList               []message_service.Message_view_item_t             `json:"message_list"`
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

	//validation.Required(f.Product.Name).Message("名称不能为空！")
	//if f.Product.Name != "" {
	//	validation.MinSize(f.Product.Name, 2).Message("名称长度不能小于2！")
	//}

	if f.Product.Name != "" {
		validation.MaxSize(f.Product.Name, 100).Message("名称长度不能大于100！")
	}

	validation.Required(f.Product.CostPrice).Message("成本价不能为空！")

	validation.Required(f.Product.Weight).Message("重量不能为空！")

	validation.Required(f.Product.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProduct) Detail(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	c.RenderArgs["can_save"] = pass_service.Is_login_user_has_pass(pass_res.KCGL_PRODUCT_DETAIL_SAVE, c.Session)

	//return c.RenderText(fmt.Sprintf("product_id %d", product_id))
	c.RenderArgs["product_id"] = product_id

	c.RenderArgs["is_create"] = product_id == 0

	return c.RenderTemplate("erp/kcgl/product/product_detail.html")
}

func (c ErpProduct) Edit(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_DETAIL_SAVE, c.Session)

	c.RenderArgs["can_save"] = pass_service.Is_login_user_has_pass(pass_res.KCGL_PRODUCT_DETAIL_SAVE, c.Session)

	c.RenderArgs["product_id"] = product_id

	c.RenderArgs["is_create"] = product_id == 0

	return c.RenderTemplate("erp/kcgl/product/product_edit.html")
}

func (c ErpProduct) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	if product_id != 0 { // 如果是不为零的时候，才会有这样的判断
		has, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品不存在！"})
		}
	}

	product.Volume = fmt.Sprintf("%0.6f", product.Length*product.Weight*product.Height/1000000)

	if product.SaleUserId != 0 {
		var sale_user auth_models.UserInfo
		_, err := db_session.Id(product.SaleUserId).Get(&sale_user)
		core.HandleError(err)

		product.SaleUser = fmt.Sprintf("%s(%s)", sale_user.LoginName, sale_user.NickName)
	}

	if product.ServiceUserId != 0 {
		var service_user auth_models.UserInfo
		_, err := db_session.Id(product.ServiceUserId).Get(&service_user)
		core.HandleError(err)

		product.ServiceUser = fmt.Sprintf("%s(%s)", service_user.LoginName, service_user.NickName)
	}

	if product.DevUserId != 0 {
		var dev_user auth_models.UserInfo
		_, err := db_session.Id(product.DevUserId).Get(&dev_user)
		core.HandleError(err)

		product.DevUser = fmt.Sprintf("%s(%s)", dev_user.LoginName, dev_user.NickName)
	}

	if product.Id == 0 {
		product.StatusId = models.PRODUCT_STATUS_ZC
	}

	var product_image_list = make([]kcgl_models.ProductImageInfo, 0)
	if product_id != 0 {
		err := db_session.Where(fmt.Sprintf("product_id = %d", product_id)).Find(&product_image_list)
		core.HandleError(err)
	}

	var supplier_list []sppz_models.SupplierInfo
	err := db_session.Asc("name").Find(&supplier_list)
	core.HandleError(err)

	var product_cate_list = product_cate_service.GetProductCateSelectItemList(db_session)

	var package_material_list [] sppz_models.PackageMaterialInfo
	err = db_session.Asc("name").Find(&package_material_list)
	core.HandleError(err)

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_PRODUCT)

	inventory_list := inventory_service.Get_product_inventory_list(db_session, product_id)
	order_item_list := order_service.Get_order_item_data_list_by_product_id_top_n(db_session, product_id, 10)

	var product_log_list = make([]kcgl_models.ProductLogInfo, 0)
	if product_id != 0 {
		err = db_session.Where(fmt.Sprintf("product_id = %d", product_id)).OrderBy("id DESC").Limit(10, 0).Find(&product_log_list)
		core.HandleError(err)
	}

	message_list := message_service.Get_message_list_data_by_product_id_top_n(db_session, product_id, 10)

	var product_brand_list = make([]sppz_models.BrandInfo, 0)
	core.HandleError(db_session.Find(&product_brand_list))

	var product_style_list = make([]sppz_models.ProductStyleInfo, 0)
	core.HandleError(db_session.Find(&product_style_list))

	var product_cert_list = make([]sppz_models.ProductCertInfo, 0)
	core.HandleError(db_session.Find(&product_cert_list))

	var product_purchase_list = c.GetProductPurchaseListData(product_id, 10, db_session)

	var detail_form ProductDetailForm
	detail_form.Product = product
	detail_form.ProductImageList = product_image_list
	detail_form.SupplierList = supplier_list
	detail_form.ProductCateSelectItemList = product_cate_list
	detail_form.PackageMaterialList = package_material_list
	detail_form.StatusList = status_list
	detail_form.InventoryList = inventory_list
	detail_form.OrderItemList = order_item_list
	detail_form.ProductLogList = product_log_list
	detail_form.MessageList = message_list
	detail_form.ProductBrandList = product_brand_list
	detail_form.ProductStyleList = product_style_list
	detail_form.ProductCertList = product_cert_list
	detail_form.ProductPurchaseList = product_purchase_list

	return c.RenderJson(detail_form)
}

func (c ErpProduct) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_DETAIL_SAVE, c.Session)

	db_session := c.DbSession

	var detail_form ProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &(detail_form.Product))
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_from_ui := detail_form.Product

	var affected int64
	is_create := detail_form.IsCreate()
	var product_id int64 = product_from_ui.Id
	if is_create {
		sku_count, err := db_session.Where("sku = ?", product_from_ui.Sku).Count(new(kcgl_models.ProductInfo))
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！"})
		}

		/*
		name_count, err := db_session.Where("name = ?", product_from_ui.Name).Count(new(kcgl_models.ProductInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}
		*/

		affected, err = db_session.Insert(&product_from_ui)
		core.HandleError(err)

		product_id = product_from_ui.Id // 这是数据库保存后得到得到新生成的 id

		product_log_service.Insert_product_log(&product_from_ui, nil, core.Db_operation_types.Insert, db_session, c.Session)
	} else {
		sku_count, err := db_session.Where("id <> ? and sku = ?", product_from_ui.Id, product_from_ui.Sku).Count(new(kcgl_models.ProductInfo))
		core.HandleError(err)
		if sku_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，SKU已存在！"})
		}

		/*
		if strings.TrimSpace(product_from_ui.Name) != "" {
			name_count, err := db_session.Where("id <> ? and name = ?", product_from_ui.Id, product_from_ui.Name).Count(new(kcgl_models.ProductInfo))
			core.HandleError(err)
			if name_count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
			}
		}
		*/

		var product_in_db kcgl_models.ProductInfo
		_, err = db_session.Id(product_from_ui.Id).Get(&product_in_db)
		core.HandleError(err)

		//update_cols := []string{"sku", "cn_name", "name", "cate_id", "supplier_id", "pkg_mtrl_id", "cost_price", "status_id", "desc", "weight", "need_invoice" }
		update_cols := []string{"sku", "cn_name", "name", "cate_id", "supplier_id", "pkg_mtrl_id", "cost_price", "status_id", "desc", "length", "width", "height", "weight",
			"brand_id", "style_id", "cert1_id", "cert2_id", "cert3_id", "stock_alert_days", "purchase_days", "tax_point", "drawback_point"}
		affected, err = db_session.Id(product_from_ui.Id).Cols(update_cols...).Update(&product_from_ui)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}

		product_id = product_in_db.Id //就是原来的 id, 这一句赋值其实是可以省掉的

		product_log_service.Insert_product_log(&product_from_ui, &product_in_db, core.Db_operation_types.Update, db_session, c.Session)
	}

	// 重新获得新增或者更新的product,用于返回
	var product kcgl_models.ProductInfo
	_, err = db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success":   true,
		"is_create": is_create,
		"product":   product,
		"message":   "数据保存成功",
	})
}

func (c ErpProduct) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	product_id_list := make([]int64, 0)
	c.Params.Bind(&product_id_list, "id_list")

	product_service.Delete_product_list(product_id_list, db_session, c.Session)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}

func (c ErpProduct) ProductOrderList(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("erp/kcgl/product/product_order_list.html")
}

func (c ErpProduct) ProductOrderListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	order_item_list := order_service.Get_order_item_data_list_by_product_id(db_session, product_id)
	return c.RenderJson(map[string]interface{}{
		"product":         product,
		"order_item_list": order_item_list,
	})
}

func (c ErpProduct) ProductMessageList(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("erp/kcgl/product/product_message_list.html")
}

func (c ErpProduct) ProductMessageListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	message_list := message_service.Get_message_list_data_by_product_id(db_session, product_id)
	return c.RenderJson(map[string]interface{}{
		"product":      product,
		"message_list": message_list,
	})
}

func (c ErpProduct) ProductLogList(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("erp/kcgl/product/product_log_list.html")
}

func (c ErpProduct) ProductLogListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var product_log_list = make([]kcgl_models.ProductLogInfo, 0)
	err = db_session.Where(fmt.Sprintf("product_id = %d", product_id)).OrderBy("id DESC").Find(&product_log_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"product":          product,
		"product_log_list": product_log_list,
	})
}

func (c ErpProduct) ProductPurchaseList(product_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("erp/kcgl/product/product_purchase_list.html")
}

func (c ErpProduct) ProductPurchaseListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var product_purchase_list = c.GetProductPurchaseListData(product_id, 0, db_session)

	return c.RenderJson(map[string]interface{}{
		"product":               product,
		"product_purchase_list": product_purchase_list,
	})
}

type ProductPurchaseListItem struct {
	cggl_models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	models.StatusInfo                 `xorm:"extends" json:"s"`
}

func (c ErpProduct) GetProductPurchaseListData(product_id int64, limit int64, db_session *xorm.Session) []ProductPurchaseListItem {

	/*
SELECT pbe.*, pb.*, p.*, s.*
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, pb.*, p.*, s.*")
	sb.From(`erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id`)

	sb.Where(fmt.Sprintf("pbe.product_id = %d", product_id))

	sb.OrderBy("pbe.id DESC")

	if limit > 0 {
		sb.Limit(0, limit)
	}

	product_purchase_list := make([]ProductPurchaseListItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_purchase_list)
	core.HandleError(err)

	return product_purchase_list
}

func (c ErpProduct) GetProductInfo() revel.Result {
	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"product": product,
	})
}

func (c ErpProduct) GetProductListSkus() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT, c.Session)

	db_session := c.DbSession

	var product_id_list = make([]int64, 0)
	c.Params.Bind(&product_id_list, "product_id_list")

	if len(product_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"product_skus": "",
		})
	}

	var product_id_str_list = make([]string, 0)
	for _, product_id := range product_id_list {
		product_id_str_list = append(product_id_str_list, fmt.Sprint(product_id))
	}

	type product_sku_t struct {
		Sku string `xorm:"'product_sku'"`
	}

	var product_sku_struct_list = make([]product_sku_t, 0)
	err := db_session.Sql(fmt.Sprintf("SELECT p.sku AS product_sku FROM erp_product p WHERE p.id IN (%s)", strings.Join(product_id_str_list, ","))).Find(&product_sku_struct_list)
	core.HandleError(err)

	var product_sku_list = make([]string, 0)
	for _, product_sku_struct := range product_sku_struct_list {
		product_sku_list = append(product_sku_list, product_sku_struct.Sku)
	}

	return c.RenderJson(map[string]interface{}{
		"product_skus": strings.Join(product_sku_list, ","),
	})
}
