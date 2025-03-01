package main

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
)

type EmailConfig struct {
	AuthEmail    string
	AuthPassword string
	Host         string
	Port         int
}

const emailTemplate = `
<html>
    <body style="font-family: Arial, sans-serif; color: #333; background-color: #f9f9f9; padding: 20px;">
        <div style="max-width: 600px; margin: auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">
            <h1 style="color: #4CAF50;">Hello, {{.Name}}!</h1>
            <p style="font-size: 16px; line-height: 1.5;">This is a test email with styled HTML content.</p>
            <p style="font-size: 14px; color: #666;">Thank you for using our service. We hope you enjoy this email layout!</p>
            <a href="{{.Link}}" style="display: inline-block; margin-top: 20px; padding: 10px 15px; color: #fff; background-color: #4CAF50; text-decoration: none; border-radius: 5px;">Learn More</a>
        </div>
    </body>
</html>
`

func main() {
	emailConfig := EmailConfig{
		AuthEmail:    "info@ryazan-market.ru",
		AuthPassword: "cN7fS0kK7tdI3mE1",
		Host:         "mail.hosting.reg.ru",
		Port:         587,
	}

	mailer := NewMailer(emailConfig, "info@ryazan-market.ru")
	body, _ := makeVerificationEmailTemplate("John Doe", "123456")
	err := mailer.SendMail("ivan.voronin.25@mail.ru", "Подтверждение почты", body)
	if err != nil {
		panic(err)
	}
}

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailer(emailConfig EmailConfig, from string) *Mailer {
	return &Mailer{
		from: from,
		dialer: gomail.NewDialer(
			emailConfig.Host,
			emailConfig.Port,
			emailConfig.AuthEmail,
			emailConfig.AuthPassword,
		),
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

func makeEmailTemplate(name string, link string) (string, error) {

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	// Execute the template with dynamic data
	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Name string
		Link string
	}{Name: name, Link: link})

	return body.String(), err
}

const verificationEmailTemplate = `
<html>
    <body style="font-family: Arial, sans-serif; color: #333; background-color: #f9f9f9; padding: 20px;">
        <div style="max-width: 600px; margin: auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); text-align: center;">
            <h1 style="color: #4CAF50;">Подтверждение почты</h1>
            <p style="font-size: 16px; line-height: 1.5;">Здравствуйте, {{.Name}}!</p>
            <p style="font-size: 16px; line-height: 1.5;">Ваш код подтверждения:</p>
            <p style="font-size: 24px; font-weight: bold; color: #4CAF50;">{{.Code}}</p>
            <p style="font-size: 14px; color: #666;">Введите этот код в соответствующее поле, чтобы подтвердить свою электронную почту.</p>
        </div>
    </body>
</html>`

func makeVerificationEmailTemplate(name string, code string) (string, error) {
	tmpl, err := template.New("email").Parse(verificationEmailTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	// Execute the template with dynamic data
	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Name string
		Code string
	}{Name: name, Code: code})

	return body.String(), err
}
