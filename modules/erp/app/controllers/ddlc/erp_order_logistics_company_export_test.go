package controllers

import (
	"fmt"
	"testing"
)

func Test_excel_row_label_to_index(t *testing.T) {
	var row_index = excel_col_label_to_index("a")
	if row_index != 0 {
		t.Error("结果不对！")
	}

	row_index = excel_col_label_to_index("C")
	if row_index != 2 {
		t.Error("结果不对！")
	}

	row_index = excel_col_label_to_index("AC")
	if row_index != 28 {
		fmt.Println(row_index)
		t.Error("结果不对！")
	}

	row_index = excel_col_label_to_index("ba")
	if row_index != 52 {
		fmt.Println(row_index)
		t.Error("结果不对！")
	}
}

func Test_split_address(t *testing.T) {
	address1, address2 := split_address("hello hello hello hello hello world abc hello hello hello hello hello hello hello hello hello hello hello", 35)
	fmt.Println("address1", address1)
	fmt.Println("address2", address2)
}


func  Test_get_state_short_name(t *testing.T) {
	var short_name = get_state_short_name("Nsw")
	fmt.Println("short_name", short_name)
}