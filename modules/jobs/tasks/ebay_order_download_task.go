package tasks


import (
	"github.com/revel/revel"
	"github.com/mithereal/go-ebay"
	"strconv"
)

type EBayOrderDownloadTask struct {
}

func (t EBayOrderDownloadTask) Run() {
	config := GetEBayConfig()

	credentials := ebay.RequesterCredentials{
		RequestToken: config.Token,
	}

	orders_request := ebay.GetOrdersRequest{
		CreateTimeFrom:       "2017-01-20T00:00:00",
		CreateTimeTo:         "2017-01-29T23:59:59",
		ModTimeFrom:          "2017-01-20T00:00:00",
		ModTimeTo:            "2017-01-29T23:59:59",
		DetailLevel:          "ReturnAll",
		OrderRole:            "Seller",
		IncludeFinalValueFee: "true",
		OrderStatus:          "Completed",
		ListingType:          "Unknown",
		NumberOfDays:         "9",
		SortingOrder:         "Ascending",
		// Credentials are used for custom config data store
		RequesterCredentials: credentials,
		WarningLevel:         "Low",
		Pagination: ebay.Pagination{
			EntriesPerPage:       "10",
			OrdersPerPage:        "10",
			PageNumber:           "10",
			TotalNumberOfEntries: 10,
			TotalNumberOfPages:   10,
		},
	}

	response := orders_request.FetchOrders(config)

	//response.OrdersArray

	//进行订单数据库操作等

	// do something with the orders //

	switch response.Ack {
	case "Success":
		revel.TRACE.Printf("Response: " + response.Ack + "\n")

		orderactualcount := strconv.Itoa(response.ReturnedOrderCountActual)
		pagecount := strconv.Itoa(response.Paginations.TotalNumberOfPages)
		pagenumber := strconv.Itoa(response.PageNumber)
		totalentries := strconv.Itoa(response.Paginations.TotalNumberOfEntries)
		hasMoreOrders := strconv.FormatBool(response.HasMoreOrders)

		revel.TRACE.Printf("Total Pages: " + pagecount + "\n")
		revel.TRACE.Printf("Current Page: " + pagenumber + "\n")
		revel.TRACE.Printf("Page Entries: " + orderactualcount + "\n")
		revel.TRACE.Printf("Total Orders: " + totalentries + "\n")
		revel.TRACE.Printf("Has more orders: " + hasMoreOrders + "\n")
	case "Failure":
		revel.ERROR.Printf("Response: " + response.Ack + "\n")
	default:
		revel.ERROR.Printf("Response: " + response.Ack + "\n")
	}
}
