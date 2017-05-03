package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"matrix/app/routes"
	"matrix/core"
	"matrix/db"
	"matrix/modules/erp/service/login_log_service"
)

func BeforeInterceptor(c *revel.Controller) revel.Result {
	if isStaticRequest(c) == false {
		if revel.DevMode {
			cache.Set("request_begin_time_"+c.Session.Id(), time.Now(), time.Minute*10) // 一般请求时间不会超过10分钟吧
			revel.TRACE.Println(blue_paint(fmt.Sprintf("### begin action %s ###", c.Action)))
		}
		//revel.TRACE.Println("session id: " + c.Session.Id())

		//在这之前的数据库访问请自行开启 db session
		dbSession := db.Engine.NewSession()
		c.DbSession = dbSession.NoAutoCondition(true)
		//c.DbSession = db.Engine.NoAutoCondition() 这个不行，xorm会有问题
		c.RenderArgs["dbsession"] = dbSession

		//if revel.DevMode == false {
		err := c.DbSession.Begin()
		core.HandleError(err)
		//}

		authResult := userAuth(c)
		if authResult != nil {
			return authResult
		}
	}

	return nil
}

func FinallyInterceptor(c *revel.Controller) revel.Result {
	if isStaticRequest(c) == false {
		//if revel.DevMode == false {
		err := c.DbSession.Commit()
		core.HandleError(err)
		//}

		c.DbSession.Close() //模板tags中的数据库访问请自行开启 db session

		if revel.DevMode {
			//revel.TRACE.Println("session id: " + c.Session.Id())
			var request_begin_time time.Time
			cache.Get("request_begin_time_"+c.Session.Id(), &request_begin_time)
			var milli_seconds float64 = float64(int64(time.Now().Sub(request_begin_time))) / float64(int64(time.Millisecond))
			revel.TRACE.Println(green_paint(fmt.Sprintf("### end   action %s ### request time: %0.3f ms", c.Action, milli_seconds)))
			revel.TRACE.Println(green_paint("----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"))
		}
	}

	return nil
}

func PanicInterceptor(c *revel.Controller) revel.Result {
	if isStaticRequest(c) == false {
		//if revel.DevMode == false {
		err := c.DbSession.Rollback()
		core.HandleError(err)
		//}
		// 这里也要 close 吧
		c.DbSession.Close();
	}

	return nil
}

func isStaticRequest(c *revel.Controller) bool {
	if c.Name != "Static" {
		return false
	} else {
		return true
	}
}

func userAuth(c *revel.Controller) revel.Result {
	/*
		Static.Serve
		Static.ServeModule
		TestRunner.Index
		TestRunner.Run
	*/
	if c.Name == "Static" || c.Name == "TestRunner" {
		return nil
	}

	/*
	if c.Action == "Login.Index" || c.Action == "Login.Login" {
		return nil
	}
	*/

	if c.Name == "Login" {
		return nil
	}

	if c.Action == "Home.Hotload" {
		return nil
	}

	if c.Action == "Home.SyncDb" || c.Action == "Home.SyncDbPost" {
		return nil
	}

	if c.Action == "Captcha.Image" {
		return nil
	}

	if c.Action == "Test.Index" {
		return nil
	}

	if strings.HasPrefix(c.Name, "Wechat") {
		return nil //暂时先 return true
	}

	// 以 Api 开头 和 以 Api 结尾 的 Controller 都不做登录限制
	if strings.HasPrefix(c.Name, "Api") || strings.HasSuffix(c.Name, "Api") {
		return nil
	}

	loginuser := core.GetLoginUser(c.Session)
	//revel.TRACE.Println("userAuth session id: " + c.Session.Id())

	if loginuser.UserId == 0 {
		revel.TRACE.Println("loginuser.UserId == 0")

		if core.IsAjaxRequest(c.Request) {
			//c.Result = c.RenderJson(JsonResult{Success: false, Message: "操作失败，未登陆或没相应权限！"})
			return c.RenderJson(core.JsonResult{Success: false, Message: "操作失败，未登陆或没相应权限！"})
		} else {
			//c.Result = c.Redirect(routes.Login.Index())
			return c.Redirect(routes.Login.Index())
		}
	} else {
		// 更新用户在线活动时间
		login_log_service.Refresh_user_active(c.Session)
	}

	return nil
}
