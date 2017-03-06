package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/franela/goreq"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/ebay"
	"matrix/modules/erp/models"
)

/*

get orders document:
http://developer.ebay.com/devzone/xml/docs/Reference/eBay/GetOrders.html

test tools:
http://developer.ebay.com/DevZone/build-test/test-tool/?index=0&api=trading&call=GetOrders&variation=xml#

request header:
X-EBAY-API-CALL-NAME:GetOrders
X-EBAY-API-DEV-NAME:***REMOVED***
X-EBAY-API-CERT-NAME:***REMOVED***
X-EBAY-API-APP-NAME:***REMOVED***
X-EBAY-API-VERSION:967
X-EBAY-API-COMPATIBILITY-LEVEL:967
X-EBAY-API-REQUEST-ENCODING:XML
X-EBAY-API-RESPONSE-ENCODING:XML
X-EBAY-API-SITEID:0
Accept:application/xml,application/xhtml+xml

request body:
<?xml version="1.0" encoding="utf-8"?>
<GetOrdersRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**uVWPWA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GiD5aDpAydj6x9nY+seQ**pRMEAA**AAMAAA**VrbSv0o7iGoYiFvprWmTAsvp1TYnEDq74jjQOsqLWlwk5eS1y2Uyrq9aYWuWdWvRuM7/L0nUP1IFac3aBTY2aanmI+37H8fN0mUrSY7gFIwEdGx1pneZFwweLzc3l0yZcV1m22tWrSqRMQ6g+A31y8vxtjGHkSZkbeMJQJhMHLo5gMOjo2veZNacluDlzL2hKrs+Z7iIPuqxqwYEz2hDGe36BeXdBY5sDUvbZBdRVV4ws4Ta5L5Bal5tJO8tH/8WO0fXu5OKlw/BIGpgtbQc8Xe5MmhotfbznjfnIgrzxoSbsky6HZqpQ7AQh3UBZ/T7i2H5BdD9Rr3sMeXas6x7E8V0xsDqfQL1N8qAfKCOhDVP8Vxg/AOxP3dVMQaqaX0gBRp8CTnMIlYvHJKzW9URZ5vipjAJ+Fcx/IW5sHdiJFGHawHlfi0yPS3cYecLcPBUAKYqr1fBYqEZAQGsqId/PMVue8hwbl/gzaLus3WAOUPME4hFXMu2n9ADkYturlMa2TYRyHufNJmziHHUUb68bALZpTMhwH6QuG4ILSIluVDz9YcM7GuRGtQzBDhbTRdim6WGcjO7dkogo665c8QqSBIW2kT1+TTT24gXs8dUkv6Zea3Ihiko/bgh6mm2vCJ+LtMNqk+/vFHhL4nMhToDGX+0fFxu2U43aq3FJGo5GEYbn6MowwY/JzRHFKDkN2wzVKo7DJ2c1+faulxDGyVGhOMeEmfdZgL3iEd24JPm0mcvT1xA/UXpDYyCdhRb6I7+</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<CreateTimeFrom>2017-01-01T00:00:00</CreateTimeFrom>
	<CreateTimeTo>2017-02-01T23:59:59</CreateTimeTo>
	<OrderRole>Seller</OrderRole>
	<OrderStatus>Active</OrderStatus>
</GetOrdersRequest>

response body:
<?xml version="1.0" encoding="UTF-8"?>
<GetOrdersResponse xmlns="urn:ebay:apis:eBLBaseComponents">
	<Timestamp>2017-01-31T12:57:53.266Z</Timestamp>
	<Ack>Success</Ack>
	<Version>989</Version>
	<Build>E989_CORE_API_18131074_R1</Build>
	<PaginationResult>
		<TotalNumberOfPages>0</TotalNumberOfPages>
		<TotalNumberOfEntries>0</TotalNumberOfEntries>
	</PaginationResult>
	<HasMoreOrders>false</HasMoreOrders>
	<OrderArray/>
	<OrdersPerPage>100</OrdersPerPage>
	<PageNumber>1</PageNumber>
	<ReturnedOrderCountActual>0</ReturnedOrderCountActual>
</GetOrdersResponse>

 */

type EBayOrderDownloadTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t EBayOrderDownloadTask) Run() {
	job_time_begin := time.Now()

	// get db driver and db url

	if t.db_driver == "" {
		t.db_driver, _ = revel.Config.String("db.driver")
	}

	if t.db_url == "" {
		t.db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(t.db_driver, t.db_url)
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	ebay_api_config_list := make([]models.EBayAPIConfigInfo, 0)

	/*
	SELECT cfg.* FROM erp_ebay_api_config cfg INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id
	*/
	err = db_session.Sql("SELECT cfg.* FROM erp_ebay_api_config cfg INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id").Find(&ebay_api_config_list)
	if err != nil {
		revel.ERROR.Println(err)
	}

	for _, config := range ebay_api_config_list {
		has_more_orders := true
		page_number := 1
		for has_more_orders == true {
			response := t.fetch_orders(config, page_number)

			order_list := response.OrdersArray.Orders

			for _, order_download := range order_list {
				order_code := order_download.OrderID

				var order_in_db models.OrderInfo
				has_order, err := db_session.Where("code = ?", order_code).Get(&order_in_db)
				if err != nil {
					revel.ERROR.Println(err)
				}

				order_in_db.Code = order_code
				order_in_db.ShopId = config.ShopId

				var shop models.ShopAccountInfo
				_, err = db_session.Id(config.ShopId).Get(&shop)
				if err != nil {
					revel.ERROR.Println(err)
				}

				// 老鼠仓的判断
				// 如果是通过PayPal支付的，而且店铺设置里面也设置了PayPal列表的
				if strings.ToLower(order_download.PaymentMethods) == strings.ToLower("PayPal") && shop.PaypalAccount != "" {
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
						order_in_db.HoldTypeId = core.NewNullInt(models.ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH, true)
					}
				}

				// 状态？
				revel.TRACE.Println("OrderStatus: ", order_download.OrderStatus)

				order_in_db.TypeId = models.ORDER_TYPE_XZDD //订单类型为下载订单

				update_col_list := []string{"type_id", "ebay_status", "status_id"}

				if order_download.CreatedTime != "" {
					created_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.CreatedTime)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					//created_time = created_time.Add(8 * time.Hour)
					created_time = created_time.Local()
					order_in_db.CreatedTime = core.NewNullLongTime(created_time, true)
					update_col_list = append(update_col_list, "created_time")
				}

				if order_download.ModifiedTime != "" {
					modified_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.ModifiedTime)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					// modified_time = modified_time.Add(8 * time.Hour)
					modified_time = modified_time.Local()
					order_in_db.ModifiedTime = core.NewNullLongTime(modified_time, true)
					update_col_list = append(update_col_list, "modified_time")
				}

				if order_download.PaidTime != "" {
					paid_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.PaidTime)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					// paid_time = paid_time.Add(8 * time.Hour)\
					paid_time = paid_time.Local()
					order_in_db.PaidTime = core.NewNullLongTime(paid_time, true)
					update_col_list = append(update_col_list, "paid_time")
				}

				if order_download.ShippedTime != "" {
					shipped_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.ShippedTime)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					// shipped_time = shipped_time.Add(8 * time.Hour)
					shipped_time = shipped_time.Local()
					order_in_db.ShippedTime = core.NewNullLongTime(shipped_time, true)
					update_col_list = append(update_col_list, "shipped_time")
				}

				order_in_db.AdjustmentAmountCurrency = order_download.AdjustmentAmount.Currency
				update_col_list = append(update_col_list, "adjustment_amount_currency")
				order_in_db.AdjustmentAmount = order_download.AdjustmentAmount.Value
				update_col_list = append(update_col_list, "adjustment_amount")

				order_in_db.AmountPaidCurrency = order_download.AmountPaid.Currency
				update_col_list = append(update_col_list, "amount_paid_currency")
				order_in_db.AmountPaid = order_download.AmountPaid.Value
				update_col_list = append(update_col_list, "amount_paid")

				order_in_db.AmountSavedCurrency = order_download.AmountSaved.Currency
				update_col_list = append(update_col_list, "amount_saved_currency")
				order_in_db.AmountSaved = order_download.AmountSaved.Value
				update_col_list = append(update_col_list, "amount_saved")

				//Total
				order_in_db.TotalAmountCurrency = order_download.Total.Currency
				update_col_list = append(update_col_list, "total_amount_currency")
				order_in_db.TotalAmount = order_download.Total.Value
				update_col_list = append(update_col_list, "total_amount")

				//SubTotal
				order_in_db.SubTotalAmountCurrency = order_download.Subtotal.Currency
				update_col_list = append(update_col_list, "subtotal_amount_currency")
				order_in_db.SubTotalAmount = order_download.Subtotal.Value
				update_col_list = append(update_col_list, "subtotal_amount")

				order_in_db.PaymentMethods = order_download.PaymentMethods
				order_in_db.SellerEmail = order_download.SellerEmail
				order_in_db.SellerUserId = order_download.SellerUserID

				order_in_db.BuyerUserId = order_download.BuyerUserID
				order_in_db.BuyerCheckoutMessage = order_download.BuyerCheckoutMessage

				// 订单状态，这部分逻辑还可以继续探讨

				// 订单状态 设置 开始
				order_status := order_download.OrderStatus
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

				// 如果是新的订单，一般是新增的订单，系统数据库中还没有的订单
				if order_status_id < models.ORDER_STATUS_YCJ && order_in_db.PaidTime.Valid == false { //前面小于号的判断要不要加呢
					order_in_db.StatusId = models.ORDER_STATUS_YCJ // 这个状态，有可能会被接下来的状态逻辑所修改，可能会，更进一步改成 “已支付”, “已完成”
				}

				// 有可能是新增的订单（最有可能），或者是原来 已创建 的状态过来的。这时支付时间不为空
				if order_status_id < models.ORDER_STATUS_YZF && order_in_db.PaidTime.Valid == true {
					order_in_db.StatusId = models.ORDER_STATUS_YZF // 这个状态，有可能会被接下来的状态逻辑所修改，可能会，更进一步改成 “已完成”
				}

				// 同步逻辑要保证中间的 待配货，配货中，已发货 等的状态在同步时不被前面的状态覆盖，只有后面的“已完成”可以覆盖这些状态

				// 如果 下载订单里面的 OrderStatus 是 Completed， 且订单的发货时间不为空，订单状态就为 已完成
				if strings.ToLower(order_status) == strings.ToLower("Completed") && order_in_db.ShippedTime.Valid == true {
					order_in_db.StatusId = models.ORDER_STATUS_YWC
				}

				//订单状态 设置 结束

				if has_order == true {
					// update order
					_, err = db_session.Id(order_in_db.Id).Cols(update_col_list...).Update(&order_in_db)
					if err != nil {
						revel.ERROR.Println(err)
					}
				} else {
					// insert order
					order_in_db.Code = order_code

					_, err = db_session.Insert(&order_in_db)
					if err != nil {
						revel.ERROR.Println(err)
					}
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
					ItemId                   string          `xorm:"nvarchar(200) notnull 'item_id'" json:"item_id"`
					ItemProductId            int64           `xorm:"bigint notnull index 'item_product_id'" json:"item_product_id"`
					ItemSku                  string          `xorm:"nvarchar(100) notnull index 'item_sku'" json:"item_sku"`
					ItemTitle                string          `xorm:"nvarchar(200) notnull 'item_title'" json:"item_title"`
					ItemConditionId          string          `xorm:"nvarchar(50) notnull 'item_condition_id'" json:"item_condition_id"`
					ItemConditionDisplayName string          `xorm:"nvarchar(100) notnull 'item_condition_display_name'" json:"item_condition_display_name"`
					CreatedDate              core.LongTime   `xorm:"datetime notnull 'created_date'" json:"created_date"`
					Desc                     core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
				 */

				for _, transaction := range order_download.TransactionArray.Transactions {
					// 获取并检查明细是否存在，根据 order_id 和 item_sku
					var order_item models.OrderItemInfo
					item_sku := transaction.Item.SKU
					has_order_item, err := db_session.Where(fmt.Sprintf("order_id = %d and item_sku = '%s'", order_id, item_sku)).Get(&order_item)
					if err != nil {
						revel.ERROR.Println(err)
					}

					// 设置明细字段
					order_item.OrderId = order_id
					order_item.OrderLineItemId = transaction.OrderLineItemID
					order_item.Quantity = transaction.QuantityPurchased
					order_item.TransactionId = transaction.TransactionID
					order_item.TransactionPriceCurrency = transaction.TransactionPrice.Currency
					order_item.TransactionPrice = transaction.TransactionPrice.Value
					order_item.ItemId = transaction.Item.ItemID

					//order_item.ItemProductId 如果这时产品不存在会怎样，下载产品
					var product_id struct {
						ProductId int64 `xorm:"'product_id'"`
					}
					_, err = db_session.Sql(fmt.Sprintf("SELECT p.id AS product_id FROM erp_product p WHERE p.sku = '%s'", transaction.Item.SKU)).Get(&product_id)
					if err != nil {
						revel.ERROR.Println(err)
					}

					order_item.ItemSku = transaction.Item.SKU
					order_item.ItemTitle = transaction.Item.Title
					order_item.ItemConditionId = transaction.Item.ConditionID
					order_item.ItemConditionDisplayName = transaction.Item.ConditionDisplayName

					transaction_created_date, err := time.Parse("2006-01-02T15:04:05.000Z", transaction.CreatedDate)
					if err != nil {
						revel.ERROR.Println(err)
					}
					transaction_created_date = transaction_created_date.Local()
					order_item.CreatedDate = core.LongTime(transaction_created_date)

					// 新增或更新明细
					if has_order_item == false {
						//新增
						_, err = db_session.Insert(&order_item)
						if err != nil {
							revel.ERROR.Println(err)
						}
					} else {
						//更新
						_, err = db_session.Id(order_item.Id).Update(&order_item)
						if err != nil {
							revel.ERROR.Println(err)
						}
					}
				}

				//新增或更新送货地址
				var order_shipping_address models.OrderShippingAddressInfo
				has_order_shipping_address, err := db_session.Where("order_id = ?", order_id).Get(&order_shipping_address)

				//设置每个字段的值，存在的每个字段都要设置，不能漏，否则更新和新增的时候会出问题
				order_shipping_address.OrderId = order_id
				order_shipping_address.Name = order_download.ShippingAddress.Name
				order_shipping_address.Street1 = order_download.ShippingAddress.Street1
				order_shipping_address.Street2 = order_download.ShippingAddress.Street2
				order_shipping_address.CityName = order_download.ShippingAddress.CityName
				order_shipping_address.Province = order_download.ShippingAddress.StateOrProvince
				order_shipping_address.Country = order_download.ShippingAddress.Country
				order_shipping_address.CountryName = order_download.ShippingAddress.CountryName
				order_shipping_address.PostalCode = order_download.ShippingAddress.PostalCode
				order_shipping_address.AddressId = order_download.ShippingAddress.AddressID
				order_shipping_address.AddressOwner = order_download.ShippingAddress.AddressOwner
				order_shipping_address.ExternalAddressId = order_download.ShippingAddress.ExternalAddressID

				if has_order_shipping_address == true {
					// do update
					_, err = db_session.Id(order_shipping_address.Id).Update(&order_shipping_address)
					if err != nil {
						revel.ERROR.Println(err)
					}
				} else {
					// do insert
					_, err = db_session.Insert(&order_shipping_address)
					if err != nil {
						revel.ERROR.Println(err)
					}
				}

				// Order Item

			} // for _, order_download := range order_list

			switch response.Ack {
			case "Success":
				revel.TRACE.Printf("eBay GetOrders Response: " + response.Ack + "\n")
				revel.TRACE.Printf("eBay GetOrders actual fetched orders count: %d", response.ReturnedOrderCountActual)
			case "Failure":
				revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
			default:
				revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
			}

			has_more_orders = response.HasMoreOrders
			page_number += 1
		} // for has_more_orders == true
	} // for _, config := range ebay_api_config_list

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("sync order time duration: ", job_duration.Seconds(), "seconds")
}

