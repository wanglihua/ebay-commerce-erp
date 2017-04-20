package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/auth/pass_res_cate"
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

	cate := new(PassResCateInfo)
	cate.Id = pass_res_cate.AUTH
	cate.Name = "权限管理"
	cate.Desc = core.MakeNullString("权限管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.CGGL
	cate.Name = "采购管理"
	cate.Desc = core.MakeNullString("采购管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.DBGL
	cate.Name = "调拨管理"
	cate.Desc = core.MakeNullString("调拨管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.DDBB
	cate.Name = "订单报表"
	cate.Desc = core.MakeNullString("订单报表", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.DDLC
	cate.Name = "订单流程"
	cate.Desc = core.MakeNullString("订单流程", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.DDPZ
	cate.Name = "订单配置"
	cate.Desc = core.MakeNullString("订单配置", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.KCGL
	cate.Name = "库存管理"
	cate.Desc = core.MakeNullString("库存管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.PPGL
	cate.Name = "PayPal管理"
	cate.Desc = core.MakeNullString("PayPal管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.SPBB
	cate.Name = "商品报表"
	cate.Desc = core.MakeNullString("商品报表", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.SPPZ
	cate.Name = "商品配置"
	cate.Desc = core.MakeNullString("商品配置", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.THGL
	cate.Name = "退货管理"
	cate.Desc = core.MakeNullString("退货管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.XSBB
	cate.Name = "销售报表"
	cate.Desc = core.MakeNullString("销售报表", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.XSGL
	cate.Name = "销售管理"
	cate.Desc = core.MakeNullString("销售管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.XXGL
	cate.Name = "消息管理"
	cate.Desc = core.MakeNullString("消息管理", true)
	db_session.Insert(cate)

	cate = new(PassResCateInfo)
	cate.Id = pass_res_cate.XTPZ
	cate.Name = "系统配置"
	cate.Desc = core.MakeNullString("系统配置", true)
	db_session.Insert(cate)
}
