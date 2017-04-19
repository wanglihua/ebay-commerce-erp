package pass

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
)

// 这里的代码，基本上没用了的，使用 pass_service中的代码就行

// 有一种做法是登录时把用户的权限全部缓存起来，这样下次检查权限的时候，就不用查询数据库了
// 这里先用不缓存的方法来做

// 检查当前登录用户是否有 res_id 的操作权限
func UserHasPass(res_id int, web_session revel.Session, db_session *xorm.Session) bool {
	//获得当前登录用户的Id
	login_user := core.GetLoginUser(web_session)
	login_user_id := login_user.UserId

	// 管理员有任何权限
	if login_user.IsAdmin == true {
		return true
	}

	type UserResCount struct {
		Count int64 `xorm:"'count'" json:"count"`
	}

	var user_res_count UserResCount
	sql := `
SELECT count(*) AS count FROM auth_pass_user_res WHERE user_id = ? AND res_id = ?;
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

func AdminOnly(web_session revel.Session) {
	login_user := core.GetLoginUser(web_session)
	if login_user.IsAdmin == false {
		panic("对不起，您没有当前操作的权限！")
	}
}

// 获取登录用户 授权资源类别列表
func GetLoginUserPassResCateList(db_session *xorm.Session, web_session revel.Session) []int64 {
	type UserPassResCateId struct {
		ResCateId int64 `xorm:"'res_cate_id'" json:"res_cate_id"`
	}

	var user_pass_res_cate_id_list = make([]UserPassResCateId, 0)
	sql := `
SELECT DISTINCT rc.id as res_cate_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id
WHERE ur.user_id = ?;`

	//获得当前登录用户的Id
	login_user := core.GetLoginUser(web_session)
	login_user_id := login_user.UserId

	err := db_session.Sql(sql, login_user_id).Find(&user_pass_res_cate_id_list)
	core.HandleError(err)

	var res_cate_id_list = make([]int64, 0)

	for _, user_pass_res_cate_id := range user_pass_res_cate_id_list {
		res_cate_id_list = append(res_cate_id_list, user_pass_res_cate_id.ResCateId)
	}

	return res_cate_id_list
}

// 获取登录用户 授权资源列表
func GetLoginUserPassResList(db_session *xorm.Session, web_session revel.Session) []int64 {
	type UserPassResId struct {
		ResId int64 `xorm:"'res_id'" json:"res_id"`
	}

	var user_pass_res_id_list = make([]UserPassResId, 0)
	sql := `
SELECT DISTINCT r.id AS res_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
WHERE ur.user_id = ?;`

	//获得当前登录用户的Id
	login_user := core.GetLoginUser(web_session)
	login_user_id := login_user.UserId

	err := db_session.Sql(sql, login_user_id).Find(&user_pass_res_id_list)
	core.HandleError(err)

	var res_id_list = make([]int64, 0)

	for _, user_pass_res_cate_id := range user_pass_res_id_list {
		res_id_list = append(res_id_list, user_pass_res_cate_id.ResId)
	}

	return res_id_list
}
