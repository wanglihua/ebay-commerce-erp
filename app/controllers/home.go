package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/db"
	auth_models "matrix/modules/auth/models"
	erp_models "matrix/modules/erp/models"
	track_modes "matrix/modules/track/models"
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

	model_list = append(model_list, &erp_models.BrandInfo{})
	model_list = append(model_list, &erp_models.ContainerTransferInfo{})
	model_list = append(model_list, &erp_models.ContainerTransferDestPortInfo{})
	model_list = append(model_list, &erp_models.ContainerTransferEntryInfo{})
	model_list = append(model_list, &erp_models.ContainerTransferStartPortInfo{})
	model_list = append(model_list, &erp_models.CurrencyInfo{})
	model_list = append(model_list, &erp_models.EBayAPIConfigInfo{})
	model_list = append(model_list, &erp_models.ECommercePlatformInfo{})
	model_list = append(model_list, &erp_models.EntityCodeInfo{})
	model_list = append(model_list, &erp_models.EnvelopeFeeRateInfo{})
	model_list = append(model_list, &erp_models.EubGoodsTypeInfo{})
	model_list = append(model_list, &erp_models.EubInsuranceTypeInfo{})
	model_list = append(model_list, &erp_models.ExchangeRateInfo{})
	model_list = append(model_list, &erp_models.ExchangeRateLogInfo{})
	model_list = append(model_list, &erp_models.HsCodeAuDutyInfo{})
	model_list = append(model_list, &erp_models.HsCodeCnReportElemInfo{})
	model_list = append(model_list, &erp_models.InventoryInfo{})
	model_list = append(model_list, &erp_models.InventoryCheckInfo{})
	model_list = append(model_list, &erp_models.InventoryCheckEntryInfo{})
	model_list = append(model_list, &erp_models.InventoryHistoryInfo{})
	model_list = append(model_list, &erp_models.LittleIconInfo{})
	model_list = append(model_list, &erp_models.LoginLogInfo{})
	model_list = append(model_list, &erp_models.LogisticsCompanyInfo{})
	model_list = append(model_list, &erp_models.LogisticsFeeRateInfo{})
	model_list = append(model_list, &erp_models.LogisticsModeInfo{})
	model_list = append(model_list, &erp_models.LogisticsModeTypeInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberBatchInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberLogInfo{})
	model_list = append(model_list, &erp_models.LogisticsSettingInfo{})
	model_list = append(model_list, &erp_models.MelTollSeviceIdInfo{})
	model_list = append(model_list, &erp_models.MessageInfo{})
	model_list = append(model_list, &erp_models.OrderCustomTypeInfo{})
	model_list = append(model_list, &erp_models.OrderHoldTypeInfo{})
	model_list = append(model_list, &erp_models.OrderInfo{})
	model_list = append(model_list, &erp_models.OrderHistoryInfo{})
	model_list = append(model_list, &erp_models.OrderItemInfo{})
	model_list = append(model_list, &erp_models.OrderLogInfo{})
	model_list = append(model_list, &erp_models.OrderLogisticsInfo{})
	model_list = append(model_list, &erp_models.OrderLogisticsLogicInfo{})
	model_list = append(model_list, &erp_models.OrderLogisticsMatchDebugInfo{})
	model_list = append(model_list, &erp_models.OrderShippingAddressInfo{})
	model_list = append(model_list, &erp_models.OrderTagInfo{})
	model_list = append(model_list, &erp_models.OrderTypeInfo{})
	model_list = append(model_list, &erp_models.PackageMaterialInfo{})
	model_list = append(model_list, &erp_models.PaypalAccountInfo{})
	model_list = append(model_list, &erp_models.PaypalReconInfo{})
	model_list = append(model_list, &erp_models.PickFeeInfo{})
	model_list = append(model_list, &erp_models.PlatformLogisticsModeInfo{})
	model_list = append(model_list, &erp_models.ProductSaleDayAnalysisInfo{})
	model_list = append(model_list, &erp_models.ProductSaleAnalysisInfo{})
	model_list = append(model_list, &erp_models.ProductSaleAnalysisSettingInfo{})
	model_list = append(model_list, &erp_models.ProductBindInfo{})
	model_list = append(model_list, &erp_models.ProductBindLogInfo{})
	model_list = append(model_list, &erp_models.ProductCateInfo{})
	model_list = append(model_list, &erp_models.ProductCateLogInfo{})
	model_list = append(model_list, &erp_models.ProductImageInfo{})
	model_list = append(model_list, &erp_models.ProductImageTypeInfo{})
	model_list = append(model_list, &erp_models.ProductImportInfo{})
	model_list = append(model_list, &erp_models.ProductInfo{})
	model_list = append(model_list, &erp_models.ProductLogInfo{})
	model_list = append(model_list, &erp_models.ProductStoreItemInfo{})
	model_list = append(model_list, &erp_models.PurchaseBillEntryInfo{})
	model_list = append(model_list, &erp_models.PurchaseBillInfo{})
	model_list = append(model_list, &erp_models.RegionInfo{})
	model_list = append(model_list, &erp_models.SalesManagerInfo{})
	model_list = append(model_list, &erp_models.SeaContainerCostInfo{})
	model_list = append(model_list, &erp_models.ShopAccountInfo{})
	model_list = append(model_list, &erp_models.ShopAccountLogInfo{})
	model_list = append(model_list, &erp_models.SpecialProductInfo{})
	model_list = append(model_list, &erp_models.StatusInfo{})
	model_list = append(model_list, &erp_models.StockBillEntryInfo{})
	model_list = append(model_list, &erp_models.StockBillInfo{})
	model_list = append(model_list, &erp_models.StockBillTypeInfo{})
	model_list = append(model_list, &erp_models.StockInfo{})
	model_list = append(model_list, &erp_models.StockPickCrewInfo{})
	model_list = append(model_list, &erp_models.StockTypeInfo{})
	model_list = append(model_list, &erp_models.StorageLocInfo{})
	model_list = append(model_list, &erp_models.SupplierInfo{})
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
