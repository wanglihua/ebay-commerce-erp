package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/auth/pass_res"
)

type PassResInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
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

	// 订单流程
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_ALL
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "所有订单"
		pass.Desc = core.MakeNullString("所有订单", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YZF
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "已支付"
		pass.Desc = core.MakeNullString("已支付", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_DPH
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "待配货"
		pass.Desc = core.MakeNullString("待配货", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_PHZ
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "配货中"
		pass.Desc = core.MakeNullString("配货中", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YFH
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "已发货"
		pass.Desc = core.MakeNullString("已发货", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YWC
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "已完成"
		pass.Desc = core.MakeNullString("已完成", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_KHB
		pass.CateId = PASS_RES_CATE_DDLC
		pass.Name = "可合并的订单"
		pass.Desc = core.MakeNullString("可合并的订单", true)
		db_session.Insert(pass)
	}

	// 订单配置
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.DDPZ_LOGISTICS_MODE
		pass.CateId = PASS_RES_CATE_DDPZ
		pass.Name = "货运方式"
		pass.Desc = core.MakeNullString("货运方式", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_PLATFORM_LOGISTICS_MODE
		pass.CateId = PASS_RES_CATE_DDPZ
		pass.Name = "电商平台物流方式"
		pass.Desc = core.MakeNullString("电商平台物流方式", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ORDER_CUSTOM_TYPE
		pass.CateId = PASS_RES_CATE_DDPZ
		pass.Name = "订单自定义类型"
		pass.Desc = core.MakeNullString("订单自定义类型", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ORDER_LOG
		pass.CateId = PASS_RES_CATE_DDPZ
		pass.Name = "订单修改日志"
		pass.Desc = core.MakeNullString("订单修改日志", true)
		db_session.Insert(pass)
	}

	// 库存管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "商品列表"
		pass.Desc = core.MakeNullString("商品列表", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_BIND
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "绑定商品"
		pass.Desc = core.MakeNullString("绑定商品", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_IMPORT
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "商品导入"
		pass.Desc = core.MakeNullString("商品导入", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_EXPORT
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "商品导出"
		pass.Desc = core.MakeNullString("商品导出", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "商品库存"
		pass.Desc = core.MakeNullString("商品库存", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_CHECK
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "库存盘点"
		pass.Desc = core.MakeNullString("库存盘点", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_CHECK_APPROVE
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "盘点审批"
		pass.Desc = core.MakeNullString("盘点审批", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_TRANSFER
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "分仓调拨"
		pass.Desc = core.MakeNullString("分仓调拨", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_QUALITY_FEEDBACK
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "商品质量反馈"
		pass.Desc = core.MakeNullString("商品质量反馈", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_CHECKIN
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "手工入库"
		pass.Desc = core.MakeNullString("手工入库", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_CHECKOUT
		pass.CateId = PASS_RES_CATE_KCGL
		pass.Name = "手工出库"
		pass.Desc = core.MakeNullString("手工出库", true)
		db_session.Insert(pass)
	}

	// 权限管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.AUTH_USER
		pass.CateId = PASS_RES_CATE_AUTH
		pass.Name = "用户管理"
		pass.Desc = core.MakeNullString("用户管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.AUTH_ADMIN
		pass.CateId = PASS_RES_CATE_AUTH
		pass.Name = "管理员组"
		pass.Desc = core.MakeNullString("管理员组", true)
		db_session.Insert(pass)
	}

	// 系统配置
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XTPZ_ECOMMERCE_PLATFORM
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "电商平台"
		pass.Desc = core.MakeNullString("电商平台", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_SHOP_ACCOUNT
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "店铺账号"
		pass.Desc = core.MakeNullString("店铺账号", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_PAYPAL_ACCOUNT
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "PayPal账号"
		pass.Desc = core.MakeNullString("PayPal账号", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_EXCHANGE_RATE
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "汇率设置"
		pass.Desc = core.MakeNullString("汇率设置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_LOGIN_LOG
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "系统日志"
		pass.Desc = core.MakeNullString("系统日志", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_SCHEDULE_TASK
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "定时任务"
		pass.Desc = core.MakeNullString("定时任务", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_LITTLE_ICON
		pass.CateId = PASS_RES_CATE_XTPZ
		pass.Name = "图标管理"
		pass.Desc = core.MakeNullString("图标管理", true)
		db_session.Insert(pass)
	}
}
