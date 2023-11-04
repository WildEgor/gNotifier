package domain

import (
	"errors"
	"strconv"
)

type SMSNotification struct {
	Phone   string `json:"phone,omitempty"`
	Message string `json:"message,omitempty"`
}

func ValidateSMSNotification(d *SMSNotification) error {
	var msg string

	if d.Phone == "" {
		msg = "[SMSNotification] Phone number must defined"
		return errors.New(msg)
	}

	if len(d.Phone) != 11 {
		msg = "[SMSNotification] Phone number must 11 digits"
		return errors.New(msg)
	}

	if _, err := strconv.ParseInt(d.Phone, 10, 64); err != nil {
		msg = "[SMSNotification] Parse error"
		return errors.New(msg)
	}

	return nil
}
