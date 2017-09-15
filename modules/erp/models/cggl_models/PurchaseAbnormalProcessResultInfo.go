package cggl_models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	ABNORMAL_PROCESS_RESULT_WCL   = 100 // 未处理
	ABNORMAL_PROCESS_RESULT_YTK   = 200 // 已退款
	ABNORMAL_PROCESS_RESULT_YCXRK = 300 // 已重新入库
	ABNORMAL_PROCESS_RESULT_ZCSH  = 400 // 正常损耗
	ABNORMAL_PROCESS_RESULT_YCXCG = 500 // 已重新采购
	ABNORMAL_PROCESS_RESULT_QT    = 600 // 其它
)

type PurchaseAbnormalProcessResultInfo struct {
	Id           int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	Name         string    `xorm:"nvarchar(100) index 'name'" json:"name" smith:"verbose_name=名称"`
	Description  string    `xorm:"nvarchar(500) 'description'" json:"description" smith:"verbose_name=备注,min=2,max=400"`
	CheckinAgain int64     `xorm:"bigint notnull 'checkin_again'" json:"checkin_again" smith:"verbose_name=是否重新入库"`
	CreateTime   core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseAbnormalProcessResultInfo) TableName() string {
	return "erp_purchase_abnormal_process_result"
}

func (e PurchaseAbnormalProcessResultInfo) ModelDesc() string {
	return "verbose_name=采购异常类型,entity_json=process_result,route=purchase/abnormal/process/result"
}

func (e PurchaseAbnormalProcessResultInfo) InitData(db_session *xorm.Session) {

	process_result := new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_WCL
	process_result.Name = "未处理"
	process_result.Description = "未处理"
	process_result.CheckinAgain = 0
	db_session.Insert(process_result)

	process_result = new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_YTK
	process_result.Name = "已退款"
	process_result.Description = "已退款"
	process_result.CheckinAgain = 0
	db_session.Insert(process_result)

	process_result = new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_YCXRK
	process_result.Name = "已重新入库"
	process_result.Description = "已重新入库"
	process_result.CheckinAgain = 1
	db_session.Insert(process_result)

	process_result = new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_ZCSH
	process_result.Name = "正常损耗"
	process_result.Description = "正常损耗"
	process_result.CheckinAgain = 0
	db_session.Insert(process_result)

	process_result = new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_YCXCG
	process_result.Name = "已重新采购"
	process_result.Description = "已重新采购"
	process_result.CheckinAgain = 0
	db_session.Insert(process_result)

	process_result = new(PurchaseAbnormalProcessResultInfo)
	process_result.Id = ABNORMAL_PROCESS_RESULT_QT
	process_result.Name = "其它"
	process_result.Description = "其它"
	process_result.CheckinAgain = 0
	db_session.Insert(process_result)
}
