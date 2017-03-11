package tasks

import (
	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/revel/revel"
)

type PaypalDownloadTask struct {
}

func (t PaypalDownloadTask) Run() {
	config := GetPaypalConfig()
	client, err := paypalsdk.NewClient(config.ClientID, config.Secret, config.APIBase)
	if err != nil {
		revel.ERROR.Println("paypalsdk NewClient error:", err.Error())
	}

	token_response, err := client.GetAccessToken()
	if err != nil {
		revel.ERROR.Println("paypalsdk client.GetAccessToken() error:", err.Error())
	}

	revel.TRACE.Println(token_response.Token)

	revel.TRACE.Println(client)

	revel.TRACE.Println("PaypalDownloadTask", "Run()")
}
