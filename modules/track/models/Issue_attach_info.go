package models

import (
	"matrix/core"
)

type Issue_attach_info struct {
	Id          int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	IssueId     int64             `xorm:"biging notnull index 'issue_id'" json:"issue_id" smith:"verbose_name=缺陷id"`
	Type        int64             `xorm:"biging notnull index default(1) 'type'" json:"type" smith:"verbose_name=类型"` // 1 表示图片，2 表示普通附件,默认为 1
	DirName     string            `xorm:"nvarchar(300) notnull 'dir_name'" json:"dir_name" smith:"verbose_name=路径名"`  // 可以按月存放等
	FileName    string            `xorm:"nvarchar(300) notnull 'file_name'" json:"file_name" smith:"verbose_name=文件名"`
	Status      int64             `xorm:"biging notnull index default(1) 'status'" json:"status" smith:"verbose_name=状态"` //状态 如果是小于零就是已删除，大于等于零就是正常。默认为 1， 一般 -1 表示已删除
	Create_time core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int               `xorm:"version notnull 'version'" json:"version"`
}

func (e Issue_attach_info) TableName() string {
	return TablePrefix + "issue_attach"
}

func (e Issue_attach_info) ModelDesc() string {
	return "verbose_name=缺陷附件,entity_json=issue_attach,route=issue/attach"
}
