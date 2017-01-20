package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	PASS_RES_USER = 10001 // 用户管理
)

type PassResInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CateId     int64           `xorm:"bigint notnull index 'cate_id'" json:"cate_id" smith:"verbose_name=类别"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e PassResInfo) TableName() string {
	return TablePrefix + "pass_res"
}

func (e PassResInfo) ModelDesc() string {
	return "verbose_name=权限,entity_json=pass_res,route=pass/res"
}

func (e PassResInfo) InitData(db_session *xorm.Session) {

}
