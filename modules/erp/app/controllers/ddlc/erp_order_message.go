package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/core"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type ErpOrderMessage struct {
	*revel.Controller
}

func (c ErpOrderMessage) show() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SEND_MESSAGE, c.Session)

	return c.RenderTemplate("erp/ddlc/order_message/order_message_send_to_customers.html")
}

func (c ErpOrderMessage) send() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_SEND_MESSAGE, c.Session)

	var orderIdList = make([]int64, 0)
	c.Params.Bind(&orderIdList, "id_list")

	var subject string
	c.Params.Bind(&subject, "subject")

	var body string
	c.Params.Bind(&body, "body")

	var orderIdStrList = make([]string, 0)
	for _, orderId := range orderIdList {
		orderIdStrList = append(orderIdStrList, fmt.Sprint(orderId))
	}

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/send/message/to/customers", url.Values{"orderIdList": orderIdStrList, "subject": {subject}, "body": {body}})
	core.HandleError(err)

	defer resp.Body.Close()
	response_body, err := ioutil.ReadAll(resp.Body)
	core.HandleError(err)

	var result core.JsonResult
	err = json.Unmarshal(response_body, &result)
	core.HandleError(err)

	if result.Success == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "消息回复失败！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "消息回复成功！",
	})
}
