package models

import "matrix/core"

type IssueRelateInfo struct {
	Id         int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	IssueId    int64             `xorm:"biging notnull index 'issue_id'" json:"issue_id" smith:"verbose_name=缺陷id"`
	RelateId   int64             `xorm:"biging notnull 'relate_id'" json:"relate_id" smith:"verbose_name=关联id"`
	CreateTime core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e IssueRelateInfo) TableName() string {
	return TablePrefix + "issue_relate"
}

func (e IssueRelateInfo) ModelDesc() string {
	return "verbose_name=关联缺陷,entity_json=issue_relate,route=issue/relate"
}
