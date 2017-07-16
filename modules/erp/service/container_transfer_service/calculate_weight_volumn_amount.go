package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/core"
)

func CalculateWeightVolumeAmount(p_container_transfer *models.ContainerTransferInfo, db_session *xorm.Session) {
	type weight_volume_amount_t struct {
		Weight float64 `xorm:"'weight'"`
		Volume float64 `xorm:"'volume'"`
		Amount float64 `xorm:"'amount'"`
	}

	var query_sql = `SELECT
    p.weight * cte.qty                      AS weight,
    p.length * p.width * p.height / 1000000 AS volume,
    cte.price * cte.qty                     AS amount
FROM erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id
WHERE cte.transfer_id = %d;`

	query_sql = fmt.Sprintf(query_sql, p_container_transfer.Id)

	var weight_volume_amount weight_volume_amount_t
	_, err := db_session.Sql(query_sql).Get(&weight_volume_amount)
	core.HandleError(err)

	p_container_transfer.TotalVolumn = weight_volume_amount.Volume
	p_container_transfer.TotalWeight = weight_volume_amount.Weight
	p_container_transfer.TotalAmount1 = weight_volume_amount.Amount * (1 + p_container_transfer.PriceUpRate)

}
