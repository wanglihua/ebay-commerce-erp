package core

import (
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

const LOGIN_USER_CACHE_KEY = "Login_User_b6db637be022550d"

//根据以后需要再在这里添加成员
type LoginUser struct {
	UserId       int64
	LoginName    string
	NickName     string
	IsAdmin      bool
	IsSuperAdmin bool
}

func GetLoginUser(session revel.Session) LoginUser {
	var loginUser LoginUser
	err := cache.Get(LOGIN_USER_CACHE_KEY+session.Id(), &loginUser)
	if err != nil {
		revel.TRACE.Println(err.Error())
	}
	return loginUser
}

func PutLoginUserToSession(session revel.Session, user LoginUser) {
	cache.Set(LOGIN_USER_CACHE_KEY+session.Id(), user, GetSessionExpire())
}

func RemoveLoginUserInSession(session revel.Session) {
	cache.Delete(LOGIN_USER_CACHE_KEY + session.Id())
}