func (t EBayOrderDownloadTask) fetch_orders(config models.EBayAPIConfigInfo, page_number int) ebay.GetOrdersResponse {

	number_of_days := 1

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetOrdersRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<NumberOfDays>%d</NumberOfDays>
	<OrderRole>Seller</OrderRole>
	<OrderStatus>All</OrderStatus>
	<Pagination>
		<EntriesPerPage>100</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetOrdersRequest>
`
	req_body = fmt.Sprintf(req_body, config.Token, number_of_days, page_number)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetOrders")
	req.AddHeader("X-EBAY-API-DEV-NAME", config.DevId)
	req.AddHeader("X-EBAY-API-CERT-NAME", config.CertId)
	req.AddHeader("X-EBAY-API-APP-NAME", config.AppId)
	req.AddHeader("X-EBAY-API-VERSION", config.GetApiVersion())
	req.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", config.GetApiVersion())
	req.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-SITEID", config.GetSiteId())
	req.AddHeader("Accept", "application/xml,application/xhtml+xml")

	res, err := req.Do()

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - processUrl -> req.Do: " + err.Error())
		return ebay.GetOrdersResponse{}
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - ioutil.ReadAll : " + err.Error())
		return ebay.GetOrdersResponse{}
	}

	if revel.DevMode {
		//core.Force_write_to_file(fmt.Sprintf(`z:\test_order_response\%s.xml`, time.Now().Format("20060102150405")), string(data))
		core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_order_response/%s.xml`, time.Now().Format("20060102150405")), string(data))
	}

	response := ebay.GetOrdersResponse{}
	err = xml.Unmarshal(data, &response)
	if err != nil {
		revel.ERROR.Println(err)
	}

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetOrders ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetOrdersResponse{}
	}

	return response
}

/*

<ShippingAddress>
	<Name>Som Chanthasenh</Name>
	<Street1>23 Graziers Way</Street1>
	<Street2/>
	<CityName>Carnes Hill</CityName>
	<StateOrProvince>New South Wales</StateOrProvince>
	<Country>AU</Country>
	<CountryName>Australia</CountryName>
	<Phone>04 19292611</Phone>
	<PostalCode>2171</PostalCode>
	<AddressID>1255435139019</AddressID>
	<AddressOwner>eBay</AddressOwner>
	<ExternalAddressID/>
</ShippingAddress>
<ShippingServiceSelected>
	<ShippingService>AU_Regular</ShippingService>
	<ShippingServiceCost currencyID="AUD">0.0</ShippingServiceCost>
</ShippingServiceSelected>

 */
