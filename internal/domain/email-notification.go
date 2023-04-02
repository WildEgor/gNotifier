package domain

import (
	"errors"
	"regexp"
)

type EmailNotification struct {
	Email   string `json:"email,omitempty"`
	Subject string `json:"subj,omitempty"`
	Message string `json:"msg,omitempty"`
}

func ValidateEmailNotification(d *EmailNotification) error {
	var msg string

	if d.Email == "" {
		msg = "[EmailNotification] Email must defined"
		return errors.New(msg)
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	submatch := re.FindStringSubmatch(d.Email)
	if len(submatch) < 2 {
		msg = "[SMSNotification] Email incorrect format"
		return errors.New(msg)
	}

	if d.Subject == "" || d.Message == "" {
		msg = "[SMSNotification] Provide Subject and Message"
		return errors.New(msg)
	}

	return nil
}
