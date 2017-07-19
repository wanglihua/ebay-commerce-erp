package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	DEST_PORT_MELBOURNE = 100
	DEST_PORT_BRISBANE  = 200
	DEST_PORT_PERTH     = 300
	DEST_PORT_SYDNEY    = 400
)

type ContainerTransferDestPortInfo struct {
	Id         int64      `xorm:"bigint notnull pk 'id'" json:"id"`
	PortName   string     `xorm:"nvarchar(200) notnull 'port_name'" json:"port_name" smith:"verbose_name=港口名称"`
	CreateTime core.Time  `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time  `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ContainerTransferDestPortInfo) TableName() string {
	return TablePrefix + "container_transfer_dest_port"
}

func (e ContainerTransferDestPortInfo) ModelDesc() string {
	return "verbose_name=集装箱调目的港口,entity_json=transfer_port,route=container/transfer/start/port"
}

func (e ContainerTransferDestPortInfo) InitData(db_session *xorm.Session) {
	var port = new(ContainerTransferDestPortInfo)
	port.Id = DEST_PORT_MELBOURNE
	port.PortName = "Melbourne Australia"
	db_session.Insert(port)

	port = new(ContainerTransferDestPortInfo)
	port.Id = DEST_PORT_BRISBANE
	port.PortName = "Brisbane Australia"
	db_session.Insert(port)

	port = new(ContainerTransferDestPortInfo)
	port.Id = DEST_PORT_PERTH
	port.PortName = "Perth Australia"
	db_session.Insert(port)

	port = new(ContainerTransferDestPortInfo)
	port.Id = DEST_PORT_SYDNEY
	port.PortName = "Sydney Australia"
	db_session.Insert(port)
}
