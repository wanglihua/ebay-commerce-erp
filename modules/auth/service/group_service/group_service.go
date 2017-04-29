package group_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Get_users_group_name_list_map(user_id_list []int64, db_session *xorm.Session) map[int64][]string {

	var user_id_str_list = core.Int_list_to_str_list(user_id_list)

	type user_group_name_t struct {
		UserId    int64  `json:"user_id"`
		GroupName string `json:"group_name"`
	}

	/*
	SELECT gu.user_id AS user_id, g.group_name AS group_name FROM auth_group g INNER JOIN auth_group_user gu ON g.id = gu.group_id WHERE gu.user_id IN (1, 2, 3);
	 */

	var sql = `SELECT gu.user_id AS user_id, g.group_name AS group_name FROM auth_group g INNER JOIN auth_group_user gu ON g.id = gu.group_id WHERE gu.user_id IN (%s);`
	sql = fmt.Sprintf(sql, strings.Join(user_id_str_list, ","))

	var users_group_name_list = make([]user_group_name_t, 0)
	err := db_session.Sql(sql).Find(&users_group_name_list)
	core.HandleError(err)

	var users_group_name_list_map = make(map[int64][]string)
	for _, user_group_name := range users_group_name_list {
		user_group_name_list, found := users_group_name_list_map[user_group_name.UserId]
		if found == false {
			user_group_name_list = make([]string, 0)
		}
		user_group_name_list = append(user_group_name_list, user_group_name.GroupName)
		users_group_name_list_map[user_group_name.UserId] = user_group_name_list
	}

	return users_group_name_list_map
}
