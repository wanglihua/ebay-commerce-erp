package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/models"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/auth_setting_service"
	"matrix/modules/auth/service/pass_service"
)

type AuthSetting struct {
	*revel.Controller
}

func (c AuthSetting) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_SETTING, c.Session)

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
	pass_service.Check_login_user_pass(pass_res.AUTH_SETTING, c.Session)

	db_session := c.DbSession

	var setting = auth_setting_service.Get_auth_setting(db_session)

	var detail_form SettingDetailForm
	detail_form.Setting = setting

	return c.RenderJson(detail_form)
}

func (c AuthSetting) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form SettingDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	setting := &(detail_form.Setting) //取地址，这样下面的对 setting 的修改都会反映到 detail_form 里面去的，返回的 json 里面也会有反映

	setting.Id = 0
	_, err = db_session.Insert(setting)
	core.HandleError(err)

	// 删除原来表里面的信息记录
	sql_delete := "DELETE FROM auth_setting WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, setting.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"detail":  detail_form,
		"message": "数据保存成功!",
	})
}
