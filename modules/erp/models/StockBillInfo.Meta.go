package models

var StockBillCols = struct {
	Id string
	Code string
	StockId string
	SrcStockId string
	Desc string
	Date string
	UserId string
	TypeId string
	StatusId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	Code: "code",
	StockId: "stock_id",
	SrcStockId: "src_stock_id",
	Desc: "desc",
	Date: "date",
	UserId: "user_id",
	TypeId: "type_id",
	StatusId: "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
