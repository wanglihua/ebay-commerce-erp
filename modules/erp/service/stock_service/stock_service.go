package stock_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

func Get_stock_id_type_by_name(stock_name string, db_session *xorm.Session) (id_return int64, type_return int64) {

	var stock_id_type struct {
		StockId int64 `xorm:"'stock_id'"`
		TypeId  int64 `xorm:"'type_id'"`
	}

	/*
SELECT s.id AS stock_id, s.stock_type_id AS type_id FROM erp_stock s WHERE s.name = 'abc'
	 */

	var sql_get_stock_id = `SELECT s.id AS stock_id, s.stock_type_id AS type_id FROM erp_stock s WHERE s.name = '%s''`
	sql_get_stock_id = fmt.Sprintf(sql_get_stock_id, stock_name)
	_, err := db_session.Sql(sql_get_stock_id).Get(&stock_id_type)
	core.HandleError(err)

	return stock_id_type.StockId, stock_id_type.TypeId
}

func Get_stock_id_name_map(stock_id_list []int64, db_session *xorm.Session) map[int64]string {
	var stock_id_name_map = make(map[int64]string)

	if len(stock_id_list) == 0 {
		return stock_id_name_map
	}

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.In("id", stock_id_list).Find(&stock_list)
	core.HandleError(err)

	for _, stock := range stock_list {
		stock_id_name_map[stock.Id] = stock.Name
	}

	return stock_id_name_map
}

// 获取人员的仓库列表，如果是拣货人员的话，只获取有权限的仓库
func Get_login_user_stock_list(db_session *xorm.Session, web_session revel.Session) []models.StockInfo {
	/*
	SELECT * FROM erp_stock WHERE stock_type_id = 1 AND id IN (SELECT spc.stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = 1)
	 */

	var stock_list = make([]models.StockInfo, 0)

	var stock_list_filter = fmt.Sprintf("stock_type_id = %d", models.STOCK_TYPE_AU) //只取澳洲的仓库
	var login_user = core.GetLoginUser(web_session)
	if login_user.IsSuperAdmin == false && pass_service.Is_user_admin(login_user.UserId) == false {
		var is_user_custom_service = pass_service.Is_user_custom_service(login_user.UserId)
		var is_user_sales = pass_service.Is_user_sales(login_user.UserId)
		if is_user_custom_service == false && is_user_sales == false {
			if pass_service.Is_user_pick_crew(login_user.UserId) { // 如果是仓库拣货人员
				stock_list_filter = core.SqlAnd(stock_list_filter, fmt.Sprintf("id IN (SELECT spc.stock_id FROM erp_stock_pick_crew spc WHERE spc.user_id = %d)", login_user.UserId))
			} else {
				stock_list_filter = core.SqlAnd(stock_list_filter, "1 = 0") //不是管理员，不是拣货人员，也不是销售和客服，就不能看任何一个仓库
			}
		}
	}

	err := db_session.Where(stock_list_filter).OrderBy("ord asc").Find(&stock_list)
	core.HandleError(err)

	return stock_list
}

//获取墨尔本仓库的 stock id，这个函数后来没采用，用了这个 。stock_pick_crew_service.Is_pick_crew_has_mel_stock
func Get_mel_stock_id(db_session *xorm.Session) int64 {

	var stock_id struct {
		StockId int64 `xorm:"'stock_id'"`
	}

	/*
	SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%墨尔本%';
	 */

	var sql = `SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%墨尔本%'`
	_, err := db_session.Sql(sql).Get(&stock_id)
	core.HandleError(err)

	return stock_id.StockId
}
