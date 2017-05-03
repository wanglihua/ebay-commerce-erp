package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/models"
)

type AuthSetting struct {
	*revel.Controller
}

func (c AuthSetting) Index() revel.Result {
	return c.RenderTemplate("auth/setting/setting_index.html")
}

type SettingDetailForm struct {
	Setting models.SettingInfo `json:"auth_setting"`
}

func (f *SettingDetailForm) IsCreate() bool {
	return f.Setting.Id == 0
}

func (f *SettingDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Setting.DefaultPassword).Message("默认密码不能为空！")
	if f.Setting.DefaultPassword != "" {
		validation.MinSize(f.Setting.DefaultPassword, 6).Message("默认密码长度不能小于6！")
	}
	if f.Setting.DefaultPassword != "" {
		validation.MaxSize(f.Setting.DefaultPassword, 20).Message("默认密码长度不能大于20！")
	}

	validation.Required(f.Setting.TimeOut).Message("超时时长不能为空！")
	validation.Min(int(f.Setting.TimeOut), 1).Message("超时时长不能小于1！")
	validation.Max(int(f.Setting.TimeOut), 10000).Message("超时时长不能大于10000！")

	return validation.HasErrors() == false
}

func (c AuthSetting) DetailData() revel.Result {
	db_session := c.DbSession

	var setting_id int64
	c.Params.Bind(&setting_id, "id")

	var setting models.SettingInfo
	if setting_id != 0 {
		has, err := db_session.Id(setting_id).Get(&setting)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的授权设置不存在！"})
		}
	}

	var detail_form SettingDetailForm
	detail_form.Setting = setting

	return c.RenderJson(detail_form)
}

func (c AuthSetting) Save() revel.Result {
	db_session := c.DbSession

	var detail_form SettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	setting := detail_form.Setting

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&setting)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(setting.Id).Update(&setting)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}
