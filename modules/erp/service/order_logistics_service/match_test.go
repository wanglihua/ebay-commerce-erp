package order_logistics_service

import (
	"testing"
	"fmt"
)

func TestMatchResult(t *testing.T) {
	var final_result match_logistics_mode_result
	var final_result_second match_logistics_mode_result

	final_result.stock_id = 1
	final_result_second = final_result

	fmt.Println(final_result_second.stock_id)
}
