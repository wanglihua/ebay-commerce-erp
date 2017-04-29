package controllers

import (
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/models"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type AuthUserPassRes struct {
	*revel.Controller
}

func (c AuthUserPassRes) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_PASS, c.Session)

	db_session := c.DbSession

	var user_id_list = make([]int64, 0)
	c.Params.Bind(&user_id_list, "user")

	if len(user_id_list) == 0 {
		//return c.RenderJson(core.JsonResult{Success: false, Message: "用户列表不能为空，请选择用户！"})
		//panic("用户列表不能为空，请选择用户！")
		return c.RenderText("用户列表不能为空，请选择用户！")
	}

	// 获取已选中的用户名称汇总 begin
	type user_name_t struct {
		UserName string `xorm:"'user_name'"`
	}

	var user_id_str_list = core.Int_list_to_str_list(user_id_list)
	sql := `SELECT concat(u.login_name, ' ', u.nick_name) AS user_name FROM auth_user u WHERE u.id IN (%s)`
	sql = fmt.Sprintf(sql, strings.Join(user_id_str_list, ","))

	var user_name_struct_list = make([]user_name_t, 0)
	err := db_session.Sql(sql).Find(&user_name_struct_list)
	core.HandleError(err)

	var user_name_list = make([]string, 0)
	for _, user_name_struct := range user_name_struct_list {
		user_name_list = append(user_name_list, user_name_struct.UserName)
	}
	// 获取已选中的用户名称汇总 end

	var pass_res_cate_list = make([]models.PassResCateInfo, 0)
	err = db_session.OrderBy("id ASC").Find(&pass_res_cate_list)
	core.HandleError(err)

	c.RenderArgs["user_id_list"] = user_id_list
	c.RenderArgs["user_names"] = strings.Join(user_name_list, ", ")
	c.RenderArgs["pass_res_cate_list"] = pass_res_cate_list
	return c.RenderTemplate("auth/pass_res/pass_res_index.html")
}

type AuthUserPassResViewItem struct {
	models.PassResInfo     `xorm:"extends" json:"r"`
	models.PassResCateInfo `xorm:"extends" json:"rc"`
	IsAuth core.NullInt `xorm:"'is_auth'" json:"is_auth"`
}

func (c AuthUserPassRes) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_PASS, c.Session)

	db_session := c.DbSession

	var user_ids_str string
	c.Params.Bind(&user_ids_str, "user_id_list")

	if len(user_ids_str) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "用户列表不能为空，请选择用户！"})
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT r.*, rc.*, (SELECT count(*) FROM auth_pass_user_res ur WHERE ur.res_id = r.id AND ur.user_id IN (1, 2)) AS is_auth
		FROM auth_pass_res r
		    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id
		ORDER BY is_auth ASC;
	*/

	query := db_session.
	Select(fmt.Sprintf("r.*, rc.*, (SELECT count(*) FROM auth_pass_user_res ur WHERE ur.res_id = r.id AND ur.user_id IN (%s)) AS is_auth", user_ids_str)).
		Table("auth_pass_res").Alias("r").
		Join("INNER", []string{"auth_pass_res_cate", "rc"}, "r.cate_id = rc.id").
		Where(filter)

	/*
        aoData.push({'name': 'pass_res_cate_id_query', 'value': $("#pass_res_cate_query").val()});
        aoData.push({'name': 'pass_res_cate_name_query', 'value': $("#pass_res_name_query").val()});
	 */
	//query extra filter here
	var pass_res_cate_id_query int64
	c.Params.Bind(&pass_res_cate_id_query, "pass_res_cate_id_query")

	var pass_res_cate_name_query string
	c.Params.Bind(&pass_res_cate_name_query, "pass_res_cate_name_query")

	if pass_res_cate_id_query != 0 {
		query = query.Where(fmt.Sprintf("rc.id = %d", pass_res_cate_id_query))
	}

	if pass_res_cate_name_query != "" {
		query = query.Where(fmt.Sprintf("r.name LIKE '%%%s%%'", pass_res_cate_name_query))
	}

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("r.id")
	}

	pass_res_list := make([]AuthUserPassResViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&pass_res_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(AuthUserPassResViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  pass_res_list,
		Total: count,
	})
}

func (c AuthUserPassRes) Grant() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_PASS, c.Session)

	db_session := c.DbSession

	user_id_list := make([]int64, 0)
	c.Params.Bind(&user_id_list, "user_id_list")

	if len(user_id_list) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "用户列表不能为空，请选择用户！"})
	}

	var res_id int64
	c.Params.Bind(&res_id, "res_id")

	if res_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "资源Id不能为空！"})
	}

	// 除了Action权限控制，还要授权时检查当前登录人是不是有当前他授权的那个操作的权限。只能授权本人本来就有的权限。
	if pass_service.Is_login_user_has_pass(res_id, c.Session) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "你没有该资源的操作权限，不能对该资源进行授权操作！",
		})
	}

	/*
		DELETE FROM auth_pass_user_res WHERE user_id IN (1, 2) AND res_id = 2;
		INSERT INTO auth_pass_user_res (user_id, res_id, create_time) VALUES (1, 1, now()), (2, 1, now());
	*/
	// 删除原来的
	user_id_str_list := make([]string, 0)
	for _, user_id := range user_id_list {
		user_id_str_list = append(user_id_str_list, fmt.Sprint(user_id))
	}
	sql := `DELETE FROM auth_pass_user_res WHERE user_id IN (%s) AND res_id = %d;`
	sql = fmt.Sprintf(sql, strings.Join(user_id_str_list, ","), res_id)
	_, err := db_session.Exec(sql)
	core.HandleError(err)

	// 插入新的数据
	value_list := make([]string, 0)
	for _, user_id := range user_id_list {
		value_list = append(value_list, fmt.Sprintf("(%d, %d, now())", user_id, res_id))
	}
	sql = `INSERT INTO auth_pass_user_res (user_id, res_id, create_time) VALUES %s;`
	sql = fmt.Sprintf(sql, strings.Join(value_list, ","))
	_, err = db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "授权成功！"})
}

func (c AuthUserPassRes) Revoke() revel.Result {
	pass_service.Check_login_user_pass(pass_res.AUTH_PASS, c.Session)

	db_session := c.DbSession

	user_id_list := make([]int64, 0)
	c.Params.Bind(&user_id_list, "user_id_list")

	if len(user_id_list) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "用户列表不能为空，请选择用户！"})
	}

	user_id_str_list := make([]string, 0)
	for _, user_id := range user_id_list {
		user_id_str_list = append(user_id_str_list, fmt.Sprint(user_id))
	}

	var res_id int64
	c.Params.Bind(&res_id, "res_id")

	if res_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "资源Id不能为空！"})
	}

	// 除了Action权限控制，还要授权时检查当前登录人是不是有当前他授权的那个操作的权限。只能授权本人本来就有的权限。
	if pass_service.Is_login_user_has_pass(res_id, c.Session) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "你没有该资源的操作权限，不能对该资源进行收权操作！",
		})
	}

	/*
		DELETE FROM auth_pass_user_res WHERE user_id IN (1, 2) AND res_id = 1;
	*/
	sql := `DELETE FROM auth_pass_user_res WHERE user_id IN (%s) AND res_id = %d;`
	sql = fmt.Sprintf(sql, strings.Join(user_id_str_list, ","), res_id)
	_, err := db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "收权成功！"})
}
