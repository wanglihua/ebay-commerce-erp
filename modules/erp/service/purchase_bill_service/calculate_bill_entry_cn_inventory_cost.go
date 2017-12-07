package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func CalculateBillEntryCnInventoryCost(purchase_bill_id int64, db_session *xorm.Session) {
	/*

	// A 更新不是分批采购的中国仓成本价

	// SELECT (pb.advance_expenses + pb.final_expenses) / sum(pbe.price * pbe.qty) AS fee_point
	// FROM erp_purchase_bill pb
	//     INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
	// WHERE pb.id = 1 AND pbe.type_id = 123
	// GROUP BY pbe.bill_id;

	// UPDATE erp_purchase_bill_entry
	// SET inventory_cost = price - if(need_invoice = 0, 0, drawback_amount) + fitting_cost + other_cost + price * 0.4
	// WHERE bill_id = 1 AND type_id = 123;

	// 1
	// 获取运费平均的计算点数。计算点数 = (1/所有明细的价格乘以数量的总和)*采购单运费 。
	// 计算点数可以用来计算每笔按采购价格平摊的运费，如果是计算SKU的运费，乘以采购价格就可以了。
	var fee_point struct {
		FeePoint float64 `xorm:"'fee_point'"`
	}
	var sql_get_fee_point = `SELECT (pb.advance_expenses + pb.final_expenses) / sum(pbe.price * pbe.qty) AS fee_point
FROM erp_purchase_bill pb
    INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
WHERE pb.id = %d AND pbe.type_id = %d
GROUP BY pbe.bill_id;`
	sql_get_fee_point = fmt.Sprintf(sql_get_fee_point, purchase_bill_id, cggl_models.PURCHASE_BILL_TYPE_NORMAL)
	_, err := db_session.Sql(sql_get_fee_point).Get(&fee_point)
	core.HandleError(err)

	// 2
	// 利用计算点数，计算每个明细的中国仓成本价
	var sql_update_inventory_cost = `UPDATE erp_purchase_bill_entry
SET inventory_cost = price - if(need_invoice = 0, 0, drawback_amount) + fitting_cost + other_cost + price * %f
WHERE bill_id = %d AND type_id = %d;`
	sql_update_inventory_cost = fmt.Sprintf(sql_update_inventory_cost, fee_point.FeePoint, purchase_bill_id, cggl_models.PURCHASE_BILL_TYPE_NORMAL)
	_, err = db_session.Exec(sql_update_inventory_cost)
	core.HandleError(err)

	// B 更新分批采购的中国仓成本价

	// UPDATE erp_purchase_bill_entry
	// SET inventory_cost = price - if(need_invoice = 0, 0, drawback_amount) + fitting_cost + other_cost + final_expenses / qty
	// WHERE bill_id = 1 AND type_id = 456;

	var sql_update_in_turn_inventory_cost = `UPDATE erp_purchase_bill_entry
SET inventory_cost = price - if(need_invoice = 0, 0, drawback_amount) + fitting_cost + other_cost + final_expenses / qty
WHERE bill_id = %d AND type_id = %d;`
	sql_update_in_turn_inventory_cost = fmt.Sprintf(sql_update_in_turn_inventory_cost, purchase_bill_id, cggl_models.PURCHASE_BILL_TYPE_IN_TURN)
	_, err = db_session.Exec(sql_update_in_turn_inventory_cost)
	core.HandleError(err)

	*/

	// 更新中国仓成本价，分批和不分批的都采用同样的算法

	// 1
	// 获取运费平均的计算点数。计算点数 = (1/所有明细的价格乘以数量的总和)*采购单运费 。
	// 计算点数可以用来计算每笔按采购价格平摊的运费，如果是计算SKU的运费，乘以采购价格就可以了。
	// 这里的 pb.final_expenses 一般都是零
	var fee_point struct {
		FeePoint float64 `xorm:"'fee_point'"`
	}
	var sql_get_fee_point = `SELECT (pb.advance_expenses + pb.final_expenses) / sum(pbe.price * pbe.qty) AS fee_point
FROM erp_purchase_bill pb
    INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
WHERE pb.id = %d
GROUP BY pbe.bill_id;`
	sql_get_fee_point = fmt.Sprintf(sql_get_fee_point, purchase_bill_id)
	_, err := db_session.Sql(sql_get_fee_point).Get(&fee_point)
	core.HandleError(err)

	// 2
	// 利用计算点数，计算每个明细的中国仓成本价
	var sql_update_inventory_cost = `UPDATE erp_purchase_bill_entry
SET inventory_cost = price - if(need_invoice = 0, 0, drawback_amount) + fitting_cost + other_cost + final_expenses / qty + price * %f
WHERE bill_id = %d;`
	sql_update_inventory_cost = fmt.Sprintf(sql_update_inventory_cost, fee_point.FeePoint, purchase_bill_id)
	_, err = db_session.Exec(sql_update_inventory_cost)
	core.HandleError(err)
}
