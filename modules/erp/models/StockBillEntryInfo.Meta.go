package models

var StockBillEntryCols = struct {
	Id string
	BillId string
	StorageLocId string
	ProductId string
	PreQty string
	Qty string
	Price string
	Amount string
	AddTime string
}{
	Id: "id",
	BillId: "bill_id",
	StorageLocId: "storage_loc_id",
	ProductId: "product_id",
	PreQty: "pre_qty",
	Qty: "qty",
	Price: "price",
	Amount: "amount",
	AddTime: "add_time",
}
