package xtpz_models

import (
	"matrix/core"
)

type ExchangeRateLogInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string    `xorm:"nvarchar(200) notnull 'name'" json:"name" smith:"verbose_name=名称"`
	Desc       string    `xorm:"nvarchar(500) notnull 'desc'" json:"desc" smith:"verbose_name=操作"`
	LoginName  string    `xorm:"nvarchar(200) notnull 'login_name'" json:"login_name" smith:"verbose_name=工号"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e ExchangeRateLogInfo) TableName() string {
	return "erp_exchange_rate_log"
}

func (e ExchangeRateLogInfo) ModelDesc() string {
	return "verbose_name=货币汇率操作日志,entity_json=rate_log,route=exchange/rate/log"
}
