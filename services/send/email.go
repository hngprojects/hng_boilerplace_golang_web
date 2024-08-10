package send

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/config"
)

type EmailRequest struct {
	ExtReq         request.ExternalRequest
	To             []string `json:"to"`
	Subject        string   `json:"subject"`
	Body           string   `json:"body"`
	AttachmentName string
	Attachment     []byte
}

func NewEmailRequest(extReq request.ExternalRequest, to []string, subject, templateFileName, baseTemplateFileName string, templateData map[string]interface{}) (*EmailRequest, error) {
	body, err := ParseTemplate(extReq, templateFileName, baseTemplateFileName, templateData)
	if err != nil {
		return &EmailRequest{}, err
	}
	return &EmailRequest{
		ExtReq:  extReq,
		To:      to,
		Subject: subject,
		Body:    body, //or parsed template
	}, nil
}

func NewSimpleEmailRequest(extReq request.ExternalRequest, to []string, subject, body string) *EmailRequest {
	return &EmailRequest{
		ExtReq:  extReq,
		To:      to,
		Subject: subject,
		Body:    body, //or parsed template
	}
}

func SendEmail(extReq request.ExternalRequest, to string, subject, templateFileName, baseTemplateFileName string, data map[string]interface{}) error {
	mailRequest, err := NewEmailRequest(extReq, []string{to}, subject, templateFileName, baseTemplateFileName, data)
	if err != nil {
		return fmt.Errorf("error getting email request, %v", err)
	}

	err = mailRequest.Send()
	if err != nil {
		return fmt.Errorf("error sending email, %v", err)
	}
	return nil
}

func (e EmailRequest) validate() error {
	if e.Subject == "" {
		return fmt.Errorf("EMAIL::validate ==> subject is required")
	}
	if e.Body == "" {
		return fmt.Errorf("EMAIL::validate ==> body is required")
	}

	if e.To == nil {
		return fmt.Errorf("receiving email is empty")
	}

	for _, v := range e.To {
		if v == "" {
			return fmt.Errorf("receiving email is empty: %s", v)
		}

		if !strings.Contains(v, "@") {
			return fmt.Errorf("receiving email is invalid: %s", v)
		}
	}

	return nil
}

func (e *EmailRequest) Send() error {

	if err := e.validate(); err != nil {
		return err
	}

	if e.ExtReq.Test {
		return nil
	}

	err := e.sendEmailViaSMTP()

	if err != nil {
		e.ExtReq.Logger.Error("error sending email: ", err.Error())
		return err
	}
	return nil
}

func (e *EmailRequest) sendEmailViaSMTP() error {
	var (
		mailConfig = config.GetConfig().Mail
	)

	auth := smtp.PlainAuth(
		"",
		mailConfig.Username,
		mailConfig.Password,
		mailConfig.Server,
	)

	sender := mailConfig.Username
	subject := e.Subject
	recipients := e.To
	mime := "\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n%s%s%s", sender, recipients[0], subject, mime, e.Body))

	conn, err := tls.Dial(
		"tcp",
		mailConfig.Server+":"+mailConfig.Port,
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         mailConfig.Server,
		})

	if err != nil {

		return fmt.Errorf("failed to connect to the server: %v", err)

	}

	defer conn.Close()

	client, err := smtp.NewClient(conn, mailConfig.Server)

	if err != nil {

		return fmt.Errorf("failed to create SMTP client: %v", err)

	}

	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)

	}

	if err = client.Mail(sender); err != nil {
		return fmt.Errorf("failed to set the sender: %v", err)

	}

	if err = client.Rcpt(recipients[0]); err != nil {
		return fmt.Errorf("failed to set the recipient: %v", err)

	}

	writer, err := client.Data()
	if err != nil {

		return fmt.Errorf("failed to write the message: %v", err)

	}

	_, err = writer.Write(body)
	if err != nil {

		return fmt.Errorf("failed to send the message: %v", err)

	}

	err = writer.Close()
	if err != nil {

		return fmt.Errorf("failed to close the writer: %v", err)

	}

	return nil

}
