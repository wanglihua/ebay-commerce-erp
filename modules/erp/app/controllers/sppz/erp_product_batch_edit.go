package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type ErpProductBatchEdit struct {
	*revel.Controller
}

const (
	EDIT_ATTR_NAME       = 1001
	EDIT_ATTR_COST_PRICE = 1002
	EDIT_ATTR_WEIGHT     = 1003
)

func (c ErpProductBatchEdit) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_BATCH, c.Session)

	return c.RenderTemplate("erp/sppz/product_batch_edit/product_batch_edit_index.html")
}

func (c ErpProductBatchEdit) EditAttrList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_BATCH, c.Session)

	type edit_attr struct {
		AttrId   int64  `json:"attr_id"`
		AttrName string `json:"attr_name"`
	}

	edit_attr_list := make([]edit_attr, 0)
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId: EDIT_ATTR_NAME, AttrName: "名称"})
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId: EDIT_ATTR_COST_PRICE, AttrName: "成本价"})
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId: EDIT_ATTR_WEIGHT, AttrName: "重量"})

	return c.RenderJson(edit_attr_list)
}

type BatchEditProductDetailForm struct {
	EditAttr       int64    `json:"edit_attr"`
	EditAttrValues string   `json:"edit_attr_values"`
}

func (f *BatchEditProductDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.EditAttr).Message("修改的元素不能为空！")
	if f.EditAttr == 0 {
		validation.Error("请选择修改的元素！")
	}
	validation.Required(f.EditAttr).Message("修改的元素不能为空！")
	validation.Required(f.EditAttrValues).Message("修改的值不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductBatchEdit) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_BATCH, c.Session)

	db_session := c.DbSession

	var detail_form BatchEditProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	attr_value_lines := strings.Split(detail_form.EditAttrValues, "\n")

	type sku_attr_value_type struct {
		sku        string
		attr_value string
	}

	sku_attr_value_list := make([]sku_attr_value_type, 0)
	for _, attr_value_line := range attr_value_lines {
		if strings.TrimSpace(attr_value_line) != "" {
			attr_value_line_splits := strings.Split(attr_value_line, ",")
			if len(attr_value_line_splits) != 2 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "请将SKU和修改的元素值用逗号分开！"})
			}
			var sku_attr_value sku_attr_value_type
			sku_attr_value.sku = strings.TrimSpace(attr_value_line_splits[0])
			sku_attr_value.attr_value = strings.TrimSpace(attr_value_line_splits[1])
			sku_attr_value_list = append(sku_attr_value_list, sku_attr_value)
		}
	}

	// 检查SKU是否存在

	type product_count struct {
		ProductCount int64 `json:"product_count"`
	}

	for _, sku_attr_value := range sku_attr_value_list {
		sku := sku_attr_value.sku

		sql_get_sku_product_count := `
SELECT count(*) AS product_count FROM erp_product WHERE sku = '%s'
`
		sql_get_sku_product_count = fmt.Sprintf(sql_get_sku_product_count, sku)

		var product_count product_count
		_, err := db_session.Sql(sql_get_sku_product_count).Get(&product_count)
		core.HandleError(err)

		if product_count.ProductCount == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: fmt.Sprintf("SKU: %s 不存在！", sku)})
		}
	}

	if detail_form.EditAttr == EDIT_ATTR_NAME {
		for _, sku_attr_value := range sku_attr_value_list {
			sku := sku_attr_value.sku
			attr_value := sku_attr_value.attr_value

			sql_update_name := `
UPDATE erp_product SET name = '%s' WHERE sku = '%s';
`
			sql_update_name = fmt.Sprintf(sql_update_name, attr_value, sku)

			_, err := db_session.Exec(sql_update_name)
			core.HandleError(err)
		}
	}

	if detail_form.EditAttr == EDIT_ATTR_COST_PRICE {
		for _, sku_attr_value := range sku_attr_value_list {
			sku := sku_attr_value.sku
			attr_value := sku_attr_value.attr_value

			sql_update_cost_price := `
UPDATE erp_product SET cost_price = %s WHERE sku = '%s';
`
			sql_update_cost_price = fmt.Sprintf(sql_update_cost_price, attr_value, sku)

			_, err := db_session.Exec(sql_update_cost_price)
			core.HandleError(err)
		}
	}

	if detail_form.EditAttr == EDIT_ATTR_WEIGHT {
		for _, sku_attr_value := range sku_attr_value_list {
			sku := sku_attr_value.sku
			attr_value := sku_attr_value.attr_value

			sql_update_cost_price := `
UPDATE erp_product SET weight = %s WHERE sku = '%s';
`
			sql_update_cost_price = fmt.Sprintf(sql_update_cost_price, attr_value, sku)

			_, err := db_session.Exec(sql_update_cost_price)
			core.HandleError(err)
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "批量更新成功！"})
}
