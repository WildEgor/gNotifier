package adapters

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/WildEgor/gNotifier/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ISMSAdapter interface {
	Send(req *domain.SMSNotification) (err error)
}

type SMSAdapter struct {
	config *config.SMSConfig
}

func NewSMSAdapter(
	config *config.SMSConfig,
) *SMSAdapter {
	return &SMSAdapter{
		config: config,
	}
}

func (s *SMSAdapter) Send(req *domain.SMSNotification) (err error) {

	err = domain.ValidateSMSNotification(req)
	if err != nil {
		log.Println("[SMSAdapter] Not valid sms notification: " + err.Error())
		return
	}

	baseURL := url.URL{
		Scheme: "https",
		Host:   s.config.BaseURL,
	}

	queryParams := url.Values{
		"action":      {"sendmessage"},
		"username":    {s.config.Username},
		"password":    {s.config.Password},
		"recipient":   {req.Phone},
		"messagetype": {"SMS:TEXT"},
		"originator":  {req.TeamID},
		"messagedata": {req.Message},
	}

	requestUrl := fmt.Sprintf("%v?%v", baseURL, queryParams.Encode())

	_, err = http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Errorf("[SMSAdapter] Creating the request failed: %w", err)
		return nil
	}

	return nil
}
