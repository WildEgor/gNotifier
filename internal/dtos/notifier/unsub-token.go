package dtos

import "time"

type UnsubTokenReqDto struct {
	SubscriberID string    `json:"sub_id"`
	Token        string    `json:"token"`
	Error        error     `json:"-"`
	TimeReqStart time.Time `json:"-"`
}

func (r *UnsubTokenReqDto) HasError() bool {
	if r.Error != nil {
		return true
	}
	return false
}
