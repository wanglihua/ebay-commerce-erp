package order_feedback_service

import "strings"

func feedback_type_to_type_str(feedback_type int64) string {
	//  // 1: negative, 2: neutral, 3: positive

	var order_feedback_type_str = ""
	if feedback_type == ORDER_FEEDBACK_NEGATIVE {
		order_feedback_type_str = strings.ToLower("negative")
	} else if feedback_type == ORDER_FEEDBACK_NEUTRAL {
		order_feedback_type_str = strings.ToLower("neutral")
	} else if feedback_type == ORDER_FEEDBACK_POSITIVE {
		order_feedback_type_str = strings.ToLower("positive")
	}

	return order_feedback_type_str
}
