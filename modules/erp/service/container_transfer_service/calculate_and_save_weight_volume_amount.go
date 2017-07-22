package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/core"
)

func CalculateAndSaveWeightVolumeAmount(transfer_id int64, db_session *xorm.Session) {
	var container_transfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&container_transfer)
	core.HandleError(err)

	/*
SELECT
    sum(p.weight * cte.qty)                                         AS total_weight,
    sum(p.length * p.width * p.height * cte.qty) / 1000000          AS total_volume,
    sum(cte.price * cte.qty)                                        AS total_amount1,
    sum(cte.price * cte.qty) * (1 + 0.1)                            AS total_amount2,
    sum(if(cte.au_hs_duty_rate <> 0, cte.price * cte.qty * 0.1, 0)) AS duty_amount_discounted,
    sum(if(cte.au_hs_duty_rate = 0, cte.price * cte.qty, 0))        AS duty_free_amount
FROM erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id
WHERE cte.transfer_id = 1;
	 */

	type weight_volume_amount_t struct {
		TotalWeight          float64 `xorm:"'total_weight'"`
		TotalVolume          float64 `xorm:"'total_volume'"`
		TotalAmount1         float64 `xorm:"'total_amount1'"`
		TotalAmount2         float64 `xorm:"'total_amount2'"`
		DutyAmountDiscounted float64 `xorm:"'duty_amount_discounted'"`
		DutyFreeAmount       float64 `xorm:"'duty_free_amount'"`
	}

	var query_sql = `SELECT
    sum(p.weight * cte.qty)                                         AS total_weight,
    sum(p.length * p.width * p.height * cte.qty) / 1000000          AS total_volume,
    sum(cte.price * cte.qty)                                        AS total_amount1,
    sum(cte.price * cte.qty) * (1 + %f)                             AS total_amount2,
    sum(if(cte.au_hs_duty_rate <> 0, cte.price * cte.qty * %f, 0))  AS duty_amount_discounted,
    sum(if(cte.au_hs_duty_rate = 0, cte.price * cte.qty, 0))        AS duty_free_amount
FROM erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id
WHERE cte.transfer_id = %d;`

	query_sql = fmt.Sprintf(query_sql, container_transfer.CostUpRate, container_transfer.PriceDownRate, container_transfer.Id)

	var weight_volume_amount weight_volume_amount_t
	_, err = db_session.Sql(query_sql).Get(&weight_volume_amount)
	core.HandleError(err)

	container_transfer.TotalWeight = weight_volume_amount.TotalWeight
	container_transfer.TotalVolumn = weight_volume_amount.TotalVolume
	container_transfer.TotalAmount1 = weight_volume_amount.TotalAmount1
	container_transfer.TotalAmount2 = weight_volume_amount.TotalAmount2

	if weight_volume_amount.DutyFreeAmount != 0 {
		container_transfer.PriceUpRate = (weight_volume_amount.TotalAmount2 - weight_volume_amount.DutyAmountDiscounted - weight_volume_amount.DutyFreeAmount) / weight_volume_amount.DutyFreeAmount
	} else {
		container_transfer.PriceUpRate = 0
	}

	_, err = db_session.Cols("total_weight", "total_volumn", "total_amount1", "total_amount2", "price_up_rate").Update(&container_transfer)
	core.HandleError(err)
}
