package adapters

import (
	"errors"
	"fmt"
	"github.com/WildEgor/gNotifier/internal/configs"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/WildEgor/gNotifier/internal/domain"
	"github.com/appleboy/go-fcm"
)

func MapToAndroidNotification(req *domain.PushNotification) *fcm.Message {
	notification := &fcm.Message{
		To:                    req.To,
		Condition:             req.Condition,
		CollapseKey:           req.CollapseKey,
		ContentAvailable:      req.ContentAvailable,
		MutableContent:        req.MutableContent,
		DelayWhileIdle:        req.DelayWhileIdle,
		TimeToLive:            req.TimeToLive,
		RestrictedPackageName: req.RestrictedPackageName,
		DryRun:                req.DryRun,
	}

	if len(req.Tokens) > 0 {
		notification.RegistrationIDs = req.Tokens
	}

	if req.Priority == "high" || req.Priority == "normal" {
		notification.Priority = req.Priority
	}

	// Add another field
	if len(req.Data) > 0 {
		notification.Data = make(map[string]interface{})
		for k, v := range req.Data {
			notification.Data[k] = v
		}
	}

	n := &fcm.Notification{}
	isNotificationSet := false
	if req.Notification != nil {
		isNotificationSet = true
		n = req.Notification
	}

	if len(req.Message) > 0 {
		isNotificationSet = true
		n.Body = req.Message
	}

	if len(req.Title) > 0 {
		isNotificationSet = true
		n.Title = req.Title
	}

	if len(req.Image) > 0 {
		isNotificationSet = true
		n.Image = req.Image
	}

	if v, ok := req.Sound.(string); ok && len(v) > 0 {
		isNotificationSet = true
		n.Sound = v
	}

	if isNotificationSet {
		notification.Notification = n
	}

	// handle iOS apns in fcm

	if len(req.Apns) > 0 {
		notification.Apns = req.Apns
	}

	return notification
}

type IFCMAdapter interface {
	Send(req *domain.PushNotification) error
}

type FCMAdapter struct {
	client *fcm.Client
	config *configs.FCMConfig
}

// NewFCMAdapter Create new FCM Client
func NewFCMAdapter(
	config *configs.FCMConfig,
) *FCMAdapter {

	FCMClient, err := fcm.NewClient(config.APIKey, fcm.WithTimeout(time.Second*5))
	if err != nil && config.Production {
		log.Fatalf("[FCMAdapter] Cannot init FCM client: %v", err)
	}

	return &FCMAdapter{
		client: FCMClient,
		config: config,
	}
}

// Send provide send notification to Android server.
func (f *FCMAdapter) Send(push *domain.PushNotification) (err error) {
	var retryCount = 0

	// Validate notification data
	err = domain.ValidatePushNotification(push)
	if err != nil {
		log.Println("[FCMAdapter] Not valid push notification: " + err.Error())
		return
	}

	f.retry(push, retryCount)

	return err
}

func (f *FCMAdapter) retry(push *domain.PushNotification, retryCount int) error {
	maxRetry := f.config.MaxRetry
	notification := MapToAndroidNotification(push)

	// Limit retry count
	if push.Retry > 0 && push.Retry < maxRetry {
		maxRetry = push.Retry
	}

	res, err := f.client.Send(notification)
	if err != nil {
		// Send Message error
		log.Println("[FCMAdapter] FCM server send message error: " + err.Error())

		// Save logs depends on topic or tokens provided
		if push.IsTopic() {
			f.saveLogs("fail_push", push.To, push, err)
		} else {
			for _, token := range push.Tokens {
				f.saveLogs("fail_push", token, push, err)
			}
		}
		return err
	}

	if !push.IsTopic() {
		log.Debugln(fmt.Sprintf("Android Success count: %d, Failure count: %d", res.Success, res.Failure))
	}

	var newTokens []string
	// result from Send messages to specific devices
	for k, result := range res.Results {
		to := ""
		if k < len(push.Tokens) {
			to = push.Tokens[k]
		} else {
			to = push.To
		}

		if result.Error != nil {
			// We should retry only "retryable" statuses. More info about response:
			// https://firebase.google.com/docs/cloud-messaging/http-server-ref#downstream-http-messages-plain-text
			if !result.Unregistered() {
				newTokens = append(newTokens, to)
			}

			f.saveLogs("fail_push", to, push, result.Error)
			continue
		}

		f.saveLogs("success_push", to, push, nil)
	}

	// result from Send messages to topics
	if push.IsTopic() {
		to := ""
		if push.To != "" {
			to = push.To
		} else {
			to = push.Condition
		}
		log.Println("Send Topic Message: ", to)
		// Success
		if res.MessageID != 0 {
			f.saveLogs("success_push", to, push, nil)
		} else {
			// failure
			f.saveLogs("fail_push", to, push, res.Error)
		}
	}

	// Device Group HTTP Response
	if len(res.FailedRegistrationIDs) > 0 {
		newTokens = append(newTokens, res.FailedRegistrationIDs...)
		f.saveLogs("fail_push", notification.To, push, errors.New("device group: partial success or all fails"))
	}

	if len(newTokens) > 0 && retryCount < maxRetry {
		retryCount++
		// resend fail token
		push.Tokens = newTokens
		f.retry(push, retryCount)
	}

	return nil
}

// TODO: save logs to storage
func (f *FCMAdapter) saveLogs(status, token string, req *domain.PushNotification, err error) error {
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
