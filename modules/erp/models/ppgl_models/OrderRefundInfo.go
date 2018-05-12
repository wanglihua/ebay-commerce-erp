package ppgl_models

import (
	"matrix/core"
)

type OrderRefundInfo struct {
	Id            int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId        int64     `xorm:"bigint notnull index 'shop_id'" json:"shop_id"`
	ShopName      string    `xorm:"nvarchar(200) notnull 'shop_name'" json:"shop_name"`
	OrderId       int64     `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	OrderCode     string    `xorm:"nvarchar(200) notnull 'order_code'" json:"order_code"`
	OrderAmount   float64   `xorm:"decimal(18,6) notnull 'order_amount'" json:"order_amount"`
	ExpressFee    float64   `xorm:"decimal(18,6) notnull 'express_fee'" json:"express_fee"`   // 收入运费
	ExpressCost   float64   `xorm:"decimal(18,6) notnull 'express_cost'" json:"express_cost"` // 支出运费
	TransactionId string    `xorm:"nvarchar(100) notnull index 'transaction_id'" json:"transaction_id"`
	Amount        float64   `xorm:"decimal(18,6) notnull 'amount'" json:"amount"`
	TypeId        int64     `xorm:"bigint notnull index 'type_id'" json:"type_id"`
	TypeName      string    `xorm:"nvarchar(200) notnull 'type_name'" json:"type_name"`
	ReasonId      int64     `xorm:"bigint notnull index 'reason_id'" json:"reason_id"`
	Reason        string    `xorm:"nvarchar(200) notnull 'reason'" json:"reason"`
	Description   string    `xorm:"nvarchar(500) notnull 'description'" json:"description"`
	ProductId     int64     `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	ProductCode   string    `xorm:"nvarchar(200) notnull 'product_code'" json:"product_code"`
	ProductName   string    `xorm:"nvarchar(200) notnull 'product_name'" json:"product_name"`
	SalesUserId   int64     `xorm:"bigint notnull index 'sales_user_id'" json:"sales_user_id"`
	SalesUserName string    `xorm:"nvarchar(200) notnull 'sales_user_name'" json:"sales_user_name"`
	StatusId      int64     `xorm:"bigint notnull index 'status_id'" json:"status_id"`
	CreateTime    core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e OrderRefundInfo) TableName() string {
	return "erp_order_refund"
}

func (e OrderRefundInfo) ModelDesc() string {
	return "verbose_name=订单退款信息,entity_json=order_refund,route=order/refund"
}
