package handlers

import (
	"encoding/json"
	"time"

	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

type NotifierHandler struct {
}

func NewNotifierHandler() *NotifierHandler {
	return &NotifierHandler{}
}

func (h *NotifierHandler) Handle(data rabbitmq.Delivery) rabbitmq.Action {
	req := h.initRequest(data.Body)
	if req.HasError() {
		h.resend(req)
	}

	// TODO: impl logic here

	return rabbitmq.Ack
}

func (h *NotifierHandler) initRequest(b []byte) *notifier_dtos.NotifierRequestDto {
	var req notifier_dtos.NotifierRequestDto
	req.TimeReqStart = time.Now()
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}
	return &req
}

func (h *NotifierHandler) resend(req *notifier_dtos.NotifierRequestDto) {
	log.Error("[NotifierHandler] Error: ", req.Error)
	// reqRes := notifier_dtos.NotifierResendRequestDto{
	// 	Req:     *parsedRequest,
	// 	Error:   parsedRequest.Error.Error(),
	// 	TimeReq: time.Now().Sub(parsedRequest.TimeReqStart).String(),
	// }
	// TODO: resend to error queue
	time.Sleep(time.Millisecond * 18)
	log.Println("execute task: ", time.Now().Sub(req.TimeReqStart).String())
}
