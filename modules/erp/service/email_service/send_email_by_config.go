package email_service

import (
	"strings"
	"net/smtp"
	"fmt"
	"matrix/core"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/xtpz_models"
)

func send_mail(host string, port string, user string, password string, from string, to string, subject string, body string) {
	auth := smtp.PlainAuth("", user, password, host)

	var content_type = "Content-Type:text/html;charset=UTF-8"

	var msg = []byte("To: " + to + "\r\nFrom: " + from + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	var send_to_mails = strings.Split(to, ";")

	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, user, send_to_mails, msg)
	core.HandleError(err)
}

func Send_mail_by_config(mail_to string, subject string, body string, email_address_id int64, db_session *xorm.Session) {
	var email_address xtpz_models.EmailAddressInfo
	_, err := db_session.Id(email_address_id).Get(&email_address)
	core.HandleError(err)

	var host = email_address.SmtpServer
	var port = email_address.SmtpPort
	var account = email_address.Account
	var password = email_address.Password
	var email = email_address.Email

	send_mail(host, port, account, password, email, mail_to, subject, body)
}
