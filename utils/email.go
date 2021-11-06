package utils

import (
	"io"
	"sea/auth/config"

	"go.uber.org/zap"
	"gopkg.in/mail.v2"
)

type EmailSender struct {
	d mail.Sender
}

func (e EmailSender) Send(from string, to []string, msg io.WriterTo) error {
	return e.d.Send(from, to, msg)
}

func NewEmailSender(conf config.AppConfig, sh *Shutdown) (*EmailSender, error) {
	d := mail.NewDialer(conf.EmailConfig.Host, conf.EmailConfig.Port, conf.EmailConfig.User, conf.EmailConfig.Password)
	if conf.EmailConfig.MandatoryTLS {
		d.StartTLSPolicy = mail.MandatoryStartTLS
	} else {
		d.StartTLSPolicy = mail.NoStartTLS
	}
	close, err := d.Dial()
	if err != nil {
		zap.S().Errorw("could not dial to email service", "error", err)
		return nil, err
	}
	sh.Add(close.Close)
	return &EmailSender{close}, nil
}

func TestEmail(d *EmailSender) {
	m := mail.NewMessage()
	m.SetHeader("From", "sea@auca.kg")
	m.SetHeader("To", "student_s@auca.kg")
	m.SetHeader("Subject", "Test mailhog email")
	m.SetBody("text/html", "Hello! <hr/> This is a test email")

	d.Send("sea@auca.kg", []string{"student_s@auca.kg"}, m)
}
