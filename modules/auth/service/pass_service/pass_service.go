package pass_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"matrix/core"
	"matrix/modules/auth/models"
)

const (
	PASS_ADMIN_LIST_CACHE_KEY                     = "pass_admin_list"
	PASS_USER_RES_CACHE_KEY                       = "pass_user_res_%d"
	PASS_USER_RES_CATE_CACHE_KEY                  = "pass_user_res_cate_%d"
	PASS_GROUP_SALES_USER_LIST_CACHE_KEY          = "pass_group_sales_user_list"
	PASS_GROUP_CUSTOM_SERVICE_USER_LIST_CACHE_KEY = "pass_group_custom_service_user_list"
	PASS_GROUP_PICK_CREW_USER_LIST_CACHE_KEY      = "pass_group_pick_crew_user_list"
)

// 加载指定用户的权限信息到缓存，包括类别和管理员列表，群组用户列表
func Load_user_pass_to_cache(user_id int64, db_session *xorm.Session) {
	// 缓存管理员列表
	load_admin_list_to_cache(db_session)
	load_group_user_list_to_cache(db_session)
	load_user_pass_res_to_cache(user_id, db_session)
	load_user_pass_cate_to_cache(user_id, db_session)
}

// 加载指定用户的权限信息到缓存
func load_user_pass_res_to_cache(user_id int64, db_session *xorm.Session) {

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

	var user_res_id_map = make(map[int64]bool)
	for _, user_res_id := range user_res_id_list {
		user_res_id_map[user_res_id.ResId] = true
	}

	err = cache.Set(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), user_res_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

// 加载所有用户的权限信息到缓存，包括类别和管理员列表
func Load_all_pass_to_cache(db_session *xorm.Session) {
	// 缓存管理员列表
	load_admin_list_to_cache(db_session)
	load_group_user_list_to_cache(db_session)
	load_all_pass_res_to_cache(db_session)
	load_all_pass_cate_to_cache(db_session)
}

// 加载所有用户的权限信息到缓存
func load_all_pass_res_to_cache(db_session *xorm.Session) {

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
	var user_res_list_map = make(map[int64](map[int64]bool))

	for _, user_res_id := range user_res_id_list {
		user_res_id_map, found := user_res_list_map[user_res_id.UserId]
		if found == false {
			user_res_id_map = make(map[int64]bool)
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

func load_user_pass_cate_to_cache(user_id int64, db_session *xorm.Session) {
	type user_pass_res_cate_id_t struct {
		ResCateId int64 `xorm:"'res_cate_id'"`
	}

	/*
SELECT DISTINCT rc.id as res_cate_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id
WHERE ur.user_id = 0
	 */
	sql_pass_cate_list := `
SELECT DISTINCT rc.id as res_cate_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id
WHERE ur.user_id = %d`
	sql_pass_cate_list = fmt.Sprintf(sql_pass_cate_list, user_id)

	var user_pass_res_cate_id_list = make([]user_pass_res_cate_id_t, 0)
	err := db_session.Sql(sql_pass_cate_list).Find(&user_pass_res_cate_id_list)
	core.HandleError(err)

	var user_res_cate_id_map = make(map[int64]bool)
	for _, user_pass_res_cate_id := range user_pass_res_cate_id_list {
		user_res_cate_id_map[user_pass_res_cate_id.ResCateId] = true
	}

	err = cache.Set(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, user_id), user_res_cate_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

func load_all_pass_cate_to_cache(db_session *xorm.Session) {
	type user_pass_res_cate_id_t struct {
		UserId    int64  `xorm:"'user_id'"`
		ResCateId int64  `xorm:"'res_cate_id'"`
	}

	/*
SELECT DISTINCT ur.user_id AS user_id, rc.id AS res_cate_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id;
	 */
	sql_pass_cate_list := `
SELECT DISTINCT ur.user_id AS user_id, rc.id AS res_cate_id
FROM auth_pass_user_res ur
    INNER JOIN auth_pass_res r ON ur.res_id = r.id
    INNER JOIN auth_pass_res_cate rc ON r.cate_id = rc.id`

	var user_pass_res_cate_id_list = make([]user_pass_res_cate_id_t, 0)
	err := db_session.Sql(sql_pass_cate_list).Find(&user_pass_res_cate_id_list)
	core.HandleError(err)

	// 每个用户各一个资源类别列表
	var user_res_cate_list_map = make(map[int64](map[int64]bool))

	for _, user_pass_cate_id := range user_pass_res_cate_id_list {
		user_pass_cate_id_map, found := user_res_cate_list_map[user_pass_cate_id.UserId]
		if found == false {
			user_pass_cate_id_map = make(map[int64]bool)
		}

		user_pass_cate_id_map[user_pass_cate_id.ResCateId] = true
		user_res_cate_list_map[user_pass_cate_id.UserId] = user_pass_cate_id_map
	}

	// 把各个用户的资源列表类别都放到缓存中去
	for user_id := range user_res_cate_list_map {
		var user_pass_cate_id_map = user_res_cate_list_map[user_id]
		err = cache.Set(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, user_id), user_pass_cate_id_map, core.GetSessionExpire())
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

	var admin_user_id_map = make(map[int64]bool)
	for _, admin_user_id := range admin_user_id_list {
		admin_user_id_map[admin_user_id.UserId] = true
	}

	err = cache.Set(PASS_ADMIN_LIST_CACHE_KEY, admin_user_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

func load_group_user_list_to_cache(db_session *xorm.Session) {
	type group_user_id_t struct {
		UserId int64 `xorm:"'user_id'"`
	}

	/*
	SELECT gu.user_id AS user_id FROM auth_group_user gu WHERE gu.group_id = 0;
	 */

	// 通过使用函数获取字符串的办法，可以保证每次获得的字符串都是一样的，都是原始的，没有格式化过的，修改过的
	var get_group_user_list_sql = func() string {
		return `SELECT gu.user_id AS user_id FROM auth_group_user gu WHERE gu.group_id = %d`
	}

	// load sales
	var sql_sales_group_user_list = fmt.Sprintf(get_group_user_list_sql(), models.GROUP_SALES)
	var sales_user_id_list = make([]group_user_id_t, 0)
	err := db_session.Sql(sql_sales_group_user_list).Find(&sales_user_id_list)
	core.HandleError(err)

	var sales_user_id_map = make(map[int64]bool)
	for _, sales_user_id := range sales_user_id_list {
		sales_user_id_map[sales_user_id.UserId] = true
	}

	err = cache.Set(PASS_GROUP_SALES_USER_LIST_CACHE_KEY, sales_user_id_map, core.GetSessionExpire())
	core.HandleError(err)

	// load custom service
	var sql_custom_service_group_user_list = fmt.Sprintf(get_group_user_list_sql(), models.GROUP_CUSTOM_SERVICE)
	var custom_service_user_id_list = make([]group_user_id_t, 0)
	err = db_session.Sql(sql_custom_service_group_user_list).Find(&custom_service_user_id_list)
	core.HandleError(err)

	var custom_service_user_id_map = make(map[int64]bool)
	for _, custom_service_user_id := range custom_service_user_id_list {
		custom_service_user_id_map[custom_service_user_id.UserId] = true
	}

	err = cache.Set(PASS_GROUP_CUSTOM_SERVICE_USER_LIST_CACHE_KEY, custom_service_user_id_map, core.GetSessionExpire())
	core.HandleError(err)

	// load pick crew
	var sql_pick_crew_group_user_list = fmt.Sprintf(get_group_user_list_sql(), models.GROUP_PICK_CREW)
	var pick_crew_user_id_list = make([]group_user_id_t, 0)
	err = db_session.Sql(sql_pick_crew_group_user_list).Find(&pick_crew_user_id_list)
	core.HandleError(err)

	var pick_crew_user_id_map = make(map[int64]bool)
	for _, pick_crew_user_id := range pick_crew_user_id_list {
		pick_crew_user_id_map[pick_crew_user_id.UserId] = true
	}

	err = cache.Set(PASS_GROUP_PICK_CREW_USER_LIST_CACHE_KEY, pick_crew_user_id_map, core.GetSessionExpire())
	core.HandleError(err)
}

func Check_login_user_pass(res_id int64, web_session revel.Session) {
	var login_user = core.GetLoginUser(web_session)
	if is_user_has_pass_in_cache(login_user.UserId, res_id) == false {
		panic(fmt.Sprintf("你没有操作权限 %d，请联系管理员！", res_id))
	}
}

// 如果列表里面的 其中一个 资源已经授权，就认为是授权的
func Check_login_user_pass_or(res_id_list []int64, web_session revel.Session) {
	var login_user = core.GetLoginUser(web_session)
	for _, res_id := range res_id_list {
		if is_user_has_pass_in_cache(login_user.UserId, res_id) == true {
			return
		}
	}

	var res_id_str_list = core.Int_list_to_str_list(res_id_list)
	panic(fmt.Sprintf("你没有操作权限 %d，请联系管理员！", strings.Join(res_id_str_list, ",")))
}

// 如果列表里面的 全部 资源已经授权，就认为是授权的
func Check_login_user_pass_and(res_id_list []int64, web_session revel.Session) {
	var login_user = core.GetLoginUser(web_session)
	var has_pass = true
	for _, res_id := range res_id_list {
		if is_user_has_pass_in_cache(login_user.UserId, res_id) == false {
			has_pass = false
		}
	}

	if has_pass == false {
		var res_id_str_list = core.Int_list_to_str_list(res_id_list)
		panic(fmt.Sprintf("你没有操作权限 %d，请联系管理员！", strings.Join(res_id_str_list, ",")))
	}
}

// 查看当前登录用户是不是有某个权限
func Is_login_user_has_pass(res_id int64, web_session revel.Session) bool {
	var login_user = core.GetLoginUser(web_session)
	return is_user_has_pass_in_cache(login_user.UserId, res_id)
}

func Get_login_user_pass_map_from_cache(web_session revel.Session) map[int64]bool {
	var login_user = core.GetLoginUser(web_session)

	var user_res_id_map = make(map[int64]bool)
	cache.Get(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, login_user.UserId), &user_res_id_map)
	// err := cache.Get(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, login_user.UserId), &user_res_id_map)
	// core.HandleError(err)

	return user_res_id_map
}

// 检查是不是管理员
func Is_user_admin(user_id int64) bool {
	var admin_user_id_map = make(map[int64]bool)
	cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// err := cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// core.HandleError(err)

	is_admin, found := admin_user_id_map[user_id]
	if found == true && is_admin == true {
		return true
	}

	return false
}

// 检查是不是销售人员
func Is_user_sales(user_id int64) bool {
	var sales_user_id_map = make(map[int64]bool)
	 cache.Get(PASS_GROUP_SALES_USER_LIST_CACHE_KEY, &sales_user_id_map)
	// err := cache.Get(PASS_GROUP_SALES_USER_LIST_CACHE_KEY, &sales_user_id_map)
	// core.HandleError(err)

	is_sales, found := sales_user_id_map[user_id]
	if found == true && is_sales == true {
		return true
	}

	return false
}

// 检查是不是客服人员
func Is_user_custom_service(user_id int64) bool {
	var custom_service_user_id_map = make(map[int64]bool)
	cache.Get(PASS_GROUP_CUSTOM_SERVICE_USER_LIST_CACHE_KEY, &custom_service_user_id_map)
	// err := cache.Get(PASS_GROUP_CUSTOM_SERVICE_USER_LIST_CACHE_KEY, &custom_service_user_id_map)
	// core.HandleError(err)

	is_custom_service, found := custom_service_user_id_map[user_id]
	if found == true && is_custom_service == true {
		return true
	}

	return false
}

// 检查是不是拣货人员
func Is_user_pick_crew(user_id int64) bool {
	var pick_crew_user_id_map = make(map[int64]bool)
	 cache.Get(PASS_GROUP_PICK_CREW_USER_LIST_CACHE_KEY, &pick_crew_user_id_map)
	// err := cache.Get(PASS_GROUP_PICK_CREW_USER_LIST_CACHE_KEY, &pick_crew_user_id_map)
	// core.HandleError(err)

	is_pick_crew, found := pick_crew_user_id_map[user_id]
	if found == true && is_pick_crew == true {
		return true
	}

	return false
}

// 通过缓存的信息检查用户权限
func is_user_has_pass_in_cache(user_id int64, res_id int64) bool {
	//检查是不是管理员
	var admin_user_id_map = make(map[int64]bool)
	cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// err := cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// core.HandleError(err)

	is_admin, found := admin_user_id_map[user_id]
	if found == true && is_admin == true {
		return true
	}

	// 如果不是管理员，就检查有没有授权
	var user_res_id_map = make(map[int64]bool)
	cache.Get(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), &user_res_id_map)
	// err = cache.Get(fmt.Sprintf(PASS_USER_RES_CACHE_KEY, user_id), &user_res_id_map)
	// core.HandleError(err)

	has_res, found := user_res_id_map[res_id]
	if found == true && has_res == true {
		return true
	}

	return false
}

func Is_login_user_has_pass_cate(cate_id int64, web_session revel.Session) bool {
	var login_user = core.GetLoginUser(web_session)
	return Is_user_has_pass_cate_in_cache(login_user.UserId, cate_id)
}

func Get_login_user_pass_cate_map_from_cache(web_session revel.Session) map[int64]bool {
	var login_user = core.GetLoginUser(web_session)

	var user_res_cate_id_map = make(map[int64]bool)
	cache.Get(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, login_user.UserId), &user_res_cate_id_map)
	// err := cache.Get(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, login_user.UserId), &user_res_cate_id_map)
	// core.HandleError(err)

	return user_res_cate_id_map
}

// 通过缓存的信息检查用户对资源类别的权限
func Is_user_has_pass_cate_in_cache(user_id int64, cate_id int64) bool {
	//检查是不是管理员
	var admin_user_id_map = make(map[int64]bool)
	cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// err := cache.Get(PASS_ADMIN_LIST_CACHE_KEY, &admin_user_id_map)
	// core.HandleError(err)

	is_admin, found := admin_user_id_map[user_id]
	if found == true && is_admin == true {
		return true
	}

	// 如果不是管理员，就检查有没有授权
	var user_res_cate_id_map = make(map[int64]bool)
	cache.Get(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, user_id), &user_res_cate_id_map)
	// err = cache.Get(fmt.Sprintf(PASS_USER_RES_CATE_CACHE_KEY, user_id), &user_res_cate_id_map)
	// core.HandleError(err)

	has_res, found := user_res_cate_id_map[cate_id]
	if found == true && has_res == true {
		return true
	}

	return false
}

// 只有管理员才可以操作的资源的检查
func Check_login_user_admin_only(web_session revel.Session) {
	var login_user = core.GetLoginUser(web_session)
	var is_admin = Is_user_admin(login_user.UserId)
	if is_admin == false {
		panic(fmt.Sprintf("你不是管理员，没有操作权限！"))
	}
}

//只有超级管理员才可以操作的资源的检查
func Check_login_user_super_admin_only(web_session revel.Session) {
	var login_user = core.GetLoginUser(web_session)
	var is_admin = Is_user_admin(login_user.UserId)
	if is_admin == false || login_user.IsSuperAdmin == false {
		panic(fmt.Sprintf("你不是超级管理员，没有操作权限！"))
	}
}

/*
func check_user_pass_in_db(user_id int64, res_id int64, db_session *xorm.Session) bool {

	// 检查是不是管理员
	var admin_count struct {
		Count int64 `xorm:"'admin_count'"`
	}

	// SELECT count(a.id) AS admin_count FROM auth_admin a WHERE a.user_id = 0;

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

	// SELECT count(ur.id) AS pass_count FROM auth_pass_user_res ur WHERE ur.user_id = 0 AND ur.res_id = 0

	var sql_pass_count = "SELECT count(ur.id) AS pass_count FROM auth_pass_user_res ur WHERE ur.user_id = %d AND ur.res_id = %d"
	sql_pass_count = fmt.Sprintf(sql_pass_count, user_id, res_id)
	_, err = db_session.Sql(sql_pass_count).Get(&pass_count)
	core.HandleError(err)

	if pass_count.Count == 0 {
		return false
	}

	return true
}
*/
