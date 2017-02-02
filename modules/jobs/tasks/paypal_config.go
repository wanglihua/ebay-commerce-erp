package tasks

import "github.com/logpacker/PayPal-Go-SDK"

/*
https://www.sandbox.paypal.com/

注册 paypal sandbox 用户
https://www.sandbox.paypal.com/signup/account
***REMOVED***
***REMOVED***

https://www.sandbox.paypal.com/c2/merchantsignup/create
***REMOVED***
alpha_142536

*/

type PaypalConfig struct {
	ClientID string
	Secret   string
	APIBase  string // paypalsdk.APIBaseSandBox or paypalsdk.APIBaseLive
}

func GetPaypalConfig() PaypalConfig {
	return PaypalConfig{
		ClientID: "AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f",
		Secret:   "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd",
		APIBase:  paypalsdk.APIBaseSandBox,
	}
}
