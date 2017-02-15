package models

var ProductImportCols = struct {
	Id          string
	Sku         string
	Name        string
	CostPrice   string
	Weight      string
	UserId      string
	ImportTime  string
}{
	Id:          "id",
	Sku:         "sku",
	Name:        "name",
	CostPrice:   "cost_price",
	Weight:      "weight",
	UserId:      "user_id",
	ImportTime:  "import_time",
}
