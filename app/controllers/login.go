package controllers

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/revel/revel"
	"matrix/app/layout"
	"matrix/core"
	"matrix/core/captcha"
	"matrix/core/requests"
	"matrix/modules/auth/forms"
	"matrix/modules/auth/models"
	"matrix/modules/auth/service/pass_service"
	erp_models "matrix/modules/erp/models"
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

	// 是不是管理员
	var admin_user models.AdminInfo
	has, err = db_session.Where("user_id = ?", user.Id).Get(&admin_user)
	core.HandleError(err)
	loginUser.IsAdmin = false
	if has == true && admin_user.Id != 0 {
		loginUser.IsAdmin = true

		// 是不是超级管理员
		if admin_user.IsSuper == 1 {
			loginUser.IsSuperAdmin = true
		}
	}

	core.PutLoginUserToSession(c.Session, loginUser)
	pass_service.Load_user_pass_to_cache(loginUser.UserId, db_session)
	//pass_service.Load_all_pass_to_cache(db_session)

	//login log
	var login_log erp_models.LoginLogInfo
	login_log.LoginName = user.LoginName

	var client_ip = get_request_client_ip(c.Request)

	login_log.IpAddress = client_ip
	login_log.Content = "员工登录"
	login_log.Desc = core.NullStringFrom(get_ip_location(client_ip))

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

func get_ip_location(ip_address string) string {
	//http://freeapi.ipip.net/120.151.39.56
	var url = fmt.Sprintf("http://freeapi.ipip.net/%s", ip_address)
	var location_script = requests.Get(url)
	var location_list = make([]string, 0)
	err := json.Unmarshal([]byte(location_script), &location_list)
	if err != nil {
		revel.ERROR.Println(err)
		revel.ERROR.Println(string(debug.Stack()))
	}

	// 去掉空白的项
	var location_list_no_empty = make([]string, 0)
	for _, loc := range location_list {
		if strings.TrimSpace(loc) != "" {
			location_list_no_empty = append(location_list_no_empty, loc)
		}
	}
	return strings.Join(location_list_no_empty, " ")
}
