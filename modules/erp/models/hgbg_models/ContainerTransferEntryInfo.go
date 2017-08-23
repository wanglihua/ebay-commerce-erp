package hgbg_models

import (
	"matrix/core"
)

type ContainerTransferEntryInfo struct {
	Id           int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TransferId   int64           `xorm:"bigint notnull index 'transfer_id'" json:"transfer_id" smith:"verbose_name=集装箱调拨Id"`
	BatchCode    string          `xorm:"nvarchar(200) notnull 'batch_code'" json:"batch_code" smith:"verbose_name=批次号"`
	ProductId    int64           `xorm:"bigint notnull 'product_id'" json:"product_id" smith:"verbose_name=商品Id"`
	ProductSku   string          `xorm:"nvarchar(200) notnull 'product_sku'" json:"product_sku" smith:"verbose_name=SKU"`
	Qty          int64           `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	BoxQty       int64           `xorm:"bigint notnull 'box_qty'" json:"box_qty" smith:"verbose_name=箱数"`
	Price        float64         `xorm:"float64 notnull 'price'" json:"price"  smith:"verbose_name=价格"`
	AuHsId       int64           `xorm:"bigint notnull 'au_hs_id'" json:"au_hs_id"`
	AuHsCode     string          `xorm:"nvarchar(100) notnull 'au_hs_code'" json:"au_hs_code"`
	AuHsDutyRate float64         `xorm:"decimal(18,6) notnull 'au_hs_duty_rate'" json:"au_hs_duty_rate"`
	AuHsCnName   string          `xorm:"nvarchar(200) notnull 'au_hs_cn_name'" json:"au_hs_cn_name"`
	AuHsEnName   string          `xorm:"nvarchar(200) notnull 'au_hs_en_name'" json:"au_hs_en_name"` // 澳洲海关申报名
	CnHsId       int64           `xorm:"bigint notnull 'cn_hs_id'" json:"cn_hs_id"`
	CnHsCode     string          `xorm:"nvarchar(100) notnull 'cn_hs_code'" json:"cn_hs_code"`
	CnHsCnName   string          `xorm:"nvarchar(200) notnull 'cn_hs_cn_name'" json:"cn_hs_cn_name"` // 中国海关申报名
	CnHsEnName   string          `xorm:"nvarchar(200) notnull 'cn_hs_en_name'" json:"cn_hs_en_name"`
	NeedInvoice  int64           `xorm:"bigint notnull 'need_invoice' default(0)" json:"need_invoice"` // 是否开票
	CreateTime   core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`

	WeightSum float64 `xorm:"-"` // 用户记录导出时SKU合并的重量
	VolumeSum float64 `xorm:"-"` // 用户记录导出时SKU合并的体积
}

func (e ContainerTransferEntryInfo) TableName() string {
	return "erp_container_transfer_entry"
}

func (e ContainerTransferEntryInfo) ModelDesc() string {
	return "verbose_name=集装箱调拨明细,entity_json=transfer_entry,route=container/transfer/entry"
}
