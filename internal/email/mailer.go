package email

import (
	"github.com/RCSE2025/backend-go/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailer(emailConfig config.EmailConfig) *Mailer {
	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthEmail,
		emailConfig.AuthPassword,
	)

	return &Mailer{
		from:   emailConfig.From,
		dialer: dialer,
	}
}

func (m *Mailer) SendMail(toEmail string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.from)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	err := m.dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
