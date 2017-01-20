package pass

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

// 有一种做法是登录时把用户的权限全部缓存起来，这样下次检查权限的时候，就不用查询数据库了
// 这里先用不缓存的方法来做

// 检查当前登录用户是否有 res_id 的操作权限
func UserHasPass(res_id int, web_session revel.Session, db_session *xorm.Session) bool {
	//获得当前登录用户的Id
	login_user := core.GetLoginUser(web_session)
	login_user_id := login_user.UserId

	type UserResCount struct {
		Count int64 `xorm:"'count'" json:"count"`
	}

	var user_res_count UserResCount
	sql := `
SELECT count(*) AS count FROM erp_pass_user_res WHERE user_id = ? AND res_id = ?;
	`
	_, err := db_session.Sql(sql, login_user_id, res_id).Get(&user_res_count)
	core.HandleError(err)

	if user_res_count.Count == 1 {
		return true
	} else {
		return false
	}
}

// 检查当前登录用户是否有 res_id 的操作权限
func CheckUserPass(res_id int, web_session revel.Session, db_session *xorm.Session) {
	if UserHasPass(res_id, web_session, db_session) == false {
		panic("对不起，您没有当前操作的权限！")
	}
}
