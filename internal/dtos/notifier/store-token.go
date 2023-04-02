package dtos

import "time"

type StoreTokenReqDto struct {
	SubscriberID string    `json:"sub_id"`
	Platform     string    `json:"platform"`
	Token        string    `json:"token"`
	Error        error     `json:"-"`
	TimeReqStart time.Time `json:"-"`
}

func (r *StoreTokenReqDto) HasError() bool {
	if r.Error != nil {
		return true
	}
	return false
}
