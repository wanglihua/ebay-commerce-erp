package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/service/message_service"
	auth_models "matrix/modules/auth/models"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {

	var db_session = c.DbSession

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	c.RenderArgs["service_user_list"] = user_list

	return c.RenderTemplate("erp/xxgl/message/message_index.html")
}

func (c ErpMessage) ListData() revel.Result {
	db_session := c.DbSession

	message_list, msg_count := message_service.Get_message_list_data(db_session, c.Request, c.Params)

	return c.RenderJson(core.GridResult{
		Data:  message_list,
		Total: msg_count,
	})
}

func (c ErpMessage) ShowReply() revel.Result {
	var db_session = c.DbSession

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var message = message_service.Get_message_data(message_id, db_session)

	c.RenderArgs["message"] = message

	return c.RenderTemplate("erp/xxgl/message/message_reply.html")
}

func (c ErpMessage) Reply() revel.Result {

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var message_body string
	c.Params.Bind(&message_body, "message_body")

	resp, err := http.PostForm("http://localhost:8080/ozosync/message/reply", url.Values{"messageId": {fmt.Sprint(message_id)}, "messageBody": {message_body}})
	core.HandleError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	core.HandleError(err)

	fmt.Println(string(body))

	var result core.JsonResult
	err = json.Unmarshal(body, &result)
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
