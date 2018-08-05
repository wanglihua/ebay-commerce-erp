package controllers

import (
	"github.com/revel/revel"
	"matrix/app/layout"
	"matrix/core"
	"matrix/core/captcha"
	"matrix/modules/auth/forms"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/login_log_service"
	"matrix/modules/erp/models"
)

type Login struct {
	*revel.Controller
}

func (c Login) Index() revel.Result {
	c.RenderArgs["sys_name"] = layout.GetSysName(c.DbSession)
	// c.RenderArgs["static_version"] = core.GetStaticVersion()
	return c.RenderTemplate("login/login_index.html")
}

func (c Login) Login() revel.Result {
	db_session := c.DbSession

	form := new(forms.LoginForm)
	c.Params.Bind(form, "form")

	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var user auth_models.UserInfo
	has, err := db_session.Where("login_name = ?", form.LoginName).Get(&user)
	core.HandleError(err)

	if has == false {
		//当用户名和密码错误的时候，应该让原来的验证码失效
		captcha.MakeCurrentExpire(c.Session)
		// return c.RenderJson(core.JsonResult{Success: false, Message: "工号不存在！"})
		return c.RenderJson(core.JsonResult{Success: false, Message: "工号或密码错误！"})
	}

	if core.CompareHashAndPassword(user.Password, form.Password) == false {
		//当用户名和密码错误的时候，应该让原来的验证码失效
		captcha.MakeCurrentExpire(c.Session)
		// return c.RenderJson(core.JsonResult{Success: false, Message: "密码错误！"})
		return c.RenderJson(core.JsonResult{Success: false, Message: "工号或密码错误！"})
	}

	if captcha.IsCorrect(form.Captcha, c.Session) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "验证码错误！"})
	}

	if user.StatusId != models.USER_STATUS_ENABLED {
		// return c.RenderJson(core.JsonResult{Success: false, Message: "用户未启用！"})
		return c.RenderJson(core.JsonResult{Success: false, Message: "工号或密码错误！"})
	}

	c.makePreviousSessionExpire()

	//将当前登录用户信息放入Session中
	var loginUser core.LoginUser
	loginUser.UserId = user.Id
	loginUser.LoginName = user.LoginName
	loginUser.NickName = user.NickName

	// 是不是管理员
	var admin_user auth_models.AdminInfo
	has, err = db_session.Where("user_id = ?", user.Id).Get(&admin_user)
	core.HandleError(err)

	// loginUser.IsAdmin = false
	if has == true && admin_user.Id != 0 {
		// loginUser.IsAdmin = true

		// 是不是超级管理员
		if admin_user.IsSuper == 1 {
			loginUser.IsSuperAdmin = true
		}
	}

	core.PutLoginUserToSession(c.Session, loginUser)
	pass_service.Load_user_pass_to_cache(loginUser.UserId, db_session)
	//pass_service.Load_all_pass_to_cache(db_session)

	var ip_address = get_request_client_ip(c.Request)
	login_log_service.Log_login(user, ip_address, c.Session, db_session)

	return c.RenderJson(map[string]interface{}{
		"success":            true,
		"require_change_pwd": user.IsDefaultPwd > 0,
		"message":            "登陆成功！",
	})
}

func (c Login) Logout() revel.Result {
	var db_session = c.DbSession

	var login_user = core.GetLoginUser(c.Session)
	if login_user.UserId != 0 {
		var session_key = core.GetLoginUserCacheKey(c.Session)
		login_log_service.Log_logout(session_key, db_session)
		login_log_service.Remove_user_from_online(c.Session)
	}

	core.RemoveLoginUserInSession(c.Session)
	//c.makePreviousSessionExpire()
	return c.RenderJson(core.JsonResult{Success: true, Message: "注销成功！"})
}

func (c Login) makePreviousSessionExpire() {
	delete(c.Session, revel.SESSION_ID_KEY)
}

func get_request_client_ip(request *revel.Request) string {
	return revel.ClientIP(request.Request)
	//return strings.Split(request.RemoteAddr, ":")[0]
}
