package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/models"
	"matrix/core"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {
	db_session := c.DbSession

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("*");
	sql_builder.Count("count(*) as user_count")
	sql_builder.From("auth_user")
	sql_builder.Where("id > 1")
	sql_builder.OrderBy("login_name asc, id asc")
	sql_builder.Limit(1, 2);
	var user_list = make([]models.UserInfo, 0)
	err := db_session.Sql(sql_builder.SqlList()).Find(&user_list)
	core.HandleError(err)

	var user_count struct {
		UserCount int64 `xorm:"'user_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&user_count)
	core.HandleError(err)

	return c.RenderText("Hello World!")
}
