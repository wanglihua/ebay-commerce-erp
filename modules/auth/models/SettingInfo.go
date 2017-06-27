package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type SettingInfo struct {
	Id              int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	DefaultPassword string          `xorm:"nvarchar(100) notnull 'default_password'" json:"default_password" smith:"verbose_name=默认密码,min=6,max=20"`
	TimeOut         int64           `xorm:"bigint notnull 'time_out'" json:"time_out" smith:"verbose_name=超时时长,min=1,max=10000"`
	CreateTime      core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime      core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e SettingInfo) TableName() string {
	return TablePrefix + "setting"
}

func (e SettingInfo) ModelDesc() string {
	return "verbose_name=授权设置,entity_json=auth_setting,route=auth/setting"
}

func (e SettingInfo) InitData(db_session *xorm.Session) {
	var setting_list = make([]SettingInfo, 0)
	db_session.Find(&setting_list)

	if len(setting_list) == 0 {
		var auth_setting SettingInfo
		auth_setting.DefaultPassword = "111111" // 默认密码 6 个 1
		auth_setting.TimeOut = 10               // 超时时长，如果这段时间内部操作，就会自动退出登录
		db_session.Insert(&auth_setting)
	}
}
