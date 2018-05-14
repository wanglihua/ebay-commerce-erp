package order_paypal_service

import (
	"matrix/core"
	"net/url"
	"github.com/revel/revel"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ppgl_models"
	"fmt"
)

func Order_refund(order_refund_id int64, db_session *xorm.Session) core.JsonResult {

	// 查询订单退款信息，如果已经存在就阻止
	var order_refund ppgl_models.OrderRefundInfo
	has, err := db_session.Id(order_refund_id).Get(&order_refund)
	core.HandleError(err)

	if has == false {
		return core.JsonResult{
			Success: false,
			Message: "订单退款申请信息不存在！",
		}
	}

	var paypal_payment_list = make([]ppgl_models.PaypalPaymentInfo, 0)
	err = db_session.Where("transaction_id = ?", order_refund.PayTransId).Find(&paypal_payment_list)
	core.HandleError(err)

	if len(paypal_payment_list) == 0 {
		return core.JsonResult{
			Success: false,
			Message: "订单PayPal付款信息不存在！",
		}
	}

	var form_data url.Values = make(url.Values)
	form_data.Add("order_refund_id", fmt.Sprint(order_refund_id))

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/payment/refund", form_data)
	core.HandleError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	core.HandleError(err)

	var result core.JsonResult
	err = json.Unmarshal(body, &result)
	core.HandleError(err)

	// 订单历史记录都在 JAVA 端完成

	//return result.Success
	return result
}
