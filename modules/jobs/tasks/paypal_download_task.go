package tasks

import (
	"os"
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/revel/revel"
)

// 好像 paypalsdk.Client 里面也有对token的过期管理，先用自己的来管理好了
var paypal_access_token string
var paypal_access_token_expired time.Time

type PaypalDownloadTask struct {
}

func (t PaypalDownloadTask) Run() {
	config := GetPaypalConfig()
	client, err := paypalsdk.NewClient(config.ClientID, config.Secret, config.APIBase)
	if err != nil {
		revel.ERROR.Println("paypalsdk NewClient error:", err.Error())
	}

	if revel.DevMode == true {
		client.Log = os.Stdout
	} else {
		client.Log = nil
	}

	t.check_paypal_client_token(client)

	payment_list, err := client.GetPayments()
	if err != nil {
		revel.ERROR.Println("paypalsdk client.GetPayments() error:", err.Error())
	}

	revel.TRACE.Println("payment_list:", payment_list)
}

func (t PaypalDownloadTask) check_paypal_client_token(client *paypalsdk.Client) {
	time_now := time.Now()

	if paypal_access_token == "" || time_now.After(paypal_access_token_expired) {
		token_response, err := client.GetAccessToken()
		if err != nil {
			revel.ERROR.Println("paypalsdk client.GetAccessToken() error:", err.Error())
		}

		paypal_access_token = token_response.Token
		paypal_access_token_expired = time_now.Add(time.Duration(token_response.ExpiresIn-180) * time.Second)
	} else {
		client.SetAccessToken(paypal_access_token)
	}
}

/*
paypal api 文档 链接：

1、
http://www.cnblogs.com/see7di/archive/2012/10/21/2732825.html

2、
http://blog.csdn.net/hipihi_0/article/details/7967767
TransactionSearch和GetTransactionDetails

3、
http://www.cnblogs.com/charje/archive/2011/12/01/paypal_ipn.html

4、
http://stackoverflow.com/questions/16611336/paypal-restful-api-pull-list-of-sale-transactions

5、
http://stackoverflow.com/questions/18860883/paypal-rest-api-list-of-refund-transactions?noredirect=1&lq=1

6、
https://www.google.com/?gws_rd=ssl#q=paypal+api+get+transactions&*


 */
