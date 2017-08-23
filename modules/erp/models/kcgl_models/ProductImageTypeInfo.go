package kcgl_models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	//橱窗图片
	PRODUCT_IMAGE_TYPE_CCTP = 1

	//内容图片
	PRODUCT_IMAGE_TYPE_NRTP = 2

	//库存图片
	PRODUCT_IMAGE_TYPE_KCTP = 3

	//质检图片
	PRODUCT_IMAGE_TYPE_ZJTP = 4
)

type ProductImageTypeInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductImageTypeInfo) TableName() string {
	return "erp_product_image_type"
}

func (e ProductImageTypeInfo) ModelDesc() string {
	return "verbose_name=商品图片类型,entity_json=product_image_type,route=product/image/type"
}

func (e ProductImageTypeInfo) InitData(db_session *xorm.Session) {

	product_image_type_id_list := []int64{
		PRODUCT_IMAGE_TYPE_CCTP,
		PRODUCT_IMAGE_TYPE_NRTP,
		PRODUCT_IMAGE_TYPE_KCTP,
		PRODUCT_IMAGE_TYPE_ZJTP,
	}

	var product_image_type_list []ProductImageTypeInfo
	err := db_session.In("id", product_image_type_id_list).
		Find(&product_image_type_list)
	core.HandleError(err)

	if len(product_image_type_list) != len(product_image_type_id_list) {
		product_image_type := new(ProductImageTypeInfo)
		//product_image_type.Id = PRODUCT_IMAGE_TYPE_CCTP
		product_image_type.Name = "橱窗图片"
		product_image_type.Desc = core.MakeNullString("橱窗图片 上传", true)
		db_session.Insert(product_image_type)

		product_image_type = new(ProductImageTypeInfo)
		//product_image_type.Id = PRODUCT_IMAGE_TYPE_NRTP
		product_image_type.Name = "内容图片"
		product_image_type.Desc = core.MakeNullString("内容图片 上传", true)
		db_session.Insert(product_image_type)

		product_image_type = new(ProductImageTypeInfo)
		//product_image_type.Id = PRODUCT_IMAGE_TYPE_KCTP
		product_image_type.Name = "库存图片"
		product_image_type.Desc = core.MakeNullString("库存图片", true)
		db_session.Insert(product_image_type)

		product_image_type = new(ProductImageTypeInfo)
		//product_image_type.Id = PRODUCT_IMAGE_TYPE_ZJTP
		product_image_type.Name = "质检图片"
		product_image_type.Desc = core.MakeNullString("质检图片", true)
		db_session.Insert(product_image_type)
	}
}
