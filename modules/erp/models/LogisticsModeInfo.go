package models

import (
	"matrix/core"
)

type LogisticsModeInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CompanyId  int64           `xorm:"bigint notnull index unique(c_s_name) 'company_id'" json:"company_id" smith:"verbose_name=物流公司"`
	StockId    int64           `xorm:"bigint notnull index unique(c_s_name) 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	Name       string          `xorm:"nvarchar(200) notnull index unique(c_s_name) 'name'" json:"name" smith:"verbose_name=物流方式名称,min=2,max=100"`
	Sfkltjz    int             `xorm:"int notnull 'sfkltjz'" json:"sfkltjz" smith:"verbose_name=是否考虑体积重"`     // 0 否， 1 是
	Tjzgs      float64         `xorm:"decimal(18,6) notnull 'tjzgs'" json:"tjzgs" smith:"verbose_name=体积重公式"` // 如果公式是:长X宽X高 / 4000, 这里就取值 4000
	Tjzxs      float64         `xorm:"decimal(18,6) notnull 'tjzxs'" json:"tjzxs" smith:"verbose_name=体积重系数"`
	Tjzysjzqx  int             `xorm:"int notnull 'tjzysjzqx'" json:"tjzysjzqx" smith:"verbose_name=体积重与实际重取向"`      // 体积重与实际重取向 取大(1) or 取小(2)
	Zlsfybzyq  int             `xorm:"int notnull 'zlsfybzyq'" json:"zlsfybzyq" smith:"verbose_name=体积重与实际重是否有比值要求"` //重量是否有比值要求  0 否， 1 是
	Tjzsjzbz   float64         `xorm:"decimal(18,6) notnull 'tjzsjzbz'" json:"tjzsjzbz" smith:"verbose_name=体积重/实际重比值"`
	Zcbsfyxz   int             `xorm:"int notnull 'zcbsfyxz'" json:"zcbsfyxz" smith:"verbose_name=最长边是否有限制"` // 0 否， 1 是
	Zcbxzcd    float64         `xorm:"decimal(18,6) notnull 'zcbxzcd'" json:"zcbxzcd" smith:"verbose_name=最长边限制 长度（cm）"`
	Zlsfyxz    int             `xorm:"int notnull 'zlsfyxz'" json:"zlsfyxz" smith:"verbose_name=重量是否有限制"` // 重量是否有限制 0 否， 1 是
	Zlbncg     float64         `xorm:"decimal(18,6) notnull 'zlbncg'" json:"zlbncg" smith:"verbose_name=重量不能超过（Kg)"`
	Sfyxcc     int             `xorm:"int notnull 'sfyxcc'" json:"sfyxcc" smith:"verbose_name=是否允许超长"`                 // 是否允许超长 0 否， 1 是
	Ccsfjt     float64         `xorm:"decimal(18,6) notnull 'ccsfjt'" json:"ccsfjt" smith:"verbose_name=超长收费阶梯"`       //超长收费阶梯，cm
	Ccjsyf     float64         `xorm:"decimal(18,6) notnull 'ccjsyf'" json:"ccjsyf" smith:"verbose_name=超长加收多少运费（澳币）"` //超长加收运费
	Sfyxcz     int             `xorm:"int notnull 'sfyxcz'" json:"sfyxcz" smith:"verbose_name=是否允许超重"`                 // 是否允许超重 0 否， 1 是
	Czjsyf     float64         `xorm:"decimal(18,6) notnull 'czjsyf'" json:"czjsyf" smith:"verbose_name=超重加收多少运费（澳币）"` //超重加收运费
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	TypeId     int64           `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=物流方式类型"` // LogisticsModeTypeInfo 的主键
	StatusId   int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LogisticsModeInfo) TableName() string {
	return TablePrefix + "logistics_mode"
}

func (e LogisticsModeInfo) ModelDesc() string {
	return "verbose_name=物流方式,entity_json=logistics_mode,route=logistics/mode"
}
