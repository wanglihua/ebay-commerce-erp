package cggl_models

import (
	"matrix/core"
)

type PurchaseBillInfo struct {
	Id              int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code            string          `xorm:"nvarchar(100) index 'code'" json:"code" smith:"verbose_name=批次号"`
	Desc            core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	SupplierId      int64           `xorm:"bigint notnull 'supplier_id'" json:"supplier_id" smith:"verbose_name=供应商"`
	TotalAmount     float64         `xorm:"decimal(18,6) notnull 'total_amount'" json:"total_amount" smith:"verbose_name=总金额"`
	AdvanceAmount   float64         `xorm:"decimal(18,6) notnull 'advance_amount'" json:"advance_amount" smith:"verbose_name=预付款"`
	FinalAmount     float64         `xorm:"decimal(18,6) notnull 'final_amount'" json:"final_amount" smith:"verbose_name=尾款"`
	AdvanceExpenses float64         `xorm:"decimal(18,6) notnull 'advance_expenses'" json:"advance_expenses" smith:"verbose_name=预付运费"` //transportation expenses
	FinalExpenses   float64         `xorm:"decimal(18,6) notnull 'final_expenses'" json:"final_expenses" smith:"verbose_name=提货费"`      // 就是运费的尾款
	ArrivalDays     int64           `xorm:"bigint notnull 'arrival_days' default(0)" json:"arrival_days"`                               // 到货天数
	PayMethodId     int64           `xorm:"bigint notnull 'pay_method_id' default(0)" json:"pay_method_id"`                             // 付款方式
	PayDays         int64           `xorm:"bigint notnull 'pay_days' default(0)" json:"pay_days"`                                       // 账期
	StatusId        int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	ApplyUserId     int64           `xorm:"bigint notnull 'apply_user_id' default(0)" json:"apply_user_id" smith:"verbose_name=申请人"`
	ApplyTime       core.NullTime   `xorm:"datetime null 'apply_time'" json:"apply_time" smith:"verbose_name=申请时间"`
	ApproveUserId   int64           `xorm:"bigint notnull 'approve_user_id' default(0)" json:"approve_user_id" smith:"verbose_name=审批人"`
	ApproveTime     core.NullTime   `xorm:"datetime null 'approve_time'" json:"approve_time" smith:"verbose_name=审批时间"` // 也就是采购时间
	PayUserId       int64           `xorm:"bigint notnull 'pay_user_id' default(0)" json:"pay_user_id" smith:"verbose_name=付款人"`
	PayTime         core.NullTime   `xorm:"datetime null 'pay_time'" json:"pay_time" smith:"verbose_name=付款时间"`  // 也就是下单时间
	CheckinUserId   int64           `xorm:"bigint notnull 'checkin_user_id' default(0)" json:"checkin_user_id" smith:"verbose_name=入库人"`
	CheckinTime     core.NullTime   `xorm:"datetime null 'checkin_time'" json:"checkin_time" smith:"verbose_name=入库时间"`
	CreateTime      core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime      core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillInfo) TableName() string {
	return "erp_purchase_bill"
}

func (e PurchaseBillInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill,route=purchase/bill"
}

/*

总金额（通过明细总计得来）、预付款、尾款、付款方式（原来这个属性是放在供应商信息里面的，应该是放在采购单里面才合适）
运费分成两项：预付运费，提货费
付款后到货天数
账期，字段可以手工输入，如果是有账期的，下单起开始算起，账期到期前1周（时间长度可设置）前提醒，提醒可以在财务界面展示（那个批次那个供应商账期快到了）。

 */
