package models

var PaypalAccountCols = struct {
	Id            string
	PaypalAccount string
	ApiAccount    string
	ApiToken      string
	Balance       string
	Desc          string
	StatusId      string
	CreateTime    string
	UpdateTime    string
	Version       string
}{
	Id:            "id",
	PaypalAccount: "paypal_account",
	ApiAccount:    "api_account",
	ApiToken:      "api_token",
	Balance:       "balance",
	Desc:          "desc",
	StatusId:      "status_id",
	CreateTime:    "create_time",
	UpdateTime:    "update_time",
	Version:       "version",
}
