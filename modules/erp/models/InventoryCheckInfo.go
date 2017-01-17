package models

import (
	"matrix/core"
)

type InventoryCheckInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code          string          `xorm:"nvarchar(100) index 'code'" json:"code" smith:"verbose_name=批次号"`
	StockId       int64           `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	CheckDate     core.Date       `xorm:"datetime notnull 'check_date'" json:"check_date" smith:"verbose_name=盘点日期"`
	CheckUserId   int64           `xorm:"bigint notnull 'check_user_id'" json:"check_user_id" smith:"verbose_name=盘点人"`
	ApproveUserId core.NullInt    `xorm:"bigint null 'approve_user_id'" json:"approve_user_id" smith:"verbose_name=审批人"`
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version       int             `xorm:"version notnull 'version'" json:"version"`
}

func (e InventoryCheckInfo) TableName() string {
	return TablePrefix + "inventory_check"
}

func (e InventoryCheckInfo) ModelDesc() string {
	return "verbose_name=库存盘点,entity_json=inventory_check,route=inventory/check"
}
