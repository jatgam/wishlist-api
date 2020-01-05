package sgmail

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type SGConfig struct {
	Client      *sendgrid.Client
	FromName    string
	FromAddress string
	Enabled     bool
	Debug       bool
}

var sgMailer *SGConfig

func SetupMail(apikey, fromName, fromAddress string, debug bool) {
	newMailer := &SGConfig{FromName: fromName, FromAddress: fromAddress, Debug: debug}
	if apikey != "" {
		newMailer.Client = sendgrid.NewSendClient(apikey)
		newMailer.Enabled = true
	} else {
		newMailer.Enabled = false
	}
	sgMailer = newMailer
}

func GetMailer() *SGConfig {
	return sgMailer
}

func (sg *SGConfig) SendMail(toAddress, subject, message string, logger *logrus.Entry) error {
	var fields logrus.Fields
	if sg.Debug {
		fields = logrus.Fields{
			"to":      toAddress,
			"subject": subject,
			"message": message,
		}
	} else {
		fields = logrus.Fields{
			"to": toAddress,
		}
	}
	if sg.Enabled {
		from := mail.NewEmail(sg.FromName, sg.FromAddress)
		to := mail.NewEmail("", toAddress)
		textContent := mail.NewContent("text/plain", message)
		mailToSend := mail.NewV3MailInit(from, subject, to, textContent)
		_, err := sg.Client.Send(mailToSend)
		if err != nil {
			logger.WithFields(fields).Error("Failed to send EMail")
			return err
		}
	} else {
		logger.WithFields(fields).Warn("SGMail Disabled: No mail sent")
	}
	logger.WithFields(fields).Info("Email Sent")
	return nil
}
