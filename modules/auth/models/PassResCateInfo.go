package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	PASS_RES_CATE_AUTH = 10100 // 权限管理
	PASS_RES_CATE_CGGL = 10200 // 采购管理
	PASS_RES_CATE_DBGL = 10300 // 调拨管理
	PASS_RES_CATE_DDBB = 10400 // 订单报表
	PASS_RES_CATE_DDLC = 10500 // 订单流程
	PASS_RES_CATE_DDPZ = 10600 // 订单配置
	PASS_RES_CATE_KCGL = 10700 // 库存管理
	PASS_RES_CATE_PPGL = 10800 // PayPal管理
	PASS_RES_CATE_SPBB = 10900 // 商品报表
	PASS_RES_CATE_SPPZ = 11000 // 商品配置
	PASS_RES_CATE_THGL = 11100 // 退货管理
	PASS_RES_CATE_XSBB = 11200 // 销售报表
	PASS_RES_CATE_XSGL = 11300 // 销售管理
	PASS_RES_CATE_XXGL = 11400 // 消息管理
	PASS_RES_CATE_XTPZ = 11500 // 系统配置
)

type PassResCateInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e PassResCateInfo) TableName() string {
	return TablePrefix + "pass_res_cate"
}

func (e PassResCateInfo) ModelDesc() string {
	return "verbose_name=权限类别,entity_json=pass_res_cate,route=pass/res/cate"
}

func (e PassResCateInfo) InitData(db_session *xorm.Session) {

	pass_res_cate := new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_AUTH
	pass_res_cate.Name = "权限管理"
	pass_res_cate.Desc = core.NewNullString("权限管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_CGGL
	pass_res_cate.Name = "采购管理"
	pass_res_cate.Desc = core.NewNullString("采购管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_DBGL
	pass_res_cate.Name = "调拨管理"
	pass_res_cate.Desc = core.NewNullString("调拨管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_DDBB
	pass_res_cate.Name = "订单报表"
	pass_res_cate.Desc = core.NewNullString("订单报表", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_DDLC
	pass_res_cate.Name = "订单流程"
	pass_res_cate.Desc = core.NewNullString("订单流程", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_DDPZ
	pass_res_cate.Name = "订单配置"
	pass_res_cate.Desc = core.NewNullString("订单配置", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_KCGL
	pass_res_cate.Name = "库存管理"
	pass_res_cate.Desc = core.NewNullString("库存管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_PPGL
	pass_res_cate.Name = "PayPal管理"
	pass_res_cate.Desc = core.NewNullString("PayPal管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_SPBB
	pass_res_cate.Name = "商品报表"
	pass_res_cate.Desc = core.NewNullString("商品报表", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_SPPZ
	pass_res_cate.Name = "商品配置"
	pass_res_cate.Desc = core.NewNullString("商品配置", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_THGL
	pass_res_cate.Name = "退货管理"
	pass_res_cate.Desc = core.NewNullString("退货管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_XSBB
	pass_res_cate.Name = "销售报表"
	pass_res_cate.Desc = core.NewNullString("销售报表", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_XSGL
	pass_res_cate.Name = "销售管理"
	pass_res_cate.Desc = core.NewNullString("销售管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_XXGL
	pass_res_cate.Name = "消息管理"
	pass_res_cate.Desc = core.NewNullString("消息管理", true)
	db_session.Insert(pass_res_cate)

	pass_res_cate = new(PassResCateInfo)
	pass_res_cate.Id = PASS_RES_CATE_XTPZ
	pass_res_cate.Name = "系统配置"
	pass_res_cate.Desc = core.NewNullString("系统配置", true)
	db_session.Insert(pass_res_cate)
}
