package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/product_sale_setting_service"
	"matrix/modules/erp/models/xsgl_models"
)

type ErpProductSaleSetting struct {
	*revel.Controller
}

func (c ErpProductSaleSetting) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	return c.RenderTemplate("erp/xsgl/product_sale_setting/product_sale_setting_index.html")
}

type ProductSaleSettingDetailForm struct {
	ProductSaleSetting xsgl_models.ProductSaleSettingInfo `json:"product_sale_setting"`
}

func (f *ProductSaleSettingDetailForm) Valid(validation *revel.Validation) bool {
	return validation.HasErrors() == false
}

func (c ErpProductSaleSetting) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	db_session := c.DbSession

	var product_sale_setting = product_sale_setting_service.Get_product_sale_setting(db_session)

	var detail_form ProductSaleSettingDetailForm
	detail_form.ProductSaleSetting = product_sale_setting

	return c.RenderJson(detail_form)
}

func (c ErpProductSaleSetting) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form ProductSaleSettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_sale_setting := &(detail_form.ProductSaleSetting) // 取地址，这样下面的对 logistics_setting 的修改都会反映到 detail_form 里面去的，返回的 json 里面也会有反映

	// 先插入再删除，这样从表里面取数据的时候，总能取到的

	// 插入当前的配置信息
	product_sale_setting.Id = 0
	_, err = db_session.Insert(product_sale_setting) // 因为前面取的是地址，所以这里就不用取地址了的
	core.HandleError(err)

	// 删除原来表里面的信息记录
	sql_delete := "DELETE FROM erp_product_sale_setting WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, product_sale_setting.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"detail":  detail_form, // detail_form 里面的数据有过优化和更改的（对 AddressSpecialWord 的优化），所以就再返回一下
		"message": "数据保存成功!",
	})
}
