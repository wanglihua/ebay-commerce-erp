package core

import (
	"strings"
	"time"

	"github.com/revel/revel"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
func GetInt64FromRequest(request *revel.Request, key string) int64 {
    value, err := strconv.ParseInt(request.Form[key][0], 10, 64)
    HandleError(err)
    return value
}
*/

func PreventSQLInjection(sqlStr string) string {
	sqlStr = strings.Replace(sqlStr, "--", "", -1)
	sqlStr = strings.Replace(sqlStr, "#", "", -1)
	sqlStr = strings.Replace(sqlStr, ";", "", -1)
	sqlStr = strings.Replace(sqlStr, "\n", "", -1)
	sqlStr = strings.Replace(sqlStr, "\r", "", -1)
	sqlStr = strings.Replace(sqlStr, "''", "", -1)

	return sqlStr
}

var expire_after_duration time.Duration = -1 // -1 表示还没获取过过期时间

func GetSessionExpire() time.Duration {
	if expire_after_duration == -1 { // -1 表示还没获取过过期时间
		expire_after_duration = 30 * 24 * time.Hour // 如果没有配置时的默认值
		if expires_string, ok := revel.Config.String("session.expires"); ok == true {
			expire_after_duration_paresed, err := time.ParseDuration(expires_string)
			if err != nil {
				expire_after_duration = expire_after_duration_paresed //如果配置了，并且值也正确，就设置为获得的值
			}

		}
	}

	return expire_after_duration
}

func IsAjaxRequest(request *revel.Request) bool {
	if len(request.Header["X-Requested-With"]) != 0 && request.Header["X-Requested-With"][0] == "XMLHttpRequest" {
		return true
	} else {
		return false
	}
}
