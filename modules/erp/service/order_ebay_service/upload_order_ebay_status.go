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
