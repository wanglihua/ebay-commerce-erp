package bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"time"
	"fmt"
	"github.com/revel/revel"
)

const (
	ENTITY_ID_INVENTORY_CHECK      = 10001 // 盘点单
	ENTITY_ID_STOCK_BILL_CHECK_IN  = 10002 // 入库单
	ENTITY_ID_STOCK_BILL_CHECK_OUT = 10003 // 出库单
	ENTITY_ID_STOCK_BILL_TRANSFER  = 10004 // 调拨单
	ENTITY_ID_PURCHASE_BILL        = 10005 // 采购单
)

func GetNextBillNo(entity_id int64, db_session *xorm.Session, web_session revel.Session) string {
	sql_proc := `
CREATE PROCEDURE erp_get_next_bill_serial_v0002(p_entity_id INT)
    BEGIN
        DECLARE v_today DATETIME;
        DECLARE v_count INT;

        SET v_today = date(now());

        SELECT count(*) INTO v_count FROM erp_entity_code WHERE entity_id = p_entity_id AND date = v_today;

        IF v_count = 0 THEN
            INSERT INTO erp_entity_code (entity_id, serial, date) VALUES (p_entity_id, 1, v_today);
        ELSE
            UPDATE erp_entity_code SET serial = erp_entity_code.serial + 1 WHERE entity_id = p_entity_id AND date = v_today;
        END IF;

        SELECT serial AS next_serial FROM erp_entity_code WHERE entity_id = p_entity_id AND date = v_today;
    END;
	`
	core.CreateProc("erp_get_next_bill_serial_v0002", sql_proc, db_session)

	type BillNextSerial struct {
		NextSerial int64 `xorm:"'next_serial'" json:"next_serial"`
	}

	var bill_next_serial BillNextSerial
	sql := `CALL erp_get_next_bill_serial_v0002();`
	_, err := db_session.Sql(sql).Get(&bill_next_serial)
	core.HandleError(err)

	var login_user = core.GetLoginUser(web_session)

	return fmt.Sprintf("%s%02d-%s", time.Now().Format("20060102"), bill_next_serial.NextSerial, login_user.LoginName)
}
