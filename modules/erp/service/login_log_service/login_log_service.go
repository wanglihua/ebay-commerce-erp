package login_log_service

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"matrix/core"
	"matrix/core/requests"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/service/auth_setting_service"
	"matrix/modules/erp/models"
)

type online_user_t struct {
	Login_user_key string
	Active_time    time.Time
}

const ONLINE_USER_LIST_CACHE_KEY = "online_user_list"

var online_user_list_expire = time.Hour * 720

func Log_login(user auth_models.UserInfo, ip_address string, web_session revel.Session, db_session *xorm.Session) {
	//login log
	var login_log models.LoginLogInfo
	login_log.UserId = user.Id
	login_log.LoginName = user.LoginName

	var client_ip = ip_address
	login_log.IpAddress = client_ip
	login_log.IpLoc = get_ip_location(client_ip)

	login_log.LoginTime = core.Time(time.Now())

	db_session.Insert(login_log)

	// 放入在线用户列表跟踪
	var online_user_list = get_online_user_list()

	var online_user online_user_t
	online_user.Login_user_key = core.GetLoginUserCacheKey(web_session)
	online_user.Active_time = time.Now()

	online_user_list = append(online_user_list, online_user)

	cache.Set(ONLINE_USER_LIST_CACHE_KEY, online_user_list, online_user_list_expire)
}

func Log_logout(login_user core.LoginUser, db_session *xorm.Session) {
	// 先获取用户最近的一条登录记录
	/*
	SELECT * FROM erp_login_log WHERE user_id = 1 ORDER BY create_time DESC LIMIT 0, 1
	 */
	var sql = `SELECT * FROM erp_login_log WHERE user_id = %d ORDER BY create_time DESC LIMIT 0, 1`
	sql = fmt.Sprintf(sql, login_user.UserId)
	var login_log models.LoginLogInfo
	has, err := db_session.Sql(sql).Get(&login_log)
	if has && err == nil {
		login_log.LogoutTime = core.Time(time.Now())
		login_log.OnlineTime = time.Time(login_log.LogoutTime).Sub(time.Time(login_log.LoginTime)).Minutes()

		db_session.Id(login_log.Id).Cols("logout_time", "online_time").Update(&login_log)
	}
}

// 每次请求的时候，刷新在线用户的活动时间
func Refresh_user_active(web_session revel.Session) {
	var login_user_key = core.GetLoginUserCacheKey(web_session)
	var online_user_list = get_online_user_list()
	for index, online_user := range online_user_list {
		if online_user.Login_user_key == login_user_key { //这样用 cache key来比较要合理一点，否则一个用户多个地方登录时，就会有问题
			online_user_list[index].Active_time = time.Now()
		}
	}

	cache.Set(ONLINE_USER_LIST_CACHE_KEY, online_user_list, online_user_list_expire)
}

func Check_user_active(db_session *xorm.Session) {
	// 加载超时配置
	var auth_setting = auth_setting_service.Get_auth_setting(db_session) // 可以考虑缓存起来

	var online_user_list = get_online_user_list()

	var to_remove_login_user_key_map = make(map[string]bool)

	for _, online_user := range online_user_list {
		if time.Now().Sub(online_user.Active_time).Minutes() > float64(auth_setting.TimeOut) {
			// 获得在线登录的用户
			var login_user core.LoginUser
			cache.Get(online_user.Login_user_key, &login_user)

			// 把用户踢出登录
			cache.Delete(online_user.Login_user_key)

			// 从在线用户列表中删除
			to_remove_login_user_key_map[online_user.Login_user_key] = true

			// 记录登出日志
			Log_logout(login_user, db_session)
		}
	}

	var new_online_user_list = make([]online_user_t, 0)
	for _, online_user := range online_user_list {
		if _, found := to_remove_login_user_key_map[online_user.Login_user_key]; found == false {
			new_online_user_list = append(new_online_user_list, online_user)
		}
	}

	cache.Set(ONLINE_USER_LIST_CACHE_KEY, new_online_user_list, online_user_list_expire)
}

// 用户自己手工登出时从在线用户列表里面删除掉
func Remove_user_from_online(web_session revel.Session) {
	var login_user_key = core.GetLoginUserCacheKey(web_session)

	var online_user_list = get_online_user_list()

	for index, online_user := range online_user_list {
		if online_user.Login_user_key == login_user_key {
			online_user_list = append(online_user_list[:index], online_user_list[index:]...)
		}
	}

	cache.Set(ONLINE_USER_LIST_CACHE_KEY, online_user_list, online_user_list_expire)
}

func get_ip_location(ip_address string) string {
	defer func() {
		// 有时网络不通什么的，课内
		if err := recover(); err != nil {
			revel.ERROR.Println(err)
			// revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	//http://freeapi.ipip.net/120.151.39.56
	var url = fmt.Sprintf("http://freeapi.ipip.net/%s", ip_address)
	var location_script = requests.Get(url)
	var location_list = make([]string, 0)
	err := json.Unmarshal([]byte(location_script), &location_list)
	core.HandleError(err)

	// 去掉空白的项
	var location_list_no_empty = make([]string, 0)
	for _, loc := range location_list {
		if strings.TrimSpace(loc) != "" {
			location_list_no_empty = append(location_list_no_empty, loc)
		}
	}
	return strings.Join(location_list_no_empty, " ")
}

func get_online_user_list() []online_user_t {
	var online_user_list = make([]online_user_t, 0)
	cache.Get(ONLINE_USER_LIST_CACHE_KEY, &online_user_list)
	return online_user_list
}
