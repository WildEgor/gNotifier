package adapters

import (
	"fmt"
	"net/smtp"

	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/WildEgor/gNotifier/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ISMTPAdapter interface {
	Send(req *domain.EmailNotification) (err error)
}

type SMTPAdapter struct {
	config *config.SMTPConfig
}

func NewSMTPAdapter(
	config *config.SMTPConfig,
) *SMTPAdapter {
	return &SMTPAdapter{
		config: config,
	}
}

func (s *SMTPAdapter) Send(req *domain.EmailNotification) (err error) {

	err = domain.ValidateEmailNotification(req)
	if err != nil {
		log.Println("[SMTPAdapter] Not valid email notification: " + err.Error())
		return
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	address := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)

	err = smtp.SendMail(
		address,
		auth,
		s.config.From,
		[]string{req.Email},
		[]byte(req.Message),
	)
	if err != nil {
		log.Println("[SMTPAdapter] Cannot send message: " + err.Error())
		return err
	}

	return nil
}
