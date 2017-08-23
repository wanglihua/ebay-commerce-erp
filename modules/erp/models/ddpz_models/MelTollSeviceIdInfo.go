package ddpz_models

type MelTollSeviceIdInfo struct {
	Postcode  int64             `xorm:"nvarchar(50) notnull pk 'postcode'" json:"postcode"`
	ServiceId string            `xorm:"nvarchar(50) notnull 'service_id'" json:"service_id"`
}

func (e MelTollSeviceIdInfo) TableName() string {
	return "erp_mel_toll_service_id"
}

func (e MelTollSeviceIdInfo) ModelDesc() string {
	return "verbose_name=墨尔本Toll服务ID,entity_json=mel_toll_service_id,route=mel/toll/service/id"
}
