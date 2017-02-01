package tasks

import (
	"github.com/revel/revel"
	"github.com/logpacker/PayPal-Go-SDK"
)

type PaypalDownloadTask struct {
}

func (t PaypalDownloadTask) Run() {
	config := GetPaypalConfig()
	client, err := paypalsdk.NewClient(config.ClientID, config.Secret, config.APIBase)

	if err != nil {
		revel.ERROR.Println("paypalsdk NewClient error:", err.Error())
	}

	revel.TRACE.Println(client)

	revel.TRACE.Println("PaypalDownloadTask", "Run()")
}
