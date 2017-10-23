package ddpz_models

import (
	"matrix/core"
)

type FavStockTimeRangeInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Weekday     int64     `xorm:"bigint notnull 'weekday'" json:"weekday" smith:"verbose_name=星期"`
	WeekdayName string    `xorm:"-" json:"weekday_name" smith:"verbose_name=星期"`
	BeginHour   int64     `xorm:"bigint notnull 'begin_hour'" json:"begin_hour" smith:"verbose_name=开始小时"`
	BeginMinute int64     `xorm:"bigint notnull 'begin_minute'" json:"begin_minute" smith:"verbose_name=开始分钟"`
	EndHour     int64     `xorm:"bigint notnull 'end_hour'" json:"end_hour" smith:"verbose_name=结束小时"`
	EndMinute   int64     `xorm:"bigint notnull 'end_minute'" json:"end_minute" smith:"verbose_name=结束分钟"`
	Description string    `xorm:"nvarchar(300) notnull 'description' default('')" json:"description" smith:"verbose_name=描述"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e FavStockTimeRangeInfo) TableName() string {
	return "erp_fav_stock_time_range"
}

func (e FavStockTimeRangeInfo) ModelDesc() string {
	return "verbose_name=优先仓库时间段,entity_json=fav_stock_time_range,route=fav/stock/time/range"
}
