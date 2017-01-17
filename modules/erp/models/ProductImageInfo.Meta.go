package models

var ProductImageCols = struct {
	Id          string
	ProductId   string
	ProductCode string
	TypeId      string
	PlatformId  string
	Url         string
	CreateTime  string
	UpdateTime  string
	Version     string
}{
	Id:          "id",
	ProductId:   "product_id",
	ProductCode: "product_code",
	TypeId:      "type_id",
	PlatformId:  "platform_id",
	Url:         "url",
	CreateTime:  "create_time",
	UpdateTime:  "update_time",
	Version:     "version",
}
