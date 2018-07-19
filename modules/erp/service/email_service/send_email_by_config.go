package email_service

import (
	"strings"
	"net/smtp"
	"fmt"
	"matrix/core"
)

func send_mail(host string, port string, user string, password string, from string, to string, subject string, body string) {
	auth := smtp.PlainAuth("", user, password, host)

	var content_type = "Content-Type:text/html;charset=UTF-8"

	var msg = []byte("To: " + to + "\r\nFrom: " + from + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	var send_to_mails = strings.Split(to, ";")

	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, user, send_to_mails, msg)
	core.HandleError(err)
}
