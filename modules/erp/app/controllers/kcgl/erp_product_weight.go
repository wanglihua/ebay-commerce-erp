package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductWeight struct {
	*revel.Controller
}

func (c ErpProductWeight) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_weight/product_weight_index.html")
}

func (c ErpProductWeight) ListData() revel.Result {
	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_weight_view_item_t struct {
		Product models.ProductInfo `xorm:"extends" json:"p"`
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

type product_weight_detail_form_t struct {
	Product models.ProductInfo `json:"product"`
}

func (c ErpProductWeight) DetailData() revel.Result {
	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_weight_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductWeight) Save() revel.Result {
	db_session := c.DbSession

	var detail_form product_weight_detail_form_t
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var product_id = detail_form.Product.Id

	var product_in_db models.ProductInfo
	has, err := db_session.Id(product_id).Get(&product_in_db)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品不存在！",
		})
	}

	product_in_db.Length = detail_form.Product.Length
	product_in_db.Width = detail_form.Product.Width
	product_in_db.Height = detail_form.Product.Height
	product_in_db.Weight = detail_form.Product.Weight

	updata_col_List := []string{"length", "width", "height", "weight"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品长宽高重量数据未能更新！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductWeight) Import() revel.Result {

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
