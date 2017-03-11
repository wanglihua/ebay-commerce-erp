package tasks

import (
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/revel/revel"
)

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

	t.check_paypal_client_token(client)

	revel.TRACE.Println("PaypalDownloadTask", "Run()")
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
