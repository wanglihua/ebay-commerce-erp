package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"matrix/core"
	"matrix/modules/auth/models"
	"strings"
)

type AuthUserPassRes struct {
	*revel.Controller
}

func (c AuthUserPassRes) Index() revel.Result {
	var user_id_list = make([]int64, 0)
	c.Params.Bind(&user_id_list, "user")

	c.RenderArgs["user_id_list"] = user_id_list
	return c.RenderTemplate("auth/pass_res/pass_res_index.html")
}

type AuthUserPassResViewItem struct {
	models.PassResInfo     `xorm:"extends" json:"r"`
	models.PassResCateInfo `xorm:"extends" json:"rc"`
	IsAuth core.NullInt    `xorm:"'is_auth'" json:"is_auth"`
}

func (c AuthUserPassRes) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT r.*, rc.*, (SELECT count(*) FROM auth_pass_user_res ur WHERE ur.res_id = r.id AND ur.user_id IN (1, 2)) AS is_auth
		FROM auth_pass_res r
		    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id
		ORDER BY is_auth ASC;
	*/

	query := db_session.
	Select(fmt.Sprintf("r.*, rc.*, (SELECT count(*) FROM auth_pass_user_res ur WHERE ur.res_id = r.id AND ur.user_id IN (%s)) AS is_auth", "1,2")).
		Table("auth_pass_res").Alias("r").
		Join("INNER", []string{"auth_pass_res_cate", "rc"}, "r.cate_id = rc.id").
		Where(filter)

	//query extra filter here

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
	db_session := c.DbSession

	user_id_list := make([]int64, 0)
	c.Params.Bind(&user_id_list, "user_id_list")

	var res_id int64
	c.Params.Bind(&res_id, "res_id")

	/*
		INSERT INTO auth_pass_user_res (user_id, res_id, create_time) VALUES (1, 1, now()), (2, 1, now());
	*/

	sql := `INSERT INTO auth_pass_user_res (user_id, res_id, create_time) VALUES %s;`

	value_list := make([]string, 0)
	for _, user_id := range user_id_list {
		value_list = append(value_list, fmt.Sprintf("(%d, %d, now())", user_id, res_id))
	}

	sql = fmt.Sprintf(sql, strings.Join(value_list, ","))

	_, err := db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "授权成功！"})
}

func (c AuthUserPassRes) Revoke() revel.Result {
	db_session := c.DbSession

	user_id_list := make([]int64, 0)
	c.Params.Bind(&user_id_list, "user_id_list")

	user_id_str_list := make([]string, 0)
	for _, user_id := range user_id_list {
		user_id_str_list = append(user_id_str_list, fmt.Sprint(user_id))
	}

	var res_id int64
	c.Params.Bind(&res_id, "res_id")

	/*
		DELETE FROM auth_pass_user_res WHERE user_id IN (1, 2) AND res_id = 1;
	*/
	sql := `DELETE FROM auth_pass_user_res WHERE user_id IN (%s) AND res_id = %d;`
	sql = fmt.Sprintf(sql, strings.Join(user_id_str_list, ","), res_id)
	_, err := db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "收权成功！"})
}
