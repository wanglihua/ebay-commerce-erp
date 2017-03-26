package tasks

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/revel/revel"
	"matrix/modules/erp/models"
)

type PaypalDownloadTask struct {
}

func (t PaypalDownloadTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Println(err)
		}
	}()

	job_time_begin := time.Now()

	// get db driver and db url

	if task_db_driver == "" {
		task_db_driver, _ = revel.Config.String("db.driver")
	}

	if task_db_url == "" {
		task_db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(task_db_driver, task_db_url)
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	paypal_account_list := make([]models.PaypalAccountInfo, 0)
	err = db_session.Find(&paypal_account_list)
	if err != nil {
		revel.ERROR.Println(err)
	}

	for _, paypal_account := range paypal_account_list {
		if strings.TrimSpace(paypal_account.ApiAccount) == "" || strings.TrimSpace(paypal_account.ApiToken) == "" {
			//如果账号里面没有设置api相关，就没有必要调用
			continue
		}

		client, err := paypalsdk.NewClient(paypal_account.ApiAccount, paypal_account.ApiToken, paypal_api_base)
		if err != nil {
			revel.ERROR.Println("paypalsdk NewClient error:", err.Error())
		}

		if revel.DevMode == true {
			client.Log = os.Stdout
		} else {
			client.Log = nil
		}

		t.check_paypal_client_token(client, paypal_account.PaypalAccount)

		/*
		payment_list, err := client.GetPayments()
		if err != nil {
			revel.ERROR.Println("paypalsdk client.GetPayments() error:", err.Error())
		}
		*/

		payment_list, err := t.get_payment_list(client)
		if err != nil {
			revel.ERROR.Println("PaypalDownloadTask.get_payment_list error:", err.Error())
		}

		revel.TRACE.Println("payment_list:", payment_list)

		//todo: insert payment_list to database
	}

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("eBay GetPayments time duration: ", job_duration.Seconds(), "seconds")
}

// 好像 paypalsdk.Client 里面也有对token的过期管理，先用自己的来管理好了
func (t PaypalDownloadTask) check_paypal_client_token(client *paypalsdk.Client, paypal_account string) {
	time_now := time.Now()

	cache_item := paypal_access_token_cache.get(paypal_account)

	if cache_item.paypal_access_token == "" || time_now.After(cache_item.paypal_access_token_expired) {
		token_response, err := client.GetAccessToken() // client.GetAccessToken() 调用后，没有必要调用 client.SetAccessToken(cache_item.paypal_access_token)
		if err != nil {
			revel.ERROR.Println("paypalsdk client.GetAccessToken() error:", err.Error())
		}

		cache_item.paypal_access_token = token_response.Token
		cache_item.paypal_access_token_expired = time_now.Add(time.Duration(token_response.ExpiresIn-180) * time.Second)
		paypal_access_token_cache.set(paypal_account, cache_item)
	} else {
		client.SetAccessToken(cache_item.paypal_access_token)
	}
}

func (t PaypalDownloadTask) get_payment_list(client *paypalsdk.Client) (payment_list []paypalsdk.Payment, err error) {
	/*
		request:

		count query_string integer
		The number of items to list in the response.
		Default: 10.

		start_id query_string string
		The ID of the starting resource in the response. When results are paged, you can use the next_id value as the start_id to continue with the next set of results.

		start_index query_string integer
		The start index of the payments to list. Typically, you use the start_index to jump to a specific position in the resource history based on its cart. For example, to start at the second item in a list of results, specify ?start_index=2.

		start_time query_string string
		The start date and time for the range to show in the response, in Internet date and time format. For example, start_time=2016-03-06T11:00:00Z.
		Format: date-time.

		end_time query_string string
		The end date and time for the range to show in the response, in Internet date and time format. For example, start_time=2016-03-06T11:00:00Z.
		Format: date-time.

		payee_id query_string string
		Filters the response by a PayPal-assigned merchant ID that identifies the payee.

		sort_by query_string enum
		Sorts the response by the create time or update time.
		Possible values: create_time, update_time.

		sort_order query_string enum
		Sorts the response in ascending or descending order.
		Possible values: asc, desc.
    */

	number_of_days := 1 //从当前算起一天时间

	end_time := time.Now().UTC()
	start_time := end_time.Add(time.Duration(-number_of_days * 24 * int(time.Hour)))

	payment_list = make([]paypalsdk.Payment, 0)

	next_payment_id := ""
	is_first_call := true

	for is_first_call == true || next_payment_id != "" { // 首次调用的时候，next_payment_id是为空的
		var query_params = make(url.Values)
		query_params.Add("count", "10")
		query_params.Add("start_time", start_time.Format("2006-01-02T15:04:05.00Z"))
		query_params.Add("end_time", end_time.Format("2006-01-02T15:04:05.00Z"))
		query_params.Add("sort_by", "create_time")
		query_params.Add("sort_order", "desc")

		if next_payment_id != "" {
			query_params.Add("start_id", next_payment_id)
		}

		request_url := fmt.Sprintf("%s%s?%s", client.APIBase, "/v1/payments/payment/", query_params.Encode())
		var response ListPaymentsResponse
		request, err := http.NewRequest("GET", request_url, nil)
		if err != nil {
			return response.Payments, err
		}

		err = client.SendWithAuth(request, &response)
		if err != nil {
			return response.Payments, err
		}

		next_payment_id = response.NextId
		is_first_call = false
		payment_list = append(payment_list, response.Payments...)
	}

	return payment_list, nil
}

/*
response:

payments array (contains the payment object)
A payment. Use this object to create, process, and manage payments.

count see description
The number of items returned in each range of results. Note that the last results range might have fewer items than the requested number of items. The maximum value is 20.
Possible types: integer

next_id string
The ID of the next element. Use to get the next range of results.
Read only.
 */

type ListPaymentsResponse struct {
	Payments []paypalsdk.Payment `json:"payments"`
	Count    int64               `json:"count"`
	NextId   string              `json:"next_id"`
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
