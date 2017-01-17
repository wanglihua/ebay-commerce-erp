package models

var StockBillEntryCols = struct {
	Id           string
	BillId       string
	StorageLocId string
	ProductId    string
	PreQty       string
	Qty          string
	Price        string
	Amount       string
	CreateTime   string
	UpdateTime   string
	Version      string
}{
	Id:           "id",
	BillId:       "bill_id",
	StorageLocId: "storage_loc_id",
	ProductId:    "product_id",
	PreQty:       "pre_qty",
	Qty:          "qty",
	Price:        "price",
	Amount:       "amount",
	CreateTime:   "create_time",
	UpdateTime:   "update_time",
	Version:      "version",
}
