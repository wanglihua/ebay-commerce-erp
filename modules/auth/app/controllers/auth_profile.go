package controllers

import (
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/forms"
	"matrix/modules/auth/models"
)

// 这里没有什么权限控制的，能登录的都能访问
type AuthProfile struct {
	*revel.Controller
}

func (c AuthProfile) Index() revel.Result {
	session := c.DbSession

	var require_change_pwd bool = false
	c.Params.Bind(&require_change_pwd, "require_change_pwd")

	loginUserId := core.GetLoginUser(c.Session).UserId
	loginUser := new(models.UserInfo)
	_, err := session.Id(loginUserId).Get(loginUser)
	core.HandleError(err)

	c.UnbindToRenderArgs(loginUser, "user")
	c.RenderArgs["require_change_pwd"] = require_change_pwd

	return c.RenderTemplate("auth/profile/profile_index.html")
}

func (c AuthProfile) Save() revel.Result {
	session := c.DbSession

	form := new(forms.ProfileForm)
	c.Params.Bind(form, "form")

	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	loginUserId := core.GetLoginUser(c.Session).UserId
	loginUser := new(models.UserInfo)
	_, err := session.Id(loginUserId).Get(loginUser)
	core.HandleError(err)

	if strings.TrimSpace(form.Password) != "" {
		loginUser.Password = core.EncryptPassword(form.Password)
		loginUser.IsDefaultPwd = -1 //下次登录就不会要求改密码了

		_, err = session.Id(loginUserId).Cols("password", "is_default_pwd").Update(loginUser)
		core.HandleError(err)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "保存成功!"})
}
