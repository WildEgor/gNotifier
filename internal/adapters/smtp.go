package adapters

import (
	"fmt"
	"strings"

	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/WildEgor/gNotifier/internal/domain"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
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

	address := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)
	auth := sasl.NewPlainClient("", s.config.Username, s.config.Password)

	to := []string{req.Email}
	msg := strings.NewReader(
		"To: " +
			req.Email +
			"\r\n" +
			"Subject: " +
			req.Subject +
			"\r\n" +
			req.Message +
			"\r\n",
	)

	err = smtp.SendMail(address, auth, s.config.From, to, msg)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
