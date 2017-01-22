package bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"time"
	"fmt"
)

func GetNextBillNo(db_session *xorm.Session) string {
	sql_proc := `
CREATE PROCEDURE erp_get_next_bill_serial_v0001()
    BEGIN
        DECLARE ENTITY_ID_BILL INT;
        DECLARE v_today DATETIME;
        DECLARE v_count INT;

        SET ENTITY_ID_BILL = 100100;
        SET v_today = date(now());

        SELECT count(*) INTO v_count FROM erp_entity_code WHERE entity_id = ENTITY_ID_BILL AND date = v_today;

        IF v_count = 0 THEN
            INSERT INTO erp_entity_code (entity_id, serial, date) VALUES (ENTITY_ID_BILL, 1, v_today);
        ELSE
            UPDATE erp_entity_code SET serial = erp_entity_code.serial + 1 WHERE entity_id = ENTITY_ID_BILL AND date = v_today;
        END IF;

        SELECT serial AS next_serial FROM erp_entity_code WHERE entity_id = ENTITY_ID_BILL AND date = v_today;
    END;
	`
	core.CreateProc("erp_get_next_bill_serial_v0001", sql_proc, db_session)

	type BillNextSerial struct {
		NextSerial int64 `xorm:"'next_serial'" json:"next_serial"`
	}

	var bill_next_serial BillNextSerial
	sql := `CALL erp_get_next_bill_serial_v0001();`
	_, err := db_session.Sql(sql).Get(&bill_next_serial)
	core.HandleError(err)

	return fmt.Sprintf("%s%04d", time.Now().Format("2006-01-02"), bill_next_serial.NextSerial)
}
