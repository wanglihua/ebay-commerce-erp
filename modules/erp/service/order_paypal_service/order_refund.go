package order_paypal_service

import (
	"matrix/core"
	"net/url"
	"fmt"
	"github.com/revel/revel"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ppgl_models"
)

func Order_refund(order_id int64, amount float64, memo string, db_session *xorm.Session) core.JsonResult {

	// todo
	// 查询订单退款信息，如果已经存在就阻止

	var sql = `SELECT pp.*
FROM erp_order_external_transaction oet
    INNER JOIN erp_paypal_payment pp ON pp.id = oet.payment_id
WHERE pp.type = 'Payment' AND order_id = %d;`

	sql = fmt.Sprintf(sql, order_id)

	var paypal_payment_list = make([]ppgl_models.PaypalPaymentInfo, 0)
	err := db_session.Sql(sql).Find(&paypal_payment_list)
	core.HandleError(err)

	if len(paypal_payment_list) == 0 {
		return core.JsonResult{
			Success: false,
			Message: "订单PayPal付款信息不存在！",
		}
	}

	var transaction_id = paypal_payment_list[0].TransactionId

	var form_data url.Values = make(url.Values)

	/*
    public ModelAndView refund(@RequestParam(value = "transactionId") Long transactionId,
                               @RequestParam(value = "amount") BigDecimal amount,
                               @RequestParam(value = "memo") String memo) throws Exception {
	 */

	form_data.Add("transactionId", transaction_id)
	form_data.Add("amount", fmt.Sprint(amount))
	form_data.Add("memo", memo)

	ozosync_app_base_url, _ := revel.Config.String("ozosync.app")
	resp, err := http.PostForm(ozosync_app_base_url+"/payment/refund", form_data)
	core.HandleError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	core.HandleError(err)

	var result core.JsonResult
	err = json.Unmarshal(body, &result)
	core.HandleError(err)

	// 记录订单退款信息 和 订单历史记录都在 JAVA 端完成

	//return result.Success
	return result
}
