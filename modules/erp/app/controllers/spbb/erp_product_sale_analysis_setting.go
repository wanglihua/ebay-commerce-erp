package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
)

type ErpProductSaleAnalysisSetting struct {
	*revel.Controller
}

func (c ErpProductSaleAnalysisSetting) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS_SETTING, c.Session)

	return c.RenderTemplate("erp/spbb/product_sale_analysis_setting/product_sale_analysis_setting_index.html")
}

type ProductSaleAnalysisSettingDetailForm struct {
	Setting models.ProductSaleAnalysisSettingInfo `json:"setting"`
}

func (f *ProductSaleAnalysisSettingDetailForm) Valid(validation *revel.Validation) bool {

	validation.Min(int(f.Setting.Level1Days), 0).Message("1档统计天数不能小于0！")
	validation.Min(int(f.Setting.Level2Days), 0).Message("2档统计天数不能小于0！")
	validation.Min(int(f.Setting.Level3Days), 0).Message("3档统计天数不能小于0！")
	validation.Min(int(f.Setting.Level4Days), 0).Message("4档统计天数不能小于0！")
	validation.Min(int(f.Setting.Level5Days), 0).Message("5档统计天数不能小于0！")

	validation.Max(int(f.Setting.Level1Days), 300).Message("1档统计天数不能大于300！")
	validation.Max(int(f.Setting.Level2Days), 300).Message("2档统计天数不能大于300！")
	validation.Max(int(f.Setting.Level3Days), 300).Message("3档统计天数不能大于300！")
	validation.Max(int(f.Setting.Level4Days), 300).Message("4档统计天数不能大于300！")
	validation.Max(int(f.Setting.Level5Days), 300).Message("5档统计天数不能大于300！")

	if f.Setting.Level1Weight < 0 || f.Setting.Level1Weight > 1 {
		validation.Error("1档统计权重只能在0和1之间！")
	}

	if f.Setting.Level2Weight < 0 || f.Setting.Level2Weight > 1 {
		validation.Error("2档统计权重只能在0和1之间！")
	}

	if f.Setting.Level3Weight < 0 || f.Setting.Level3Weight > 1 {
		validation.Error("3档统计权重只能在0和1之间！")
	}

	if f.Setting.Level4Weight < 0 || f.Setting.Level4Weight > 1 {
		validation.Error("4档统计权重只能在0和1之间！")
	}

	if f.Setting.Level5Weight < 0 || f.Setting.Level5Weight > 1 {
		validation.Error("5档统计权重只能在0和1之间！")
	}

	return validation.HasErrors() == false
}

func (c ErpProductSaleAnalysisSetting) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS_SETTING, c.Session)

	db_session := c.DbSession

	var setting = product_sale_analysis_setting_service.Get_product_sale_analysis_setting(db_session)

	var detail_form ProductSaleAnalysisSettingDetailForm
	detail_form.Setting = setting

	return c.RenderJson(detail_form)
}

func (c ErpProductSaleAnalysisSetting) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPBB_PRODUCT_SALE_ANALYSIS_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form ProductSaleAnalysisSettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	setting := &(detail_form.Setting) // 取地址，这样下面的对 logistics_setting 的修改都会反映到 detail_form 里面去的，返回的 json 里面也会有反映

	// 先插入再删除，这样从表里面取数据的时候，总能取到的

	// 插入当前的配置信息
	setting.Id = 0
	setting.Name = "default"
	_, err = db_session.Insert(setting) // 因为前面取的是地址，所以这里就不用取地址了的
	core.HandleError(err)

	// 删除原来表里面的信息记录
	sql_delete := "DELETE FROM erp_product_sale_analysis_setting WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, setting.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"detail":  detail_form,
		"message": "数据保存成功!",
	})
}
