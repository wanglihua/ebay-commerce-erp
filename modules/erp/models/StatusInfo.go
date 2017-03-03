package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	//状态类型 begin

	// 仓库状态
	STATUS_TYPE_STOCK = 100
	// 仓库状态
	STATUS_TYPE_EXCHANGE_RATE = 200
	// 商品目录状态
	STATUS_TYPE_PRODUCT_CATE = 300
	// 商品状态
	STATUS_TYPE_PRODUCT = 400
	// 店铺账号状态
	STATUS_TYPE_SHOP_ACCOUNT = 500
	// Paypal账号状态
	STATUS_TYPE_PAYPAL_ACCOUNT = 600
	// 包材 状态
	STATUS_TYPE_PACKAGE_MTRL = 700
	// 库存盘点单状态
	STATUS_TYPE_INVENTORY_CHECK = 800
	// Paypal对账结果
	STATUS_TYPE_PAYPAL_RECON = 900
	// 订单自定义类型
	STATUS_TYPE_ORDER_CUSTOM_TYPE = 1000
	// 不规则订单运单号类型
	STATUS_TYPE_LOGISTICS_NUMBER = 1100
	// 不规则订单运单号类型
	STATUS_TYPE_LOGISTICS_NUMBER_BATCH = 1200
	// 绑定商品状态
	STATUS_TYPE_PRODUCT_BIND = 1300
	// 入库单状态
	STATUS_TYPE_STOCK_BILL_CHECKIN = 1400
	// 出库单状态
	STATUS_TYPE_STOCK_BILL_CHECKOUT = 1500
	// 调拨单状态
	STATUS_TYPE_STOCK_BILL_TRANSFER = 1600
	// 订单状态
	STATUS_TYPE_ORDER = 1700
	// 采购单状态
	STATUS_TYPE_PURCHASE_BILL = 1800

	//状态类型 end

	//状态 id begin

	// 仓库状态 正常
	STOCK_STATUS_ENABLED = 100
	// 仓库状态 停用
	STOCK_STATUS_DISABLED = 200

	// 货币汇率状态 正常
	EXCHANGE_RATE_STATUS_ENABLED = 300
	// 货币汇率状态 停用
	EXCHANGE_RATE_STATUS_DISABLED = 400

	// 商品目录状态 正常
	PRODUCT_CATE_STATUS_ENABLED = 500
	// 商品目录状态 停用
	PRODUCT_CATE_STATUS_DISABLED = 600

	// 商品状态 重点维护
	PRODUCT_STATUS_ZDWH = 700
	// 商品状态 待开发
	PRODUCT_STATUS_DKF = 800
	// 商品状态 正常
	PRODUCT_STATUS_ZC = 900
	// 商品状态 清仓
	PRODUCT_STATUS_QC = 1000
	// 商品状态 停产
	PRODUCT_STATUS_TC = 1100
	// 商品状态 自动创建
	PRODUCT_STATUS_ZDCJ = 1200
	// 商品状态 暂停销售
	PRODUCT_STATUS_ZTXC = 1300

	// 店铺账号状态 正常
	SHOP_ACCOUNT_STATUS_ENABLED = 1400
	// 店铺账号状态 停用
	SHOP_ACCOUNT_STATUS_DISABLED = 1500

	// Paypal账号状态 正常
	PAYPAL_ACCOUNT_STATUS_ENABLED = 1600
	// Paypal账号状态 停用
	PAYPAL_ACCOUNT_STATUS_DISABLED = 1700

	// 包材状态 正常
	PACKAGE_MTRL_STATUS_ENABLED = 1800
	// 包材状态 停用
	PACKAGE_MTRL_STATUS_DISABLED = 1900

	// 盘点单状态 已创建
	INVENTORY_CHECK_STATUS_YCJ = 2000
	// 盘点单状态 未通过
	INVENTORY_CHECK_STATUS_WTG = 2050
	// 盘点单状态 审核中
	INVENTORY_CHECK_STATUS_SPZ = 2100
	// 盘点单状态 已完成
	INVENTORY_CHECK_STATUS_YWC = 2200

	// Paypal对账 成功
	PAYPAL_RECON_STATUS_SUCCESS = 2300
	// Paypal对账 失败
	PAYPAL_RECON_STATUS_FAILED = 2400

	// 订单自定义类型 正常
	ORDER_CUSTOM_TYPE_STATUS_ENABLED = 2500
	// 订单自定义类型 停用
	ORDER_CUSTOM_TYPE_STATUS_DISABLED = 2600

	// 不规则订单运单号 已使用
	LOGISTICS_NUMBER_STATUS_USED = 2700
	// 不规则订单运单号 未使用
	LOGISTICS_NUMBER_STATUS_NOT_USED = 2800

	// 订单运单号 正常
	LOGISTICS_NUMBER_BATCH_STATUS_ENABLED = 2900
	// 订单运单号 停用
	LOGISTICS_NUMBER_BATCH_STATUS_DISABLED = 3000

	// 绑定商品状态 正常
	PRODUCT_BIND_STATUS_NORMAL = 3100
	// 绑定商品状态 停产
	PRODUCT_BIND_STATUS_STOPPED = 3200

	// 入库单状态 暂存
	STOCK_BILL_CHECKIN_STATUS_SAVE = 3300
	// 入库单状态 提交
	STOCK_BILL_CHECKIN_STATUS_SUBMIT = 3400

	// 出库单状态 暂存
	STOCK_BILL_CHECKOUT_STATUS_SAVE = 3500
	// 出库单状态 提交
	STOCK_BILL_CHECKOUT_STATUS_SUBMIT = 3600

	// 调拨单状态 暂存
	STOCK_BILL_TRANSFER_STATUS_SAVE = 3700
	// 调拨单状态 提交
	STOCK_BILL_TRANSFER_STATUS_SUBMIT = 3800

	// 订单状态 已创建
	ORDER_STATUS_YCJ = 3850
	// 订单状态 已支付
	ORDER_STATUS_YZF = 3900
	// 订单状态 待配货
	ORDER_STATUS_DPH = 4000
	// 订单状态 配货中
	ORDER_STATUS_PHZ = 4100
	// 订单状态 已发货
	ORDER_STATUS_YFH = 4200
	// 订单状态 已完成
	ORDER_STATUS_YWC = 4300

	// 采购单状态 已创建
	PURCHASE_BILL_STATUS_YCJ = 4400
	// 采购单状态 已下单
	PURCHASE_BILL_STATUS_YXD = 4500
	// 采购单状态 采购中
	PURCHASE_BILL_STATUS_CGZ = 4600
	// 采购单状态 已入库
	PURCHASE_BILL_STATUS_YRK = 4700

	//状态 id end
)

type StatusInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Type       int64           `xorm:"bigint notnull index 'type'" json:"type" smith:"verbose_name=类型"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e StatusInfo) TableName() string {
	return TablePrefix + "status"
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
		status.Desc = core.NewNullString("仓库状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_STATUS_DISABLED
		status.Type = STATUS_TYPE_STOCK
		status.Name = "停用"
		status.Desc = core.NewNullString("仓库状态 停用", true)
		db_session.Insert(status)
	}

	//货币汇率状态
	{
		status := new(StatusInfo)
		status.Id = EXCHANGE_RATE_STATUS_ENABLED
		status.Type = STATUS_TYPE_EXCHANGE_RATE
		status.Name = "正常"
		status.Desc = core.NewNullString("货币汇率状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = EXCHANGE_RATE_STATUS_DISABLED
		status.Type = STATUS_TYPE_EXCHANGE_RATE
		status.Name = "停用"
		status.Desc = core.NewNullString("货币汇率状态 停用", true)
		db_session.Insert(status)
	}

	//商品目录状态
	{
		status := new(StatusInfo)
		status.Id = PRODUCT_CATE_STATUS_ENABLED
		status.Type = STATUS_TYPE_PRODUCT_CATE
		status.Name = "正常"
		status.Desc = core.NewNullString("商品目录状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_CATE_STATUS_DISABLED
		status.Type = STATUS_TYPE_PRODUCT_CATE
		status.Name = "停用"
		status.Desc = core.NewNullString("商品目录状态 停用", true)
		db_session.Insert(status)
	}

	//商品状态

	{
		status := new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZDWH
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "重点维护"
		status.Desc = core.NewNullString("商品状态 重点维护", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_DKF
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "待开发"
		status.Desc = core.NewNullString("商品状态 待开发", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "正常"
		status.Desc = core.NewNullString("商品状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_QC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "清仓"
		status.Desc = core.NewNullString("商品状态 清仓", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_TC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "停产"
		status.Desc = core.NewNullString("商品状态 停产", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZDCJ
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "自动创建"
		status.Desc = core.NewNullString("商品状态 自动创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_STATUS_ZTXC
		status.Type = STATUS_TYPE_PRODUCT
		status.Name = "暂停销售"
		status.Desc = core.NewNullString("商品状态 暂停销售", true)
		db_session.Insert(status)
	}

	//店铺账号状态
	{
		status := new(StatusInfo)
		status.Id = SHOP_ACCOUNT_STATUS_ENABLED
		status.Type = STATUS_TYPE_SHOP_ACCOUNT
		status.Name = "正常"
		status.Desc = core.NewNullString("店铺账号状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = SHOP_ACCOUNT_STATUS_DISABLED
		status.Type = STATUS_TYPE_SHOP_ACCOUNT
		status.Name = "停用"
		status.Desc = core.NewNullString("店铺账号状态 停用", true)
		db_session.Insert(status)
	}

	//Paypal账号状态

	{
		status := new(StatusInfo)
		status.Id = PAYPAL_ACCOUNT_STATUS_ENABLED
		status.Type = STATUS_TYPE_PAYPAL_ACCOUNT
		status.Name = "正常"
		status.Desc = core.NewNullString("Paypal账号状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PAYPAL_ACCOUNT_STATUS_DISABLED
		status.Type = STATUS_TYPE_PAYPAL_ACCOUNT
		status.Name = "停用"
		status.Desc = core.NewNullString("Paypal账号状态 停用", true)
		db_session.Insert(status)
	}

	//包材状态
	{
		status := new(StatusInfo)
		status.Id = PACKAGE_MTRL_STATUS_ENABLED
		status.Type = STATUS_TYPE_PACKAGE_MTRL
		status.Name = "正常"
		status.Desc = core.NewNullString("包材状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PACKAGE_MTRL_STATUS_DISABLED
		status.Type = STATUS_TYPE_PACKAGE_MTRL
		status.Name = "停用"
		status.Desc = core.NewNullString("包材状态 停用", true)
		db_session.Insert(status)
	}

	// 盘点单状态
	{
		status := new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_YCJ
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "已创建"
		status.Desc = core.NewNullString("盘点单状态 已创建", true)
		db_session.Insert(status)


		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_WTG
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "未通过"
		status.Desc = core.NewNullString("盘点单状态 未通过", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_SPZ
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "审核中"
		status.Desc = core.NewNullString("盘点单状态 审核中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = INVENTORY_CHECK_STATUS_YWC
		status.Type = STATUS_TYPE_INVENTORY_CHECK
		status.Name = "已完成"
		status.Desc = core.NewNullString("盘点单状态 已完成", true)
		db_session.Insert(status)

		// Paypal对账结果
		status = new(StatusInfo)
		status.Id = PAYPAL_RECON_STATUS_SUCCESS
		status.Type = STATUS_TYPE_PAYPAL_RECON
		status.Name = "匹配成功"
		status.Desc = core.NewNullString("Paypal对账结果 匹配成功", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PAYPAL_RECON_STATUS_FAILED
		status.Type = STATUS_TYPE_PAYPAL_RECON
		status.Name = "匹配失败"
		status.Desc = core.NewNullString("Paypal对账结果 匹配失败", true)
		db_session.Insert(status)
	}

	// 订单自定义类型
	{
		status := new(StatusInfo)
		status.Id = ORDER_CUSTOM_TYPE_STATUS_ENABLED
		status.Type = STATUS_TYPE_ORDER_CUSTOM_TYPE
		status.Name = "正常"
		status.Desc = core.NewNullString("订单自定义类型 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_CUSTOM_TYPE_STATUS_DISABLED
		status.Type = STATUS_TYPE_ORDER_CUSTOM_TYPE
		status.Name = "停用"
		status.Desc = core.NewNullString("订单自定义类型 停用", true)
		db_session.Insert(status)

		//规则订单运单号状态
		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_BATCH_STATUS_ENABLED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER_BATCH
		status.Name = "正常"
		status.Desc = core.NewNullString("订单运单号 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_BATCH_STATUS_DISABLED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER_BATCH
		status.Name = "停用"
		status.Desc = core.NewNullString("订单运单号 停用", true)
		db_session.Insert(status)
	}

	//不规则订单运单号状态
	{
		status := new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_STATUS_NOT_USED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER
		status.Name = "未使用"
		status.Desc = core.NewNullString("订单运单号 未使用", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = LOGISTICS_NUMBER_STATUS_USED
		status.Type = STATUS_TYPE_LOGISTICS_NUMBER
		status.Name = "已使用"
		status.Desc = core.NewNullString("订单运单号 已使用", true)
		db_session.Insert(status)
	}

	// 绑定商品状态 正常 停产
	{
		status := new(StatusInfo)
		status.Id = PRODUCT_BIND_STATUS_NORMAL
		status.Type = STATUS_TYPE_PRODUCT_BIND
		status.Name = "正常"
		status.Desc = core.NewNullString("绑定商品状态 正常", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PRODUCT_BIND_STATUS_STOPPED
		status.Type = STATUS_TYPE_PRODUCT_BIND
		status.Name = "停产"
		status.Desc = core.NewNullString("绑定商品状态 停产", true)
		db_session.Insert(status)
	}

	// 入库单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_CHECKIN_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKIN
		status.Name = "临时保存"
		status.Desc = core.NewNullString("入库单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_CHECKIN_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKIN
		status.Name = "已提交入库"
		status.Desc = core.NewNullString("入库单状态 已提交入库", true)
		db_session.Insert(status)
	}

	// 出库单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_CHECKOUT_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKOUT
		status.Name = "临时保存"
		status.Desc = core.NewNullString("出库单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_CHECKOUT_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_CHECKOUT
		status.Name = "已提交出库"
		status.Desc = core.NewNullString("出库单状态 已提交出库", true)
		db_session.Insert(status)
	}

	// 调拨单状态  暂存 提交
	{
		status := new(StatusInfo)
		status.Id = STOCK_BILL_TRANSFER_STATUS_SAVE
		status.Type = STATUS_TYPE_STOCK_BILL_TRANSFER
		status.Name = "临时保存"
		status.Desc = core.NewNullString("调拨单状态 临时保存", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = STOCK_BILL_TRANSFER_STATUS_SUBMIT
		status.Type = STATUS_TYPE_STOCK_BILL_TRANSFER
		status.Name = "已提交调拨"
		status.Desc = core.NewNullString("调拨单状态 已提交调拨", true)
		db_session.Insert(status)
	}

	// 订单状态 已支付 待配货 配货中 已发货 已完成
	{
		status := new(StatusInfo)
		status.Id = ORDER_STATUS_YCJ
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已创建"
		status.Desc = core.NewNullString("订单状态 已创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YZF
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已支付"
		status.Desc = core.NewNullString("订单状态 已支付", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_DPH
		status.Type = STATUS_TYPE_ORDER
		status.Name = "待配货"
		status.Desc = core.NewNullString("订单状态 待配货", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_PHZ
		status.Type = STATUS_TYPE_ORDER
		status.Name = "配货中"
		status.Desc = core.NewNullString("订单状态 配货中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YFH
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已发货"
		status.Desc = core.NewNullString("订单状态 已发货", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = ORDER_STATUS_YWC
		status.Type = STATUS_TYPE_ORDER
		status.Name = "已完成"
		status.Desc = core.NewNullString("订单状态 已完成", true)
		db_session.Insert(status)
	}

	// 采购单状态 已创建，已下单，采购中，已入库
	{
		status := new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YCJ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已创建"
		status.Desc = core.NewNullString("采购单状态 已创建", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YXD
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已下单"
		status.Desc = core.NewNullString("采购单状态 已下单", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_CGZ
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "采购中"
		status.Desc = core.NewNullString("采购单状态 采购中", true)
		db_session.Insert(status)

		status = new(StatusInfo)
		status.Id = PURCHASE_BILL_STATUS_YRK
		status.Type = STATUS_TYPE_PURCHASE_BILL
		status.Name = "已入库"
		status.Desc = core.NewNullString("采购单状态 已入库", true)
		db_session.Insert(status)
	}
}
