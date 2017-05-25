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
