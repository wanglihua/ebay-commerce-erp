package ddpz_models

import "matrix/core"

type LogisticsFeeRate2Info struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ModeId      int64     `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`
	Postcode    string    `xorm:"nvarchar(200) notnull index 'postcode'" json:"postcode" smith:"verbose_name=邮编,min=2,max=100"`
	Fee01       float64   `xorm:"decimal(18,6) notnull 'fee01'" json:"fee01"`
	Fee02       float64   `xorm:"decimal(18,6) notnull 'fee02'" json:"fee02"`
	Fee03       float64   `xorm:"decimal(18,6) notnull 'fee03'" json:"fee03"`
	Fee04       float64   `xorm:"decimal(18,6) notnull 'fee04'" json:"fee04"`
	Fee05       float64   `xorm:"decimal(18,6) notnull 'fee05'" json:"fee05"`
	Fee06       float64   `xorm:"decimal(18,6) notnull 'fee06'" json:"fee06"`
	Fee07       float64   `xorm:"decimal(18,6) notnull 'fee07'" json:"fee07"`
	Fee08       float64   `xorm:"decimal(18,6) notnull 'fee08'" json:"fee08"`
	Fee09       float64   `xorm:"decimal(18,6) notnull 'fee09'" json:"fee09"`
	Fee10       float64   `xorm:"decimal(18,6) notnull 'fee10'" json:"fee10"`
	Fee11       float64   `xorm:"decimal(18,6) notnull 'fee11'" json:"fee11"`
	Fee12       float64   `xorm:"decimal(18,6) notnull 'fee12'" json:"fee12"`
	Sfpydq      int       `xorm:"int notnull 'sfpydq'" json:"sfpydq" smith:"verbose_name=是否为偏远地区"` // 0 否， 1 是
	Pydqjsyf    float64   `xorm:"decimal(18,6) notnull 'pydqjsyf'" json:"pydqjsyf" smith:"verbose_name=偏远地区加收运费"`
	Description string    `xorm:"nvarchar(200) notnull 'description'" json:"description" smith:"verbose_name=描述,max=200"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LogisticsFeeRate2Info) TableName() string {
	return "erp_logistics_fee_rate2"
}

func (e LogisticsFeeRate2Info) ModelDesc() string {
	return "verbose_name=物流费率,entity_json=logistics_fee_rate2,route=logistics/fee/rate2"
}

/*

CREATE TABLE erp_logistics_fee_rate2 (
    id          BIGINT(20)     NOT NULL AUTO_INCREMENT,
    mode_id     BIGINT(20)     NOT NULL,
    postcode    VARCHAR(200)   NOT NULL,
    fee01       DECIMAL(18, 6) NOT NULL,
    fee02       DECIMAL(18, 6) NOT NULL,
    fee03       DECIMAL(18, 6) NOT NULL,
    fee04       DECIMAL(18, 6) NOT NULL,
    fee05       DECIMAL(18, 6) NOT NULL,
    fee06       DECIMAL(18, 6) NOT NULL,
    fee07       DECIMAL(18, 6) NOT NULL,
    fee08       DECIMAL(18, 6) NOT NULL,
    fee09       DECIMAL(18, 6) NOT NULL,
    fee10       DECIMAL(18, 6) NOT NULL,
    fee11       DECIMAL(18, 6) NOT NULL,
    fee12       DECIMAL(18, 6) NOT NULL,
    sfpydq      INT(11)        NOT NULL,
    pydqjsyf    DECIMAL(18, 6) NOT NULL,
    description VARCHAR(500)   NOT NULL,
    create_time DATETIME       NOT NULL,
    update_time DATETIME       NOT NULL,
    PRIMARY KEY (id),
    KEY IDX_erp_logistics_fee_rate2_postcode (postcode),
    KEY idx_lfr2_mode_id (mode_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8


 */
