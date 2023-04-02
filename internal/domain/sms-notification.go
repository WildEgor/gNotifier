package domain

import (
	"errors"
	"regexp"
)

type SMSNotification struct {
	TeamID  string `json:"team_id"`
	Phone   string `json:"phone,omitempty"`
	Message string `json:"message,omitempty"`
}

func ValidateSMSNotification(d *SMSNotification) error {
	var msg string

	if d.Phone == "" {
		msg = "[SMSNotification] Phone number must defined"
		return errors.New(msg)
	}

	re := regexp.MustCompile(`(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])`)
	submatch := re.FindStringSubmatch(d.Phone)
	if len(submatch) < 2 {
		msg = "[SMSNotification] Phone number incorrect format"
		return errors.New(msg)
	}

	return nil
}
