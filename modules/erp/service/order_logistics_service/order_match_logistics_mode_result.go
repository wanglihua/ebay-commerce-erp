package order_logistics_service

import (
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
)

type match_logistics_mode_result struct {
	order_id      int64
	order_code    string
	address_id    int64
	address       ddlc_models.OrderShippingAddressInfo
	company_id    int64
	company       ddpz_models.LogisticsCompanyInfo
	stock_id      int64
	stock         sppz_models.StockInfo
	mode_id       int64
	mode          ddpz_models.LogisticsModeInfo
	fee_rate_id   int64
	fee_rate      ddpz_models.LogisticsFeeRateInfo
	fee_rate2_id  int64
	fee_rate2     ddpz_models.LogisticsFeeRate2Info
	length        float64
	width         float64
	height        float64
	actual_weight float64
	dim_weight    float64
	sfkltjz       int64
	weight        float64 //计算时用的重量体积重或者实际重量, kg
	total_fee     float64
	pick_fee      float64
	cost_amount   float64
	total_amount  float64
	has_stock     bool
}

//排序相关的函数实现

func (r match_logistics_mode_result) can_deliver() bool {
	if r.fee_rate_id == 0 && r.fee_rate2_id == 0 {
		return false
	} else {
		return true
	}
}

type match_logistics_mode_result_slice []match_logistics_mode_result

func (s match_logistics_mode_result_slice) Len() int {
	return len(s)
}

func (s match_logistics_mode_result_slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 费用少的排在最前面
func (s match_logistics_mode_result_slice) Less(i, j int) bool {
	return s[i].total_fee < s[j].total_fee
}

/*

package main

import (
    "fmt"
    "sort"
)

type Person struct {
    Name string
    Age  int
}

// 按照 Person.Age 从大到小排序
type PersonSlice [] Person

func (a PersonSlice) Len() int {         // 重写 Len() 方法
    return len(a)
}
func (a PersonSlice) Swap(i, j int){     // 重写 Swap() 方法
    a[i], a[j] = a[j], a[i]
}
func (a PersonSlice) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
    return a[j].Age < a[i].Age
}

func main() {
    people := [] Person{
        {"zhang san", 12},
        {"li si", 30},
        {"wang wu", 52},
        {"zhao liu", 26},
    }

    fmt.Println(people)

    sort.Sort(PersonSlice(people))    // 按照 Age 的逆序排序
    fmt.Println(people)

    sort.Sort(sort.Reverse(PersonSlice(people)))    // 按照 Age 的升序排序
    fmt.Println(people)

}
 */
