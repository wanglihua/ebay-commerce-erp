package cggl_models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	ABNORMAL_TYPE_ZC    = 50  // 正常，一般默认是正常
	ABNORMAL_TYPE_HWDQ  = 100 // 货未到全
	ABNORMAL_TYPE_YTHCP = 200 // 有退货/次品
	ABNORMAL_TYPE_ZXPS  = 300 // 纸箱破损
)

type PurchaseAbnormalTypeInfo struct {
	Id          int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string    `xorm:"nvarchar(100) index 'name'" json:"name" smith:"verbose_name=名称"`
	Description string    `xorm:"nvarchar(500) 'description'" json:"description" smith:"verbose_name=备注,min=2,max=400"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseAbnormalTypeInfo) TableName() string {
	return "erp_purchase_abnormal_type"
}

func (e PurchaseAbnormalTypeInfo) ModelDesc() string {
	return "verbose_name=采购异常类型,entity_json=purchase_abnormal_type,route=purchase/abnormal/type"
}

func (e PurchaseAbnormalTypeInfo) InitData(db_session *xorm.Session) {
	abnormal_type := new(PurchaseAbnormalTypeInfo)
	abnormal_type.Id = ABNORMAL_TYPE_ZC
	abnormal_type.Name = "正常"
	abnormal_type.Description = "正常"
	db_session.Insert(abnormal_type)

	abnormal_type = new(PurchaseAbnormalTypeInfo)
	abnormal_type.Id = ABNORMAL_TYPE_HWDQ
	abnormal_type.Name = "货未到全"
	abnormal_type.Description = "货未到全"
	db_session.Insert(abnormal_type)

	abnormal_type = new(PurchaseAbnormalTypeInfo)
	abnormal_type.Id = ABNORMAL_TYPE_YTHCP
	abnormal_type.Name = "有退货/次品"
	abnormal_type.Description = "有退货/次品"
	db_session.Insert(abnormal_type)

	abnormal_type = new(PurchaseAbnormalTypeInfo)
	abnormal_type.Id = ABNORMAL_TYPE_ZXPS
	abnormal_type.Name = "纸箱破损"
	abnormal_type.Description = "纸箱破损"
	db_session.Insert(abnormal_type)
}
