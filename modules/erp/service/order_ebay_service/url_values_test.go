package order_ebay_service

import (
	"fmt"
	"testing"
	"net/url"
)

func TestUrlValuesAdd(t *testing.T) {

	var order_id_list []int64 = []int64{1, 2, 3, 4}

	var form_data url.Values = make(url.Values)

	for _, order_id := range order_id_list {
		form_data.Add("orderIdList[]", fmt.Sprint(order_id))
	}

	fmt.Println(form_data.Encode())

}
