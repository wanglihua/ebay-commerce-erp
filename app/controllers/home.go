package controllers

import (
	"matrix/app/models"
	"matrix/core"
	"matrix/db"
	erp_models "matrix/modules/erp/models"
	auth_models "matrix/modules/auth/models"

	"github.com/revel/revel"
)

type Home struct {
	*revel.Controller
}

func (c Home) Index() revel.Result {
	return c.RenderTemplate("home/home_index.html")
}

func (c Home) SyncDb() revel.Result {
	return c.RenderTemplate("home/home_syncdb.html")
}

func (c Home) SyncDbPost() revel.Result {

	model_list := make([]interface{}, 0)

	model_list = append(model_list, &models.EntityCodeInfo{})

	model_list = append(model_list, &auth_models.UserInfo{})
	model_list = append(model_list, &auth_models.AdminInfo{})
	model_list = append(model_list, &auth_models.GroupInfo{})
	model_list = append(model_list, &auth_models.GroupUserInfo{})

	model_list = append(model_list, &erp_models.StatusInfo{})
	model_list = append(model_list, &erp_models.BrandInfo{})
	model_list = append(model_list, &erp_models.StockInfo{})
	model_list = append(model_list, &erp_models.StorageLocInfo{})
	model_list = append(model_list, &erp_models.SupplierInfo{})
	model_list = append(model_list, &erp_models.ExchangeRateInfo{})
	model_list = append(model_list, &erp_models.ProductCateInfo{})
	model_list = append(model_list, &erp_models.EubGoodsTypeInfo{})
	model_list = append(model_list, &erp_models.EubInsuranceTypeInfo{})
	model_list = append(model_list, &erp_models.LoginLogInfo{})

	model_list = append(model_list, &erp_models.ShopAccountInfo{})
	model_list = append(model_list, &erp_models.PaypalAccountInfo{})
	model_list = append(model_list, &erp_models.LogisticsModeInfo{})

	model_list = append(model_list, &erp_models.PackageMaterialInfo{})
	model_list = append(model_list, &erp_models.ProductInfo{})
	model_list = append(model_list, &erp_models.ProductImageInfo{})
	model_list = append(model_list, &erp_models.ProductImageTypeInfo{})
	model_list = append(model_list, &erp_models.ProductLogInfo{})
	model_list = append(model_list, &erp_models.OrderLogInfo{})
	model_list = append(model_list, &erp_models.InventoryInfo{})
	model_list = append(model_list, &erp_models.InventoryCheckInfo{})
	model_list = append(model_list, &erp_models.StockBillTypeInfo{})
	model_list = append(model_list, &erp_models.StockBillInfo{})
	model_list = append(model_list, &erp_models.StockBillEntryInfo{})

	model_list = append(model_list, &erp_models.ECommercePlatformInfo{})
	model_list = append(model_list, &erp_models.PaypalReconInfo{})
	model_list = append(model_list, &erp_models.OrderCustomTypeInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberBatchInfo{})
	model_list = append(model_list, &erp_models.LogisticsNumberLogInfo{})
	model_list = append(model_list, &erp_models.ProductBindInfo{})
	model_list = append(model_list, &erp_models.ProductImportInfo{})

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
