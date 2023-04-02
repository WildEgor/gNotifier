package domain

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/appleboy/go-fcm"
)

const (
	PlatFormIos = iota + 1
	PlatFormAndroid
)

type AnyData map[string]interface{}

// Alert is APNs payload
type Alert struct {
	Action          string   `json:"action,omitempty"`
	ActionLocKey    string   `json:"action-loc-key,omitempty"`
	Body            string   `json:"body,omitempty"`
	LaunchImage     string   `json:"launch-image,omitempty"`
	LocArgs         []string `json:"loc-args,omitempty"`
	LocKey          string   `json:"loc-key,omitempty"`
	Title           string   `json:"title,omitempty"`
	Subtitle        string   `json:"subtitle,omitempty"`
	TitleLocArgs    []string `json:"title-loc-args,omitempty"`
	TitleLocKey     string   `json:"title-loc-key,omitempty"`
	SummaryArg      string   `json:"summary-arg,omitempty"`
	SummaryArgCount int      `json:"summary-arg-count,omitempty"`
}

type PushNotification struct {
	// Common
	ID               string      `json:"notif_id,omitempty"`
	Tokens           []string    `json:"tokens" binding:"required"`
	Platform         int         `json:"platform" binding:"required"`
	Message          string      `json:"message,omitempty"`
	Title            string      `json:"title,omitempty"`
	Image            string      `json:"image,omitempty"`
	Priority         string      `json:"priority,omitempty"`
	ContentAvailable bool        `json:"content_available,omitempty"`
	MutableContent   bool        `json:"mutable_content,omitempty"`
	Sound            interface{} `json:"sound,omitempty"`
	Data             AnyData     `json:"data,omitempty"`
	Retry            int         `json:"retry,omitempty"`

	// Android
	To                    string            `json:"to,omitempty"`
	CollapseKey           string            `json:"collapse_key,omitempty"`
	DelayWhileIdle        bool              `json:"delay_while_idle,omitempty"`
	TimeToLive            *uint             `json:"time_to_live,omitempty"`
	RestrictedPackageName string            `json:"restricted_package_name,omitempty"`
	DryRun                bool              `json:"dry_run,omitempty"`
	Condition             string            `json:"condition,omitempty"`
	Notification          *fcm.Notification `json:"notification,omitempty"`

	// iOS
	Expiration  *int64   `json:"expiration,omitempty"`
	ApnsID      string   `json:"apns_id,omitempty"`
	CollapseID  string   `json:"collapse_id,omitempty"`
	Topic       string   `json:"topic,omitempty"`
	PushType    string   `json:"push_type,omitempty"`
	Badge       *int     `json:"badge,omitempty"`
	Category    string   `json:"category,omitempty"`
	ThreadID    string   `json:"thread-id,omitempty"`
	URLArgs     []string `json:"url-args,omitempty"`
	Alert       Alert    `json:"alert,omitempty"`
	Production  bool     `json:"production,omitempty"`
	Development bool     `json:"development,omitempty"`
	SoundName   string   `json:"name,omitempty"`
	SoundVolume float32  `json:"volume,omitempty"`
	Apns        AnyData  `json:"apns,omitempty"`

	// ref: https://github.com/sideshow/apns2/blob/54928d6193dfe300b6b88dad72b7e2ae138d4f0a/payload/builder.go#L7-L24
	InterruptionLevel string `json:"interruption_level,omitempty"`
}

// Bytes for queue message
func (p *PushNotification) ToBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return b
}

// IsTopic check if message format is topic for FCM
// ref: https://firebase.google.com/docs/cloud-messaging/send-message#topic-http-post-request
func (p *PushNotification) IsTopic() bool {
	if p.Platform == PlatFormAndroid {
		return p.To != "" && strings.HasPrefix(p.To, "/topics/") || p.Condition != ""
	}

	return false
}

// Static method
// ValidateNotification for check request message
func ValidateNotification(d *PushNotification) error {
	var msg string

	// ignore send topic mesaage from FCM
	if !d.IsTopic() && len(d.Tokens) == 0 && d.To == "" {
		msg = "[PushNotification] The message must specify at least one registration ID"
		return errors.New(msg)
	}

	if len(d.Tokens) == PlatFormIos && d.Tokens[0] == "" {
		msg = "[PushNotification] The token must not be empty"
		return errors.New(msg)
	}

	if d.Platform == PlatFormAndroid && len(d.Tokens) > 1000 {
		msg = "[PushNotification] The message may specify at most 1000 registration IDs"
		return errors.New(msg)
	}

	// ref: https://firebase.google.com/docs/cloud-messaging/http-server-ref
	if d.Platform == PlatFormAndroid && d.TimeToLive != nil && *d.TimeToLive > uint(2419200) {
		msg = "[PushNotification] The message's TimeToLive field must be an integer " +
			"between 0 and 2419200 (4 weeks)"
		return errors.New(msg)
	}

	return nil
}
