package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var templatePath = "internal/resources/views/email/reset-password.html"

type EmailHelper struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPasswd string
}

func New() EmailHelper {
	return EmailHelper{
		SMTPHost:     os.Getenv("MAIL_HOST"),
		SMTPPort:     os.Getenv("MAIL_PORT"),
		SenderEmail:  os.Getenv("MAIL_USERNAME"),
		SenderPasswd: os.Getenv("MAIL_PASSWORD"),
	}
}

type PasswordResetEmailInput struct {
	To       string
	Subject  string
	UserName string
	Url      string
}

func (e EmailHelper) SendPasswordResetEmail(input PasswordResetEmailInput) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	htmlContent := string(content)
	htmlContent = strings.ReplaceAll(htmlContent, "{{ .Username }}", input.UserName)
	htmlContent = strings.ReplaceAll(htmlContent, "{{ .Url }}", input.Url)

	err = e.sendHtmlEmail(input.To, input.Subject, htmlContent)
	if err != nil {
		return err
	}

	return nil
}

func (e EmailHelper) sendHtmlEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", e.SenderEmail, e.SenderPasswd, e.SMTPHost)

	headers := make(map[string]string)
	headers["From"] = e.SenderEmail
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message string
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	message += "\r\n" + htmlBody

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort),
		auth,
		e.SenderEmail,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		return err
	}

	return nil
}
