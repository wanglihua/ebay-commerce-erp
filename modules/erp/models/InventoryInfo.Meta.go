package models

var InventoryCols = struct {
	Id           string
	StockId      string
	StorageLocId string
	ProductId    string
	Qty          string
	CreateTime   string
	UpdateTime   string
	Version      string
}{
	Id:           "id",
	StockId:      "stock_id",
	StorageLocId: "storage_loc_id",
	ProductId:    "product_id",
	Qty:          "qty",
	CreateTime:   "create_time",
	UpdateTime:   "update_time",
	Version:      "version",
}
