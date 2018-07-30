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
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/service/email_service"
	"matrix/modules/erp/models/xxgl_models"
	"github.com/go-xorm/xorm"
)

const REPLY_METHOD_MESSAGE int64 = 1
const REPLY_METHOD_MAIL int64 = 2

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var db_session = c.DbSession

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	var shop_list = make([]xtpz_models.ShopAccountInfo, 0)
	core.HandleError(db_session.Find(&shop_list))

	c.RenderArgs["user_list"] = user_list
	c.RenderArgs["shop_list"] = shop_list

	c.RenderArgs["tag_list"] = get_message_tag_list()
	c.RenderArgs["status_list"] = get_message_status_list(db_session)

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

	var db_session = c.DbSession

	var reply_method int64
	c.Params.Bind(&reply_method, "reply_method")

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var reply_body string
	c.Params.Bind(&reply_body, "reply_body")

	if reply_method == REPLY_METHOD_MESSAGE {
		ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
		resp, err := http.PostForm(ozosync_app_base_url+"/message/reply", url.Values{"messageId": {fmt.Sprint(message_id)}, "messageBody": {reply_body}})
		core.HandleError(err)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		core.HandleError(err)

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

	if reply_method == REPLY_METHOD_MAIL {

		var recipient_email string
		c.Params.Bind(&recipient_email, "recipient_email")

		var reply_subject string
		c.Params.Bind(&reply_subject, "reply_subject")

		var reply_email_address_id int64
		c.Params.Bind(&reply_email_address_id, "reply_email_address_id")

		email_service.Send_mail_by_config(recipient_email, reply_subject, reply_body, reply_email_address_id, db_session)

		return c.RenderJson(map[string]interface{}{
			"success": true,
			"message": "邮件回复成功！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": false,
		"message": "未知回复方式！",
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

	c.RenderArgs["reply_method_message"] = REPLY_METHOD_MESSAGE
	c.RenderArgs["reply_method_mail"] = REPLY_METHOD_MAIL

	var message = message_service.Get_message_data(message_id, db_session)
	c.RenderArgs["message"] = message

	var email_address_list = make([]xtpz_models.EmailAddressInfo, 0)
	err := db_session.Where("shop_id = ?", message.Message.ShopId).Find(&email_address_list)
	core.HandleError(err)
	c.RenderArgs["email_address_list"] = email_address_list

	var product_status models.StatusInfo
	_, err = db_session.Id(message.Product.StatusId).Get(&product_status)
	core.HandleError(err)
	c.RenderArgs["product_status"] = product_status

	var sender_id = message.Message.SenderId
	var sender_message_list = message_service.Get_message_list_data_by_sender_shop_id(db_session, sender_id, message.Message.ShopId)
	c.RenderArgs["sender_message_list"] = sender_message_list

	var sender_item_order_item_list = service.Get_user_item_order_item_list(message.Message.SenderId, message.Message.ItemId, db_session)
	c.RenderArgs["sender_item_order_item_list"] = sender_item_order_item_list

	var sender_order_item_list = service.Get_user_order_item_list(message.Message.SenderId, db_session)
	c.RenderArgs["sender_order_item_list"] = sender_order_item_list

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))
	c.RenderArgs["user_list"] = user_list

	c.RenderArgs["message_tag_list"] = get_message_tag_list()

	return c.RenderTemplate("erp/xxgl/message/message_detail.html")
}

func (c ErpMessage) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XXGL_MESSAGE_ALL, c.Session)

	var db_session = c.DbSession

	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var assigned_user_id int64
	c.Params.Bind(&assigned_user_id, "assigned_user_id")

	var tag_id int64
	c.Params.Bind(&tag_id, "tag_id")

	var message xxgl_models.MessageInfo
	_, err := db_session.Id(message_id).Get(&message)
	core.HandleError(err)

	var assigned_user auth_models.UserInfo
	_, err = db_session.Id(assigned_user_id).Get(&assigned_user)
	core.HandleError(err)

	var message_tag = xxgl_models.Get_message_tag_by_id(tag_id)

	message.AssignedUserId = assigned_user_id
	message.AssignedUser = fmt.Sprintf("%s(%s)", assigned_user.NickName, assigned_user.LoginName)
	message.TagId = tag_id
	message.TagName = message_tag.Name

	_, err = db_session.Id(message_id).Cols("assigned_user_id", "assigned_user", "tag_id", "tag_name").Update(&message)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func get_message_tag_list() []xxgl_models.MessageTagInfo {
	var message_tag_list = make([]xxgl_models.MessageTagInfo, 0);
	message_tag_list = append(message_tag_list, xxgl_models.MESSAGE_TAG_NORMAL)
	message_tag_list = append(message_tag_list, xxgl_models.MESSAGE_TAG_NOREPLY)
	message_tag_list = append(message_tag_list, xxgl_models.MESSAGE_TAG_OTHER)

	return message_tag_list
}

func get_message_status_list(db_session *xorm.Session) []models.StatusInfo {
	var message_status_list = make([]models.StatusInfo, 0)

	err := db_session.Where("type = ?", models.STATUS_TYPE_MESSAGE).Find(&message_status_list)
	core.HandleError(err)

	return message_status_list
}
