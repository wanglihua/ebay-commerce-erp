package ddlc_models

import (
	"matrix/core"
)

type CombineOrderInfo struct {
	Id           int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId      int64     `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	OrderCode    string    `xorm:"nvarchar(200) notnull 'order_code'" json:"order_code"`
	SubOrderId   int64     `xorm:"bigint notnull index 'sub_order_id'" json:"sub_order_id"`
	SubOrderCode string    `xorm:"nvarchar(200) notnull 'sub_order_code'" json:"sub_order_code"`
	CreateTime   core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e CombineOrderInfo) TableName() string {
	return "erp_combine_order"
}

func (e CombineOrderInfo) ModelDesc() string {
	return "verbose_name=合并订单,entity_json=combine_order,route=combine/order"
}
