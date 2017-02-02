package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"encoding/json"
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
	return c.RenderTemplate("erp/sppz/product_batch_edit/product_batch_edit_index.html")
}

func (c ErpProductBatchEdit) EditAttrList() revel.Result {
	type edit_attr struct {
		AttrId   int64  `json:"attr_id"`
		AttrName string `json:"attr_name"`
	}

	edit_attr_list := make([]edit_attr, 0)
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId:EDIT_ATTR_NAME, AttrName:"名称"})
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId:EDIT_ATTR_COST_PRICE, AttrName:"成本价"})
	edit_attr_list = append(edit_attr_list, edit_attr{AttrId:EDIT_ATTR_WEIGHT, AttrName:"重量"})

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
	//db_session := c.DbSession

	var detail_form BatchEditProductDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	revel.TRACE.Println("EditAttr: ", detail_form.EditAttr)
	revel.TRACE.Println("EditAttrValues: ", detail_form.EditAttrValues)

	return c.RenderJson(core.JsonResult{Success: false, Message: "功能未实现！"})
}
