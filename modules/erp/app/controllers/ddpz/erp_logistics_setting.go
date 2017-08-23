package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/logistics_setting_service"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpLogisticsSetting struct {
	*revel.Controller
}

func (c ErpLogisticsSetting) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	return c.RenderTemplate("erp/ddpz/logistics_setting/logistics_setting_index.html")
}

type LogisticsSettingDetailForm struct {
	LogisticsSetting ddpz_models.LogisticsSettingInfo `json:"logistics_setting"`
}

func (f *LogisticsSettingDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsSetting.AddressSpecialWord).Message("地址特殊字段不能为空！")
	if f.LogisticsSetting.AddressSpecialWord != "" {
		validation.MaxSize(f.LogisticsSetting.AddressSpecialWord, 300).Message("地址特殊字段长度不能大于300！")
	}

	validation.Required(f.LogisticsSetting.MaxDeficitPrice).Message("最大亏损价格不能为空！")
	validation.Required(f.LogisticsSetting.MaxDeficitRate).Message("最大亏损率不能为空！")
	validation.Required(f.LogisticsSetting.EnvelopeMaxLength).Message("信封最大长度不能为空！")
	validation.Required(f.LogisticsSetting.EnvelopeMaxWidth).Message("信封最大宽度不能为空！")
	validation.Required(f.LogisticsSetting.EnvelopeMaxHeight).Message("信封最大高度不能为空！")
	validation.Required(f.LogisticsSetting.EnvelopeMaxWeight).Message("信封最大重量不能为空！")
	validation.Required(f.LogisticsSetting.MinExpressFee).Message("最少快递费不能为空！")

	return validation.HasErrors() == false
}

func (c ErpLogisticsSetting) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	db_session := c.DbSession

	var logistics_setting = logistics_setting_service.Get_logistics_setting(db_session)

	var detail_form LogisticsSettingDetailForm
	detail_form.LogisticsSetting = logistics_setting

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsSetting) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form LogisticsSettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_setting := &(detail_form.LogisticsSetting) // 取地址，这样下面的对 logistics_setting 的修改都会反映到 detail_form 里面去的，返回的 json 里面也会有反映

	// 检查地址特殊字段的格式
	if strings.Index(logistics_setting.AddressSpecialWord, "，") != -1 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "地址特殊字段不能用中文的逗号分隔！",
		})
	}

	address_special_word_list := strings.Split(logistics_setting.AddressSpecialWord, ",")
	for index, address_special_word := range address_special_word_list {
		if strings.TrimSpace(address_special_word) == "" {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "地址特殊字段请用单个的英文逗号分隔！",
			})
		}

		address_special_word_list[index] = strings.TrimSpace(address_special_word)
	}

	logistics_setting.AddressSpecialWord = strings.Join(address_special_word_list, ",") //这里面保存的已经是去掉 字段 前后 多余的空白了的

	// 先插入再删除，这样从表里面取数据的时候，总能取到的

	// 插入当前的配置信息
	logistics_setting.Id = 0
	_, err = db_session.Insert(logistics_setting) // 因为前面取的是地址，所以这里就不用取地址了的
	core.HandleError(err)

	// 删除原来表里面的信息记录
	sql_delete := "DELETE FROM erp_logistics_setting WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, logistics_setting.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"detail":  detail_form, // detail_form 里面的数据有过优化和更改的（对 AddressSpecialWord 的优化），所以就再返回一下
		"message": "数据保存成功!",
	})
}
