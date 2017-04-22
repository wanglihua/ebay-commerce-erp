package product_bind_service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

// 暂时先写在这里
// bind_produt_list: bind_product_expression
func Check_product_bind_valid(bind_product_list string, db_session *xorm.Session) (result bool, err error) {
	if strings.Index(bind_product_list, "，") != -1 { //中文的逗号，看看有没有输入错误
		return false, errors.New("多个商品时不能用中文的逗号（，）分割。")
	}

	product_sku_with_number_list := strings.Split(bind_product_list, ",")
	for _, product_sku_with_number := range product_sku_with_number_list {
		product_sku_number_slice := strings.Split(product_sku_with_number, "*")
		if len(product_sku_number_slice) != 2 {
			return false, errors.New("请用 SKU*数量 表示商品的数量，单个数量就用 *1 的方式！")
		}

		product_sku := strings.TrimSpace(product_sku_number_slice[0]);
		// 检查SKU的正确性
		if check_product_sku(product_sku, db_session) == false {
			return false, errors.New("请用 SKU*数量 表示商品的数量，商品SKU不存在！")
		}

		product_number_str := strings.TrimSpace(product_sku_number_slice[1]);
		// 检查数量的格式
		_, err = strconv.ParseInt(product_number_str, 10, 64)
		if err != nil {
			return false, errors.New("请用 SKU*数量 表示商品的数量，数量应为正整数！")
		}
	}

	return true, nil
}

// 是否是绑定商品的sku，就是虚拟的那个sku
func Is_product_bind_sku(product_bind_sku string, db_session *xorm.Session) bool {
	sql := `SELECT count(pb.id) AS product_bind_count FROM erp_product_bind pb WHERE pb.code = '%s'`
	sql = fmt.Sprintf(sql, product_bind_sku)

	var product_bind_count struct {
		Count int64 `xorm:"'product_bind_count'"`
	}
	_, err := db_session.Sql(sql).Get(&product_bind_count)
	core.HandleError(err)

	if product_bind_count.Count == 0 {
		return false
	}

	return true
}

type Product_bind_sku_t struct {
	Sku   string
	Count int64
}

// 解析绑定 sku， product_bind_code 就是绑定商品的虚拟 sku
func Get_product_bind_skus(product_bind_code string, db_session *xorm.Session) []Product_bind_sku_t {
	sql := `SELECT pb.bind_list FROM erp_product_bind pb WHERE pb.code = '%s'`
	sql = fmt.Sprint(sql, product_bind_code)

	var product_bind_list struct {
		BindList string  `xorm:"'bind_list'"`
	}

	_, err := db_session.Sql(sql).Get(&product_bind_list)
	core.HandleError(err)

	var product_bind_sku_list = make([]Product_bind_sku_t, 0)

	product_sku_with_number_list := strings.Split(product_bind_list.BindList, ",")
	for _, product_sku_with_number := range product_sku_with_number_list {
		product_sku_number_slice := strings.Split(product_sku_with_number, "*")

		product_sku := strings.TrimSpace(product_sku_number_slice[0]);
		product_number_str := strings.TrimSpace(product_sku_number_slice[1]);
		product_number, _ := strconv.ParseInt(product_number_str, 10, 64)

		var product_bind_sku Product_bind_sku_t
		product_bind_sku.Sku = product_sku
		product_bind_sku.Count = product_number

		product_bind_sku_list = append(product_bind_sku_list, product_bind_sku)
	}

	return nil
}

// 检查SKU的正确性，是否是合法的数据库中存在的SKU
func check_product_sku(product_sku string, db_session *xorm.Session) bool {

	sql := `SELECT count(p.id) AS product_count FROM erp_product p WHERE p.sku = '%s'`
	sql = fmt.Sprintf(sql, product_sku)
	var product_count struct {
		Count int64 `xorm:"'product_count'"`
	}
	_, err := db_session.Sql(sql).Get(&product_count)
	core.HandleError(err)

	if product_count.Count == 0 {
		return false
	}

	return true
}
