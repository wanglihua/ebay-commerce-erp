package ppgl_models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

type OrderRefundTypeInfo struct {
	Id          int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string    `xorm:"nvarchar(200) notnull index 'name'" json:"name"`
	Description string    `xorm:"nvarchar(500) 'description'" json:"description"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderRefundTypeInfo) TableName() string {
	return "erp_order_refund_type"
}

func (e OrderRefundTypeInfo) ModelDesc() string {
	return "verbose_name=订单退款类型,entity_json=order_refund_type,route=order/refund/type"
}

func (e OrderRefundTypeInfo) InitData(db_session *xorm.Session) {
}
