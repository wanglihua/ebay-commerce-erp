package models

var InventoryCheckCols = struct {
	Id string
	Code string
	StockId string
	Desc string
	CheckDate string
	CheckUserId string
	ApproveUserId string
	StatusId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	Code: "code",
	StockId: "stock_id",
	Desc: "desc",
	CheckDate: "check_date",
	CheckUserId: "check_user_id",
	ApproveUserId: "approve_user_id",
	StatusId: "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
