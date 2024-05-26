package mail

import (
	"github.com/yunsonggo/helper/v3/types"
	"gopkg.in/gomail.v2"
)

func SendEmail(conf *types.Mail, fromName, toEmail, subject, body string) (err error) {
	fromEmail := conf.MailUser
	smtpAddr := conf.MailHost
	mailHeader := map[string][]string{
		"To":      {toEmail},
		"Subject": {subject},
	}
	m := gomail.NewMessage()
	m.SetHeaders(mailHeader)
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpAddr, conf.MailPort, fromEmail, conf.MailPass)
	return d.DialAndSend(m)
}
