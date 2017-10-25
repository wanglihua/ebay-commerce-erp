package kcgl_models

import (
	"matrix/core"
)

type ProductInfo struct {
	Id               int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Sku              string          `xorm:"nvarchar(200) notnull index unique 'sku'" json:"sku" smith:"verbose_name=SKU,min=2,max=100"`
	CnName           string          `xorm:"nvarchar(200) notnull index 'cn_name'" json:"cn_name" smith:"verbose_name=名称,min=2,max=100"`
	Name             string          `xorm:"nvarchar(400) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	CateId           core.NullInt    `xorm:"bigint index 'cate_id'" json:"cate_id" smith:"verbose_name=商品目录"`
	CostPrice        float64         `xorm:"float notnull 'cost_price'" json:"cost_price" smith:"verbose_name=成本价,min=0"` //其实就是采购价，是 RMB
	SupplierId       core.NullInt    `xorm:"bigint index 'supplier_id'" json:"supplier_id" smith:"verbose_name=供应商"`
	PkgMtrlId        core.NullInt    `xorm:"bigint index 'pkg_mtrl_id'" json:"pkg_mtrl_id" smith:"verbose_name=包材"`
	ImageUrl         core.NullString `xorm:"nvarchar(500) 'image_url'" json:"image_url" smith:"verbose_name=图片地址"` // 许多图片中的一个图片，用于商品列表中的图片显示，存储的是某个店铺的图片地址
	Desc             core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,max=500"`
	StatusId         int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	Length           float64         `xorm:"float notnull 'length'" json:"length" smith:"verbose_name=长度,min=0"`          //cm
	Width            float64         `xorm:"float notnull 'width'" json:"width" smith:"verbose_name=宽度,min=0"`            //cm
	Height           float64         `xorm:"float notnull 'height'" json:"height" smith:"verbose_name=高度,min=0"`          //cm
	Weight           float64         `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量,min=0"`          //kg
	BatchCode        int64           `xorm:"bigint null index 'batch_code'" json:"batch_code" smith:"verbose_name=导入批次号"` //导入的批次号
	CnDeclName       string          `xorm:"nvarchar(300) notnull index 'cn_decl_name'" json:"cn_decl_name"`              // 中国海关申报名
	AuDeclName       string          `xorm:"nvarchar(300) notnull index 'au_decl_name'" json:"au_decl_name"`              // 澳洲海关申报名
	CnHsCode         string          `xorm:"nvarchar(100) notnull index 'cn_hs_code'" json:"cn_hs_code" smith:"verbose_name=中国海关编码,min=1,max=10"`
	AuHsCode         string          `xorm:"nvarchar(100) notnull index 'au_hs_code'" json:"au_hs_code" smith:"verbose_name=澳洲海关编码,min=1,max=8"`
	NeedInvoice      int64           `xorm:"bigint notnull 'need_invoice'" json:"need_invoice"` // 是否开票
	BrandId          int64           `xorm:"bigint notnull 'brand_id' default(0)" json:"brand_id"`
	StyleId          int64           `xorm:"bigint notnull 'style_id' default(0)" json:"style_id"`
	Cert1Id          int64           `xorm:"bigint notnull 'cert1_id' default(0)" json:"cert1_id"`
	Cert2Id          int64           `xorm:"bigint notnull 'cert2_id' default(0)" json:"cert2_id"`
	Cert3Id          int64           `xorm:"bigint notnull 'cert3_id' default(0)" json:"cert3_id"`
	StockAlertQty    int64           `xorm:"bigint notnull 'stock_alert_qty' default(0)" json:"stock_alert_qty"`                                           // 警戒库存，计算得到值
	StockAlertDays   int64           `xorm:"bigint notnull 'stock_alert_days' default(0)" json:"stock_alert_days"`                                         // 库存警戒天数
	DailyMeanSales   float64         `xorm:"decimal(18,6) notnull 'daily_mean_sales' default(0)" json:"daily_mean_sales"`                                  // 日均销量
	PurchaseDays     int64           `xorm:"bigint notnull 'purchase_days' default(0)" json:"purchase_days"`                                               // 采购天数
	Volume           string          `xorm:"-" json:"volume"`                                                                                              // 体积， 立方米， 只是用于界面显示
	DimWeight        string          `xorm:"-" json:"dim_weight"`                                                                                          // 体积重， Kg， 只是用于界面显示
	AvgFee           string          `xorm:"-" json:"avg_fee"`                                                                                             // 平均费用
	PostageCalWeight float64         `xorm:"decimal(18,6) notnull 'postage_cal_weight' default(0)" json:"postage_cal_weight" smith:"verbose_name=客户端运费权重"` // 客户端运费计算权重
	SaleUserId       int64           `xorm:"bigint notnull index 'sale_user_id' default(0)" json:"sale_user_id" smith:"verbose_name=销售人员"`
	SaleUser         string          `xorm:"-" json:"sale_user"`
	ServiceUserId    int64           `xorm:"bigint notnull index 'service_user_id' default(0)" json:"service_user_id" smith:"verbose_name=客服人员"`
	ServiceUser      string          `xorm:"-" json:"service_user"`
	DevUserId        int64           `xorm:"bigint notnull index 'dev_user_id' default(0)" json:"dev_user_id" smith:"verbose_name=开发人员"`
	DevUser          string          `xorm:"-" json:"dev_user"`
	CreateTime       core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime       core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductInfo) TableName() string {
	return "erp_product"
}

func (e ProductInfo) ModelDesc() string {
	return "verbose_name=商品,entity_json=product,route=product"
}
