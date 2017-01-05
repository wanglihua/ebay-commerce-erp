package models

var StockBillCols = struct {
	Id string
	Code string
	StockId string
	Desc string
	Date string
	UserId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	Code: "code",
	StockId: "stock_id",
	Desc: "desc",
	Date: "date",
	UserId: "check_user_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
