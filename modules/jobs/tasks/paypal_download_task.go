package tasks

import (
	"github.com/revel/revel"
	"github.com/logpacker/PayPal-Go-SDK"
)

type PaypalDownloadTask struct {
}

func (t PaypalDownloadTask) Run() {
	client, err := paypalsdk.NewClient("", "", "")

	if err != nil {
		revel.ERROR.Println("paypalsdk NewClient error:", err.Error())
	}

	revel.TRACE.Println(client)

	revel.TRACE.Println("PaypalDownloadTask", "Run()")
}
