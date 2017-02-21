package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	auth_models "matrix/modules/auth/models"
	"fmt"
)

type ErpSalesManager struct {
	*revel.Controller
}

func (c ErpSalesManager) Index() revel.Result {
	return c.RenderTemplate("erp/xsgl/sales_manager/sales_manager_index.html")
}

func (c ErpSalesManager) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpSalesManagerListViewItem struct {
		models.SalesManagerInfo `xorm:"extends" json:"sm"`
		models.ShopAccountInfo  `xorm:"extends" json:"sa"`
		models.ProductInfo      `xorm:"extends" json:"p"`
		auth_models.UserInfo    `xorm:"extends" json:"u"`
	}

	/*
	SELECT sm.*, sa.*, p.*, u.*
	FROM erp_sales_manager sm
		INNER JOIN erp_shop_account sa ON sm.shop_id = sa.id
		INNER JOIN erp_product p ON sm.product_id = p.id
		INNER JOIN auth_user u ON sm.user_id = u.id
	*/

	sb := core.NewSqlBuilder()
	sb.Select("sm.*, sa.*, p.*, u.*")
	sb.Count("count(sm.id) as sales_manager_count")
	sb.From(`erp_sales_manager sm
		INNER JOIN erp_shop_account sa ON sm.shop_id = sa.id
		INNER JOIN erp_product p ON sm.product_id = p.id
		INNER JOIN auth_user u ON sm.user_id = u.id`)
	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sm.id ASC")
	}

	sb.Limit(offset, limit)

	sales_manager_list := make([]ErpSalesManagerListViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&sales_manager_list)
	core.HandleError(err)

	var sales_manager_count struct {
		Count int64 `xorm:"'sales_manager_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&sales_manager_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  sales_manager_list,
		Total: sales_manager_count.Count,
	})
}

type SalesManagerDetailForm struct {
	SalesManager    models.SalesManagerInfo `json:"sales_manager"`
	ShopAccountName string                  `json:"shop_account_name"`
	ProductName     string                  `json:"product_name"`
	UserName        string                  `json:"user_name"`
}

func (f *SalesManagerDetailForm) IsCreate() bool {
	return f.SalesManager.Id == 0
}

func (f *SalesManagerDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.SalesManager.ShopId).Message("店铺不能为空！")

	validation.Required(f.SalesManager.ProductId).Message("商品不能为空！")

	validation.Required(f.SalesManager.UserId).Message("销售不能为空！")

	return validation.HasErrors() == false
}

func (c ErpSalesManager) DetailData() revel.Result {
	db_session := c.DbSession

	var sales_manager_id int64
	c.Params.Bind(&sales_manager_id, "id")

	var shop_account models.ShopAccountInfo
	var product models.ProductInfo
	var user auth_models.UserInfo
	var sales_manager models.SalesManagerInfo
	if sales_manager_id != 0 {
		has, err := db_session.Id(sales_manager_id).Get(&sales_manager)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的销售负责人不存在！"})
		}

		_, err = db_session.Id(sales_manager.ShopId).Get(&shop_account)
		core.HandleError(err)
		_, err = db_session.Id(sales_manager.ProductId).Get(&product)
		core.HandleError(err)
		_, err = db_session.Id(sales_manager.UserId).Get(&user)
		core.HandleError(err)
	}

	var detail_form SalesManagerDetailForm
	detail_form.SalesManager = sales_manager
	detail_form.ShopAccountName = fmt.Sprintf("%s %s", shop_account.Platform, shop_account.Name)
	detail_form.ProductName = fmt.Sprintf("%s %s", product.Sku, product.Name)
	detail_form.UserName = fmt.Sprintf("%s %s", user.LoginName, user.NickName)

	return c.RenderJson(detail_form)
}

func (c ErpSalesManager) Save() revel.Result {
	db_session := c.DbSession

	var detail_form SalesManagerDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	sales_manager := detail_form.SalesManager

	var affected int64
	if detail_form.IsCreate() {
		shopId_count, err := db_session.Where("shop_id = ? and product_id = ? and user_id = ?", sales_manager.ShopId, sales_manager.ProductId, sales_manager.UserId).Count(&sales_manager)
		core.HandleError(err)
		if shopId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，店铺、商品、销售已存在！" })
		}

		affected, err = db_session.Insert(&sales_manager)
		core.HandleError(err)
	} else {
		shopId_count, err := db_session.Where("id <> ? and shop_id = ? and product_id = ? and user_id = ?", sales_manager.Id, sales_manager.ShopId, sales_manager.ProductId, sales_manager.UserId).Count(&sales_manager)
		core.HandleError(err)
		if shopId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，店铺、商品、销售已存在！" })
		}

		affected, err = db_session.Id(sales_manager.Id).Update(&sales_manager)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpSalesManager) Delete() revel.Result {
	db_session := c.DbSession

	sales_manager_id_list := make([]int64, 0)
	c.Params.Bind(&sales_manager_id_list, "id_list")

	affected, err := db_session.In("id", sales_manager_id_list).Delete(new(models.SalesManagerInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpSalesManager) ShopAccountSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpShopAccountViewItem struct {
		models.ShopAccountInfo `xorm:"extends" json:"a"`
		models.StatusInfo      `xorm:"extends" json:"s"`
	}

	/*
		SELECT a.*, s.* FROM erp_shop_account a
		    INNER JOIN erp_status s ON a.status_id = s.id
	*/

	query := db_session.
	Select("a.*, s.*").
		Table("erp_shop_account").Alias("a").
		Join("INNER", []string{"erp_status", "s"}, "a.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("a.id")
	}

	shop_account_list := make([]ErpShopAccountViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&shop_account_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpShopAccountViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  shop_account_list,
		Total: count,
	})
}

func (c ErpSalesManager) ProductSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpProductViewItem struct {
		models.ProductInfo          `xorm:"extends" json:"prd"`
		models.ProductCateInfo      `xorm:"extends" json:"cate"`
		models.SupplierInfo         `xorm:"extends" json:"sup"`
		models.PackageMaterialInfo  `xorm:"extends" json:"pm"`
		models.StatusInfo           `xorm:"extends" json:"st"`
	}

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

func (c ErpSalesManager) UserSelectListData() revel.Result {
	session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := session.Where(filter)

	//query extra filter here

	dataQuery := *query
	if order != "" {
		dataQuery = *dataQuery.OrderBy(order)
	} else {
		dataQuery = *dataQuery.Asc("id")
	}

	userList := make([]auth_models.UserInfo, 0, limit)
	err := dataQuery.Limit(limit, offset).Find(&userList)
	core.HandleError(err)

	countQuery := *query
	count, err := countQuery.Count(new(auth_models.UserInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  userList,
		Total: count,
	})
}
