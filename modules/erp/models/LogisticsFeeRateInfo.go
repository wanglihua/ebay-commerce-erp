package models

import "matrix/core"

type LogisticsFeeRateInfo struct {
	Id       int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ModeId   int64           `xorm:"bigint notnull 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`
	Postcode string          `xorm:"nvarchar(200) notnull index 'postcode'" json:"postcode" smith:"verbose_name=邮编,min=2,max=100"`
	Qbj      float64         `xorm:"decimal(18,6) notnull 'qbj'" json:"qbj" smith:"verbose_name=起步价"`
	Qbjzl    float64         `xorm:"decimal(18,6) notnull 'qbjzl'" json:"qbjzl" smith:"verbose_name=起步价重量"`
	Zdsf     float64         `xorm:"decimal(18,6) notnull 'zdsf'" json:"zdsf" smith:"verbose_name=最低收费"`
	Xzjt     float64         `xorm:"decimal(18,6) notnull 'xzjt'" json:"xzjt" smith:"verbose_name=续重阶梯"`
	Xzjtsf     float64         `xorm:"decimal(18,6) notnull 'xzjtsf'" json:"xzjtsf" smith:"verbose_name=续重阶梯收费"`
	Sfpydq     int             `xorm:"int notnull 'sfpydq'" json:"sfpydq" smith:"verbose_name=是否为偏远地区"` // 0 否， 1 是
	Pydqjsyf   float64         `xorm:"decimal(18,6) notnull 'pydqjsyf'" json:"pydqjsyf" smith:"verbose_name=偏远地区加收运费"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsFeeRateInfo) TableName() string {
	return TablePrefix + "logistics_fee_rate"
}

func (e LogisticsFeeRateInfo) ModelDesc() string {
	return "verbose_name=物流费率,entity_json=logistics_fee_rate,route=logistics/fee/rate"
}
