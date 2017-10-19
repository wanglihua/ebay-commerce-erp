package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Calculate_product_price(db_session *xorm.Session) {
	var sql = `CALL calculate_product_price();`

	_, err := db_session.Exec(sql)
	core.HandleError(err)
}
