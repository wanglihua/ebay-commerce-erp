package tasks

import (
	"runtime/debug"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/ebay_service"
)

type EBayTokenStatusCheckTask struct {
}

func (t EBayTokenStatusCheckTask) Name() string {
	return "店铺Token状态检查"
}

func (t EBayTokenStatusCheckTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("检查Token状态 出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	job_time_begin := time.Now()

	if task_db_driver == "" {
		task_db_driver, _ = revel.Config.String("db.driver")
	}

	if task_db_url == "" {
		task_db_url, _ = revel.Config.String("db.url")
	}

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
		response := ebay_service.Get_token_status(config)

		var expiration_time time.Time
		if response.TokenStatus.ExpirationTime != "" {

			expiration_time, err = time.ParseInLocation("2006-01-02T15:04:05.000Z", response.TokenStatus.ExpirationTime, time.UTC)
			core.HandleError(err)

			expiration_time = expiration_time.Local()
		}

		// update shop account token_result and token_expire

		var shop_account models.ShopAccountInfo
		_, err := db_session.Id(config.ShopId).Get(&shop_account)
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

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("eBay GetTokenStatus time duration: ", job_duration.Seconds(), "seconds")
}
