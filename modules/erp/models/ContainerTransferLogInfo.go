package models

import (
	"matrix/core"
)

type ContainerTransferLogInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TransferId int64           `xorm:"bigint notnull index 'transfer_id'" json:"transfer_id" smith:"verbose_name=调拨id"`
	BatchCode  string          `xorm:"nvarchar(200) notnull index 'batch_code'" json:"batch_code" smith:"verbose_name=调拨批次号"`
	UserId     int64           `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName  string          `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=用户工号"`
	Operation  string          `xorm:"nvarchar(200) notnull 'operation'" json:"operation"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ContainerTransferLogInfo) TableName() string {
	return TablePrefix + "container_transfer_log"
}

func (e ContainerTransferLogInfo) ModelDesc() string {
	return "verbose_name=集装箱调拨日志,entity_json=transfer_log,route=container/transfer/log"
}
