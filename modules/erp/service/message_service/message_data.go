package message_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/sppz_models"
	"strings"
)

type Message_view_item_t struct {
	models.MessageInfo                `xorm:"extends" json:"msg"`
	xtpz_models.ShopAccountInfo       `xorm:"extends" json:"sa"`
	xtpz_models.ECommercePlatformInfo `xorm:"extends" json:"pf"`
	kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	sppz_models.ProductStoreItemInfo  `xorm:"extends" json:"psi"`
	models.StatusInfo                 `xorm:"extends" json:"s"`
}

var msg_count struct {
	Count int64 `xorm:"msg_count"`
}

/*
SELECT msg.*, sa.*, pf.*, p.*, psi.*, s.*
FROM erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    LEFT JOIN erp_product_store_item psi on psi.product_id = p.id and psi.item_id = msg.item_id
    INNER JOIN erp_status s ON s.id = msg.status_id;
 */

func get_message_list_sql_builder() *core.SqlBuilder {
	sb := core.NewSqlBuilder()
	sb.Select("msg.*, sa.*, pf.*, p.*, psi.*, s.*")
	sb.Count("count(msg.id) as msg_count")
	sb.From(`
erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    LEFT JOIN erp_product_store_item psi ON psi.product_id = p.id AND psi.item_id = msg.item_id
    INNER JOIN erp_status s ON s.id = msg.status_id`)

	sb.OrderBy("msg.id DESC")

	return sb
}

func Get_message_list_data(db_session *xorm.Session, request *revel.Request, params *revel.Params) (message_list []Message_view_item_t, count int64) {

	filter, order, offset, limit := core.GetGridRequestParam(request)

	var shop_id_query int64
	params.Bind(&shop_id_query, "shop_id_query")

	var item_id_query string
	params.Bind(&item_id_query, "item_id_query")

	var sku_query string
	params.Bind(&sku_query, "sku_query")

	var sender_id_query string
	params.Bind(&sender_id_query, "sender_id_query")

	var recipient_id_query string
	params.Bind(&recipient_id_query, "recipient_id_query")

	if shop_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.shop_id = %d", shop_id_query))
	}

	if strings.TrimSpace(item_id_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.item_id = '%s'", item_id_query))
	}

	if strings.TrimSpace(sku_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.sku = '%s'", sku_query))
	}

	if strings.TrimSpace(sender_id_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.sender_id = '%s'", sender_id_query))
	}

	if strings.TrimSpace(recipient_id_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.recipient_id = '%s'", recipient_id_query))
	}

	sb := get_message_list_sql_builder()

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("msg.id DESC")
	}

	sb.Limit(offset, limit)

	message_list = make([]Message_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	_, err = db_session.Sql(sb.SqlCount()).Get(&msg_count)
	core.HandleError(err)

	return message_list, msg_count.Count
}

func Get_message_list_data_by_product_id(db_session *xorm.Session, product_id int64) []Message_view_item_t {
	sb := get_message_list_sql_builder()
	sb.Where(fmt.Sprintf("msg.product_id = %d", product_id))

	var message_list = make([]Message_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	return message_list
}

func Get_message_list_data_by_product_id_top_n(db_session *xorm.Session, product_id int64, top_n int64) []Message_view_item_t {
	sb := get_message_list_sql_builder()
	sb.Where(fmt.Sprintf("msg.product_id = %d", product_id))
	sb.Limit(0, top_n)

	var message_list = make([]Message_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	return message_list
}
