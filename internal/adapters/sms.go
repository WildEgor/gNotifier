package adapters

import (
	"fmt"
	"github.com/WildEgor/gNotifier/internal/configs"
	"net/http"
	"net/url"

	"github.com/WildEgor/gNotifier/internal/domain"
	log "github.com/sirupsen/logrus"
)

type ISMSAdapter interface {
	Send(req *domain.SMSNotification) (err error)
}

type SMSAdapter struct {
	config *configs.SMSConfig
}

func NewSMSAdapter(
	config *configs.SMSConfig,
) *SMSAdapter {
	return &SMSAdapter{
		config: config,
	}
}

// Send implement own logic here
func (s *SMSAdapter) Send(notification *domain.SMSNotification) (err error) {

	err = domain.ValidateSMSNotification(notification)
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
		"recipient":   {notification.Phone},
		"messagetype": {"SMS:TEXT"},
		"originator":  {""},
		"messagedata": {notification.Message},
	}

	requestUrl := fmt.Sprintf("%v?%v", baseURL, queryParams.Encode())

	_, err = http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Errorf("[SMSAdapter] Creating the request failed: %w", err)
		return nil
	}

	return nil
}
