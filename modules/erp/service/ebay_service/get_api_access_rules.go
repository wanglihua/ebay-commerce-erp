package ebay_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/core"
	"fmt"
	"github.com/franela/goreq"
	"github.com/revel/revel"
	"io/ioutil"
)

func Get_api_access_rules(shop_id int64, db_session *xorm.Session) string {

	var config xtpz_models.EBayAPIConfigInfo
	_, err := db_session.Where("shop_id = ?", shop_id).Limit(1).Get(&config)
	core.HandleError(err)

	/*

<?xml version="1.0" encoding="utf-8"?>
<GetApiAccessRulesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
  <!-- (No call-specific Input fields) -->

  <!-- Standard Input Fields -->
  <ErrorLanguage> string </ErrorLanguage>
  <MessageID> string </MessageID>
  <Version> string </Version>
  <WarningLevel> WarningLevelCodeType </WarningLevel>
</GetApiAccessRulesRequest>
	 */

	var req_body = `<?xml version="1.0" encoding="utf-8"?>
<GetApiAccessRulesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
</GetApiAccessRulesRequest>
`
	req_body = fmt.Sprintf(req_body, config.Token)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	add_request_header("GetApiAccessRules", &req, config)

	res, err := req.Do()

	if err != nil {
		revel.ERROR.Println("eBay GetApiAccessRules ERROR - processUrl -> req.Do: " + err.Error())
		return err.Error()
	}

	data, err := ioutil.ReadAll(res.Body)

	return string(data)
}
