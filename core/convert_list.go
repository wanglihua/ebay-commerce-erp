package core

import (
	"fmt"
	"strconv"
)

func Int_list_to_str_list(int_list []int64) []string {
	var str_list = make([]string, 0)
	for _, int_ := range int_list {
		str_list = append(str_list, fmt.Sprint(int_))
	}
	return str_list
}

func Str_list_to_int_list(str_list []string) []int64 {
	var int_list = make([]int64, 0)
	for _, str := range str_list {
		int_, err := strconv.ParseInt(str, 10, 64)
		HandleError(err)
		/*
		if (err != nil) {
			revel.ERROR.Println(err)
		}
		*/
		int_list = append(int_list, int_)
	}
	return int_list
}
