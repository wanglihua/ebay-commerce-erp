package service

import "fmt"

func calculate_order_item_list_extra_field(message_order_item_list []MessageOrderItem) {
	if message_order_item_list == nil || len(message_order_item_list) == 0 {
		return
	}

	for index, message_order_item := range message_order_item_list {
		message_order_item_list[index].OrdNum = fmt.Sprint(index + 1)
		message_order_item_list[index].Price = fmt.Sprintf("%0.2f", message_order_item.OrderItem.TransactionPrice)
		message_order_item_list[index].Amount = fmt.Sprintf("%0.2f", message_order_item.OrderItem.TransactionPrice*float64(message_order_item.OrderItem.Quantity))

		if message_order_item.Order.PaidTime.Valid == true {
			message_order_item_list[index].PaidTime = message_order_item.Order.PaidTime.Time.Format("2006-01-02 15:04")
		}

		if message_order_item.Order.ShippedTime.Valid == true {
			message_order_item_list[index].ShippedTime = message_order_item.Order.ShippedTime.Time.Format("2006-01-02 15:04")
		}
	}
}
