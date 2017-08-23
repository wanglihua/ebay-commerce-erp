package container_transfer_service

import (
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/erp/models/kcgl_models"
)

type ContainerTransferEntryItem struct {
	ContainerTransferEntry hgbg_models.ContainerTransferEntryInfo   `xorm:"extends" json:"cte"`
	Product                kcgl_models.ProductInfo                       `xorm:"extends" json:"p"`
}
