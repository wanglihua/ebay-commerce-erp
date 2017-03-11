package tasks

import (
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

		payment_list, err := client.GetPayments()
		if err != nil {
			revel.ERROR.Println("paypalsdk client.GetPayments() error:", err.Error())
		}

		revel.TRACE.Println("payment_list:", payment_list)
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
