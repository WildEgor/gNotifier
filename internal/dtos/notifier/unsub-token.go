package dtos

import "time"

type UnsubTokenReqDto struct {
	SubscriberID string    `json:"sub_id"`
	Token        string    `json:"token"`
	Error        error     `json:"-"`
	TimeReqStart time.Time `json:"-"`
}

func (r *UnsubTokenReqDto) HasError() bool {
	return r.Error != nil
}
