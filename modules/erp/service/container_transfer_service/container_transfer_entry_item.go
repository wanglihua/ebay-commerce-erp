package container_transfer_service

import (
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/hgbg_models"
)

type ContainerTransferEntryItem struct {
	ContainerTransferEntry hgbg_models.ContainerTransferEntryInfo   `xorm:"extends" json:"cte"`
	Product                models.ProductInfo                       `xorm:"extends" json:"p"`
}
