package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	// 权限管理
	PASS_RES_AUTH_USER  = 10100 // 用户管理
	PASS_RES_AUTH_ADMIN = 10200 // 管理员组

	// 采购管理

	// 调拨管理

	// 订单报表

	// 订单流程
	PASS_RES_DDLC_ORDER_ALL = 11100 // 所有订单
	PASS_RES_DDLC_ORDER_YZF = 11200 // 已支付
	PASS_RES_DDLC_ORDER_DPH = 11300 // 待配货
	PASS_RES_DDLC_ORDER_PHZ = 11400 // 配货中
	PASS_RES_DDLC_ORDER_YFH = 11500 // 已发货
	PASS_RES_DDLC_ORDER_YWC = 11600 // 已完成
	PASS_RES_DDLC_ORDER_KHB = 11700 // 可合并的订单

	// 订单配置
	PASS_RES_DDPZ_LOGISTICS_MODE          = 11800 // 货运方式
	PASS_RES_DDPZ_PLATFORM_LOGISTICS_MODE = 11900 // 电商平台物流方式
	PASS_RES_DDPZ_ORDER_CUSTOM_TYPE       = 12000 // 订单自定义类型
	PASS_RES_DDPZ_ORDER_LOG               = 12100 // 订单修改日志

	// 库存管理
	PASS_RES_KCGL_PRODUCT                  = 12200 // 商品列表
	PASS_RES_KCGL_PRODUCT_BIND             = 12300 // 绑定商品
	PASS_RES_KCGL_PRODUCT_IMPORT           = 12400 // 商品导入
	PASS_RES_KCGL_PRODUCT_EXPORT           = 12500 // 商品导出
	PASS_RES_KCGL_INVENTORY                = 12600 // 商品库存
	PASS_RES_KCGL_INVENTORY_CHECK          = 12700 // 库存盘点
	PASS_RES_KCGL_INVENTORY_CHECK_APPROVE  = 12700 // 盘点审批
	PASS_RES_KCGL_STOCKBILL_TRANSFER       = 12800 // 分仓调拨
	PASS_RES_KCGL_PRODUCT_QUALITY_FEEDBACK = 12900 // 商品质量反馈
	PASS_RES_KCGL_STOCKBILL_CHECKIN        = 13000 // 手工入库
	PASS_RES_KCGL_STOCKBILL_CHECKOUT       = 13100 // 手工出库

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

	// 权限管理
	{
		pass_res := new(PassResInfo)
		pass_res.Id = PASS_RES_AUTH_USER
		pass_res.CateId = PASS_RES_CATE_AUTH
		pass_res.Name = "用户管理"
		pass_res.Desc = core.MakeNullString("用户管理", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_AUTH_ADMIN
		pass_res.CateId = PASS_RES_CATE_AUTH
		pass_res.Name = "管理员组"
		pass_res.Desc = core.MakeNullString("管理员组", true)
		db_session.Insert(pass_res)
	}

	// 订单流程
	{
		pass_res := new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_ALL
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "所有订单"
		pass_res.Desc = core.MakeNullString("所有订单", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_YZF
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "已支付"
		pass_res.Desc = core.MakeNullString("已支付", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_DPH
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "待配货"
		pass_res.Desc = core.MakeNullString("待配货", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_PHZ
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "配货中"
		pass_res.Desc = core.MakeNullString("配货中", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_YFH
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "已发货"
		pass_res.Desc = core.MakeNullString("已发货", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_YWC
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "已完成"
		pass_res.Desc = core.MakeNullString("已完成", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDLC_ORDER_KHB
		pass_res.CateId = PASS_RES_CATE_DDLC
		pass_res.Name = "可合并的订单"
		pass_res.Desc = core.MakeNullString("可合并的订单", true)
		db_session.Insert(pass_res)
	}

	// 订单配置
	{
		pass_res := new(PassResInfo)
		pass_res.Id = PASS_RES_DDPZ_LOGISTICS_MODE
		pass_res.CateId = PASS_RES_CATE_DDPZ
		pass_res.Name = "货运方式"
		pass_res.Desc = core.MakeNullString("货运方式", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDPZ_PLATFORM_LOGISTICS_MODE
		pass_res.CateId = PASS_RES_CATE_DDPZ
		pass_res.Name = "电商平台物流方式"
		pass_res.Desc = core.MakeNullString("电商平台物流方式", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDPZ_ORDER_CUSTOM_TYPE
		pass_res.CateId = PASS_RES_CATE_DDPZ
		pass_res.Name = "订单自定义类型"
		pass_res.Desc = core.MakeNullString("订单自定义类型", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_DDPZ_ORDER_LOG
		pass_res.CateId = PASS_RES_CATE_DDPZ
		pass_res.Name = "订单修改日志"
		pass_res.Desc = core.MakeNullString("订单修改日志", true)
		db_session.Insert(pass_res)
	}

	// 库存管理
	{
		pass_res := new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_PRODUCT
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "商品列表"
		pass_res.Desc = core.MakeNullString("商品列表", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_PRODUCT_BIND
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "绑定商品"
		pass_res.Desc = core.MakeNullString("绑定商品", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_PRODUCT_IMPORT
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "商品导入"
		pass_res.Desc = core.MakeNullString("商品导入", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_PRODUCT_EXPORT
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "商品导出"
		pass_res.Desc = core.MakeNullString("商品导出", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_INVENTORY
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "商品库存"
		pass_res.Desc = core.MakeNullString("商品库存", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_INVENTORY_CHECK
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "库存盘点"
		pass_res.Desc = core.MakeNullString("库存盘点", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_INVENTORY_CHECK_APPROVE
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "盘点审批"
		pass_res.Desc = core.MakeNullString("盘点审批", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_STOCKBILL_TRANSFER
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "分仓调拨"
		pass_res.Desc = core.MakeNullString("分仓调拨", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_PRODUCT_QUALITY_FEEDBACK
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "商品质量反馈"
		pass_res.Desc = core.MakeNullString("商品质量反馈", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_STOCKBILL_CHECKIN
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "手工入库"
		pass_res.Desc = core.MakeNullString("手工入库", true)
		db_session.Insert(pass_res)

		pass_res = new(PassResInfo)
		pass_res.Id = PASS_RES_KCGL_STOCKBILL_CHECKOUT
		pass_res.CateId = PASS_RES_CATE_KCGL
		pass_res.Name = "手工出库"
		pass_res.Desc = core.MakeNullString("手工出库", true)
		db_session.Insert(pass_res)
	}
}
