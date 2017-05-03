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
	Login_time     time.Time
	Active_time    time.Time
}

const ONLINE_USER_LIST_CACHE_KEY = "online_user_list"

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
	online_user.Login_time = time.Now()
	online_user.Active_time = time.Now()

	online_user_list = append(online_user_list, online_user)

	cache.Set(ONLINE_USER_LIST_CACHE_KEY, online_user_list, time.Hour*720)
}

func Log_logout(login_user core.LoginUser, db_session *xorm.Session) {
}

// 每次请求的时候，刷新在线用户的活动时间
func Refresh_user_active(web_session revel.Session) {

}

func Check_user_active(db_session *xorm.Session) {
	// 加载超时配置
	var auth_setting = auth_setting_service.Get_auth_setting(db_session) // todo 可以缓存起来

	var online_user_list = get_online_user_list()

	for _, online_user := range online_user_list {
		if time.Now().Sub(online_user.Active_time).Minutes() > float64(auth_setting.TimeOut) {
			// 把用户踢出登录
		}
	}
}

func get_ip_location(ip_address string) string {
	//http://freeapi.ipip.net/120.151.39.56
	var url = fmt.Sprintf("http://freeapi.ipip.net/%s", ip_address)
	var location_script = requests.Get(url)
	var location_list = make([]string, 0)
	err := json.Unmarshal([]byte(location_script), &location_list)
	if err != nil {
		revel.ERROR.Println(err)
		revel.ERROR.Println(string(debug.Stack()))
	}

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
