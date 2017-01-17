package models

var LoginLogCols = struct {
	Id         string
	LoginName  string
	IpAddress  string
	Content    string
	Desc       string
	CreateTime string
}{
	Id:         "id",
	LoginName:  "login_name",
	IpAddress:  "ip_address",
	Content:    "content",
	Desc:       "desc",
	CreateTime: "create_time",
}
