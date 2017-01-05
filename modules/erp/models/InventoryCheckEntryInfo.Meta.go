package models

var InventoryCheckEntryCols = struct {
	Id string
	CheckId string
	StorageLocId string
	ProductId string
	PreQty string
	CheckQty string
	QtyDiff string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	CheckId: "check_id",
	StorageLocId: "storage_loc_id",
	ProductId: "product_id",
	PreQty: "pre_qty",
	CheckQty: "check_qty",
	QtyDiff: "qty_diff",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
