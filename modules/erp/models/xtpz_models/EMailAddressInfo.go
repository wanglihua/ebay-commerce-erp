package xtpz_models

import (
	"matrix/core"
)

type EmailAddressInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Email      string    `xorm:"nvarchar(300) notnull 'email'" json:"email" smith:"verbose_name=邮件地址"`
	Account    string    `xorm:"nvarchar(300) notnull 'account'" json:"account" smith:"verbose_name=账号"`
	Password   string    `xorm:"nvarchar(30) notnull 'password'" json:"password" smith:"verbose_name=密码"`
	SmtpServer string    `xorm:"nvarchar(300) notnull 'smtp_server'" json:"smtp_server" smith:"verbose_name=发件服务器"`
	SmtpPort   string    `xorm:"nvarchar(10) notnull 'smtp_port'" json:"smtp_port" smith:"verbose_name=发件服务器端口"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e EmailAddressInfo) TableName() string {
	return "erp_email_address"
}

func (e EmailAddressInfo) ModelDesc() string {
	return "verbose_name=电子邮件地址,entity_json=email_address,route=email/address"
}
