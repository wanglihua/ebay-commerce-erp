package ebay_token_status_check_service

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"time"
	"matrix/modules/erp/service/ebay_service"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models"

	"matrix/core"
	"fmt"
)

func Check_ebay_token_status() {
	task_db_driver, _ := revel.Config.String("db.driver")

	task_db_url, _ := revel.Config.String("db.url")

	db_engine, err := xorm.NewEngine(task_db_driver, task_db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	ebay_api_config_list := get_ebay_api_config_list(db_session)

	//todo: 如果发现 token 有问题 或者 即将过期， 加个提示消息进行提示

	for _, config := range ebay_api_config_list {
		response := ebay_service.Get_token_status(config.APIConfig)

		var expiration_time time.Time
		if response.TokenStatus.ExpirationTime != "" {

			expiration_time, err = time.ParseInLocation("2006-01-02T15:04:05.000Z", response.TokenStatus.ExpirationTime, time.UTC)
			core.HandleError(err)

			expiration_time = expiration_time.Local()
		}

		// update shop account token_result and token_expire

		var shop_account xtpz_models.ShopAccountInfo
		_, err := db_session.Id(config.APIConfig.ShopId).Get(&shop_account)
		core.HandleError(err)

		var update_col_list = make([]string, 0)
		if expiration_time.IsZero() == false {
			shop_account.TokenExpire = core.MakeNullLongTime(expiration_time, true)
			update_col_list = append(update_col_list, "token_expire")
		}

		if response.TokenStatus.Status != "" {
			shop_account.TokenResult = core.MakeNullString(response.TokenStatus.Status, true)
			update_col_list = append(update_col_list, "token_result")
		}

		if len(update_col_list) != 0 {
			_, err := db_session.Id(shop_account.Id).Cols(update_col_list...).Update(&shop_account)
			core.HandleError(err)
		}

		switch response.Ack {
		case "Success":
			revel.TRACE.Printf("eBay GetTokenStatus Response: " + response.Ack + "\n")
		case "Failure":
			revel.ERROR.Printf("eBay GetTokenStatus Response: " + response.Ack + "\n")
			continue
		default:
			revel.ERROR.Printf("eBay GetTokenStatus Response: " + response.Ack + "\n")
			continue
		}
	}
}

func get_ebay_api_config_list(db_session *xorm.Session) []Shop_account_api_info {
	shop_account_api_list := make([]Shop_account_api_info, 0)

	/*
SELECT cfg.*, sa.*
FROM erp_ebay_api_config cfg
    INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id
WHERE sa.status_id = 1400
	*/
	sql := `
SELECT cfg.*, sa.*
FROM erp_ebay_api_config cfg
    INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id
WHERE sa.status_id = %d`
	err := db_session.Sql(fmt.Sprintf(sql, models.SHOP_ACCOUNT_STATUS_ENABLED)).Find(&shop_account_api_list)
	core.HandleError(err)

	return shop_account_api_list
}

type Shop_account_api_info struct {
	APIConfig xtpz_models.EBayAPIConfigInfo `xorm:"extends"`
	Shop      xtpz_models.ShopAccountInfo   `xorm:"extends"`
}
