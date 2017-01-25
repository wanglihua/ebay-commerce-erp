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

	var user_list = make([]models.UserInfo, 0)
	err := db_session.Sql("SELECT * FROM auth_user").Where("id > 1").Limit(10, 0).Find(&user_list)
	core.HandleError(err)

	return c.RenderText("Hello World!")
}
