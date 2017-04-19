package pass_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel/cache"
	"matrix/core"
)

const (
	PASS_ADMIN_LIST_CACHE_KEY = "pass_admin_list"
	PASS_USER_RES_CACHE_KEY   = "pass_user_res_%d"
)

// 加载指定用户的权限信息到缓存，包括管理员列表
func Load_user_pass_to_cache(user_id int64, db_session *xorm.Session) {

	// 缓存管理员列表
	load_admin_list_to_cache(db_session)

	// 缓存用户权限
	type user_res_id_t struct {
		ResId int64 `xorm:"'res_id'"`
	}

	/*
		SELECT ur.res_id AS res_id FROM auth_pass_user_res ur WHERE ur.user_id = 0
	 */

	var sql_user_res_id_list = `SELECT ur.res_id AS res_id FROM auth_pass_user_res ur WHERE ur.user_id = %d`
	sql_user_res_id_list = fmt.Sprintf(sql_user_res_id_list, user_id)

	var user_res_id_list = make([]user_res_id_t, 0)
	err := db_session.Sql(sql_user_res_id_list).Find(&user_res_id_list)
	core.HandleError(err)

	var user_res_id_map = make(map[int64]interface{})
	for _, user_res_id := range user_res_id_list {
		user_res_id_map[user_res_id.ResId] = true
	}

	err = cache.Set(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), user_res_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

// 加载所有用户的权限信息到缓存，包括管理员列表
func Load_all_pass_to_cache(db_session *xorm.Session) {

	// 缓存管理员列表
	load_admin_list_to_cache(db_session)

	// 缓存用户权限
	type user_res_id_t struct {
		UserId int64 `xorm:"'user_id'"`
		ResId  int64 `xorm:"'res_id'"`
	}

	/*
		SELECT ur.user_id AS user_id, ur.res_id AS res_id FROM auth_pass_user_res ur
	 */

	var sql_user_res_id_list = `SELECT ur.user_id AS user_id, ur.res_id AS res_id FROM auth_pass_user_res ur`

	var user_res_id_list = make([]user_res_id_t, 0)
	err := db_session.Sql(sql_user_res_id_list).Find(&user_res_id_list)
	core.HandleError(err)

	// 每个用户各一个资源列表
	var user_res_list_map = make(map[int64](map[int64]interface{}))

	for _, user_res_id := range user_res_id_list {
		user_res_id_map, found := user_res_list_map[user_res_id.UserId]
		if found == false {
			user_res_id_map = make(map[int64]interface{})
		}

		user_res_id_map[user_res_id.ResId] = true
		user_res_list_map[user_res_id.UserId] = user_res_id_map
	}

	// 把各个用户的资源列表都放到缓存中去
	for user_id := range user_res_list_map {
		var user_res_id_map = user_res_list_map[user_id]
		err = cache.Set(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), user_res_id_map, core.GetSessionExpire())
		core.HandleError(err)
	}
}

// 加载管理员列表到缓存
func load_admin_list_to_cache(db_session *xorm.Session) {
	type admin_user_id_t struct {
		UserId int64 `xorm:"'user_id'"`
	}

	/*
		SELECT a.user_id AS user_id FROM auth_admin a
	 */

	var sql_admin_list = `SELECT a.user_id AS user_id FROM auth_admin a`

	var admin_user_id_list = make([]admin_user_id_t, 0)
	err := db_session.Sql(sql_admin_list).Find(&admin_user_id_list)
	core.HandleError(err)

	var admin_user_id_map = make(map[int64]interface{})
	for _, admin_user_id := range admin_user_id_list {
		admin_user_id_map[admin_user_id.UserId] = true
	}

	err = cache.Set(PASS_ADMIN_LIST_CACHE_KEY, admin_user_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

// 检查用户权限
func Check_user_pass(user_id int64) bool {
	return check_user_pass_in_cache(user_id)
}

// 通过缓存的信息检查用户权限
func check_user_pass_in_cache(user_id int64) bool {
	//检查是不是管理员
	var admin_user_id_map = make(map[int64]interface{})
	err := cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	core.HandleError(err)

	is_admin, found := admin_user_id_map[user_id]
	if found == true && is_admin == true {
		return true
	}

	// 如果不是管理员，就检查有没有授权
	var user_res_id_map = make(map[int64]interface{})
	err = cache.Get(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), &user_res_id_map)
	core.HandleError(err)

	has_res, found := user_res_id_map[user_id]
	if found == true && has_res == true {
		return true
	}

	return false
}

func check_user_pass_in_db(user_id int64, res_id int64, db_session *xorm.Session) bool {

	// 检查是不是管理员
	var admin_count struct {
		Count int64 `xorm:"'admin_count'"`
	}

	/*
		SELECT count(a.id) AS admin_count FROM auth_admin a WHERE a.user_id = 0;
	 */

	var sql_admin_count = `SELECT count(a.id) AS admin_count FROM auth_admin a WHERE a.user_id = %d`
	sql_admin_count = fmt.Sprintf(sql_admin_count, user_id)
	_, err := db_session.Sql(sql_admin_count).Get(&admin_count)
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
