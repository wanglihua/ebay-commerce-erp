package models


import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	ORDER_TYPE_SGCJ = 1010 // 手工创建
	ORDER_TYPE_XZDD = 1020 // 下载订单
	ORDER_TYPE_CFDD = 1030 // 拆分订单
	ORDER_TYPE_HBDD = 1040 // 合并订单
)

type OrderTypeInfo struct {
	Id          int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,max=500"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderTypeInfo) TableName() string {
	return TablePrefix + "order_type"
}

func (e OrderTypeInfo) ModelDesc() string {
	return "verbose_name=订单类型,entity_json=order_type,route=order_type"
}

func (e OrderTypeInfo) InitData(db_session *xorm.Session) {
	order_type := new(OrderTypeInfo)
	order_type.Id = ORDER_TYPE_SGCJ
	order_type.Name = "手工创建"
	order_type.Desc = core.NewNullString("订单类型 手工创建", true)
	db_session.Insert(order_type)

	order_type = new(OrderTypeInfo)
	order_type.Id = ORDER_TYPE_XZDD
	order_type.Name = "下载订单"
	order_type.Desc = core.NewNullString("订单类型 下载订单", true)
	db_session.Insert(order_type)

	order_type = new(OrderTypeInfo)
	order_type.Id = ORDER_TYPE_CFDD
	order_type.Name = "拆分订单"
	order_type.Desc = core.NewNullString("订单类型 拆分订单", true)
	db_session.Insert(order_type)

	order_type = new(OrderTypeInfo)
	order_type.Id = ORDER_TYPE_HBDD
	order_type.Name = "合并订单"
	order_type.Desc = core.NewNullString("订单类型 合并订单", true)
	db_session.Insert(order_type)
}
