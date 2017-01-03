package models

var ProductLogCols = struct {
	Id string
	ProductId string
	ProductCode string
	UserId string
	LoginName string
	Operation string
	CreateTime string
}{
	Id: "id",
	ProductId: "product_id",
	ProductCode: "product_code",
	UserId: "user_id",
	LoginName: "login_name",
	Operation: "operation",
	CreateTime: "create_time",
}
