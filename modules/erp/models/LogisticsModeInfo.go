package models

import (
	"matrix/core"
)

type LogisticsModeInfo struct {
	Id        int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CompanyId int64           `xorm:"bigint notnull 'company_id'" json:"company_id" smith:"verbose_name=物流公司"`
	Name      string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=物流方式名称,min=2,max=100"`
	Sfkltjz   bool            `xorm:"bool notnull 'sfkltjz'" json:"sfkltjz" smith:"verbose_name=是否考虑体积重"`
	Tjzgs     float64         `xorm:"decimal(18,6) notnull 'tjzgs'" json:"tjzgs" smith:"verbose_name=体积重公式"` // 如果公式是:长X宽X高, 这里就取值 250
	Tjzxs     float64         `xorm:"decimal(18,6) notnull 'tjzxs'" json:"tjzxs" smith:"verbose_name=体积重系数"`
	Tjzysjzqx int64           `xorm:"bigint notnull 'tjzysjzqx'" json:"tjzysjzqx" smith:"verbose_name=体积重与实际重取向"` //取大(1) or 取小(2)
	Zlsfybzyq bool            `xorm:"bool notnull 'sfkltjz'" json:"sfkltjz" smith:"verbose_name=体积重与实际重是否有比值要求"`  //重量是否有比值要求
	Tjzsjzbz  float64         `xorm:"decimal(18,6) notnull 'tjzsjzbz'" json:"tjzsjzbz" smith:"verbose_name=体积重/实际重比值"`
	Zcbsfyxz  bool            `xorm:"bool notnull 'zcbsfyxz'" json:"zcbsfyxz" smith:"verbose_name=最长边是否有限制"`
	Zlsfyxz   bool            `xorm:"bool notnull 'zlsfyxz'" json:"zlsfyxz" smith:"verbose_name=重量是否有限制"`
	Zlbncg    float64         `xorm:"decimal(18,6) notnull 'zlbncg'" json:"zlbncg" smith:"verbose_name=重量不能超过（Kg)"`
	Sfyxcc    bool            `xorm:"bool notnull 'sfyxcc'" json:"sfyxcc" smith:"verbose_name=是否允许超长"`
	Ccjsyf    float64         `xorm:"decimal(18,6) notnull 'ccjsyf'" json:"ccjsyf" smith:"verbose_name=超长加收多少运费（澳币）"` //超长加收运费
	Sfyxcz    bool            `xorm:"bool notnull 'sfyxcz'" json:"sfyxcz" smith:"verbose_name=是否允许超重"`
	Czjsyf    float64         `xorm:"decimal(18,6) notnull 'czjsyf'" json:"czjsyf" smith:"verbose_name=超重加收多少运费（澳币）"` //超重加收运费
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsModeInfo) TableName() string {
	return TablePrefix + "logistics_mode"
}

func (e LogisticsModeInfo) ModelDesc() string {
	return "verbose_name=物流方式,entity_json=logistics_mode,route=logistics/mode"
}
