package email_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"gopkg.in/gomail.v2"
)

type (
	EmailRequest struct {
		Receiver string
		HTMLPath string
		Subject  string
	}

	Email interface {
		Send(ctx context.Context, request *EmailRequest) (err error)
	}

	emailClient struct {
		gomail *gomail.Dialer
	}
)

func NewEmail(config *configs.Config) (*emailClient, error) {
	emailConfig := config.Config.Email
	gomail := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.Sender.Email,
		emailConfig.Sender.Password,
	)
	return &emailClient{
		gomail: gomail,
	}, nil
}

func (e *emailClient) Send(ctx context.Context, request *EmailRequest) (err error) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "")
	mailer.SetHeader("To", request.Receiver)
	mailer.SetHeader("Subject", request.Subject)
	mailer.SetBody("text/html", request.HTMLPath)

	err = e.gomail.DialAndSend(mailer)
	if err != nil {
		return err
	}
	return nil
}
