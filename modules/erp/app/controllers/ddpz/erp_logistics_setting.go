
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsSetting struct {
	*revel.Controller
}

func (c ErpLogisticsSetting) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_setting/logistics_setting_index.html")
}

type LogisticsSettingDetailForm struct {
	LogisticsSetting models.LogisticsSettingInfo `json:"logistics_setting"`
}

func (f *LogisticsSettingDetailForm) IsCreate() bool {
	return f.LogisticsSetting.Id == 0
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
	db_session := c.DbSession

	var logistics_setting_id int64
	c.Params.Bind(&logistics_setting_id, "id")

	var logistics_setting models.LogisticsSettingInfo
	if logistics_setting_id != 0 {
		has, err := db_session.Id(logistics_setting_id).Get(&logistics_setting)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流配置不存在！"})
		}
	}

	var detail_form LogisticsSettingDetailForm
	detail_form.LogisticsSetting = logistics_setting

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsSetting) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsSettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_setting := detail_form.LogisticsSetting

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&logistics_setting)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(logistics_setting.Id).Update(&logistics_setting)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}
