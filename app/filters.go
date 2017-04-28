package app

import (
	"fmt"
	"runtime/debug"

	"github.com/revel/revel"
	"matrix/core"
)

//捕获和处理整站异常
func PanicFilter(c *revel.Controller, fc []revel.Filter) {
	defer func() {
		if err := recover(); err != nil {
			handleInvocationPanic(c, err)
		}
	}()
	fc[0](c, fc[1:])
}

func handleInvocationPanic(c *revel.Controller, err interface{}) {
	//这个函数里面不能出错（panic），这里出错，程序就挂了
	defer func() {
		recover() //不做任何事，保证程序不退出就行
	}()

	errorDesc := fmt.Sprint(err)
	revel.ERROR.Print(errorDesc, "\n", string(debug.Stack()))

	c.Result = c.RenderJson(core.JsonResult{Success: false, Message: "操作失败！详细:" + errorDesc})

	/*
	// 这个判断无法区别，是上传文件的请求还是普通的请求
	if core.IsAjaxRequest(c.Request) {
		c.Result = c.RenderJson(core.JsonResult{Success: false, Message: "操作失败！详细:" + errorDesc})
	} else {
		if revel.DevMode == true {
			c.Result = c.RenderError(fmt.Errorf("%s \r\n %s", errorDesc, string(debug.Stack())))
		} else {
			c.Result = c.RenderText(fmt.Sprintf("%s", errorDesc))
		}
	}
	*/
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	var response_header = c.Response.Out.Header()

	// Add some common security headers
	response_header.Add("X-Frame-Options", "SAMEORIGIN")
	response_header.Add("X-XSS-Protection", "1; mode=block")
	response_header.Add("X-Content-Type-Options", "nosniff")

	//add copyright
	response_header.Add("dev-company", "Jiaxing Huadu Information Technolony Co.,Ltd")
	//add email and mobile

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
