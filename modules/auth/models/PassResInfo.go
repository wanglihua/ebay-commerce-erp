package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/pass_res_cate"
)

type PassResInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	CateId     int64           `xorm:"bigint notnull index 'cate_id'" json:"cate_id" smith:"verbose_name=类别"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
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
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "所有订单"
		pass.Desc = core.MakeNullString("所有订单", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_NEW
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "新增订单"
		pass.Desc = core.MakeNullString("新增订单", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YCJ
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已创建"
		pass.Desc = core.MakeNullString("已创建", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YZF
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已支付"
		pass.Desc = core.MakeNullString("已支付", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_HOLD
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "HOLD"
		pass.Desc = core.MakeNullString("HOLD", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_PROCESSED_HOLD
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已处理HOLD"
		pass.Desc = core.MakeNullString("已处理HOLD", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_DPH
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "待配货"
		pass.Desc = core.MakeNullString("待配货", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_PHZ
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "配货中"
		pass.Desc = core.MakeNullString("配货中", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YFH
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已发货"
		pass.Desc = core.MakeNullString("已发货", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YWC
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已完成"
		pass.Desc = core.MakeNullString("已完成", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_YQX
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "已取消"
		pass.Desc = core.MakeNullString("已取消", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_VIEW_DETAIL
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "查看订单详细"
		pass.Desc = core.MakeNullString("查看订单详细", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_IMPORT
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单导入"
		pass.Desc = core.MakeNullString("订单导入", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_LOGISTICS
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单物流"
		pass.Desc = core.MakeNullString("订单物流", true)
		db_session.Insert(pass)

		// DDLC_ORDER_LOGISTICS_FEE_DIFF = 131000 // 订单物流费用偏差

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单物流费用偏差"
		pass.Desc = core.MakeNullString("订单物流费用理论和实际偏差", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_QUXIAO
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单取消"
		pass.Desc = core.MakeNullString("订单取消", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_LOGISTICS_MATCH_DEBUG
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单物流匹配调试"
		pass.Desc = core.MakeNullString("订单物流匹配调试", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_LOGISTICS_MATCH
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单物流匹配"
		pass.Desc = core.MakeNullString("订单物流匹配", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_FORCE_PEIHUO
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单流程取消波次限制配货"
		pass.Desc = core.MakeNullString("订单流程取消波次限制配货", true)
		db_session.Insert(pass)
	}

	// 订单配置
	{

		pass := new(PassResInfo)
		pass.Id = pass_res.DDPZ_LOGISTICS_COMPANY
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "物流公司"
		pass.Desc = core.MakeNullString("物流公司", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_LOGISTICS_MODE
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "物流方式"
		pass.Desc = core.MakeNullString("物流方式", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_LOGISTICS_FEE_RATE
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "物流费率"
		pass.Desc = core.MakeNullString("物流费率", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_LOGISTICS_SETTING
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "物流配置"
		pass.Desc = core.MakeNullString("物流配置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ENVELOPE_FEE_RATE
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "信封费率"
		pass.Desc = core.MakeNullString("信封费率", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_PICK_FEE
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "拣货费用"
		pass.Desc = core.MakeNullString("拣货费用", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ORDER_HOLD_TYPE
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "订单HOLD类型"
		pass.Desc = core.MakeNullString("订单HOLD类型", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ORDER_LOG
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "订单修改日志"
		pass.Desc = core.MakeNullString("订单修改日志", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDPZ_ORDER_TAG
		pass.CateId = pass_res_cate.DDPZ
		pass.Name = "订单标签"
		pass.Desc = core.MakeNullString("订单标签", true)
		db_session.Insert(pass)
	}

	// 库存管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品列表"
		pass.Desc = core.MakeNullString("商品列表", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_BIND
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "绑定商品"
		pass.Desc = core.MakeNullString("绑定商品", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_IMPORT
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品导入"
		pass.Desc = core.MakeNullString("商品导入", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_EXPORT
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品导出"
		pass.Desc = core.MakeNullString("商品导出", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_LWHW
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品长宽高重量"
		pass.Desc = core.MakeNullString("商品长宽高重量", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存"
		pass.Desc = core.MakeNullString("商品库存", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_LOC
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 商品库位"
		pass.Desc = core.MakeNullString("商品库存 商品库位", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_QTY
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 库存数量"
		pass.Desc = core.MakeNullString("商品库存 库存数量", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_PRICE
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 库存成本"
		pass.Desc = core.MakeNullString("商品库存 库存成本", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_CHECK
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "库存盘点"
		pass.Desc = core.MakeNullString("库存盘点", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_CHECK_APPROVE
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "盘点审批"
		pass.Desc = core.MakeNullString("盘点审批", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_TRANSFER
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "分仓调拨"
		pass.Desc = core.MakeNullString("分仓调拨", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_PRODUCT_QUALITY_FEEDBACK
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品质量反馈"
		pass.Desc = core.MakeNullString("商品质量反馈", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_CHECKIN
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "手工入库"
		pass.Desc = core.MakeNullString("手工入库", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_CHECKOUT
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "手工出库"
		pass.Desc = core.MakeNullString("手工出库", true)
		db_session.Insert(pass)

		// KCGL_SPECIAL_SKU
		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_SPECIAL_SKU
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "特殊SKU"
		pass.Desc = core.MakeNullString("特殊SKU", true)
		db_session.Insert(pass)
	}

	// 商品配置
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_BATCH
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品批量修改"
		pass.Desc = core.MakeNullString("商品批量修改", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PACKAGE_MATERIAL
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "包材类型"
		pass.Desc = core.MakeNullString("包材类型", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_STOCK
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "仓库"
		pass.Desc = core.MakeNullString("仓库", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_STOCK_PICK_CREW
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "仓库拣货人员"
		pass.Desc = core.MakeNullString("仓库拣货人员", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_STORAGE_LOC
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "仓位"
		pass.Desc = core.MakeNullString("仓位", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_CATE
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品目录"
		pass.Desc = core.MakeNullString("商品目录", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PSL_BATCH
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "批量修改仓位"
		pass.Desc = core.MakeNullString("批量修改仓位", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_SUPPLIER
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "供应商"
		pass.Desc = core.MakeNullString("供应商", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_STATUS
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品状态"
		pass.Desc = core.MakeNullString("商品状态", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_SEA_CONTAINER_COST
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "海运集装箱成本"
		pass.Desc = core.MakeNullString("海运集装箱成本", true)
		db_session.Insert(pass)
	}

	// 商品报表
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.SPBB_PRODUCT_SALE_ANALYSIS_SETTING
		pass.CateId = pass_res_cate.SPBB
		pass.Name = "商品销售分析设置"
		pass.Desc = core.MakeNullString("商品销售分析设置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPBB_PRODUCT_SALE_ANALYSIS
		pass.CateId = pass_res_cate.SPBB
		pass.Name = "商品销售分析"
		pass.Desc = core.MakeNullString("商品销售分析", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPBB_PRODUCT_SALE_DAY_ANALYSIS
		pass.CateId = pass_res_cate.SPBB
		pass.Name = "商品日销售分析"
		pass.Desc = core.MakeNullString("商品日销售分析", true)
		db_session.Insert(pass)
	}

	// 海关报关
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.HGBG_PRODUCT_HS_CODE
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "商品海关编码"
		pass.Desc = core.MakeNullString("商品海关编码", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_HS_CODE_CN_REPORT_ELEM
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "中国海关编码"
		pass.Desc = core.MakeNullString("中国海关编码", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_HS_CODE_AU_DUTY
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "澳洲海关编码"
		pass.Desc = core.MakeNullString("澳洲海关编码", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_CONTAINER_TRANSFER
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "集装箱调拨"
		pass.Desc = core.MakeNullString("集装箱调拨", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_CONTAINER_TRANSFER_TO_ZT_STOCK
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "集装箱调拨 调拨至在途仓"
		pass.Desc = core.MakeNullString("集装箱调拨 调拨至在途仓", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_CONTAINER_TRANSFER_TO_AU_STOCK
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "集装箱调拨 调拨至澳洲仓"
		pass.Desc = core.MakeNullString("集装箱调拨 调拨至澳洲仓", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.HGBG_PRODUCT_INVOICE
		pass.CateId = pass_res_cate.HGBG
		pass.Name = "商品是否开票"
		pass.Desc = core.MakeNullString("商品是否开票", true)
		db_session.Insert(pass)
	}

	// 权限管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.AUTH_USER
		pass.CateId = pass_res_cate.AUTH
		pass.Name = "用户管理"
		pass.Desc = core.MakeNullString("用户管理", true)
		db_session.Insert(pass)

		// 这个没有意义，只有超级管理员才能用的功能
		/*
		pass = new(PassResInfo)
		pass.Id = pass_res.AUTH_ADMIN
		pass.CateId = pass_res_cate.AUTH
		pass.Name = "管理员组"
		pass.Desc = core.MakeNullString("管理员组", true)
		db_session.Insert(pass)
		*/

		pass = new(PassResInfo)
		pass.Id = pass_res.AUTH_PASS
		pass.CateId = pass_res_cate.AUTH
		pass.Name = "授权管理"
		pass.Desc = core.MakeNullString("授权管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.AUTH_GROUP
		pass.CateId = pass_res_cate.AUTH
		pass.Name = "群组管理"
		pass.Desc = core.MakeNullString("群组管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.AUTH_SETTING
		pass.CateId = pass_res_cate.AUTH
		pass.Name = "授权配置"
		pass.Desc = core.MakeNullString("授权配置", true)
		db_session.Insert(pass)
	}

	// 系统配置
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XTPZ_ECOMMERCE_PLATFORM
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "电商平台"
		pass.Desc = core.MakeNullString("电商平台", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_SHOP_ACCOUNT
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "店铺账号"
		pass.Desc = core.MakeNullString("店铺账号", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_PAYPAL_ACCOUNT
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "PayPal账号"
		pass.Desc = core.MakeNullString("PayPal账号", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_EXCHANGE_RATE
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "汇率设置"
		pass.Desc = core.MakeNullString("汇率设置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_LOGIN_LOG
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "系统日志"
		pass.Desc = core.MakeNullString("系统日志", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_SCHEDULE_TASK
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "定时任务"
		pass.Desc = core.MakeNullString("定时任务", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XTPZ_LITTLE_ICON
		pass.CateId = pass_res_cate.XTPZ
		pass.Name = "图标管理"
		pass.Desc = core.MakeNullString("图标管理", true)
		db_session.Insert(pass)
	}
}
