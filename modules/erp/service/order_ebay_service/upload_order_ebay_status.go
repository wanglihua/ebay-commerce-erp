package order_ebay_service

import (
	"github.com/revel/revel"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"matrix/core"
)

func Upload_order_ebay_status(order_id_list []int64) core.JsonResult {

	/*
	var message_id int64
	c.Params.Bind(&message_id, "message_id")

	var message_body string
	c.Params.Bind(&message_body, "message_body")

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/message/reply", url.Values{"messageId": {fmt.Sprint(message_id)}, "messageBody": {message_body}})
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
	 */

	var form_data url.Values = make(url.Values)

	for _, order_id := range order_id_list {
		form_data.Add("orderIdList[]", fmt.Sprint(order_id))
	}

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/order/ship", form_data)
	core.HandleError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	core.HandleError(err)

	var result core.JsonResult
	err = json.Unmarshal(body, &result)
	core.HandleError(err)

	//return result.Success
	return result
}
