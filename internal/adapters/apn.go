package adapters

import (
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/WildEgor/gNotifier/internal/configs"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mitchellh/mapstructure"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"

	"github.com/WildEgor/gNotifier/internal/domain"
)

var (
	idleConnTimeout = 90 * time.Second
	tlsDialTimeout  = 20 * time.Second
	tcpKeepAlive    = 60 * time.Second
	doOnce          sync.Once
	// DialTLS is the default dial function for creating TLS connections for
	// non-proxied HTTPS requests.
	DialTLS = func(cfg *tls.Config) func(network, addr string) (net.Conn, error) {
		return func(network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   tlsDialTimeout,
				KeepAlive: tcpKeepAlive,
			}
			return tls.DialWithDialer(dialer, network, addr, cfg)
		}
	}
)

const (
	dotP8  = ".p8"
	dotPEM = ".pem"
	dotP12 = ".p12"
)

var maxConcurrentIOSPushes chan struct{}

// Sound sets the aps sound on the payload.
type Sound struct {
	Critical int     `json:"critical,omitempty"`
	Name     string  `json:"name,omitempty"`
	Volume   float32 `json:"volume,omitempty"`
}

type IAPNAdapter interface {
	Send(domain *domain.PushNotification) error
}

type APNAdapter struct {
	client *apns2.Client
	config *configs.APNConfig
}

func NewAPNAdapter(
	config *configs.APNConfig,
) *APNAdapter {
	var err error
	var authKey *ecdsa.PrivateKey
	var certificateKey tls.Certificate
	var ext string
	var client *apns2.Client

	// Try read config file (cert)
	if config.KeyPath != "" {
		ext = filepath.Ext(config.KeyPath)

		// Serve formats
		switch ext {
		case dotP12:
			certificateKey, err = certificate.FromP12File(config.KeyPath, config.Password)
		case dotPEM:
			certificateKey, err = certificate.FromPemFile(config.KeyPath, config.Password)
		case dotP8:
			authKey, err = token.AuthKeyFromFile(config.KeyPath)
		default:
			err = errors.New("[APNAdapter] Wrong certificate key extension")
		}

		if err != nil && config.Production {
			log.Fatal("[APNAdapter] Cert Error:", err.Error())
		}

		// Else try parse base64 content
	} else if config.KeyBase64 != "" {
		ext = "." + config.KeyType

		key, err := base64.StdEncoding.DecodeString(config.KeyBase64)
		if err != nil && config.Production {
			log.Fatal("[APNAdapter] Base64 decode error:", err.Error())
		}

		// Serve extension
		switch ext {
		case dotP12:
			certificateKey, err = certificate.FromP12Bytes(key, config.Password)
		case dotPEM:
			certificateKey, err = certificate.FromPemBytes(key, config.Password)
		case dotP8:
			authKey, err = token.AuthKeyFromBytes(key)
		default:
			err = errors.New("[APNAdapter] Wrong certificate key type")
		}

		if err != nil && config.Production {
			log.Fatal("Cert Error:", err.Error())
		}
	}

	// If provede .p8 cert
	if ext == dotP8 {
		if config.KeyID == "" || config.TeamID == "" {
			log.Fatal("[APNAdapter] You should provide KeyID from developer account (Certificates, Identifiers & Profiles -> Keys) and TeamID from developer account (View Account -> Membership) for p8 token")
		}
		token := &token.Token{
			AuthKey: authKey,
			// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
			KeyID: config.KeyID,
			// TeamID from developer account (View Account -> Membership)
			TeamID: config.TeamID,
		}

		// Build clien
		client, err = buildApnsTokenClient(config, token)
		if err != nil && config.Production {
			log.Fatal("[APNAdapter] Failed when init new APNClient", err)
		}

	} else {
		client, err = buildApnsClient(config, certificateKey)
		if err != nil && config.Production {
			log.Fatal("[APNAdapter] Failed when init new APNClient", err)
		}
	}

	if h2Transport, ok := client.HTTPClient.Transport.(*http2.Transport); ok {
		configureHTTP2ConnHealthCheck(h2Transport)
	}

	if err != nil && config.Production {
		log.Fatal("[APNAdapter] Transport Error:", err.Error())
	}

	doOnce.Do(func() {
		maxConcurrentIOSPushes = make(chan struct{}, 5) // WARN!
	})

	return &APNAdapter{
		client: client,
		config: config,
	}
}

