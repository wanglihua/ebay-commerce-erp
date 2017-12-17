package ebay

import "encoding/xml"

/*

<?xml version="1.0" encoding="utf-8"?>
<GetTokenStatusResponse xmlns="urn:ebay:apis:eBLBaseComponents">
  <!-- Call-specific Output Fields -->
  <TokenStatus> TokenStatusType
    <EIASToken> string </EIASToken>
    <ExpirationTime> dateTime </ExpirationTime>
    <RevocationTime> dateTime </RevocationTime>
    <Status> TokenStatusCodeType </Status>
  </TokenStatus>
  <!-- Standard Output Fields -->
  <Ack> AckCodeType </Ack>
  <Build> string </Build>
  <CorrelationID> string </CorrelationID>
  <Errors> ErrorType
    <ErrorClassification> ErrorClassificationCodeType </ErrorClassification>
    <ErrorCode> token </ErrorCode>
    <ErrorParameters ParamID="string"> ErrorParameterType
      <Value> string </Value>
    </ErrorParameters>
    <!-- ... more ErrorParameters nodes allowed here ... -->
    <LongMessage> string </LongMessage>
    <SeverityCode> SeverityCodeType </SeverityCode>
    <ShortMessage> string </ShortMessage>
  </Errors>
  <!-- ... more Errors nodes allowed here ... -->
  <HardExpirationWarning> string </HardExpirationWarning>
  <Timestamp> dateTime </Timestamp>
  <Version> string </Version>
</GetTokenStatusResponse>

 */

type GetTokenStatusResponse struct {
	Xmlns       string          `xml:"xmlns,attr"`
	Timestamp   string          `xml:"Timestamp"`
	Ack         string          `xml:"Ack"`
	Build       string          `xml:"Build"`
	Errors      Errors          `xml:"Errors"`
	TokenStatus TokenStatus     `xml:"TokenStatus"`
}

type TokenStatus struct {
	XMLName        xml.Name   `xml:"TokenStatus"`
	EIASToken      string     `xml:"EIASToken"`
	ExpirationTime string     `xml:"ExpirationTime"`
	RevocationTime string     `xml:"RevocationTime"`
	Status         string     `xml:"Status"`
}
