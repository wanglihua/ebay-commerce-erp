package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	//状态类型 begin

	STATUS_TYPE_STOCK                            = 100  // 仓库状态
	STATUS_TYPE_EXCHANGE_RATE                    = 200  // 仓库状态
	STATUS_TYPE_PRODUCT_CATE                     = 300  // 商品目录状态
	STATUS_TYPE_PRODUCT                          = 400  // 商品状态
	STATUS_TYPE_SHOP_ACCOUNT                     = 500  // 店铺账号状态
	STATUS_TYPE_PAYPAL_ACCOUNT                   = 600  // Paypal账号状态
	STATUS_TYPE_PACKAGE_MTRL                     = 700  // 包材 状态
	STATUS_TYPE_INVENTORY_CHECK                  = 800  // 库存盘点单状态
	STATUS_TYPE_PAYPAL_RECON                     = 900  // Paypal对账结果
	STATUS_TYPE_ORDER_CUSTOM_TYPE                = 1000 // 订单自定义类型
	STATUS_TYPE_LOGISTICS_NUMBER                 = 1100 // 不规则订单运单号类型
	STATUS_TYPE_LOGISTICS_NUMBER_BATCH           = 1200 // 不规则订单运单号类型
	STATUS_TYPE_PRODUCT_BIND                     = 1300 // 绑定商品状态
	STATUS_TYPE_STOCK_BILL_CHECKIN               = 1400 // 入库单状态
	STATUS_TYPE_STOCK_BILL_CHECKOUT              = 1500 // 出库单状态
	STATUS_TYPE_STOCK_BILL_TRANSFER              = 1600 // 调拨单状态
	STATUS_TYPE_ORDER                            = 1700 // 订单状态
	STATUS_TYPE_PURCHASE_BILL                    = 1800 // 采购单状态
	STATUS_TYPE_MESSAGE                          = 1900 // 消息状态
	STATUS_TYPE_LOGISTICS_COMPANY                = 2000 // 物流公司状态
	STATUS_TYPE_LOGISTICS_MODE                   = 2100 // 物流模式状态
	STATUS_TYPE_CONTAINER_TRANSFER               = 2200 // 集装箱调拨状态
	STATUS_TYPE_INTEL_REPL                       = 2300 // 智能补货
	STATUS_TYPE_PURCHASE_BILL_ENTRY_FINAL_AMOUNT = 2400 // 采购单明细尾款状态

	//状态类型 end

	//状态 id begin

	STOCK_STATUS_ENABLED  = 100 // 仓库状态 正常
	STOCK_STATUS_DISABLED = 200 // 仓库状态 停用

	EXCHANGE_RATE_STATUS_ENABLED  = 300 // 货币汇率状态 正常
	EXCHANGE_RATE_STATUS_DISABLED = 400 // 货币汇率状态 停用

	PRODUCT_CATE_STATUS_ENABLED  = 500 // 商品目录状态 正常
	PRODUCT_CATE_STATUS_DISABLED = 600 // 商品目录状态 停用

	PRODUCT_STATUS_ZDWH = 700  // 商品状态 重点维护
	PRODUCT_STATUS_DKF  = 800  // 商品状态 待开发
	PRODUCT_STATUS_ZC   = 900  // 商品状态 正常
	PRODUCT_STATUS_QC   = 1000 // 商品状态 清仓
	PRODUCT_STATUS_TC   = 1100 // 商品状态 停产
	PRODUCT_STATUS_ZDCJ = 1200 // 商品状态 自动创建
	PRODUCT_STATUS_ZTXC = 1300 // 商品状态 暂停销售

	SHOP_ACCOUNT_STATUS_ENABLED  = 1400 // 店铺账号状态 正常
	SHOP_ACCOUNT_STATUS_DISABLED = 1500 // 店铺账号状态 停用

	PAYPAL_ACCOUNT_STATUS_ENABLED  = 1600 // Paypal账号状态 正常
	PAYPAL_ACCOUNT_STATUS_DISABLED = 1700 // Paypal账号状态 停用

	PACKAGE_MTRL_STATUS_ENABLED  = 1800 // 包材状态 正常
	PACKAGE_MTRL_STATUS_DISABLED = 1900 // 包材状态 停用

	INVENTORY_CHECK_STATUS_YCJ = 2000 // 盘点单状态 已创建
	INVENTORY_CHECK_STATUS_WTG = 2050 // 盘点单状态 未通过
	INVENTORY_CHECK_STATUS_SPZ = 2100 // 盘点单状态 审核中
	INVENTORY_CHECK_STATUS_YWC = 2200 // 盘点单状态 已完成

	PAYPAL_RECON_STATUS_SUCCESS = 2300 // Paypal对账 成功
	PAYPAL_RECON_STATUS_FAILED  = 2400 // Paypal对账 失败

	ORDER_CUSTOM_TYPE_STATUS_ENABLED  = 2500 // 订单自定义类型 正常
	ORDER_CUSTOM_TYPE_STATUS_DISABLED = 2600 // 订单自定义类型 停用

	LOGISTICS_NUMBER_STATUS_USED     = 2700 // 不规则订单运单号 已使用
	LOGISTICS_NUMBER_STATUS_NOT_USED = 2800 // 不规则订单运单号 未使用

	LOGISTICS_NUMBER_BATCH_STATUS_ENABLED  = 2900 // 订单运单号 正常
	LOGISTICS_NUMBER_BATCH_STATUS_DISABLED = 3000 // 订单运单号 停用

	PRODUCT_BIND_STATUS_NORMAL  = 3100 // 绑定商品状态 正常
	PRODUCT_BIND_STATUS_STOPPED = 3200 // 绑定商品状态 停产

	STOCK_BILL_CHECKIN_STATUS_SAVE   = 3300 // 入库单状态 暂存
	STOCK_BILL_CHECKIN_STATUS_SUBMIT = 3400 // 入库单状态 提交

	STOCK_BILL_CHECKOUT_STATUS_SAVE   = 3500 // 出库单状态 暂存
	STOCK_BILL_CHECKOUT_STATUS_SUBMIT = 3600 // 出库单状态 提交

	STOCK_BILL_TRANSFER_STATUS_SAVE   = 3700 // 调拨单状态 暂存
	STOCK_BILL_TRANSFER_STATUS_SUBMIT = 3800 // 调拨单状态 提交

	ORDER_STATUS_YCJ            = 3850 // 订单状态 已创建
	ORDER_STATUS_YZF            = 3900 // 订单状态 已支付
	ORDER_STATUS_HOLD           = 3950 // 订单状态 Hold
	ORDER_STATUS_PROCESSED_HOLD = 3975 // 订单状态 已处理Hold
	ORDER_STATUS_DPH            = 4000 // 订单状态 待配货
	ORDER_STATUS_PHZ            = 4100 // 订单状态 配货中
	ORDER_STATUS_YFH            = 4200 // 订单状态 已发货
	ORDER_STATUS_YWC            = 4300 // 订单状态 已完成
	ORDER_STATUS_YQX            = 4350 // 订单状态 已取消

	PURCHASE_BILL_STATUS_YCJ = 4400 // 采购单状态 已创建
	PURCHASE_BILL_STATUS_YTJ = 4450 // 采购单状态 已提交
	PURCHASE_BILL_STATUS_YSP = 4500 // 采购单状态 已审批
	PURCHASE_BILL_STATUS_CGZ = 4600 // 采购单状态 采购中
	PURCHASE_BILL_STATUS_YRK = 4700 // 采购单状态 已入库
	PURCHASE_BILL_STATUS_YC  = 4710 // 采购单状态 异常
	PURCHASE_BILL_STATUS_YCL = 4720 // 采购单状态 异常已处理
	PURCHASE_BILL_STATUS_YWJ = 4730 // 已完结

	MESSAGE_STATUS_ANSWERED   = 4800 // 消息状态 已回复
	MESSAGE_STATUS_UNANSWERED = 4900 // 消息状态 未回复

	LOGISTICS_COMPANY_ENABLED  = 5000 // 物流公司启用
	LOGISTICS_COMPANY_DISABLED = 5100 // 物流公司禁用

	LOGISTICS_MODE_ENABLED  = 5200 // 物流模式启用
	LOGISTICS_MODE_DISABLED = 5300 // 物流模式禁用

	CONTAINER_TRANSFER_STATUS_NEW    = 5400 // 集装箱调拨 已创建
	CONTAINER_TRANSFER_STATUS_DEPART = 5500 // 集装箱调拨 已出发
	CONTAINER_TRANSFER_STATUS_STOCK  = 5600 // 集装箱调拨 已到仓

	INTEL_REPL_STATUS_XSC = 5700 // 新生成
	INTEL_REPL_STATUS_YCN = 5800 // 已采纳

	PURCHASE_BILL_ENTRY_FINAL_AMOUNT_STATUS_WZF = 5900 // 采购单明细尾款 未支付
	PURCHASE_BILL_ENTRY_FINAL_AMOUNT_STATUS_YZF = 6000 // 采购单明细尾款 已支付

	//状态 id end
)

type StatusInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Type       int64           `xorm:"bigint notnull index 'type'" json:"type" smith:"verbose_name=类型"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e StatusInfo) TableName() string {
	return "erp_status"
}

func (e StatusInfo) ModelDesc() string {
	return "verbose_name=状态,entity_json=status,route=status"
}

func (e StatusInfo) InitData(db_session *xorm.Session) {
	//仓库状态
	{
		status := new(StatusInfo)
		status.Id = STOCK_STATUS_ENABLED
		status.Type = STATUS_TYPE_STOCK
		status.Name = "正常"
		status.Desc = core.MakeNullString("仓库状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_STATUS_DISABLED
		status.Type = STATUS_TYPE_STOCK
		status.Name = "停用"
		status.Desc = core.MakeNullString("仓库状态 停用", true)
		db_session.Insert(status)
	}

	//货币汇率状态
	{
		status := new(StatusInfo)
		status.Id = EXCHANGE_RATE_STATUS_ENABLED
		status.Type = STATUS_TYPE_EXCHANGE_RATE
		status.Name = "正常"
		status.Desc = core.MakeNullString("货币汇率状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = EXCHANGE_RATE_STATUS_DISABLED
		status.Type = STATUS_TYPE_EXCHANGE_RATE
		status.Name = "停用"
		status.Desc = core.MakeNullString("货币汇率状态 停用", true)
		db_session.Insert(status)
	}

	//商品目录状态
	{
		status := new(StatusInfo)
		status.Id = PRODUCT_CATE_STATUS_ENABLED
		status.Type = STATUS_TYPE_PRODUCT_CATE
		status.Name = "正常"
		status.Desc = core.MakeNullString("商品目录状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_CATE_STATUS_DISABLED
		status.Type = STATUS_TYPE_PRODUCT_CATE
		status.Name = "停用"
		status.Desc = core.MakeNullString("商品目录状态 停用", true)
		db_session.Insert(status)
	}

	//商品状态

	{
		status := new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZDWH
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "重点维护"
		status.Desc = core.MakeNullString("商品状态 重点维护", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_DKF
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "待开发"
		status.Desc = core.MakeNullString("商品状态 待开发", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "正常"
		status.Desc = core.MakeNullString("商品状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_QC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "清仓"
		status.Desc = core.MakeNullString("商品状态 清仓", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_TC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "停产"
		status.Desc = core.MakeNullString("商品状态 停产", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZDCJ
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "自动创建"
		status.Desc = core.MakeNullString("商品状态 自动创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZTXC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "暂停销售"
		status.Desc = core.MakeNullString("商品状态 暂停销售", true)
		db_session.Insert(status)
	}

	//店铺账号状态
	{
		status := new(StatusInfo)
		status.Id = SHOP_ACCOUNT_STATUS_ENABLED
		status.Type = STATUS_TYPE_SHOP_ACCOUNT
		status.Name = "正常"
		status.Desc = core.MakeNullString("店铺账号状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = SHOP_ACCOUNT_STATUS_DISABLED
		status.Type = STATUS_TYPE_SHOP_ACCOUNT
		status.Name = "停用"
		status.Desc = core.MakeNullString("店铺账号状态 停用", true)
		db_session.Insert(status)
	}

	//Paypal账号状态

	{
		status := new(StatusInfo)
		status.Id = PAYPAL_ACCOUNT_STATUS_ENABLED
		status.Type = STATUS_TYPE_PAYPAL_ACCOUNT
		status.Name = "正常"
		status.Desc = core.MakeNullString("Paypal账号状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PAYPAL_ACCOUNT_STATUS_DISABLED
		status.Type = STATUS_TYPE_PAYPAL_ACCOUNT
		status.Name = "停用"
		status.Desc = core.MakeNullString("Paypal账号状态 停用", true)
		db_session.Insert(status)
	}

	//包材状态
	{
		status := new(StatusInfo)
		status.Id = PACKAGE_MTRL_STATUS_ENABLED
		status.Type = STATUS_TYPE_PACKAGE_MTRL
		status.Name = "正常"
		status.Desc = core.MakeNullString("包材状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PACKAGE_MTRL_STATUS_DISABLED
		status.Type = STATUS_TYPE_PACKAGE_MTRL
		status.Name = "停用"
		status.Desc = core.MakeNullString("包材状态 停用", true)
		db_session.Insert(status)
	}

	// 盘点单状态
	{
		status := new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_YCJ
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "已创建"
		status.Desc = core.MakeNullString("盘点单状态 已创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_WTG
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "未通过"
		status.Desc = core.MakeNullString("盘点单状态 未通过", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_SPZ
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "审核中"
		status.Desc = core.MakeNullString("盘点单状态 审核中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_YWC
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "已完成"
		status.Desc = core.MakeNullString("盘点单状态 已完成", true)
		db_session.Insert(status)

		// Paypal对账结果
		status = new(StatusInfo)
		status.Id = PAYPAL_RECON_STATUS_SUCCESS
		status.Type = STATUS_TYPE_PAYPAL_RECON
		status.Name = "匹配成功"
		status.Desc = core.MakeNullString("Paypal对账结果 匹配成功", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PAYPAL_RECON_STATUS_FAILED
		status.Type = STATUS_TYPE_PAYPAL_RECON
		status.Name = "匹配失败"
		status.Desc = core.MakeNullString("Paypal对账结果 匹配失败", true)
		db_session.Insert(status)
	}

	// 订单自定义类型
	{
		status := new(StatusInfo)
		status.Id = ORDER_CUSTOM_TYPE_STATUS_ENABLED
		status.Type = STATUS_TYPE_ORDER_CUSTOM_TYPE
		status.Name = "正常"
		status.Desc = core.MakeNullString("订单自定义类型 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_CUSTOM_TYPE_STATUS_DISABLED
		status.Type = STATUS_TYPE_ORDER_CUSTOM_TYPE
		status.Name = "停用"
		status.Desc = core.MakeNullString("订单自定义类型 停用", true)
		db_session.Insert(status)

		//规则订单运单号状态
		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_BATCH_STATUS_ENABLED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER_BATCH
		status.Name = "正常"
		status.Desc = core.MakeNullString("订单运单号 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_BATCH_STATUS_DISABLED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER_BATCH
		status.Name = "停用"
		status.Desc = core.MakeNullString("订单运单号 停用", true)
		db_session.Insert(status)
	}

	//不规则订单运单号状态
	{
		status := new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_STATUS_NOT_USED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER
		status.Name = "未使用"
		status.Desc = core.MakeNullString("订单运单号 未使用", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_STATUS_USED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER
		status.Name = "已使用"
		status.Desc = core.MakeNullString("订单运单号 已使用", true)
		db_session.Insert(status)
	}

	// 绑定商品状态 正常 停产
	{
		status := new(StatusInfo)
		status.Id = PRODUCT_BIND_STATUS_NORMAL
		status.Type = STATUS_TYPE_PRODUCT_BIND
		status.Name = "正常"
		status.Desc = core.MakeNullString("绑定商品状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_BIND_STATUS_STOPPED
		status.Type = STATUS_TYPE_PRODUCT_BIND
		status.Name = "停产"
		status.Desc = core.MakeNullString("绑定商品状态 停产", true)
		db_session.Insert(status)
	}

	// 入库单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_CHECKIN_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKIN
		status.Name = "临时保存"
		status.Desc = core.MakeNullString("入库单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_CHECKIN_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKIN
		status.Name = "已提交入库"
		status.Desc = core.MakeNullString("入库单状态 已提交入库", true)
		db_session.Insert(status)
	}

	// 出库单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_CHECKOUT_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKOUT
		status.Name = "临时保存"
		status.Desc = core.MakeNullString("出库单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_CHECKOUT_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKOUT
		status.Name = "已提交出库"
		status.Desc = core.MakeNullString("出库单状态 已提交出库", true)
		db_session.Insert(status)
	}

	// 调拨单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_TRANSFER_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_TRANSFER
		status.Name = "临时保存"
		status.Desc = core.MakeNullString("调拨单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_TRANSFER_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_TRANSFER
		status.Name = "已提交调拨"
		status.Desc = core.MakeNullString("调拨单状态 已提交调拨", true)
		db_session.Insert(status)
	}

	// 订单状态 已创建 已支付 Hold 待配货 配货中 已发货 已完成
	{
		status := new(StatusInfo)
		status.Id = ORDER_STATUS_YCJ
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已创建"
		status.Desc = core.MakeNullString("订单状态 已创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YZF
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已支付"
		status.Desc = core.MakeNullString("订单状态 已支付", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_HOLD
		status.Type = STATUS_TYPE_ORDER
		status.Name = "Hold"
		status.Desc = core.MakeNullString("订单状态 Hold", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_PROCESSED_HOLD
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已处理Hold"
		status.Desc = core.MakeNullString("订单状态 已处理Hold", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_DPH
		status.Type = STATUS_TYPE_ORDER
		status.Name = "待配货"
		status.Desc = core.MakeNullString("订单状态 待配货", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_PHZ
		status.Type = STATUS_TYPE_ORDER
		status.Name = "配货中"
		status.Desc = core.MakeNullString("订单状态 配货中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YFH
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已发货"
		status.Desc = core.MakeNullString("订单状态 已发货", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YWC
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已完成"
		status.Desc = core.MakeNullString("订单状态 已完成", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YQX
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已取消"
		status.Desc = core.MakeNullString("订单状态 已取消", true)
		db_session.Insert(status)
	}

	// 采购单状态 已创建，已审批，采购中，已入库
	{
		status := new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YCJ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已创建"
		status.Desc = core.MakeNullString("采购单状态 已创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YTJ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已提交"
		status.Desc = core.MakeNullString("采购单状态 已提交", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YSP
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已审批"
		status.Desc = core.MakeNullString("采购单状态 已审批", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_CGZ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "采购中"
		status.Desc = core.MakeNullString("采购单状态 采购中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YRK
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已入库"
		status.Desc = core.MakeNullString("采购单状态 已入库", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YC
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "异常"
		status.Desc = core.MakeNullString("采购单状态 异常，只用于明细的状态", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YCL
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已处理"
		status.Desc = core.MakeNullString("采购单状态 异常已处理，只用于明细的状态", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YWJ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已完结"
		status.Desc = core.MakeNullString("采购单状态 已完结", true)
		db_session.Insert(status)
	}

	// 消息状态 已回复 未回复
	{
		status := new(StatusInfo)
		status.Id = MESSAGE_STATUS_ANSWERED
		status.Type = STATUS_TYPE_MESSAGE
		status.Name = "已回复"
		status.Desc = core.MakeNullString("消息状态 已回复", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = MESSAGE_STATUS_UNANSWERED
		status.Type = STATUS_TYPE_MESSAGE
		status.Name = "未回复"
		status.Desc = core.MakeNullString("消息状态 未回复", true)
		db_session.Insert(status)
	}

	// 物流公司状态 启用 停用
	{
		status := new(StatusInfo)
		status.Id = LOGISTICS_COMPANY_ENABLED
		status.Type = STATUS_TYPE_LOGISTICS_COMPANY
		status.Name = "正常"
		status.Desc = core.MakeNullString("物流公司状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_COMPANY_DISABLED
		status.Type = STATUS_TYPE_LOGISTICS_COMPANY
		status.Name = "停用"
		status.Desc = core.MakeNullString("物流公司状态 停用", true)
		db_session.Insert(status)
	}

	// 物流模式状态 启用 停用
	{
		status := new(StatusInfo)
		status.Id = LOGISTICS_MODE_ENABLED
		status.Type = STATUS_TYPE_LOGISTICS_MODE
		status.Name = "正常"
		status.Desc = core.MakeNullString("物流模式状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_MODE_DISABLED
		status.Type = STATUS_TYPE_LOGISTICS_MODE
		status.Name = "停用"
		status.Desc = core.MakeNullString("物流模式状态 停用", true)
		db_session.Insert(status)
	}

	// 集装箱调拨状态
	{
		status := new(StatusInfo)
		status.Id = CONTAINER_TRANSFER_STATUS_NEW
		status.Type = STATUS_TYPE_CONTAINER_TRANSFER
		status.Name = "已新建"
		status.Desc = core.MakeNullString("集装箱调拨状态 已新建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = CONTAINER_TRANSFER_STATUS_DEPART
		status.Type = STATUS_TYPE_CONTAINER_TRANSFER
		status.Name = "已出发"
		status.Desc = core.MakeNullString("集装箱调拨状态 已出发", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = CONTAINER_TRANSFER_STATUS_STOCK
		status.Type = STATUS_TYPE_CONTAINER_TRANSFER
		status.Name = "已到仓"
		status.Desc = core.MakeNullString("集装箱调拨状态 已到仓", true)
		db_session.Insert(status)
	}

	// 智能补货状态
	{
		status := new(StatusInfo)
		status.Id = INTEL_REPL_STATUS_XSC
		status.Type = STATUS_TYPE_INTEL_REPL
		status.Name = "新生成"
		status.Desc = core.MakeNullString("智能补货状态 新生成", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INTEL_REPL_STATUS_YCN
		status.Type = STATUS_TYPE_INTEL_REPL
		status.Name = "已采纳"
		status.Desc = core.MakeNullString("智能补货状态 已采纳", true)
		db_session.Insert(status)
	}

	// 采购单明细尾款状态
	{
		status := new(StatusInfo)
		status.Id = PURCHASE_BILL_ENTRY_FINAL_AMOUNT_STATUS_WZF
		status.Type = STATUS_TYPE_PURCHASE_BILL_ENTRY_FINAL_AMOUNT
		status.Name = "未支付"
		status.Desc = core.MakeNullString("采购单明细尾款状态 未支付", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_ENTRY_FINAL_AMOUNT_STATUS_YZF
		status.Type = STATUS_TYPE_PURCHASE_BILL_ENTRY_FINAL_AMOUNT
		status.Name = "已支付"
		status.Desc = core.MakeNullString("采购单明细尾款状态 已支付", true)
		db_session.Insert(status)
	}
}
