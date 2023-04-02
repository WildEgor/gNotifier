package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	"github.com/wagslane/go-rabbitmq"
)

type NotifierHandler struct {
}

func NewNotifierHandler() *NotifierHandler {
	return &NotifierHandler{}
}

func (h *NotifierHandler) Handle(d rabbitmq.Delivery) rabbitmq.Action {
	fmt.Printf("[NotifierHandler] consumed: %v\n", string(d.Body))

	notifierRequest := h.parseReq(d.Body)
	if notifierRequest.HasError() {
		return h.tryResend(notifierRequest)
	}

	// TODO: impl logic here
	fmt.Print(notifierRequest)

	return rabbitmq.Ack
}

func (h *NotifierHandler) parseReq(b []byte) *notifier_dtos.NotifierReqDto {
	req := notifier_dtos.NotifierReqDto{
		TimeReqStart: time.Now(),
	}
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}

	return &req
}

func (h *NotifierHandler) tryResend(req *notifier_dtos.NotifierReqDto) rabbitmq.Action {
	fmt.Printf("[NotifierHandler] Error: %v\n", req.Error)
	// reqRes := notifier_dtos.NotifierResendRequestDto{
	// 	Req:     *parsedRequest,
	// 	Error:   parsedRequest.Error.Error(),
	// 	TimeReq: time.Now().Sub(parsedRequest.TimeReqStart).String(),
	// }
	// TODO: resend to error queue
	time.Sleep(time.Millisecond * 18)

	fmt.Println("[NotifierHandler] execute task: ", time.Now().Sub(req.TimeReqStart).String())

	return rabbitmq.Ack
}
