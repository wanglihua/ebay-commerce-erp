package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	PORT_NINGBO   = 100
	PORT_SHANGHAI = 200
	PORT_SHENZHEN = 300
)

type ContainerTransferStartPortInfo struct {
	Id         int64      `xorm:"bigint notnull pk 'id'" json:"id"`
	PortName   string     `xorm:"nvarchar(200) notnull 'port_name'" json:"port_name" smith:"verbose_name=港口名称"`
	TradeItem  string     `xorm:"nvarchar(200) notnull 'trade_item'" json:"trade_item" smith:"verbose_name=贸易协议"`
	CreateTime core.Time  `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time  `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ContainerTransferStartPortInfo) TableName() string {
	return TablePrefix + "container_transfer_start_port"
}

func (e ContainerTransferStartPortInfo) ModelDesc() string {
	return "verbose_name=集装箱调拨起始港口,entity_json=transfer_port,route=container/transfer/start/port"
}

func (e ContainerTransferStartPortInfo) InitData(db_session *xorm.Session) {
	var port = new(ContainerTransferStartPortInfo)
	port.Id = PORT_NINGBO
	port.PortName = "NINGBO CHINA"
	port.TradeItem = "FOB NINGBO"
	db_session.Insert(port)

	port = new(ContainerTransferStartPortInfo)
	port.Id = PORT_SHANGHAI
	port.PortName = "SHANGHAI CHINA"
	port.TradeItem = "FOB SHANGHAI"
	db_session.Insert(port)

	port = new(ContainerTransferStartPortInfo)
	port.Id = PORT_SHENZHEN
	port.PortName = "SHENZHEN CHINA"
	port.TradeItem = "FOB SHENZHEN"
	db_session.Insert(port)
}
