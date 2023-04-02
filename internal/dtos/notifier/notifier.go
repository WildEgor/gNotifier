package dtos

import (
	"errors"
	"time"
)

type NotifierResendReqDto struct {
	Req     NotifierReqDto `json:"request"`
	Error   string         `json:"error"`
	TimeReq string         `json:"time_req"`
}

type NotifierReqDto struct {
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
	PushSetting struct {
		To       string `json:"to"`
		Platform string `json:"platform"`
		Image    string `json:"image"`
		Template string `json:"template"`
		Title    string `json:"title"`
		Message  string `json:"message"`
	} `json:"push_settings"`
	Data         interface{} `json:"data"`
	Error        error       `json:"-"`
	TimeReqStart time.Time   `json:"-"`
}

func (r *NotifierReqDto) IsSms() bool {
	if r.Type == "sms" {
		return true
	}
	return false
}

func (r *NotifierReqDto) IsEmail() bool {
	if r.Type == "email" {
		return true
	}
	return false
}

func (r *NotifierReqDto) IsPush() bool {
	if r.Type == "push" {
		return true
	}
	return false
}

func (r *NotifierReqDto) IsForAndroid() bool {
	if r.Type == "push" && r.PushSetting.Platform == "ANDROID" {
		return true
	}
	return false
}

func (r *NotifierReqDto) IsForIOS() bool {
	if r.Type == "push" && r.PushSetting.Platform == "IOS" {
		return true
	}
	return false
}

func (r *NotifierReqDto) ValidateType() bool {
	if r.Type != "sms" && r.Type != "email" && r.Type != "push" {
		r.Error = errors.New("[NotifierReqDto] Error type - " + r.Type)
		return false
	}
	return true
}

func (r *NotifierReqDto) ValidateEmail() bool {
	if len(r.EmailSetting.Email) == 0 {
		r.Error = errors.New("[NotifierReqDto] Error pass email param")
		return false
	}
	if len(r.EmailSetting.Subject) == 0 {
		r.Error = errors.New("[NotifierReqDto] Error pass subject param")
		return false
	}
	if len(r.EmailSetting.Template) == 0 {
		r.Error = errors.New("[NotifierReqDto] Error pass email template param")
		return false
	}

	return true
}

func (r *NotifierReqDto) ValidateSms() bool {
	if len(r.PhoneSetting.Number) == 0 || len(r.PhoneSetting.Text) == 0 {
		r.Error = errors.New("[NotifierReqDto] Error pass sms params")
		return false
	}
	return true
}

func (r *NotifierReqDto) ValidatePush() bool {
	if len(r.PushSetting.To) == 0 {
		r.Error = errors.New("[NotifierReqDto] Empty param in PushSetting")
		return false
	}

	return true
}

func (r *NotifierReqDto) HasError() bool {
	if r.Error != nil {
		return true
	}
	return false
}