func buildApnsClient(cfg *configs.APNConfig, certificate tls.Certificate) (*apns2.Client, error) {
	var client *apns2.Client

	if cfg.Production {
		client = apns2.NewClient(certificate).Production()
	} else {
		client = apns2.NewClient(certificate).Development()
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS:         DialTLS(tlsConfig),
		Proxy:           http.DefaultTransport.(*http.Transport).Proxy,
		IdleConnTimeout: idleConnTimeout,
	}

	h2Transport, err := http2.ConfigureTransports(transport)
	if err != nil {
		return nil, err
	}

	configureHTTP2ConnHealthCheck(h2Transport)

	client.HTTPClient.Transport = transport

	return client, nil
}

func buildApnsTokenClient(cfg *configs.APNConfig, token *token.Token) (*apns2.Client, error) {
	var client *apns2.Client

	if cfg.Production {
		client = apns2.NewTokenClient(token).Production()
	} else {
		client = apns2.NewTokenClient(token).Development()
	}

	transport := &http.Transport{
		DialTLS:         DialTLS(nil),
		Proxy:           http.DefaultTransport.(*http.Transport).Proxy,
		IdleConnTimeout: idleConnTimeout,
	}

	h2Transport, err := http2.ConfigureTransports(transport)
	if err != nil {
		return nil, err
	}

	configureHTTP2ConnHealthCheck(h2Transport)

	client.HTTPClient.Transport = transport

	return client, nil
}

func configureHTTP2ConnHealthCheck(h2Transport *http2.Transport) {
	h2Transport.ReadIdleTimeout = 1 * time.Second
	h2Transport.PingTimeout = 1 * time.Second
}

func convertToIOSDictionary(notificationPayload *payload.Payload, req *domain.PushNotification) *payload.Payload {
	// Alert dictionary

	if len(req.Title) > 0 {
		notificationPayload.AlertTitle(req.Title)
	}

	if len(req.InterruptionLevel) > 0 {
		notificationPayload.InterruptionLevel(payload.EInterruptionLevel(req.InterruptionLevel))
	}

	if len(req.Message) > 0 && len(req.Title) > 0 {
		notificationPayload.AlertBody(req.Message)
	}

	if len(req.Alert.Title) > 0 {
		notificationPayload.AlertTitle(req.Alert.Title)
	}

	// Apple Watch & Safari display this string as part of the notification interface.
	if len(req.Alert.Subtitle) > 0 {
		notificationPayload.AlertSubtitle(req.Alert.Subtitle)
	}

	if len(req.Alert.TitleLocKey) > 0 {
		notificationPayload.AlertTitleLocKey(req.Alert.TitleLocKey)
	}

	if len(req.Alert.LocArgs) > 0 {
		notificationPayload.AlertLocArgs(req.Alert.LocArgs)
	}

	if len(req.Alert.TitleLocArgs) > 0 {
		notificationPayload.AlertTitleLocArgs(req.Alert.TitleLocArgs)
	}

	if len(req.Alert.Body) > 0 {
		notificationPayload.AlertBody(req.Alert.Body)
	}

	if len(req.Alert.LaunchImage) > 0 {
		notificationPayload.AlertLaunchImage(req.Alert.LaunchImage)
	}

	if len(req.Alert.LocKey) > 0 {
		notificationPayload.AlertLocKey(req.Alert.LocKey)
	}

	if len(req.Alert.Action) > 0 {
		notificationPayload.AlertAction(req.Alert.Action)
	}

	if len(req.Alert.ActionLocKey) > 0 {
		notificationPayload.AlertActionLocKey(req.Alert.ActionLocKey)
	}

	// General
	if len(req.Category) > 0 {
		notificationPayload.Category(req.Category)
	}

	if len(req.Alert.SummaryArg) > 0 {
		notificationPayload.AlertSummaryArg(req.Alert.SummaryArg)
	}

	if req.Alert.SummaryArgCount > 0 {
		notificationPayload.AlertSummaryArgCount(req.Alert.SummaryArgCount)
	}

	return notificationPayload
}

