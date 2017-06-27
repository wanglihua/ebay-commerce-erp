package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type OrderHoldTypeInfo struct {
	Id         int64               `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string              `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称"`
	Desc       string              `xorm:"nvarchar(500) notnull 'desc'" json:"desc" smith:"verbose_name=备注"`
	IconFile   string              `xorm:"nvarchar(200) null 'icon_file'" json:"icon_file" smith:"verbose_name=图标"`
	Ord        int64               `xorm:"bigint null 'ord'" json:"ord" smith:"verbose_name=排序"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderHoldTypeInfo) TableName() string {
	return TablePrefix + "order_hold_type"
}

func (e OrderHoldTypeInfo) ModelDesc() string {
	return "verbose_name=订单Hold类型,entity_json=order_hold_type,route=order/hold/type"
}

const (
	ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH int64 = 10000 // 卖家邮件和PayPal的地址不匹配
	ORDER_HOLD_TYPE_SPECIAL_ADDRESS        int64 = 10010 // 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能
	ORDER_HOLD_TYPE_BUYER_CHECKOUT_MESSAGE int64 = 10020 // 有买家留言
	ORDER_HOLD_TYPE_PLATFORM_LEAVE_WORD    int64 = 10030 // 有平台留言
	ORDER_HOLD_TYPE_PICK_UP                int64 = 10040 // PICK UP
	ORDER_HOLD_TYPE_OOS_ALL_STOCK          int64 = 10050 // 系统显示所有仓都缺货（中国仓除外） OOS: Out Of Stock
	ORDER_HOLD_TYPE_OOS_AU_STOCK_ZT_HAS    int64 = 10055 // 澳仓无货，在途才有货
	ORDER_HOLD_TYPE_OOS_BEST_STOCK         int64 = 10060 // 匹配运费后，最佳仓库无货，运费偏差范围能的次佳仓库有货
	ORDER_HOLD_TYPE_OOS_ACTUAL_STOCK       int64 = 10070 // 实际仓库缺货
	ORDER_HOLD_TYPE_OOS_SOME_PRODUCT       int64 = 10080 // 多数量或多产品，但部分缺货（是否拆分订单）
	ORDER_HOLD_TYPE_MULTI_PRODUCT          int64 = 10075 // 多商品，这是个临时的状态，下面三个的多商品判断实在太麻烦，先用这个临时解决方案
	ORDER_HOLD_TYPE_MULTI_STOCK_LOW_FEE    int64 = 10090 // 多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜
	ORDER_HOLD_TYPE_ONO_STOCK_EASY         int64 = 10100 // 多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便
	ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH     int64 = 10110 // 无匹配的物流方式（长宽高重量没有物流满足）
	ORDER_HOLD_TYPE_PRODUCT_INFO_MISS      int64 = 10120 // 产品信息不全（长宽高重量不全或者有误）
	ORDER_HOLD_TYPE_GREAT_DEFICIT          int64 = 10130 // 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
	ORDER_HOLD_TYPE_ENVELOPE               int64 = 10140 // 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
	ORDER_HOLD_TYPE_EXPRESS                int64 = 10150 // 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
	ORDER_HOLD_TYPE_INTERNATIONAL          int64 = 10160 // 国际件。对于非澳大利亚的订单，直接到国际件模块
	ORDER_HOLD_ADDRESS_INFO_MISS           int64 = 10170 // 收货人信息不全的，比如没有手机号，没有邮编等等
)

func (e OrderHoldTypeInfo) InitData(db_session *xorm.Session) {
	// ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH int64 = 10000 // 卖家邮件和PayPal的地址不匹配
	order_hold_type := new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH
	order_hold_type.Name = "卖家paypal不匹配"
	order_hold_type.Desc = "店铺配置的paypal账号和订单中卖家的paypal账号不匹配"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_SPECIAL_ADDRESS        int64 = 10010 // 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_SPECIAL_ADDRESS
	order_hold_type.Name = "地址有特殊字段"
	order_hold_type.Desc = "地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_BUYER_LEAVE_WORD       int64 = 10020 // 有买家留言
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_BUYER_CHECKOUT_MESSAGE
	order_hold_type.Name = "有买家留言"
	order_hold_type.Desc = "有买家留言"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_SELLER_LEAVE_WORD      int64 = 10030 // 有卖家留言
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_PLATFORM_LEAVE_WORD
	order_hold_type.Name = "有平台留言"
	order_hold_type.Desc = "有平台留言"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_PICK_UP                int64 = 10040 // PICK UP
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_PICK_UP
	order_hold_type.Name = "PICK UP"
	order_hold_type.Desc = "PICK UP"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_OOS_ALL_STOCK          int64 = 10050 // 系统显示所有仓都缺货（中国仓除外） OOS: Out Of Stock
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_OOS_ALL_STOCK
	order_hold_type.Name = "缺货"
	order_hold_type.Desc = "系统显示所有仓都缺货（中国仓除外）"
	db_session.Insert(order_hold_type)

	//ORDER_HOLD_TYPE_OOS_AU_STOCK_ZT_HAS
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_OOS_AU_STOCK_ZT_HAS
	order_hold_type.Name = "在途才有货"
	order_hold_type.Desc = "澳仓没货，在途仓才有货"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_OOS_BEST_STOCK         int64 = 10060 // 系统显示最佳仓缺货（次佳、海上有货）
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_OOS_BEST_STOCK
	order_hold_type.Name = "最佳仓缺货"
	order_hold_type.Desc = "系统显示最佳仓缺货（次佳、海上有货）"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_OOS_ACTUAL_STOCK       int64 = 10070 // 实际仓库缺货
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_OOS_ACTUAL_STOCK
	order_hold_type.Name = "实际仓库缺货"
	order_hold_type.Desc = "实际仓库缺货"
	db_session.Insert(order_hold_type)

	/*
	// ORDER_HOLD_TYPE_MULTI_PRODUCT          int64 = 10075 // 多商品，这是个临时的状态，下面三个的多商品判断实在太麻烦，先用这个临时解决方案
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_MULTI_PRODUCT
	order_hold_type.Name = "多商品"
	order_hold_type.Desc = "多商品的一个临时解决方案"
	db_session.Insert(order_hold_type)
	*/

	// ORDER_HOLD_TYPE_OOS_SOME_PRODUCT       int64 = 10080 // 多数量或多产品，但部分缺货（是否拆分订单）
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_OOS_SOME_PRODUCT
	order_hold_type.Name = "部分商品缺货"
	order_hold_type.Desc = "多数量或多产品，但部分缺货（是否拆分订单）"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_MULTI_STOCK_LOW_FEE    int64 = 10090 // 多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_MULTI_STOCK_LOW_FEE
	order_hold_type.Name = "拆单发货便宜"
	order_hold_type.Desc = "多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_ONO_STOCK_EASY         int64 = 10100 // 多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_ONO_STOCK_EASY
	order_hold_type.Name = "合单发货方便"
	order_hold_type.Desc = "多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH     int64 = 10110 // 无匹配的物流方式（长宽高重量没有物流满足）
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH
	order_hold_type.Name = "无匹配的物流方式"
	order_hold_type.Desc = "无匹配的物流方式（长宽高重量没有物流满足）"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_PRODUCT_INFO_MISS      int64 = 10120 // 产品信息不全（长宽高重量不全或者有误）
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_PRODUCT_INFO_MISS
	order_hold_type.Name = "产品信息不全"
	order_hold_type.Desc = "产品信息不全（长宽高重量不全或者有误）"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_GREAT_DEFICIT          int64 = 10130 // 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_GREAT_DEFICIT
	order_hold_type.Name = "订单亏损"
	order_hold_type.Desc = "匹配最佳仓库及物流后订单亏损"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_ENVELOPE               int64 = 10140 // 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_ENVELOPE
	order_hold_type.Name = "信封类"
	order_hold_type.Desc = "信封类"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_EXPRESS                int64 = 10150 // 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_EXPRESS
	order_hold_type.Name = "快递类"
	order_hold_type.Desc = "快递类，客户付了快递费的"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_TYPE_INTERNATIONAL          int64 = 10160 // 国际件。对于非澳大利亚的订单，直接到国际件模块
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_INTERNATIONAL
	order_hold_type.Name = "国际件"
	order_hold_type.Desc = "国际件。对于非澳大利亚的订单，直接到国际件模块"
	db_session.Insert(order_hold_type)

	// ORDER_HOLD_ADDRESS_INFO_MISS           int64 = 10170 // 收货人信息不全的，比如没有手机号，没有邮编等等
	order_hold_type = new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_ADDRESS_INFO_MISS
	order_hold_type.Name = "收货人信息不全"
	order_hold_type.Desc = "收货人信息不全的，比如没有手机号，没有邮编等等"
	db_session.Insert(order_hold_type)
}
