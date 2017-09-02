package pass_res

const (
	// 权限管理
	AUTH_USER = 10100 // 用户管理
	// AUTH_ADMIN = 10200 // 管理员组 这个不用授权啥的，只有超级管理员有权限的
	AUTH_PASS    = 10300 // 授权管理
	AUTH_GROUP   = 10400 // 角色管理
	AUTH_SETTING = 10500 // 授权配置

	// 新的编号办法： 前两位是 模块编号，11 开始；后两位是菜单编号 01 开始；最后两位保留用于中间插入新的号码

	// 系统配置        模块编号 11
	// 订单配置        模块编号 12
	// 订单流程        模块编号 13
	// 库存管理        模块编号 14
	// 商品配置        模块编号 15
	// 海关报关        模块编号 16
	// 商品报表        模块编号 17
	// 商品质量反馈    模块编号 18
	// 系统管理        模块编号 19
	// 采购管理        模块编号 20
	// 销售管理        模块编号 21

	// 采购管理 20
	CGGL_PRODUCT_SETTING = 200100 // 商品采购设置

	// 销售管理 21
	XSGL_PRODUCT_SALE_SETTING = 210100 // 商品销售设置

	// 调拨管理

	// 订单报表

	// 系统配置，模块编号 11
	XTPZ_ECOMMERCE_PLATFORM = 110100 // 电商平台
	XTPZ_SHOP_ACCOUNT       = 110200 // 店铺账号
	XTPZ_PAYPAL_ACCOUNT     = 110300 // PAYPAL账号
	XTPZ_EXCHANGE_RATE      = 110400 // 汇率设置
	XTPZ_LOGIN_LOG          = 110500 // 系统日志
	XTPZ_SCHEDULE_TASK      = 110600 // 定时任务
	XTPZ_LITTLE_ICON        = 110700 // 图标管理

	// 订单配置, 模块编号 12
	DDPZ_LOGISTICS_COMPANY  = 120100 // 物流公司
	DDPZ_LOGISTICS_MODE     = 120200 // 物流方式
	DDPZ_LOGISTICS_FEE_RATE = 120300 // 物流费率
	DDPZ_LOGISTICS_SETTING  = 120400 // 物流配置
	DDPZ_ENVELOPE_FEE_RATE  = 120500 // 信封费率
	DDPZ_PICK_FEE           = 120600 // 拣货费用
	DDPZ_ORDER_HOLD_TYPE    = 120700 // 订单HOLD类型
	DDPZ_ORDER_LOG          = 120800 // 订单修改日志
	DDPZ_ORDER_TAG          = 120900 // 订单标签

	// 订单流程, 模块编号 13
	DDLC_ORDER_ALL                   = 130100 // 所有订单
	DDLC_ORDER_NEW                   = 130130 // 新增订单
	DDLC_ORDER_YCJ                   = 130150 // 已创建
	DDLC_ORDER_YZF                   = 130200 // 已支付
	DDLC_ORDER_HOLD                  = 130300 // HOLD
	DDLC_ORDER_PROCESSED_HOLD        = 130350 // 已处理HOLD
	DDLC_ORDER_DPH                   = 130400 // 待配货
	DDLC_ORDER_PHZ                   = 130500 // 配货中
	DDLC_ORDER_YFH                   = 130600 // 已发货
	DDLC_ORDER_YWC                   = 130700 // 已完成
	DDLC_ORDER_YQX                   = 130730 // 已取消
	DDLC_ORDER_VIEW_DETAIL           = 130750 // 查看订单详细
	DDLC_ORDER_IMPORT                = 130800 // 订单导入
	DDLC_ORDER_LOGISTICS             = 130900 // 订单物流
	DDLC_ORDER_LOGISTICS_FEE_DIFF    = 131000 // 订单物流费用偏差
	DDLC_ORDER_QUXIAO                = 131100 // 订单取消，是指取消订单这个动作， 不是指查看已取消的订单
	DDLC_ORDER_LOGISTICS_MATCH_DEBUG = 131200 // 订单物流匹配调试
	DDLC_ORDER_LOGISTICS_MATCH       = 131300 // 订单物流匹配
	DDLC_ORDER_FORCE_PEIHUO          = 131400 // 强制配货

	// 库存管理 模块编号 14
	KCGL_PRODUCT              = 140100 // 商品列表
	KCGL_PRODUCT_DETAIL_SAVE  = 140150 // 商品详细保存
	KCGL_PRODUCT_BIND         = 140200 // 绑定商品
	KCGL_PRODUCT_IMPORT       = 140300 // 商品导入
	KCGL_PRODUCT_EXPORT       = 140400 // 商品导出
	KCGL_PRODUCT_LWHW         = 140401 // 商品长宽高重量
	KCGL_INVENTORY            = 140500 // 商品库存
	KCGL_INVENTORY_IMPORT_QTY = 140531 // 商品库存 导入库存数量
	KCGL_INVENTORY_EXPORT_QTY = 140532 // 商品库存 导出库存数量
	//KCGL_INVENTORY_IMPORT_QTY_OVER = 140510 // 商品库存 导入覆盖库存数量
	//KCGL_INVENTORY_IMPORT_QTY_DIFF = 140520 // 商品库存 导入修正库存数量
	KCGL_INVENTORY_SKU_LOC = 140501 // 商品库存 商品库位
	//KCGL_INVENTORY_SKU_QTY         = 140502 // 商品库存 库存数量
	KCGL_INVENTORY_SKU_PRICE     = 140503 // 商品库存 库存成本
	KCGL_INVENTORY_CHECK         = 140600 // 库存盘点
	KCGL_INVENTORY_CHECK_APPROVE = 140700 // 盘点审批
	KCGL_STOCKBILL_TRANSFER      = 140800 // 分仓调拨
	KCGL_STOCKBILL_CHECKIN       = 141000 // 手工入库
	KCGL_STOCKBILL_CHECKOUT      = 141100 // 手工出库
	KCGL_SPECIAL_SKU             = 141200 // 特殊SKU

	// 商品配置 模块编号 15
	SPPZ_PRODUCT_BATCH       = 150100 // 商品批量修改
	SPPZ_PACKAGE_MATERIAL    = 150200 // 包材类型
	SPPZ_STOCK               = 150300 // 仓库
	SPPZ_STOCK_PICK_CREW     = 150400 // 仓库拣货人员
	SPPZ_STORAGE_LOC         = 150500 // 仓位
	SPPZ_PRODUCT_CATE        = 150600 // 商品目录
	SPPZ_PSL_BATCH           = 150700 // 批量修改仓位
	SPPZ_SUPPLIER            = 150800 // 供应商
	SPPZ_SUPPLIER_DELETE     = 150850 // 供应商删除
	SPPZ_PRODUCT_STATUS      = 150900 // 商品状态
	SPPZ_SEA_CONTAINER_COST  = 151000 // 商品状态
	SPPZ_BRAND               = 151100 // 商品品牌
	SPPZ_PRODUCT_CERT        = 151200 // 商品认证
	SPPZ_PRODUCT_STYLE       = 151300 // 商品款式
	SPPZ_SUPPLIER_PAY_METHOD = 151400 // 供应商付款方式

	// 海关报关
	HGBG_PRODUCT_HS_CODE                = 160100 // 商品海关编码
	HGBG_HS_CODE_CN_REPORT_ELEM         = 160200 // 中国海关编码
	HGBG_HS_CODE_AU_DUTY                = 160300 // 澳洲海关编码
	HGBG_CONTAINER_TRANSFER             = 160400 // 集装箱调拨
	HGBG_CONTAINER_TRANSFER_TO_ZT_STOCK = 160410 // 集装箱调拨 调拨至在途仓
	HGBG_CONTAINER_TRANSFER_TO_AU_STOCK = 160420 // 集装箱调拨 调拨至澳洲仓
	HGBG_PRODUCT_INVOICE                = 160500 // 商品是否开票

	// 商品报表
	SPBB_PRODUCT_SALE_ANALYSIS_SETTING = 170100 // 商品销售分析设置
	SPBB_PRODUCT_SALE_ANALYSIS         = 170200 // 商品销售分析
	SPBB_PRODUCT_SALE_DAY_ANALYSIS     = 170300 // 商品日销售分析

	// 商品质量反馈 18
	SPZL_PRODUCT_FEEDBACK_CATE   = 180100 // 商品质量反馈类别
	SPZL_PRODUCT_FEEDBACK        = 180200 // 商品质量反馈
	SPZL_PRODUCT_FEEDBACK_DELETE = 180300 // 商品质量反馈删除

	// 系统管理 19
	XTGL_NOTICE = 190100 // 通知管理
)
