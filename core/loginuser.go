package core

import (
	"fmt"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

const LOGIN_USER_CACHE_KEY = "Login_User_%s"

//根据以后需要再在这里添加成员
type LoginUser struct {
	UserId       int64
	LoginName    string
	NickName     string
	IsSuperAdmin bool

	// IsAdmin      bool // 这个不使用
}

func GetLoginUserCacheKey(session revel.Session) string {
	return fmt.Sprintf(LOGIN_USER_CACHE_KEY, session.Id())
}

func GetLoginUser(session revel.Session) LoginUser {
	var loginUser LoginUser
	err := cache.Get(GetLoginUserCacheKey(session), &loginUser)
	if err != nil {
		revel.TRACE.Println(err.Error())
	}
	return loginUser
}

func PutLoginUserToSession(session revel.Session, user LoginUser) {
	cache.Set(GetLoginUserCacheKey(session), user, GetSessionExpire())
}

func RemoveLoginUserInSession(session revel.Session) {
	cache.Delete(GetLoginUserCacheKey(session))
}
