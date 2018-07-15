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
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
	"matrix/modules/erp/app/controllers/xxgl/service"
)

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var db_session = c.DbSession

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	c.RenderArgs["user_list"] = user_list

	return c.RenderTemplate("erp/xxgl/message/message_index.html")
}

func (c ErpMessage) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	db_session := c.DbSession

	message_list, msg_count := message_service.Get_message_list_data(db_session, c.Request, c.Params)

	return c.RenderJson(core.GridResult{
		Data:  message_list,
		Total: msg_count,
	})
}

func (c ErpMessage) ShowReply() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var db_session = c.DbSession

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var message = message_service.Get_message_data(message_id, db_session)

	c.RenderArgs["message"] = message

	return c.RenderTemplate("erp/xxgl/message/message_reply.html")
}

func (c ErpMessage) Reply() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var message_body string
	c.Params.Bind(&message_body, "message_body")

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/message/reply", url.Values{"messageId": {fmt.Sprint(message_id)}, "messageBody": {message_body}})
	core.HandleError(err)

	fmt.Println("Url: " + ozosync_app_base_url + "/message/reply")

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

func (c ErpMessage) Detail(message_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var db_session = c.DbSession

	/*
	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	c.RenderArgs["user_list"] = user_list
	*/

	var message = message_service.Get_message_data(message_id, db_session)
	c.RenderArgs["message"] = message

	var product_status models.StatusInfo
	_, err := db_session.Id(message.Product.StatusId).Get(&product_status)
	core.HandleError(err)
	c.RenderArgs["product_status"] = product_status

	var sender_id = message.Message.SenderId
	var sender_message_list = message_service.Get_message_list_data_by_sender(db_session, sender_id)
	c.RenderArgs["sender_message_list"] = sender_message_list

	var sender_item_order_item_list = service.Get_user_item_order_item_list(message.Message.SenderId, message.Message.ItemId, db_session)
	c.RenderArgs["sender_item_order_item_list"] = sender_item_order_item_list

	var sender_order_item_list = service.Get_user_order_item_list(message.Message.SenderId, db_session)
	c.RenderArgs["sender_order_item_list"] = sender_order_item_list

	return c.RenderTemplate("erp/xxgl/message/message_detail.html")
}
