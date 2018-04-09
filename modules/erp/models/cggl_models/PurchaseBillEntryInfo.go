package cggl_models

import (
	"matrix/core"
)

type PurchaseBillEntryInfo struct {
	Id                      int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId                  int64           `xorm:"bigint notnull index 'bill_id'" json:"bill_id" smith:"verbose_name=批次"`
	ProductId               int64           `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty                     int64           `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=总数量"`                             // 所有仓库的采购数量，各个仓库加起来的总和
	MelQty                  int64           `xorm:"bigint notnull default(0) 'mel_qty'" json:"mel_qty" smith:"verbose_name=墨尔本数量"`        // 墨尔本仓库采购数量
	BneQty                  int64           `xorm:"bigint notnull default(0) 'bne_qty'" json:"bne_qty" smith:"verbose_name=布里斯班数量"`       // 墨尔本仓库采购数量
	CheckinQty              int64           `xorm:"bigint notnull 'checkin_qty' default(0)" json:"checkin_qty" smith:"verbose_name=入库数量"` // 最终入库的数量，或者分批采购时从界面上获取的数量s
	Price                   float64         `xorm:"decimal(18,6) notnull 'price'" json:"price" smith:"verbose_name=单价"`
	Amount                  float64         `xorm:"decimal(18,6) notnull 'amount'" json:"amount" smith:"verbose_name=金额"`
	Desc                    core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	StatusId                int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	AbnormalTypeId          int64           `xorm:"bigint notnull index 'abnormal_type_id' default(0)" json:"abnormal_type_id" smith:"verbose_name=异常类型"`
	AbnormalProcessResultId int64           `xorm:"bigint notnull index 'abnormal_process_result_id' default(0)" json:"abnormal_process_result_id" smith:"verbose_name=异常处理结果"`
	CheckinTime             core.NullTime   `xorm:"datetime null 'checkin_time'" json:"checkin_time" smith:"verbose_name=入库时间"`
	CheckinUserId           int64           `xorm:"bigint notnull 'checkin_user_id' default(0)" json:"checkin_user_id" smith:"verbose_name=入库人"`
	NeedInvoice             int64           `xorm:"bigint notnull 'need_invoice' default(0)" json:"need_invoice" smith:"verbose_name=是否开票"`
	TaxPoint                float64         `xorm:"decimal(18,6) notnull 'tax_point' default(0)" json:"tax_point" smith:"verbose_name=收取税点"`
	DrawbackPoint           float64         `xorm:"decimal(18,6) notnull 'drawback_point' default(0)" json:"drawback_point" smith:"verbose_name=退税点"`
	DrawbackAmount          float64         `xorm:"decimal(18,6) notnull 'drawback_amount' default(0)" json:"drawback_amount" smith:"verbose_name=退税额"` // 这是单个商品的退税额
	FittingCost             float64         `xorm:"decimal(18,6) notnull 'fitting_cost' default(0)" json:"fitting_cost" smith:"verbose_name=配件费用"`
	OtherCost               float64         `xorm:"decimal(18,6) notnull 'other_cost' default(0)" json:"other_cost" smith:"verbose_name=其它费用"`
	InventoryCost           float64         `xorm:"decimal(18,6) notnull 'inventory_cost' default(0)" json:"inventory_cost" smith:"verbose_name=中国仓成本"`
	PayAmount               float64         `xorm:"decimal(18,6) null 'pay_amount' default(0)" json:"pay_amount" smith:"verbose_name=付款金额"`                     // 付款金额
	PayTime                 core.NullTime   `xorm:"datetime null 'pay_time'" json:"pay_time" smith:"verbose_name=付款时间"`                                         // 也就是下单时间
	FinalAmount             float64         `xorm:"decimal(18,6) null 'final_amount' default(0)" json:"final_amount" smith:"verbose_name=尾款"`                   // 尾款金额
	FinalAmountStatusId     core.NullInt    `xorm:"bigint null 'final_amount_status_id' default(5900)" json:"final_amount_status_id" smith:"verbose_name=尾款状态"` // 尾款状态： 未付 5900 或者 已付 6000
	FinalExpenses           float64         `xorm:"decimal(18,6) notnull 'final_expenses' default(0)" json:"final_expenses" smith:"verbose_name=提货费"`           // 就是运费的尾款，但分批采购时就是运费了，这个是总数量的运费，不是单个商品的运费
	FreightCharge           float64         `xorm:"decimal(18,6) notnull 'freight_charge' default(0)" json:"freight_charge" smith:"verbose_name=运费"`            // 搞不清，为什么要加个运费，前面都已经有提货费
	TypeId                  int64           `xorm:"bigint notnull index 'type_id' default(100)" json:"type_id" smith:"verbose_name=采购明细类型"`                     // 分批 或 正常，默认是正常明细类型
	ArrivalDays             int64           `xorm:"bigint null 'arrival_days'" json:"arrival_days"`                                                             // 到货天数
	DiffAmount              float64         `xorm:"decimal(18,6) null 'diff_amount'" json:"diff_amount" smith:"verbose_name=金额"`                                // 异常处理是的差异金额，正数表示退款
	CreateTime              core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime              core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillEntryInfo) TableName() string {
	return "erp_purchase_bill_entry"
}

func (e PurchaseBillEntryInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill_entry,route=purchase/bill/entry"
}

/*
已创建，已审批，采购中，已入库
 */

/*

是否开票 need_invoice
收取税点 tax_point
退税点   drawback_point
退税额   drawback_amount
配件费用  fitting_cost
其它费用  other_cost

*/

/*

是否开票，收税税点，退税税点
采购价、退税额、配件费用、其它费用，运费平摊 可以用于计算改商品的中国仓成本价，只有开票的才考虑退税额，收税点对计算中国仓成本价不起作用，只是一个记录。
运费平摊采用采购价按比例平摊
退税额=采购价/1.17*退税点
计算中国仓成本价的时候才有最近的平均6次计算，不足六次按实际次数平均。采购单入库后才有机会参与计算。可以入库时计算和定时任务计算。

 */