// ConvertToIOSNotification use for define iOS notification.
// The iOS Notification Payload (Payload Key Reference)
// Ref: https://apple.co/2VtH6Iu
func ConvertToIOSNotification(req *domain.PushNotification) *apns2.Notification {
	notification := &apns2.Notification{
		ApnsID:     req.ApnsID,
		Topic:      req.Topic,
		CollapseID: req.CollapseID,
	}

	if req.Expiration != nil {
		notification.Expiration = time.Unix(*req.Expiration, 0)
	}

	if len(req.Priority) > 0 {
		if req.Priority == "normal" {
			notification.Priority = apns2.PriorityLow
		} else if req.Priority == "high" {
			notification.Priority = apns2.PriorityHigh
		}
	}

	if len(req.PushType) > 0 {
		notification.PushType = apns2.EPushType(req.PushType)
	}

	payload := payload.NewPayload()

	// add alert object if message length > 0 and title is empty
	if len(req.Message) > 0 && req.Title == "" {
		payload.Alert(req.Message)
	}

	// zero value for clear the badge on the app icon.
	if req.Badge != nil && *req.Badge >= 0 {
		payload.Badge(*req.Badge)
	}

	if req.MutableContent {
		payload.MutableContent()
	}

	switch req.Sound.(type) {
	// from http request binding
	case map[string]interface{}:
		result := &Sound{}
		_ = mapstructure.Decode(req.Sound, &result)
		payload.Sound(result)
	// from http request binding for non critical alerts
	case string:
		payload.Sound(&req.Sound)
	case Sound:
		payload.Sound(&req.Sound)
	}

	if len(req.SoundName) > 0 {
		payload.SoundName(req.SoundName)
	}

	if req.SoundVolume > 0 {
		payload.SoundVolume(req.SoundVolume)
	}

	if req.ContentAvailable {
		payload.ContentAvailable()
	}

	if len(req.URLArgs) > 0 {
		payload.URLArgs(req.URLArgs)
	}

	if len(req.ThreadID) > 0 {
		payload.ThreadID(req.ThreadID)
	}

	for k, v := range req.Data {
		payload.Custom(k, v)
	}

	payload = convertToIOSDictionary(payload, req)

	notification.Payload = payload

	return notification
}

func (s *APNAdapter) getApnsClient(req *domain.PushNotification) (client *apns2.Client) {
	switch {
	case req.Production:
		client = s.client.Production()
	case req.Development:
		client = s.client.Development()
	default:
		if s.config.Production {
			client = s.client.Production()
		} else {
			client = s.client.Development()
		}
	}

	return client
}

// Send provide send notification to APNs server.
func (s *APNAdapter) Send(req *domain.PushNotification) (err error) {
	log.Debug("Start push notification for iOS")

	var (
		retryCount = 0
		maxRetry   = 5
	)

	if req.Retry > 0 && req.Retry < maxRetry {
		maxRetry = req.Retry
	}

	s.retry(req, retryCount)

	return err
}

func (s *APNAdapter) retry(req *domain.PushNotification, retryCount int) int {
	var newTokens []string
	var maxRetry = 5

	notification := ConvertToIOSNotification(req)
	client := s.getApnsClient(req)

	var wg sync.WaitGroup
	for _, token := range req.Tokens {
		// occupy push slot
		maxConcurrentIOSPushes <- struct{}{}
		wg.Add(1)
		go func(notification apns2.Notification, token string) {
			notification.DeviceToken = token

			// send ios notification
			res, err := client.Push(&notification)
			if err != nil || (res != nil && !res.Sent()) {
				if err == nil {
					// error message:
					// ref: https://github.com/sideshow/apns2/blob/master/response.go#L14-L65
					err = errors.New(res.Reason)
				}

				// apns server error
				s.saveLogs("fail_push", token, req, err)

				// We should retry only "retryable" statuses. More info about response:
				// See https://apple.co/3AdNane (Handling Notification Responses from APNs)
				if res != nil && res.StatusCode >= http.StatusInternalServerError {
					newTokens = append(newTokens, token)
				}
			}

			if res != nil && res.Sent() {
				s.saveLogs("success_push", token, req, nil)
			}

			// free push slot
			<-maxConcurrentIOSPushes
			wg.Done()
		}(*notification, token)
	}

	wg.Wait()

	if len(newTokens) > 0 && retryCount < maxRetry {
		retryCount++

		// resend fail token
		req.Tokens = newTokens
		s.retry(req, retryCount)
	}

	return retryCount
}

// TODO: send logs to storage
func (s *APNAdapter) saveLogs(status, token string, req *domain.PushNotification, err error) error {
	log.Error(map[string]interface{}{
		"ID":       req.ID,
		"Status":   status,
		"Token":    token,
		"Message":  req.Message,
		"Platform": req.Platform,
		"Error":    err,
	})
	return nil
}
