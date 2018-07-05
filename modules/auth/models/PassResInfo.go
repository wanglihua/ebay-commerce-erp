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
		pass.Id = pass_res.DDLC_ORDER_HUIFU
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单恢复"
		pass.Desc = core.MakeNullString("订单恢复", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_FHSC
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "发货上传"
		pass.Desc = core.MakeNullString("发货上传", true)
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

		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_LOGISTICS_DIFF
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "最佳仓次佳仓偏差"
		pass.Desc = core.MakeNullString("最佳仓次佳仓偏差", true)
		db_session.Insert(pass)

		// INSERT INTO auth_pass_res (id, cate_id, name, `desc`, create_time, update_time) VALUES (131600, 10500, '订单分割', '订单分割', now(), now());
		pass = new(PassResInfo)
		pass.Id = pass_res.DDLC_ORDER_SPLIT
		pass.CateId = pass_res_cate.DDLC
		pass.Name = "订单拆分"
		pass.Desc = core.MakeNullString("订单拆分", true)
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

	// 订单报表
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.DDBB_LOGISTICS_REPORT
		pass.CateId = pass_res_cate.DDBB
		pass.Name = "物流报表"
		pass.Desc = core.MakeNullString("物流报表", true)
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
		pass.Id = pass_res.KCGL_PRODUCT_DETAIL_SAVE
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品详细保存"
		pass.Desc = core.MakeNullString("商品详细保存", true)
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
		pass.Id = pass_res.KCGL_INVENTORY_EXPORT_QTY
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 导出库存数量"
		pass.Desc = core.MakeNullString("商品库存 导出库存数量", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_IMPORT_QTY
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 导入库存数量"
		pass.Desc = core.MakeNullString("商品库存 导入库存数量", true)
		db_session.Insert(pass)

		/*
		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_IMPORT_QTY_OVER
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 导入覆盖库存数量"
		pass.Desc = core.MakeNullString("商品库存 导入覆盖库存数量", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_IMPORT_QTY_DIFF
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 导入修正库存数量"
		pass.Desc = core.MakeNullString("商品库存 导入修正库存数量", true)
		db_session.Insert(pass)
		*/

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_LOC
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 商品库位"
		pass.Desc = core.MakeNullString("商品库存 商品库位", true)
		db_session.Insert(pass)

		/*
		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_QTY
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 库存数量"
		pass.Desc = core.MakeNullString("商品库存 库存数量", true)
		db_session.Insert(pass)
		*/

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_INVENTORY_SKU_PRICE
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "商品库存 库存成本"
		pass.Desc = core.MakeNullString("商品库存 库存成本", true)
		db_session.Insert(pass)

		/*
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
		*/

		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_STOCKBILL_TRANSFER
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "分仓调拨"
		pass.Desc = core.MakeNullString("分仓调拨", true)
		db_session.Insert(pass)

		/*
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
		*/

		// KCGL_SPECIAL_SKU
		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_SPECIAL_SKU
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "特殊SKU"
		pass.Desc = core.MakeNullString("特殊SKU", true)
		db_session.Insert(pass)

		// KCGL_MEL_TEMP_SKU
		pass = new(PassResInfo)
		pass.Id = pass_res.KCGL_MEL_TEMP_SKU
		pass.CateId = pass_res_cate.KCGL
		pass.Name = "墨尔本临时SKU"
		pass.Desc = core.MakeNullString("墨尔本临时SKU", true)
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
		pass.Id = pass_res.SPPZ_BRAND
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品品牌"
		pass.Desc = core.MakeNullString("商品品牌", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_CERT
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品认证"
		pass.Desc = core.MakeNullString("商品认证", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_STYLE
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品款式"
		pass.Desc = core.MakeNullString("商品款式", true)
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
		pass.Id = pass_res.SPPZ_SUPPLIER_PAY_METHOD
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "供应商付款方式"
		pass.Desc = core.MakeNullString("供应商付款方式", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_SUPPLIER_DELETE
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "供应商删除"
		pass.Desc = core.MakeNullString("供应商删除", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_SUPPLIER_EXPORT
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "供应商导出"
		pass.Desc = core.MakeNullString("供应商导出", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_STATUS
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品状态"
		pass.Desc = core.MakeNullString("商品状态", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_USER
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品负责人"
		pass.Desc = core.MakeNullString("商品负责人", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_PRODUCT_SUPPLIER
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "商品供应商"
		pass.Desc = core.MakeNullString("商品供应商", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPPZ_SEA_CONTAINER_COST
		pass.CateId = pass_res_cate.SPPZ
		pass.Name = "海运集装箱成本"
		pass.Desc = core.MakeNullString("海运集装箱成本", true)
		db_session.Insert(pass)
	}

	// 商品质量反馈
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.SPZL_PRODUCT_FEEDBACK_CATE
		pass.CateId = pass_res_cate.SPZL
		pass.Name = "商品质量反馈类别"
		pass.Desc = core.MakeNullString("商品质量反馈类别", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPZL_PRODUCT_FEEDBACK
		pass.CateId = pass_res_cate.SPZL
		pass.Name = "商品质量反馈"
		pass.Desc = core.MakeNullString("商品质量反馈", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPZL_PRODUCT_FEEDBACK_DELETE
		pass.CateId = pass_res_cate.SPZL
		pass.Name = "商品质量反馈删除"
		pass.Desc = core.MakeNullString("商品质量反馈删除", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.SPZL_PRODUCT_FEEDBACK_QUERY
		pass.CateId = pass_res_cate.SPZL
		pass.Name = "商品质量反馈查询导出"
		pass.Desc = core.MakeNullString("商品质量反馈查询导出", true)
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

		pass = new(PassResInfo)
		pass.Id = pass_res.SPBB_PRODUCT_SALE_RANK
		pass.CateId = pass_res_cate.SPBB
		pass.Name = "商品销售排名"
		pass.Desc = core.MakeNullString("商品销售排名", true)
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
		pass.Name = "用户权限"
		pass.Desc = core.MakeNullString("用户权限", true)
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
		pass.Name = "授权设置"
		pass.Desc = core.MakeNullString("授权设置", true)
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

	// 系统管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XTGL_NOTICE
		pass.CateId = pass_res_cate.XTGL
		pass.Name = "通知公告"
		pass.Desc = core.MakeNullString("通知公告", true)
		db_session.Insert(pass)
	}

	// 采购管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.CGGL_PRODUCT_SETTING
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "商品采购设置"
		pass.Desc = core.MakeNullString("商品采购设置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_INTEL_REPL
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "智能补货"
		pass.Desc = core.MakeNullString("智能补货", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_BILL
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "采购管理"
		pass.Desc = core.MakeNullString("采购管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_FINANCE
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "财务管理"
		pass.Desc = core.MakeNullString("财务管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_CHECKIN
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "入库管理"
		pass.Desc = core.MakeNullString("入库管理", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_ABNORMAL
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "异常采购单"
		pass.Desc = core.MakeNullString("异常采购单", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_FLOW
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "采购流水"
		pass.Desc = core.MakeNullString("采购流水", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.CGGL_PURCHASE_APPROVE
		pass.CateId = pass_res_cate.CGGL
		pass.Name = "采购审批"
		pass.Desc = core.MakeNullString("采购审批", true)
		db_session.Insert(pass)
	}

	// 销售管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XSGL_PRODUCT_SALE_SETTING
		pass.CateId = pass_res_cate.XSGL
		pass.Name = "商品销售设置"
		pass.Desc = core.MakeNullString("商品销售设置", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XSGL_PRODUCT_PRICE_ESTIMATE
		pass.CateId = pass_res_cate.XSGL
		pass.Name = "商品售价预估"
		pass.Desc = core.MakeNullString("商品售价预估", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XSGL_PRODUCT_LOGISTICS
		pass.CateId = pass_res_cate.XSGL
		pass.Name = "商品运费估算"
		pass.Desc = core.MakeNullString("商品运费估算", true)
		db_session.Insert(pass)
	}

	// 销售报表
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XSBB_SALE_INCOME_REPORT_DAY
		pass.CateId = pass_res_cate.XSBB
		pass.Name = "日收入报表"
		pass.Desc = core.MakeNullString("日收入报表", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XSBB_SALE_AMOUNT_REPORT_DAY
		pass.CateId = pass_res_cate.XSBB
		pass.Name = "日销售额报表"
		pass.Desc = core.MakeNullString("日销售额报表", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XSBB_SALE_AMOUNT_REPORT_MONTH
		pass.CateId = pass_res_cate.XSBB
		pass.Name = "月销售额报表"
		pass.Desc = core.MakeNullString("月销售额报表", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT
		pass.CateId = pass_res_cate.XSBB
		pass.Name = "热销商品排名"
		pass.Desc = core.MakeNullString("热销商品排名", true)
		db_session.Insert(pass)
	}

	// 消息管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.XXGL_MESSAGE_ALL
		pass.CateId = pass_res_cate.XXGL
		pass.Name = "所有消息"
		pass.Desc = core.MakeNullString("所有消息", true)
		db_session.Insert(pass)
	}

	// PayPal管理
	{
		pass := new(PassResInfo)
		pass.Id = pass_res.PPGL_PAYPAL_PAYMENT
		pass.CateId = pass_res_cate.PPGL
		pass.Name = "PayPal付款"
		pass.Desc = core.MakeNullString("PayPal付款", true)
		db_session.Insert(pass)

		// INSERT INTO auth_pass_res (id, cate_id, name, `desc`, create_time, update_time) VALUES (131700, 10500, '订单退款', '订单退款', now(), now());
		pass = new(PassResInfo)
		pass.Id = pass_res.PPGL_ORDER_REFUND_APPLY
		pass.CateId = pass_res_cate.PPGL
		pass.Name = "订单退款申请"
		pass.Desc = core.MakeNullString("订单退款申请", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.PPGL_ORDER_REFUND_APPROVE
		pass.CateId = pass_res_cate.PPGL
		pass.Name = "订单退款审核"
		pass.Desc = core.MakeNullString("订单退款审核", true)
		db_session.Insert(pass)

		pass = new(PassResInfo)
		pass.Id = pass_res.PPGL_REFUND_REPORT
		pass.CateId = pass_res_cate.PPGL
		pass.Name = "退款报表"
		pass.Desc = core.MakeNullString("退款报表", true)
		db_session.Insert(pass)
	}
}
