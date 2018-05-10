package order_shipping_address_service

import "strings"

// 整理电话号码，使固定的格式
func Normalize_phone(phone string) string {
	if phone != "" {
		phone = strings.Replace(phone, " ", "", -1)
		var phone_in_rune = []rune(phone)
		if len(phone_in_rune) > 4 {
			return string(phone_in_rune[0:4]) + " " + string(phone_in_rune[4:])
		}
		/*
		if len(phone) > 4 {
			return phone[0:4] + " " + phone[4:]
		}
		*/
	}

	return phone
}
