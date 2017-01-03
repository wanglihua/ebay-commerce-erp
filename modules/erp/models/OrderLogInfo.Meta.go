package models

var OrderLogCols = struct {
	Id string
	OrderId string
	OrderCode string
	TransCode string
	UserId string
	LoginName string
	Operation string
	CreateTime string
}{
	Id: "id",
	OrderId: "order_id",
	OrderCode: "order_code",
	TransCode: "trans_code",
	UserId: "user_id",
	LoginName: "login_name",
	Operation: "operation",
	CreateTime: "create_time",
}
