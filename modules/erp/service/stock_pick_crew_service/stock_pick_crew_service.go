package stock_pick_crew_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
)

// 用户是不是墨尔本仓库的拣货人员
func Is_pick_crew_has_mel_stock(user_id int64, db_session *xorm.Session) bool {
	var stock_pick_crew_count struct {
		Count int64 `xorm:"'stock_pick_crew_count'"`
	}
	/*
SELECT count(spc.id) AS stock_pick_crew_count
FROM erp_stock_pick_crew spc
WHERE spc.user_id = 0
      AND spc.stock_id = (SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%墨尔本%')
	 */
	var sql_is_mel_pick_screw = `
SELECT count(spc.id) AS stock_pick_crew_count
FROM erp_stock_pick_crew spc
WHERE spc.user_id = %d
      AND spc.stock_id = (SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%%墨尔本%%')`
	sql_is_mel_pick_screw = fmt.Sprintf(sql_is_mel_pick_screw, user_id)

	_, err := db_session.Sql(sql_is_mel_pick_screw).Get(&stock_pick_crew_count)
	core.HandleError(err)

	if stock_pick_crew_count.Count == 0 {
		return false
	} else {
		return true
	}
}

func Get_login_user_stock_id_str_list(db_session *xorm.Session, web_session revel.Session) []string {
	// 不同的仓库拣货人员可以看到不同的订单，管理员、销售、客服 除外
	var login_user = core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)                   // 是不是管理员
	var is_sales = pass_service.Is_user_sales(login_user.UserId)                   // 是不是销售
	var is_custom_service = pass_service.Is_user_custom_service(login_user.UserId) // 是不是客服

	var get_stock_id_list_sql = ``
	if login_user.IsSuperAdmin == false && is_admin == false && is_custom_service == false && is_sales == false { // 管理员，客服，销售 可以看到所有仓库的订单
		get_stock_id_list_sql = fmt.Sprintf(`SELECT spc.stock_id AS stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = %d`, login_user.UserId)
	} else {
		get_stock_id_list_sql = `SELECT s.id AS stock_id FROM erp_stock s`
	}

	type StockId struct {
		StockId int64 `xorm:"'stock_id'" json:"stock_id"`
	}

	var stock_id_list = make([]StockId, 0)
	err := db_session.Sql(get_stock_id_list_sql).Find(&stock_id_list)
	core.HandleError(err)

	var stock_id_str_list = make([]string, 0)
	for _, stock_id := range stock_id_list {
		stock_id_str_list = append(stock_id_str_list, fmt.Sprint(stock_id.StockId))
	}

	return stock_id_str_list
}
