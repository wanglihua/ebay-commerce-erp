package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/db"
	auth_models "matrix/modules/auth/models"
	erp_models "matrix/modules/erp/models"
	track_modes "matrix/modules/track/models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/spzl_models"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/spbb_models"
	"matrix/modules/erp/models/xtgl_models"
	"matrix/modules/erp/models/xsgl_models"
	"matrix/modules/erp/models/cggl_models"
)

type Home struct {
	*revel.Controller
}

func (c Home) Index() revel.Result {
	return c.RenderTemplate("home/home_index.html")
}

func (c Home) Loading() revel.Result {
	return c.RenderTemplate("home/home_loading.html")
}

func (c Home) Hotload() revel.Result {
	db_session := c.DbSession

	sql := `SELECT now() AS time_now`
	var time_now struct {
		Now core.Time `xorm:"'time_now'"`
	}
	_, err := db_session.Sql(sql).Get(&time_now)
	core.HandleError(err)

	revel.TRACE.Println(time_now.Now)

	return c.RenderTemplate("home/home_hotload.html")
}

func (c Home) SyncDb() revel.Result {
	return c.RenderTemplate("home/home_syncdb.html")
}

func (c Home) SyncDbPost() revel.Result {

	model_list := make([]interface{}, 0)

	model_list = append(model_list, &auth_models.AdminInfo{})
	model_list = append(model_list, &auth_models.GroupInfo{})
	model_list = append(model_list, &auth_models.GroupUserInfo{})
	model_list = append(model_list, &auth_models.PassResCateInfo{})
	model_list = append(model_list, &auth_models.PassResInfo{})
	model_list = append(model_list, &auth_models.PassUserResInfo{})
	model_list = append(model_list, &auth_models.SettingInfo{})
	model_list = append(model_list, &auth_models.UserInfo{})

	model_list = append(model_list, &ddpz_models.LogisticsCompanyInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsFeeRateInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsModeInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsModeTypeInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsNumberBatchInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsNumberInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsNumberLogInfo{})
	model_list = append(model_list, &ddpz_models.LogisticsSettingInfo{})
	model_list = append(model_list, &ddpz_models.MelTollSeviceIdInfo{})
	model_list = append(model_list, &ddpz_models.OrderCustomTypeInfo{})
	model_list = append(model_list, &ddpz_models.OrderHoldTypeInfo{})
	model_list = append(model_list, &ddpz_models.PickFeeInfo{})
	model_list = append(model_list, &ddpz_models.OrderTagInfo{})
	model_list = append(model_list, &ddpz_models.OrderTypeInfo{})
	model_list = append(model_list, &ddpz_models.EnvelopeFeeRateInfo{})
	model_list = append(model_list, &ddpz_models.PickupPayMethodInfo{})
	model_list = append(model_list, &ddpz_models.FavStockTimeRangeInfo{})

	model_list = append(model_list, &sppz_models.BrandInfo{})
	model_list = append(model_list, &sppz_models.ProductStyleInfo{})
	model_list = append(model_list, &sppz_models.RegionInfo{})
	model_list = append(model_list, &sppz_models.StockInfo{})
	model_list = append(model_list, &sppz_models.StockPickCrewInfo{})
	model_list = append(model_list, &sppz_models.StockTypeInfo{})
	model_list = append(model_list, &sppz_models.StorageLocInfo{})
	model_list = append(model_list, &sppz_models.SupplierInfo{})
	model_list = append(model_list, &sppz_models.SupplierPayMethodInfo{})
	model_list = append(model_list, &sppz_models.ProductCertInfo{})

	model_list = append(model_list, &hgbg_models.ContainerTransferInfo{})
	model_list = append(model_list, &hgbg_models.ContainerTransferDestPortInfo{})
	model_list = append(model_list, &hgbg_models.ContainerTransferEntryInfo{})
	model_list = append(model_list, &hgbg_models.ContainerTransferLogInfo{})
	model_list = append(model_list, &hgbg_models.ContainerTransferStartPortInfo{})
	model_list = append(model_list, &hgbg_models.HsCodeAuDutyInfo{})
	model_list = append(model_list, &hgbg_models.HsCodeCnReportElemInfo{})

	model_list = append(model_list, &kcgl_models.InventoryInfo{})
	model_list = append(model_list, &kcgl_models.InventoryCheckInfo{})
	model_list = append(model_list, &kcgl_models.InventoryCheckEntryInfo{})
	model_list = append(model_list, &kcgl_models.InventoryHistoryInfo{})
	model_list = append(model_list, &kcgl_models.ProductBindInfo{})
	model_list = append(model_list, &kcgl_models.ProductBindLogInfo{})
	model_list = append(model_list, &kcgl_models.ProductCateInfo{})
	model_list = append(model_list, &kcgl_models.ProductCateLogInfo{})
	model_list = append(model_list, &kcgl_models.ProductImageInfo{})
	model_list = append(model_list, &kcgl_models.ProductImageTypeInfo{})
	model_list = append(model_list, &kcgl_models.ProductImportInfo{})
	model_list = append(model_list, &kcgl_models.ProductInfo{})
	model_list = append(model_list, &kcgl_models.ProductLogInfo{})
	model_list = append(model_list, &kcgl_models.StockBillEntryInfo{})
	model_list = append(model_list, &kcgl_models.StockBillInfo{})
	model_list = append(model_list, &kcgl_models.StockBillTypeInfo{})
	model_list = append(model_list, &kcgl_models.SpecialProductInfo{})
	model_list = append(model_list, &kcgl_models.ProductMelTempSkuInfo{})
	model_list = append(model_list, &spzl_models.ProductFeedbackCateInfo{})
	model_list = append(model_list, &spzl_models.ProductFeedbackSkuInfo{})
	model_list = append(model_list, &spzl_models.ProductFeedbackInfo{})

	model_list = append(model_list, &xtpz_models.ExchangeRateInfo{})
	model_list = append(model_list, &xtpz_models.ExchangeRateLogInfo{})
	model_list = append(model_list, &xtpz_models.CurrencyInfo{})
	model_list = append(model_list, &xtpz_models.EBayAPIConfigInfo{})
	model_list = append(model_list, &xtpz_models.ECommercePlatformInfo{})
	model_list = append(model_list, &xtpz_models.PaypalAccountInfo{})
	model_list = append(model_list, &xtpz_models.ShopAccountInfo{})
	model_list = append(model_list, &xtpz_models.ShopAccountLogInfo{})
	model_list = append(model_list, &xtpz_models.PlatformLogisticsModeInfo{})

	model_list = append(model_list, &ddlc_models.OrderInfo{})
	model_list = append(model_list, &ddlc_models.OrderHistoryInfo{})
	model_list = append(model_list, &ddlc_models.OrderItemInfo{})
	model_list = append(model_list, &ddlc_models.OrderLogInfo{})
	model_list = append(model_list, &ddlc_models.OrderLogisticsInfo{})
	model_list = append(model_list, &ddlc_models.OrderLogisticsLogicInfo{})
	model_list = append(model_list, &ddlc_models.OrderLogisticsMatchDebugInfo{})
	model_list = append(model_list, &ddlc_models.OrderShippingAddressInfo{})
	model_list = append(model_list, &ddlc_models.PickupOrderInfo{})

	model_list = append(model_list, &xsgl_models.SalesManagerInfo{})
	model_list = append(model_list, &xsgl_models.ProductSaleSettingInfo{})

	model_list = append(model_list, &xtgl_models.EntityCodeInfo{})
	model_list = append(model_list, &xtgl_models.LittleIconInfo{})
	model_list = append(model_list, &xtgl_models.LoginLogInfo{})
	model_list = append(model_list, &xtgl_models.NoticeInfo{})
	model_list = append(model_list, &xtgl_models.NoticeReadInfo{})

	model_list = append(model_list, &sppz_models.EubGoodsTypeInfo{})
	model_list = append(model_list, &sppz_models.EubInsuranceTypeInfo{})
	model_list = append(model_list, &sppz_models.PackageMaterialInfo{})
	model_list = append(model_list, &sppz_models.SeaContainerCostInfo{})
	model_list = append(model_list, &sppz_models.ProductStoreItemInfo{})

	model_list = append(model_list, &cggl_models.IntelReplInfo{})
	model_list = append(model_list, &cggl_models.PurchaseBillEntryInfo{})
	model_list = append(model_list, &cggl_models.PurchaseBillInfo{})
	model_list = append(model_list, &cggl_models.PurchaseBillLogInfo{})
	model_list = append(model_list, &cggl_models.PurchaseAbnormalTypeInfo{})
	model_list = append(model_list, &cggl_models.PurchaseAbnormalProcessResultInfo{})

	model_list = append(model_list, &spbb_models.ProductSaleDayAnalysisInfo{})
	model_list = append(model_list, &spbb_models.ProductSaleAnalysisInfo{})
	model_list = append(model_list, &spbb_models.ProductSaleAnalysisSettingInfo{})
	model_list = append(model_list, &erp_models.StatusInfo{})
	model_list = append(model_list, &erp_models.MessageInfo{})
	model_list = append(model_list, &erp_models.PaypalReconInfo{})

	//model_list = append(model_list, &erp_models.TestInfo{}) // 为什么这一句打开就不行

	model_list = append(model_list, &track_modes.IssueAttachInfo{})
	model_list = append(model_list, &track_modes.IssueCateInfo{})
	model_list = append(model_list, &track_modes.IssueCommentInfo{})
	model_list = append(model_list, &track_modes.IssueInfo{})
	model_list = append(model_list, &track_modes.IssueLogInfo{})
	//model_list = append(model_list, &track_modes.IssueRelateInfo{})
	model_list = append(model_list, &track_modes.IssueStatusInfo{})
	model_list = append(model_list, &track_modes.IssueTypeInfo{})

	err := db.Engine.Sync2(model_list...)
	core.HandleError(err)

	//执行数据初始化
	db_session := c.DbSession
	for _, model := range model_list {
		if table_initer, ok := model.(core.TableDataIniter); ok {
			table_initer.InitData(db_session)
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据库同步成功!"})
}
