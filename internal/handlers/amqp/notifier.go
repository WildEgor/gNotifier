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
	notifierRequest := h.initRequest(d.Body)
	if notifierRequest.HasError() {
		return h.resend(notifierRequest)
	}

	// TODO: impl logic here
	fmt.Print(notifierRequest)

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

func (h *NotifierHandler) resend(req *notifier_dtos.NotifierRequestDto) rabbitmq.Action {
	fmt.Printf("[NotifierHandler] Error: %v", req.Error)
	// reqRes := notifier_dtos.NotifierResendRequestDto{
	// 	Req:     *parsedRequest,
	// 	Error:   parsedRequest.Error.Error(),
	// 	TimeReq: time.Now().Sub(parsedRequest.TimeReqStart).String(),
	// }
	// TODO: resend to error queue
	time.Sleep(time.Millisecond * 18)
	fmt.Printf("[NotifierHandler] execute task: ", time.Now().Sub(req.TimeReqStart).String())
	return rabbitmq.Ack
}
