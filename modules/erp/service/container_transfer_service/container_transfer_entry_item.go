package container_transfer_service

import (
	"matrix/modules/erp/models"
)

type ContainerTransferEntryItem struct {
	ContainerTransferEntry models.ContainerTransferEntryInfo   `xorm:"extends" json:"cte"`
	Product                models.ProductInfo                  `xorm:"extends" json:"p"`
}
