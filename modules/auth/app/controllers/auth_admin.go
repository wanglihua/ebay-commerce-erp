package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"matrix/core"
	"matrix/modules/auth/models"
	"strings"
)

type AuthAdmin struct {
	*revel.Controller
}

func (c AuthAdmin) Index() revel.Result {
	return c.RenderTemplate("auth/admin/admin_index.html")
}

type AdminViewItem struct {
	models.AdminInfo `xorm:"extends" json:"a"`
	models.UserInfo  `xorm:"extends" json:"u"`
}

func (c AuthAdmin) ListData() revel.Result {
	session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := session.
	Select("a.*, u.*").
		Table(models.TablePrefix + "user").Alias("u").
		Join("inner", []string{models.TablePrefix + "admin", "a"}, "u.id = a.user_id").
		Where(filter)

	//query extra filter here

	dataQuery := *query
	if order != "" {
		dataQuery = *dataQuery.OrderBy(order)
	} else {
		dataQuery = *dataQuery.Asc("a.id")
	}

	adminList := make([]AdminViewItem, 0, limit)
	err := dataQuery.Limit(limit, offset).Find(&adminList)
	core.HandleError(err)

	countQuery := *query
	count, err := countQuery.Count(new(AdminViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  adminList,
		Total: count,
	})
}

func (c AuthAdmin) Add() revel.Result {
	//session := c.DbSession
	return c.RenderTemplate("auth/admin/admin_add.html")
}

func (c AuthAdmin) AddListData() revel.Result {
	session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := session.Where(filter).
		And(fmt.Sprintf("id NOT IN (SELECT user_id FROM %sadmin)", models.TablePrefix))

	//query extra filter here

	dataQuery := *query
	if order != "" {
		dataQuery = *dataQuery.OrderBy(order)
	} else {
		dataQuery = *dataQuery.Asc("id")
	}

	userList := make([]models.UserInfo, 0, limit)
	err := dataQuery.Limit(limit, offset).Find(&userList)
	core.HandleError(err)

	countQuery := *query
	count, err := countQuery.Count(new(models.UserInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  userList,
		Total: count,
	})
}

func (c AuthAdmin) AddSave() revel.Result {
	session := c.DbSession

	userIdList := make([]int64, 0)
	c.Params.Bind(&userIdList, "id_list")

	if len(userIdList) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "添加失败，请选择用户再添加!"})
	}

	for _, userId := range userIdList {
		admin := new(models.AdminInfo)
		admin.UserId = userId
		_, err := session.Insert(admin)
		core.HandleError(err)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "添加成功!"})
}

func (c AuthAdmin) Remove() revel.Result {
	db_session := c.DbSession

	adminIdList := make([]int, 0)
	c.Params.Bind(&adminIdList, "id_list")

	admin_id_str_list := make([]string, 0)
	for _, admin_id := range adminIdList {
		admin_id_str_list = append(admin_id_str_list, fmt.Sprintf("%d", admin_id))
	}

	admin_id_strs := strings.Join(admin_id_str_list, ",")

	/*
		SELECT count(*) AS count
		FROM auth_admin a
		    INNER JOIN auth_user u ON a.user_id = u.id
		WHERE a.id IN (1, 2, 3) AND u.login_name = '0000';
	*/

	sql := `
SELECT count(*) AS count
FROM auth_admin a
    INNER JOIN auth_user u ON a.user_id = u.id
WHERE a.id IN (%s) AND u.login_name = '0000';
	`
	sql = fmt.Sprintf(sql, admin_id_strs)

	var admin_count struct {
		Count int64 `xorm:"'count'" json:"count"`
	}

	_, err := db_session.Sql(sql).Get(&admin_count)
	core.HandleError(err)

	if admin_count.Count != 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "移除失败，内置管理员不能移除!"})
	}

	/*
	admin := new(models.AdminInfo)
	count, err := db_session.Count(admin)
	if int64(len(adminIdList)) == count {
		return c.RenderJson(core.JsonResult{Success: false, Message: "移除失败，管理员不能全部移除!"})
	}
	*/

	admin := new(models.AdminInfo)
	_, err = db_session.In("id", adminIdList).Delete(admin)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "移除成功!"})
}
