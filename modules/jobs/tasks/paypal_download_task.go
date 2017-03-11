package tasks

import (
	"time"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/revel/revel"
)

type PaypalDownloadTask struct {
	access_token         string
	access_token_expired time.Time
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

	if t.access_token == "" || time_now.After(t.access_token_expired) {
		token_response, err := client.GetAccessToken()
		if err != nil {
			revel.ERROR.Println("paypalsdk client.GetAccessToken() error:", err.Error())
		}

		t.access_token = token_response.Token
		t.access_token_expired = time_now.Add(time.Duration(token_response.ExpiresIn-180) * time.Second)
	} else {
		client.SetAccessToken(t.access_token)
	}
}
