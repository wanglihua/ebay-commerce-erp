package models

var ProductCols = struct {
	Id string
	Sku string
	Name string
	EnglishName string
	CateId string
	CostPrice string
	Weight string
	SupplierId string
	Desc string
	StatusId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	Sku: "sku",
	Name: "name",
	EnglishName: "english_name",
	CateId: "cate_id",
	CostPrice: "cost_price",
	Weight: "weight",
	SupplierId: "supplier_id",
	Desc: "desc",
	StatusId: "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
