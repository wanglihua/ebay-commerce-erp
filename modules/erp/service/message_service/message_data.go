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
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models/xxgl_models"
)

// todo: 不采用多表联接的方式，采用单表的方式来做，这个需要改造
// shop_id 设置的时候同时设置名称
// stauts_id 设置的时候设置状态名称
type Message_view_item_t struct {
	Message           xxgl_models.MessageInfo           `xorm:"extends" json:"msg"`
	ShopAccount       xtpz_models.ShopAccountInfo       `xorm:"extends" json:"sa"`
	ECommercePlatform xtpz_models.ECommercePlatformInfo `xorm:"extends" json:"pf"`
	Product           kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	ServiceUser       auth_models.UserInfo              `xorm:"extends" json:"psu"`
	ProductStoreItem  sppz_models.ProductStoreItemInfo  `xorm:"extends" json:"psi"`
	Status            models.StatusInfo                 `xorm:"extends" json:"s"`
	MsgLastModifyDate string                            `xorm:"-" json:"message_last_modify_date"`
	OrdNum            int64                             `xorm:"-"`
	LastModifiedDate  string                            `xorm:"-"`
}

var msg_count struct {
	Count int64 `xorm:"msg_count"`
}

/*
SELECT msg.*, sa.*, pf.*, p.*, psu.*, psi.*, s.*
FROM erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    LEFT JOIN auth_user psu ON psu.id = p.service_user_id
    LEFT JOIN erp_product_store_item psi ON psi.product_id = p.id AND psi.item_id = msg.item_id
    INNER JOIN erp_status s ON s.id = msg.status_id;
 */

func get_message_list_sql_builder() *core.SqlBuilder {
	sb := core.NewSqlBuilder()
	sb.Select("msg.*, sa.*, pf.*, p.*, psu.*, psi.*, s.*")
	sb.Count("count(msg.id) as msg_count")
	sb.From(`
erp_message msg
    INNER JOIN erp_shop_account sa ON sa.id = msg.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = msg.product_id
    LEFT JOIN auth_user psu ON psu.id = p.service_user_id
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

	var service_user_id_query int64
	params.Bind(&service_user_id_query, "service_user_id_query")

	var sale_user_id_query int64
	params.Bind(&sale_user_id_query, "sale_user_id_query")

	var dev_user_id_query int64
	params.Bind(&dev_user_id_query, "dev_user_id_query")

	var message_content_query string
	params.Bind(&message_content_query, "message_content_query")
	message_content_query = strings.TrimSpace(message_content_query)

	var query_begin_time string
	params.Bind(&query_begin_time, "query_begin_time")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	params.Bind(&query_end_time, "query_end_time")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

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

	if service_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.service_user_id = %d", service_user_id_query))
	}

	if sale_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.sale_user_id = %d", sale_user_id_query))
	}

	if dev_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.dev_user_id = %d", dev_user_id_query))
	}

	if message_content_query != "" {
		filter = core.SqlAnd(filter, core.SqlOr(
			fmt.Sprintf("msg.subject LIKE '%%%s%%'", message_content_query),
			fmt.Sprintf("msg.body LIKE '%%%s%%'", message_content_query)))
	}

	if query_begin_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.last_modified_date >= '%s'", query_begin_date.Format("2006-01-02 15:04")))
	}

	if query_end_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("msg.last_modified_date <= '%s'", query_end_date.Format("2006-01-02 15:04")))
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

	for index, message := range message_list {
		message_list[index].MsgLastModifyDate = message.Message.LastModifiedDate.Time.Format("2006-01-02 15:04:05")
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&msg_count)
	core.HandleError(err)

	return message_list, msg_count.Count
}

func Get_message_data(message_id int64, db_session *xorm.Session) Message_view_item_t {

	var filter = fmt.Sprintf("msg.id = %d", message_id)

	sb := get_message_list_sql_builder()
	sb.Where(filter)

	var message_view_item Message_view_item_t

	_, err := db_session.Sql(sb.SqlGet()).Get(&message_view_item)
	core.HandleError(err)

	message_view_item.MsgLastModifyDate = message_view_item.Message.LastModifiedDate.Time.Format("2006-01-02 15:04:05")

	return message_view_item
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

func Get_message_list_data_by_sender(db_session *xorm.Session, sender_id string) []Message_view_item_t {
	sb := get_message_list_sql_builder()
	sb.Where(fmt.Sprintf("msg.sender_id = '%s'", sender_id))

	var message_list = make([]Message_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	for index, message := range message_list {
		message_list[index].OrdNum = int64(index + 1)
		message_list[index].LastModifiedDate = message.Message.LastModifiedDate.Time.Format("2006-01-02 15:04:05")
	}

	return message_list
}

func Get_message_list_data_by_sender_shop_id(db_session *xorm.Session, sender_id string, shop_id int64) []Message_view_item_t {
	sb := get_message_list_sql_builder()
	sb.Where(fmt.Sprintf("msg.sender_id = '%s' AND msg.shop_id = %d", sender_id, shop_id))

	var message_list = make([]Message_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&message_list)
	core.HandleError(err)

	for index, message := range message_list {
		message_list[index].OrdNum = int64(index + 1)
		message_list[index].LastModifiedDate = message.Message.LastModifiedDate.Time.Format("2006-01-02 15:04:05")
	}

	return message_list
}
