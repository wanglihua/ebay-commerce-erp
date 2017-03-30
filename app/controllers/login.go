package controllers

import (
	"strings"

	"github.com/revel/revel"
	"matrix/app/layout"
	"matrix/core"
	"matrix/core/captcha"
	"matrix/modules/auth/forms"
	"matrix/modules/auth/models"
	erp_models "matrix/modules/erp/models"
)

type Login struct {
	*revel.Controller
}

func (c Login) Index() revel.Result {
	c.RenderArgs["sys_name"] = layout.GetSysName(c.DbSession)
	c.RenderArgs["static_version"] = core.GetStaticVersion()
	return c.RenderTemplate("login/login_index.html")
}

func (c Login) Login() revel.Result {
	db_session := c.DbSession

	form := new(forms.LoginForm)
	c.Params.Bind(form, "form")

	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	user := new(models.UserInfo)
	has, err := db_session.Where("login_name = ?", form.LoginName).Get(user)
	core.HandleError(err)

	if has == false {
		//当用户名和密码错误的时候，应该让原来的验证码失效
		captcha.MakeCurrentExpire(c.Session)
		return c.RenderJson(core.JsonResult{Success: false, Message: "工号不存在！"})
	}

	if core.CompareHashAndPassword(user.Password, form.Password) == false {
		//当用户名和密码错误的时候，应该让原来的验证码失效
		captcha.MakeCurrentExpire(c.Session)
		return c.RenderJson(core.JsonResult{Success: false, Message: "密码错误！"})
	}

	if captcha.IsCorrect(form.Captcha, c.Session) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "验证码错误！"})
	}

	c.makePreviousSessionExpire()

	//将当前登录用户信息放入Session中
	var loginUser core.LoginUser
	loginUser.UserId = user.Id
	loginUser.LoginName = user.LoginName
	loginUser.NickName = user.NickName

	admin_count, err := db_session.Where("user_id = ?", user.Id).Count(new(models.AdminInfo))
	core.HandleError(err)
	loginUser.IsAdmin = false
	if admin_count != 0 {
		loginUser.IsAdmin = true
	}

	core.PutLoginUserToSession(c.Session, loginUser)

	//login log
	var login_log erp_models.LoginLogInfo
	login_log.LoginName = user.LoginName
	login_log.IpAddress = get_request_client_ip(c.Request)
	login_log.Content = "员工登录"
	db_session.Insert(login_log)

	return c.RenderJson(core.JsonResult{Success: true, Message: "登陆成功！"})
}

func (c Login) Logout() revel.Result {
	core.RemoveLoginUserInSession(c.Session)
	//c.makePreviousSessionExpire()
	return c.RenderJson(core.JsonResult{Success: true, Message: "注销成功！"})
}

func (c Login) makePreviousSessionExpire() {
	delete(c.Session, revel.SESSION_ID_KEY)
}

func get_request_client_ip(request *revel.Request) string {
	return strings.Split(request.RemoteAddr, ":")[0]
}
