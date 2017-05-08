package controllers

import (
	"fmt"
	"time"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {
	/*
	var db_session = c.DbSession

	_, err := db_session.Exec("delete from erp_test")
	core.HandleError(err)

	paid_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", "2017-05-08T18:31:43.000Z", time.Local)
	if err != nil {
		fmt.Println(err)
	}

	var test_info models.TestInfo

	test_info.PaidTime = core.MakeNullLongTime(paid_time, true)

	_, err = db_session.Insert(&test_info)
	core.HandleError(err)

	var test_info_list = make([]models.TestInfo, 0)
	err = db_session.Find(&test_info_list)
	core.HandleError(err)

	fmt.Println(test_info_list)

	return c.RenderJson(test_info_list)
	*/

	var db_session = c.DbSession

	_, err := db_session.Exec("delete from erp_test")
	core.HandleError(err)

	paid_time, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", "2017-05-08T18:31:43.000Z", time.UTC)
	if err != nil {
		fmt.Println(err)
	}

	var test_info models.TestInfo

	test_info.PaidTime = core.MakeNullLongTime(paid_time, true)

	_, err = db_session.Insert(&test_info)
	core.HandleError(err)

	var test_info_list = make([]models.TestInfo, 0)
	err = db_session.Find(&test_info_list)
	core.HandleError(err)

	fmt.Println(test_info_list)

	return c.RenderJson(test_info_list)
}
