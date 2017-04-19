package pass_res

const (
	// 权限管理
	AUTH_USER  = 10100 // 用户管理
	AUTH_ADMIN = 10200 // 管理员组

	// 采购管理

	// 调拨管理

	// 订单报表

	// 订单流程
	DDLC_ORDER_ALL = 11100 // 所有订单
	DDLC_ORDER_YZF = 11200 // 已支付
	DDLC_ORDER_DPH = 11300 // 待配货
	DDLC_ORDER_PHZ = 11400 // 配货中
	DDLC_ORDER_YFH = 11500 // 已发货
	DDLC_ORDER_YWC = 11600 // 已完成
	DDLC_ORDER_KHB = 11700 // 可合并的订单

	// 订单配置
	DDPZ_LOGISTICS_MODE          = 11800 // 货运方式
	DDPZ_PLATFORM_LOGISTICS_MODE = 11900 // 电商平台物流方式
	DDPZ_ORDER_CUSTOM_TYPE       = 12000 // 订单自定义类型
	DDPZ_ORDER_LOG               = 12100 // 订单修改日志

	// 库存管理
	KCGL_PRODUCT                  = 12200 // 商品列表
	KCGL_PRODUCT_BIND             = 12300 // 绑定商品
	KCGL_PRODUCT_IMPORT           = 12400 // 商品导入
	KCGL_PRODUCT_EXPORT           = 12500 // 商品导出
	KCGL_INVENTORY                = 12600 // 商品库存
	KCGL_INVENTORY_CHECK          = 12700 // 库存盘点
	KCGL_INVENTORY_CHECK_APPROVE  = 12700 // 盘点审批
	KCGL_STOCKBILL_TRANSFER       = 12800 // 分仓调拨
	KCGL_PRODUCT_QUALITY_FEEDBACK = 12900 // 商品质量反馈
	KCGL_STOCKBILL_CHECKIN        = 13000 // 手工入库
	KCGL_STOCKBILL_CHECKOUT       = 13100 // 手工出库

	// 新的编号办法： 前两位是 模块编号，11 开始；后两位是菜单编号 01 开始；最后两位保留用于中间插入新的号码

	// 系统配置，模块编号 11

	XTPZ_ECOMMERCE_PLATFORM = 110100 // 电商平台
	XTPZ_SHOP_ACCOUNT       = 110200 // 店铺账号
	XTPZ_PAYPAL_ACCOUNT     = 110300 // PAYPAL账号
	XTPZ_EXCHANGE_RATE      = 110400 // 汇率设置
	XTPZ_LOGIN_LOG          = 110500 // 系统日志
	XTPZ_SCHEDULE_TASK      = 110600 // 定时任务
	XTPZ_LITTLE_ICON        = 110700 // 图标管理
)
