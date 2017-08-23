package ebay_service

import (
	"github.com/franela/goreq"
	"matrix/modules/erp/models/xtpz_models"
)

const log_response_xml_to_file bool = false

func add_request_header(call_name string, request *goreq.Request, config xtpz_models.EBayAPIConfigInfo) {
	request.AddHeader("X-EBAY-API-CALL-NAME", call_name)
	request.AddHeader("X-EBAY-API-DEV-NAME", config.DevId)
	request.AddHeader("X-EBAY-API-CERT-NAME", config.CertId)
	request.AddHeader("X-EBAY-API-APP-NAME", config.AppId)
	request.AddHeader("X-EBAY-API-VERSION", config.GetApiVersion())
	request.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", config.GetApiVersion())
	request.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	request.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	request.AddHeader("X-EBAY-API-SITEID", config.GetSiteId())
	request.AddHeader("Accept", "application/xml,application/xhtml+xml")
}
