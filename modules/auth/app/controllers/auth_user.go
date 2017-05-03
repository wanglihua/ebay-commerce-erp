package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/forms"
	"matrix/modules/auth/models"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/auth_setting_service"
	"matrix/modules/auth/service/group_service"
	"matrix/modules/auth/service/pass_service"
)

type AuthUser struct {
	*revel.Controller
}

func (c AuthUser) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	var has_auth_pass_res = pass_service.Is_login_user_has_pass(pass_res.AUTH_PASS, c.Session)

	c.RenderArgs["has_auth_pass_res"] = has_auth_pass_res
	return c.RenderTemplate("auth/user/user_index.html")
}

func (c AuthUser) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type user_t struct {
		User       models.UserInfo `xorm:"extends" json:"u"`
		GroupNames string          `xorm:"-" json:"group_names"`
	}

	/*
	SELECT u.* FROM auth_user u
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("u.*")
	sb.Count("count(u.id) AS user_count")
	sb.From("auth_user u")
	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("u.id ASC")
	}

	sb.Limit(offset, limit)

	user_list := make([]user_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&user_list)
	core.HandleError(err)

	var user_id_list = make([]int64, 0)
	for _, user := range user_list {
		user_id_list = append(user_id_list, user.User.Id)
	}

	var users_group_name_list_map = group_service.Get_users_group_name_list_map(user_id_list, db_session)
	for index, user := range user_list {
		user_list[index].GroupNames = strings.Join(users_group_name_list_map[user.User.Id], "， ")
	}

	var user_count struct {
		Count int64 `xorm:"'user_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&user_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  user_list,
		Total: user_count.Count,
	})
}

func (c AuthUser) RefreshPassCache() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	var db_session = c.DbSession

	pass_service.Load_all_pass_to_cache(db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "刷新权限缓存成功！",
	})
}

func (c AuthUser) Detail() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	session := c.DbSession

	var userId int64
	c.Params.Bind(&userId, "id")

	user := new(models.UserInfo)
	if userId != 0 {
		has, err := session.Id(userId).Get(user)
		core.HandleError(err)
		if has == false {
			panic("指定的用户不存在！")
		}
	}

	form := new(forms.UserDetailForm)
	form.User = *user

	c.UnbindToRenderArgs(form, "form")
	c.RenderArgs["is_create"] = form.IsCreate()
	return c.RenderTemplate("auth/user/user_detail.html")
}

func (c AuthUser) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	session := c.DbSession

	form := new(forms.UserDetailForm)
	c.Params.Bind(form, "form")

	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	user := &form.User

	var affected int64
	if form.IsCreate() {
		count, err := session.Where("login_name = ?", user.LoginName).Count(new(models.UserInfo))
		core.HandleError(err)
		if count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，工号已存在！"})
		}

		user.Password = core.EncryptPassword(user.Password)

		affected, err = session.Insert(user)
		core.HandleError(err)
	} else {
		count, err := session.Where("id <> ? and login_name = ?", user.Id, user.LoginName).Count(new(models.UserInfo))
		core.HandleError(err)
		if count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，工号已存在！"})
		}

		if strings.TrimSpace(form.NewPassword) != "" {
			user.Password = core.EncryptPassword(form.NewPassword)
		}

		affected, err = session.Id(user.Id).Update(user)
		//affected, err := session.Id(form.Id).Cols("nick_name").Update(user)
		//affected, err := session.Table(new(User)).Id(form.Id).Update(map[string]interface{}{"password":"123456"})
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c AuthUser) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	session := c.DbSession

	userIdList := make([]int64, 0)
	c.Params.Bind(&userIdList, "id_list")

	count, err := session.In("user_id", userIdList).Count(new(models.AdminInfo))
	core.HandleError(err)
	if count != 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "不能删除管理员！"})
	}

	user := new(models.UserInfo)
	affected, err := session.In("id", userIdList).Delete(user)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c AuthUser) ResetPassword() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_USER, c.Session)

	db_session := c.DbSession

	user_id_list := make([]int64, 0)
	c.Params.Bind(&user_id_list, "id_list")

	if len(user_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "用户列表不能为空，请选择用户！",
		})
	}

	var login_user = core.GetLoginUser(c.Session)

	count, err := db_session.In("user_id", user_id_list).Count(new(models.AdminInfo))
	core.HandleError(err)

	if count != 0 && login_user.IsSuperAdmin == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "不能重置管理员的密码，需要超级管理员权限！"})
	}

	var auth_setting = auth_setting_service.Get_auth_setting(db_session)
	var default_password = auth_setting.DefaultPassword

	/*
	UPDATE auth_user SET password = '123' WHERE id IN (-1, -2);
	 */

	var sql_reset_pwd = fmt.Sprintf(`UPDATE auth_user SET password = '%s' WHERE id IN (%s)`, core.EncryptPassword(default_password), core.Int_list_to_str_list(user_id_list))
	_, err = db_session.Exec(sql_reset_pwd)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "用户密码重置成功！",
	})
}

// 这里就不进行权限检查了，别的地方要调用
func (c AuthUser) GetUserInfo() revel.Result {
	db_session := c.DbSession

	var user_id int64
	c.Params.Bind(&user_id, "user_id")

	var user models.UserInfo
	_, err := db_session.Id(user_id).Get(&user)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"user": user,
	})
}
