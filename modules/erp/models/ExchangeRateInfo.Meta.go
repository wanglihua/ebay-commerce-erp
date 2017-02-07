package models

var ExchangeRateCols = struct {
	Id string
	CurrencyId string
	Rate string
	StatusId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	CurrencyId: "currency_id",
	Rate: "rate",
	StatusId: "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
