package tasks

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/ebay_service"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/order_shipping_address_service"
	"matrix/modules/erp/service/product_bind_service"
	"matrix/modules/erp/service/product_service"
	"matrix/modules/erp/service/product_store_item_service"
	"matrix/modules/erp/service/order_history_service"
)

type EBayOrderDownloadTask struct {
}

func (t EBayOrderDownloadTask) Name() string {
	return "订单下载"
}

func (t EBayOrderDownloadTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("下载订单出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	//todo: ebay_order_download_task.go:86: 下载订单出现错误： Error 1364: Field 'express_fee' doesn't have a default value 错误来之服务器日志

	job_time_begin := time.Now()

	// get db driver and db url

	if task_db_driver == "" {
		task_db_driver, _ = revel.Config.String("db.driver")
	}

	if task_db_url == "" {
		task_db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(task_db_driver, task_db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	ebay_api_config_list := get_ebay_api_config_list(db_session)

	for _, shop_account_config := range ebay_api_config_list {
		has_more_orders := true
		page_number := 1
		for has_more_orders == true {
			response := ebay_service.Get_orders(shop_account_config.APIConfig, page_number)

			switch response.Ack {
			case "Success":
				revel.TRACE.Printf("eBay GetOrders Response: " + response.Ack + "\n")
				revel.TRACE.Printf("eBay GetOrders actual fetched orders count: %d", response.ReturnedOrderCountActual)
			case "Failure":
				revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
				continue
			default:
				revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
				continue
			}

			order_list := response.OrdersArray.Orders

			for _, order_download := range order_list {
				order_status := order_download.OrderStatus

				// 如果是已完成的订单， 就忽略
				/* todo: 临时去掉这个检查
				if order_download.ShippedTime != "" {
					shipped_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", order_download.ShippedTime, time.UTC)
					core.HandleError(err)
					shipped_time = shipped_time.Local()
					if shipped_time.IsZero() == false && err == nil {
						// 已发货的订单，
						continue
					}
				}
				*/

				order_code := order_download.OrderID // 已付款的 282359220529-1651978360018, 未付款的 282387576804-0

				order_code_split_str_list := strings.Split(order_code, "-")
				order_code_min_length := 13
				if len(order_code_split_str_list) > 1 && len(order_code_split_str_list[1]) >= order_code_min_length {
					order_code = order_code_split_str_list[1]
				}

				var order_in_db models.OrderInfo
				has_order, err := db_session.Where("code = ?", order_code).Get(&order_in_db)
				core.HandleError(err)

				order_in_db.Code = order_code
				order_in_db.ShopId = shop_account_config.APIConfig.ShopId
				order_in_db.TranId = order_download.ShippingDetails.SellingManagerSalesRecordNumber
				order_in_db.FullTranId = shop_account_config.Shop.Code + order_download.ShippingDetails.SellingManagerSalesRecordNumber

				var shop models.ShopAccountInfo
				_, err = db_session.Id(shop_account_config.APIConfig.ShopId).Get(&shop)
				core.HandleError(err)

				order_in_db.TypeId = models.ORDER_TYPE_ZCDD //订单类型为正常下载订单

				update_col_list := []string{"type_id", "ebay_status", "status_id"}

				if order_download.CreatedTime != "" {
					created_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", order_download.CreatedTime, time.UTC)
					core.HandleError(err)
					created_time = created_time.Local()
					order_in_db.CreatedTime = core.MakeNullLongTime(created_time, true)
					update_col_list = append(update_col_list, "created_time")
				}

				if order_download.ModifiedTime != "" {
					modified_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", order_download.ModifiedTime, time.UTC)
					core.HandleError(err)
					modified_time = modified_time.Local()
					order_in_db.ModifiedTime = core.MakeNullLongTime(modified_time, true)
					update_col_list = append(update_col_list, "modified_time")
				}

				if order_download.PaidTime != "" {
					paid_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", order_download.PaidTime, time.UTC)
					core.HandleError(err)
					paid_time = paid_time.Local()
					order_in_db.PaidTime = core.MakeNullLongTime(paid_time, true)
					update_col_list = append(update_col_list, "paid_time")
				}

				if order_download.ShippedTime != "" {
					shipped_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", order_download.ShippedTime, time.UTC)
					core.HandleError(err)
					shipped_time = shipped_time.Local()
					order_in_db.ShippedTime = core.MakeNullLongTime(shipped_time, true)
					update_col_list = append(update_col_list, "shipped_time")
				}

				// 获取当前的汇率设置
				var exchange_rate_list = exchange_rate_service.Get_exchange_rate_list(db_session)

				var exchange_rate_map = make(map[string]models.ExchangeRateInfo)
				for _, exchange_rate := range exchange_rate_list {
					exchange_rate_map[exchange_rate.Currency.Code] = exchange_rate.ExchangeRate
				}

				order_in_db.AdjustmentAmountCurrency = order_download.AdjustmentAmount.Currency
				update_col_list = append(update_col_list, "adjustment_amount_currency")
				order_in_db.AdjustmentAmount = order_download.AdjustmentAmount.Value
				update_col_list = append(update_col_list, "adjustment_amount")
				if exchange_rate, found := exchange_rate_map[order_download.AdjustmentAmount.Currency]; found == true {
					order_in_db.AdjustmentAmountRmb = order_download.AdjustmentAmount.Value * exchange_rate.Rate
					update_col_list = append(update_col_list, "adjustment_amount_rmb")
				}

				order_in_db.AmountPaidCurrency = order_download.AmountPaid.Currency
				update_col_list = append(update_col_list, "amount_paid_currency")
				order_in_db.AmountPaid = order_download.AmountPaid.Value
				update_col_list = append(update_col_list, "amount_paid")
				if exchange_rate, found := exchange_rate_map[order_download.AmountPaid.Currency]; found == true {
					order_in_db.ExchangeRate = exchange_rate.Rate //把 amount_paid 的汇率设置成 订单的汇率记录
					update_col_list = append(update_col_list, "exchange_rate")
					order_in_db.AmountPaidRmb = order_download.AmountPaid.Value * exchange_rate.Rate
					update_col_list = append(update_col_list, "amount_paid_rmb")
				}

				order_in_db.AmountSavedCurrency = order_download.AmountSaved.Currency
				update_col_list = append(update_col_list, "amount_saved_currency")
				order_in_db.AmountSaved = order_download.AmountSaved.Value
				update_col_list = append(update_col_list, "amount_saved")
				if exchange_rate, found := exchange_rate_map[order_download.AmountSaved.Currency]; found == true {
					order_in_db.AmountSavedRmb = order_download.AmountSaved.Value * exchange_rate.Rate
					update_col_list = append(update_col_list, "amount_saved_rmb")
				}

				//Total
				order_in_db.TotalAmountCurrency = order_download.Total.Currency
				update_col_list = append(update_col_list, "total_amount_currency")
				order_in_db.TotalAmount = order_download.Total.Value
				update_col_list = append(update_col_list, "total_amount")
				if exchange_rate, found := exchange_rate_map[order_download.Total.Currency]; found == true {
					order_in_db.TotalAmountRmb = order_download.Total.Value * exchange_rate.Rate
					update_col_list = append(update_col_list, "total_amount_rmb")
				}

				//SubTotal
				order_in_db.SubTotalAmountCurrency = order_download.Subtotal.Currency
				update_col_list = append(update_col_list, "subtotal_amount_currency")
				order_in_db.SubTotalAmount = order_download.Subtotal.Value
				update_col_list = append(update_col_list, "subtotal_amount")
				if exchange_rate, found := exchange_rate_map[order_download.Subtotal.Currency]; found == true {
					order_in_db.SubTotalAmountRmb = order_download.Subtotal.Value * exchange_rate.Rate
					update_col_list = append(update_col_list, "subtotal_amount_rmb")
				}

				order_in_db.SellerEmail = order_download.SellerEmail
				order_in_db.SellerUserId = order_download.SellerUserID

				order_in_db.BuyerUserId = order_download.BuyerUserID
				if len(order_download.TransactionArray.Transactions) > 0 {
					order_in_db.BuyerEmail = order_download.TransactionArray.Transactions[0].Buyer.Email
				}
				order_in_db.BuyerCheckoutMessage = order_download.BuyerCheckoutMessage

				// 订单状态，这部分逻辑还可以继续探讨，判断状态的时候用到了前面已经设置的 order_in_db的状态，所以这段代码不能放到前面去的

				// 订单状态 设置 开始
				order_in_db.EBayStatus = order_status
				order_status_id := order_in_db.StatusId

				/*
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
				*/

				// """同步逻辑要保证中间的 待配货，配货中，已发货 等的状态在同步时不被前面的状态覆盖，只有后面的“已完成”可以覆盖这些状态""
				// 所以下面的两个条件中需要加 < 号的状态判断

				var status_id_before = order_in_db.StatusId

				// 如果是新的订单，一般是新增的订单，系统数据库中还没有的订单
				if order_status_id < models.ORDER_STATUS_YCJ && //这个状态检查不能少，否则 待配货，已发货的状态可能会被覆盖
					order_in_db.PaidTime.Valid == false {
					order_in_db.StatusId = models.ORDER_STATUS_YCJ // 这个状态，有可能会被接下来的状态逻辑所修改，可能会，更进一步改成 “已支付”, “已完成”
				}

				// 有可能是新增的订单（最有可能），或者是原来 已创建 的状态过来的。这时支付时间不为空
				if order_status_id < models.ORDER_STATUS_YZF && //这个状态检查不能少，否则 待配货，已发货的状态可能会被覆盖
					order_in_db.PaidTime.Valid == true &&
					strings.ToLower(order_download.CheckoutStatuss.Status) == strings.ToLower("Complete") {
					order_in_db.StatusId = models.ORDER_STATUS_YZF // 这个状态，有可能会被接下来的状态逻辑所修改，可能会，更进一步改成 “已完成”
				}

				// 如果 下载订单里面的 OrderStatus 是 Completed， 且订单的发货时间不为空，订单状态就为 已完成
				if order_in_db.ShippedTime.Valid == true {
					order_in_db.StatusId = models.ORDER_STATUS_YWC
				}

				//todo cancel_status

				//订单状态 设置 结束

				var payment_method = order_download.CheckoutStatuss.PaymentMethod
				order_in_db.PaymentMethods = payment_method

				/* 这部分判断放到物流匹配中去
				// 老鼠仓的判断
				// 如果是通过PayPal支付的，而且店铺设置里面也设置了PayPal列表的
				if strings.ToLower(payment_method) == strings.ToLower("PayPal") && shop.PaypalAccount != "" {
					shop_paypal_email_list := strings.Split(shop.PaypalAccount, ",")

					seller_email := order_download.SellerEmail //是用SellerEmail来判断老鼠仓吗？

					seller_email_in_list := false
					for _, shop_paypal_email := range shop_paypal_email_list {
						if strings.ToLower(seller_email) == strings.ToLower(shop_paypal_email) {
							seller_email_in_list = true
						}
					}

					// 设置的账户里面和卖家账户匹配的，就设置成为Hold订单
					if seller_email_in_list == false {
						order_in_db.HoldTypeId = core.MakeNullInt(models.ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH, true)
					}
				}
				*/

				if has_order == true {
					// update order
					_, err = db_session.Id(order_in_db.Id).Cols(update_col_list...).Update(&order_in_db)
					core.HandleError(err)

					if status_id_before != order_in_db.StatusId {
						order_history_service.Insert_order_history(order_in_db, "订单下载状态更新", db_session)
					}
				} else {
					// insert order
					order_in_db.Code = order_code

					_, err = db_session.Insert(&order_in_db)
					core.HandleError(err)

					order_history_service.Insert_order_history(order_in_db, "订单下载新增订单", db_session)
				}

				order_id := order_in_db.Id

				//订单明细信息
				/*
					Id                       int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
					OrderId                  int64           `xorm:"bigint notnull index 'order_id'" json:"order_id"`
					OrderLineItemId          string          `xorm:"nvarchar(200) notnull 'order_line_item_id'" json:"order_line_item_id"`
					Quantity                 int             `xorm:"int notnull 'quantity'" json:"quantity"`
					TransactionId            string          `xorm:"nvarchar(200) notnull 'transaction_id'" json:"transaction_id"`
					TransactionPriceCurrency string          `xorm:"nvarchar(50) notnull 'transaction_price_currency'" json:"transaction_price_currency"`
					TransactionPrice         float64         `xorm:"decimal(18,6) notnull 'transaction_price'" json:"transaction_price"`
					TransactionPriceRmb      float64         `xorm:"decimal(18,6) notnull 'transaction_price_rmb'" json:"transaction_price_rmb"`
					ProductStoreItemId       int64           `xorm:"bigint notnull index 'product_store_item_id'" json:"product_store_item_id"` // 对应的 ProductItemInfo对象的数据库id
					ItemId                   string          `xorm:"nvarchar(50) notnull 'item_id'" json:"item_id"`                             // item 对应的在 ebay系统里面的id
					ItemProductId            int64           `xorm:"bigint notnull index 'item_product_id'" json:"item_product_id"`             // item 对应的 ProductInfo 对象的数据库id
					ItemSku                  string          `xorm:"nvarchar(100) notnull index 'item_sku'" json:"item_sku"`
					ItemTitle                string          `xorm:"nvarchar(200) notnull 'item_title'" json:"item_title"`
					ItemConditionId          string          `xorm:"nvarchar(50) notnull 'item_condition_id'" json:"item_condition_id"`
					ItemConditionDisplayName string          `xorm:"nvarchar(100) notnull 'item_condition_display_name'" json:"item_condition_display_name"`
					CreatedDate              core.LongTime   `xorm:"datetime notnull 'created_date'" json:"created_date"`
					Desc                     core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
				 */

				for _, transaction := range order_download.TransactionArray.Transactions {
					// 获取并检查明细是否存在，根据 order_id 和 item_sku
					item_sku := transaction.Item.SKU

					// 考虑是不是绑定商品的情况，先查一下这个SKU是不是在 绑定商品列表里面的SKU,
					is_product_bind_sku := product_bind_service.Is_product_bind_sku(item_sku, db_session)
					if is_product_bind_sku == true {
						product_bind_sku_list := product_bind_service.Get_product_bind_skus(item_sku, db_session)

						//todo: 这里还要做一个根据各个商品的加权平均，计算商品利润的情况
						for _, product_bind_sku := range product_bind_sku_list {
							product_sku := product_bind_sku.Sku
							product_sku_count := product_bind_sku.Count

							var order_item models.OrderItemInfo

							has_order_item, err := db_session.Where(fmt.Sprintf("order_id = %d and item_sku = '%s'", order_id, product_sku)).Get(&order_item)
							core.HandleError(err)

							// 设置明细字段
							order_item.OrderId = order_id
							order_item.OrderLineItemId = transaction.OrderLineItemID
							order_item.Quantity = transaction.QuantityPurchased * int(product_sku_count)
							order_item.TransactionId = transaction.TransactionID
							order_item.TransactionPriceCurrency = transaction.TransactionPrice.Currency
							order_item.TransactionPrice = transaction.TransactionPrice.Value
							if exchange_rate, found := exchange_rate_map[transaction.TransactionPrice.Currency]; found == true {
								order_item.TransactionPriceRmb = transaction.TransactionPrice.Value * exchange_rate.Rate
							}
							order_item.ItemId = transaction.Item.ItemID

							// 要根据 SKU 下载 商品吗，觉得解析出来的SKU作为商品不一定在eBay上卖的，所以就觉得没必要了的

							//order_item.ItemProductId 如果这时产品不存在会怎样，下载产品, 这时下载的是 绑定的虚拟商品，不是解析出来的sku的商品
							has_product_store_item, _, product_store_item_id := product_store_item_service.Query_product_store_item(transaction.Item.ItemID, db_session)
							if has_product_store_item == false {
								//download product
								_, _, product_store_item_id = fetch_product_item_to_db(shop_account_config.APIConfig, db_session, transaction.Item.ItemID)
							}

							product_id := product_service.Get_product_id_by_sku(product_sku, db_session)

							order_item.ItemProductId = product_id
							order_item.ProductStoreItemId = product_store_item_id
							order_item.ItemSku = transaction.Item.SKU
							order_item.ItemTitle = transaction.Item.Title
							order_item.ItemConditionId = transaction.Item.ConditionID
							order_item.ItemConditionDisplayName = transaction.Item.ConditionDisplayName

							transaction_created_date, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", transaction.CreatedDate, time.UTC)
							core.HandleError(err)
							transaction_created_date = transaction_created_date.Local()
							order_item.CreatedDate = core.LongTime(transaction_created_date)

							// 新增或更新明细
							if has_order_item == false {
								//新增
								_, err = db_session.Insert(&order_item)
								core.HandleError(err)
							} else {
								//更新
								_, err = db_session.Id(order_item.Id).Update(&order_item)
								core.HandleError(err)
							}
						}

					} else {
						var order_item models.OrderItemInfo

						has_order_item, err := db_session.Where(fmt.Sprintf("order_id = %d and item_sku = '%s'", order_id, item_sku)).Get(&order_item)
						core.HandleError(err)

						// 设置明细字段
						order_item.OrderId = order_id
						order_item.OrderLineItemId = transaction.OrderLineItemID
						order_item.Quantity = transaction.QuantityPurchased
						order_item.TransactionId = transaction.TransactionID
						order_item.TransactionPriceCurrency = transaction.TransactionPrice.Currency
						order_item.TransactionPrice = transaction.TransactionPrice.Value
						if exchange_rate, found := exchange_rate_map[transaction.TransactionPrice.Currency]; found == true {
							order_item.TransactionPriceRmb = transaction.TransactionPrice.Value * exchange_rate.Rate
						}
						order_item.ItemId = transaction.Item.ItemID

						//order_item.ItemProductId 如果这时产品不存在会怎样，下载产品
						has_product_store_item, product_id, product_store_item_id := product_store_item_service.Query_product_store_item(transaction.Item.ItemID, db_session)
						if has_product_store_item == false {
							//download product
							product_id, _, product_store_item_id = fetch_product_item_to_db(shop_account_config.APIConfig, db_session, transaction.Item.ItemID)
						}

						order_item.ItemProductId = product_id
						order_item.ProductStoreItemId = product_store_item_id
						order_item.ItemSku = transaction.Item.SKU
						order_item.ItemTitle = transaction.Item.Title
						order_item.ItemConditionId = transaction.Item.ConditionID
						order_item.ItemConditionDisplayName = transaction.Item.ConditionDisplayName

						transaction_created_date, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", transaction.CreatedDate, time.UTC)
						core.HandleError(err)
						transaction_created_date = transaction_created_date.Local()
						order_item.CreatedDate = core.LongTime(transaction_created_date)

						// 新增或更新明细
						if has_order_item == false {
							//新增
							_, err = db_session.Insert(&order_item)
							core.HandleError(err)
						} else {
							//更新
							_, err = db_session.Id(order_item.Id).Update(&order_item)
							core.HandleError(err)
						}
					}
				}

				if order_download.ShippingAddress.AddressOwner != "" && order_download.ShippingAddress.AddressID != "" {
					//新增或更新送货地址
					var order_shipping_address models.OrderShippingAddressInfo
					has_order_shipping_address, err := db_session.Where(fmt.Sprintf("order_id = %d AND address_owner = '%s'", order_id, order_download.ShippingAddress.AddressOwner)).Get(&order_shipping_address) // 收货地址
					core.HandleError(err)

					//设置每个字段的值，存在的每个字段都要设置，不能漏，否则更新和新增的时候会出问题
					order_shipping_address.OrderId = order_id
					order_shipping_address.Name = order_download.ShippingAddress.Name
					order_shipping_address.Street1 = order_download.ShippingAddress.Street1
					order_shipping_address.Street2 = order_download.ShippingAddress.Street2
					order_shipping_address.Phone = order_shipping_address_service.Normalize_phone(order_download.ShippingAddress.Phone) // 整理电话号码
					order_shipping_address.CityName = order_download.ShippingAddress.CityName
					order_shipping_address.Province = order_download.ShippingAddress.StateOrProvince
					order_shipping_address.Country = order_download.ShippingAddress.Country
					order_shipping_address.CountryName = order_download.ShippingAddress.CountryName
					order_shipping_address.PostalCode = order_download.ShippingAddress.PostalCode
					order_shipping_address.AddressId = order_download.ShippingAddress.AddressID
					order_shipping_address.AddressOwner = order_download.ShippingAddress.AddressOwner
					order_shipping_address.ExternalAddressId = order_download.ShippingAddress.ExternalAddressID

					var address_id = order_shipping_address.Id
					if has_order_shipping_address == true {
						// do update
						_, err = db_session.Id(order_shipping_address.Id).Update(&order_shipping_address)
						core.HandleError(err)

						address_id = order_shipping_address.Id
					} else {
						// do insert
						_, err = db_session.Insert(&order_shipping_address)
						core.HandleError(err)

						address_id = order_shipping_address.Id
					}

					//回填地址id到订单里面去
					order_in_db.AddressId = address_id
					_, err = db_session.Id(order_in_db.Id).Cols("address_id").Update(&order_in_db)
					core.HandleError(err)
				}

			} // for _, order_download := range order_list

			has_more_orders = response.HasMoreOrders
			page_number += 1
		} // for has_more_orders == true
	} // for _, config := range ebay_api_config_list

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("eBay GetOrders time duration: ", job_duration.Seconds(), "seconds")
}
