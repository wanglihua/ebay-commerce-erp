package pass_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Check_user_pass(user_id int64, res_id int64, db_session *xorm.Session) bool {

	// 检查是不是管理员
	var admin_count struct{
		Count int64 `xorm:"'admin_count'"`
	}

	/*
SELECT count(a.id) AS admin_count FROM auth_admin a WHERE a.user_id = 0;
	 */

	var sql_admin_count = `SELECT count(a.id) AS admin_count FROM auth_admin a WHERE a.user_id = %d`
	sql_admin_count = fmt.Sprintf(sql_admin_count, user_id)
	_,err:=db_session.Sql(sql_admin_count).Get(&admin_count)
	core.HandleError(err)

	// 如果是管理员，就返回有权限
	if admin_count.Count > 0 {
		return true
	}

	//如果不是管理员，就检查是不是被授权了
	var pass_count struct {
		Count int64 `xorm:"'pass_count'"`
	}

	/*
SELECT count(ur.id) AS pass_count FROM auth_pass_user_res ur WHERE ur.user_id = 0 AND ur.res_id = 0
	 */

	var sql_pass_count = "SELECT count(ur.id) AS pass_count FROM auth_pass_user_res ur WHERE ur.user_id = %d AND ur.res_id = %d"
	sql_pass_count = fmt.Sprintf(sql_pass_count, user_id, res_id)
	_, err = db_session.Sql(sql_pass_count).Get(&pass_count)
	core.HandleError(err)

	if pass_count.Count == 0 {
		return false
	}

	return true
}
