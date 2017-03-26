package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/shop_service"
)

type ErpProductImage struct {
	*revel.Controller
}

func (c ErpProductImage) Index() revel.Result {
	return c.RenderTemplate("erp/xsgl/product_image/product_image_index.html")
}

func (c ErpProductImage) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_image_view_item_t struct {
		Id         string    `xorm:"'id'" json:"id"`
		Sku        string    `xorm:"'sku'" json:"sku"`
		ImgUrl     string    `xorm:"'img_url'" json:"img_url"`
		TypeName   string    `xorm:"'type_name'" json:"type_name"`
		ShopName   string    `xorm:"'shop_name'" json:"shop_name"`
		CreateTime core.Time `xorm:"'create_time'" json:"create_time"`
		UpdateTime core.Time `xorm:"'update_time'" json:"update_time"`
	}

	/*
		SELECT *
		FROM (
		         SELECT img.id AS id, img.product_code AS sku, img.url AS img_url, type.name AS type_name, concat(platform.name, ' ', sa.name) AS shop_name, img.create_time, img.update_time
		         FROM erp_product_image img
		             INNER JOIN erp_product_image_type type ON img.type_id = type.id
		             INNER JOIN erp_shop_account sa ON img.shop_id = sa.id
		             INNER JOIN erp_ecommerce_platform platform ON sa.platform_id = platform.id
		     ) AS product_image_view
	*/

	sb := core.NewSqlBuilder()
	sb.Select("*")
	sb.Count("count(id) AS img_count")
	sb.From(`
	     (
	         SELECT img.id AS id, img.product_code AS sku, img.url AS img_url, type.name AS type_name, concat(platform.name, ' ', sa.name) AS shop_name, img.create_time, img.update_time
	         FROM erp_product_image img
	             INNER JOIN erp_product_image_type type ON img.type_id = type.id
	             INNER JOIN erp_shop_account sa ON img.shop_id = sa.id
	             INNER JOIN erp_ecommerce_platform platform ON sa.platform_id = platform.id
	     ) AS product_image_view`)
	sb.Where(filter)

	//query extra filter here

	if order == "" {
		sb.OrderBy("id ASC")
	} else {
		sb.OrderBy(order)
	}

	sb.Limit(offset, limit)

	product_image_list := make([]product_image_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_image_list)
	core.HandleError(err)

	var img_count struct {
		Count int64 `xorm:"'img_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&img_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_image_list,
		Total: img_count.Count,
	})
}

type shop_select_item_t struct {
	Id   int64  `xorm:"'shop_id'" json:"shop_id"`
	Name string `xorm:"'shop_name'" json:"shop_name"`
}

type ProductImageDetailForm struct {
	ProductImage         models.ProductImageInfo            `json:"product_image"`
	ProductImageTypeList []models.ProductImageTypeInfo      `json:"product_image_type_list"`
	ShopList             []shop_service.Shop_select_item_t  `json:"shop_list"`
}

func (f *ProductImageDetailForm) IsCreate() bool {
	return f.ProductImage.Id == 0
}

func (f *ProductImageDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductImage.ProductCode).Message("商品编码不能为空！")
	if f.ProductImage.ProductCode != "" {
		validation.MinSize(f.ProductImage.ProductCode, 2).Message("商品编码长度不能小于2！")
	}
	if f.ProductImage.ProductCode != "" {
		validation.MaxSize(f.ProductImage.ProductCode, 100).Message("商品编码长度不能大于100！")
	}

	validation.Required(f.ProductImage.TypeId).Message("图片属性不能为空！")

	validation.Required(f.ProductImage.ShopId).Message("店铺不能为空！")

	validation.Required(f.ProductImage.Url).Message("图片链接不能为空！")
	if f.ProductImage.Url != "" {
		validation.MinSize(f.ProductImage.Url, 2).Message("图片链接长度不能小于2！")
	}
	if f.ProductImage.Url != "" {
		validation.MaxSize(f.ProductImage.Url, 300).Message("图片链接长度不能大于300！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductImage) DetailData() revel.Result {
	db_session := c.DbSession

	var product_image_id int64
	c.Params.Bind(&product_image_id, "id")

	var product_image models.ProductImageInfo
	if product_image_id != 0 {
		has, err := db_session.Id(product_image_id).Get(&product_image)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品图片不存在！"})
		}
	}

	var product_image_type_list []models.ProductImageTypeInfo
	err := db_session.Find(&product_image_type_list)
	core.HandleError(err)

	var shop_list = shop_service.Get_shop_select_item_list(db_session)

	var detail_form ProductImageDetailForm
	detail_form.ProductImage = product_image
	detail_form.ProductImageTypeList = product_image_type_list
	detail_form.ShopList = shop_list

	return c.RenderJson(detail_form)
}

func (c ErpProductImage) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductImageDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_image := detail_form.ProductImage

	var product models.ProductInfo
	product_cols := models.ProductCols
	found, err := db_session.Where(fmt.Sprintf("%s = ?", product_cols.Sku), product_image.ProductCode).Get(&product)
	core.HandleError(err)

	if found == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品SKU编号不存在，请重试！"})
	}

	product_image.ProductId = product.Id

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_image)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(product_image.Id).Update(&product_image)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductImage) Delete() revel.Result {
	db_session := c.DbSession

	product_image_id_list := make([]int64, 0)
	c.Params.Bind(&product_image_id_list, "id_list")

	affected, err := db_session.In("id", product_image_id_list).Delete(new(models.ProductImageInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
