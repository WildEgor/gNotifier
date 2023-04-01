package dtos

import (
	"errors"
	"time"
)

type NotifierResendRequestDto struct {
	Req     NotifierRequestDto `json:"request"`
	Error   string             `json:"error"`
	TimeReq string             `json:"time_req"`
}

type NotifierRequestDto struct {
	Type         string `json:"type"`
	EmailSetting struct {
		Email    string `json:"email"`
		Subject  string `json:"subject"`
		Template string `json:"template"`
		Text     string `json:"text"`
	} `json:"email_setting"`
	PhoneSetting struct {
		Number string `json:"number"`
		Text   string `json:"text"`
	} `json:"phone_setting"`
	Data         interface{} `json:"data"`
	Error        error       `json:"-"`
	TimeReqStart time.Time   `json:"-"`
}

func (r *NotifierRequestDto) IsSms() bool {
	if r.Type == "sms" {
		return true
	}
	return false
}

func (r *NotifierRequestDto) IsEmail() bool {
	if r.Type == "email" {
		return true
	}
	return false
}

func (r *NotifierRequestDto) ValidateType() bool {
	if r.Type != "sms" && r.Type != "email" {
		r.Error = errors.New("[NotifierRequestDto] Error type - " + r.Type)
		return false
	}
	return true
}

func (r *NotifierRequestDto) ValidateEmail() bool {
	if len(r.EmailSetting.Email) == 0 {
		r.Error = errors.New("[NotifierRequestDto] Error pass email param")
		return false
	}
	if len(r.EmailSetting.Subject) == 0 {
		r.Error = errors.New("[NotifierRequestDto] Error pass subject param")
		return false
	}
	if len(r.EmailSetting.Template) == 0 {
		r.Error = errors.New("[NotifierRequestDto] Error pass email template param")
		return false
	}

	return true
}

func (r *NotifierRequestDto) ValidateSms() bool {
	if len(r.PhoneSetting.Number) == 0 || len(r.PhoneSetting.Text) == 0 {
		r.Error = errors.New("[NotifierRequestDto] Error pass sms params")
		return false
	}
	return true
}

func (r *NotifierRequestDto) HasError() bool {
	if r.Error != nil {
		return true
	}
	return false
}
