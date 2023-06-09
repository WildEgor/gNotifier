package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/WildEgor/gNotifier/internal/adapters"
	"github.com/WildEgor/gNotifier/internal/domain"
	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

type NotifierHandler struct {
	smtpAdapter adapters.ISMTPAdapter
	smsAdapter  adapters.ISMSAdapter
	fcmAdapter  adapters.IFCMAdapter
	apnAdapter  adapters.IAPNAdapter
}

func NewNotifierHandler(
	smtpAdapter adapters.ISMTPAdapter,
	smsAdapter adapters.ISMSAdapter,
	fcmAdapter adapters.IFCMAdapter,
	apnAdapter adapters.IAPNAdapter,
) *NotifierHandler {
	return &NotifierHandler{
		smtpAdapter: smtpAdapter,
		smsAdapter:  smsAdapter,
		fcmAdapter:  fcmAdapter,
		apnAdapter:  apnAdapter,
	}
}

func (h *NotifierHandler) Handle(d rabbitmq.Delivery) rabbitmq.Action {
	notifierRequest := h.parseReq(d.Body)
	if notifierRequest.HasError() {
		return h.tryResend(notifierRequest)
	}

	fmt.Printf("[NotifierHandler] consumed: %v\n", notifierRequest)

	if notifierRequest.IsEmail() {
		notification := domain.EmailNotification{
			Email:   notifierRequest.EmailSetting.Email,
			Message: notifierRequest.EmailSetting.Text,
		}

		if notifierRequest.WithTemplate() {
			msg, err := h.parseTemplate(notifierRequest)
			if err != nil {
				//
			}
			notification.Message = msg
		}

		if err := h.smtpAdapter.Send(&notification); err != nil {
			// TODO
			log.Errorf("[] Failed send to: ", notifierRequest.EmailSetting.Email)
		}
	}

	if notifierRequest.IsSms() {
		notification := domain.SMSNotification{
			Phone:   notifierRequest.PhoneSetting.Number,
			Message: notifierRequest.PhoneSetting.Text,
			TeamID:  "",
		}

		if err := h.smsAdapter.Send(&notification); err != nil {
			// TODO
		}
	}

	if notifierRequest.IsPush() {

		// TODO: find tokens by To and Platform as sub_id in mongodb and use only tokens updated_at or created_at >= now() - 30 days!
		notification := domain.PushNotification{
			ID:      "", // TODO
			To:      "",
			Tokens:  []string{""},
			Topic:   notifierRequest.PushSetting.To,
			Message: notifierRequest.PushSetting.Message,
			Title:   notifierRequest.PushSetting.Title,
		}

		if notifierRequest.WithTemplate() {
			msg, err := h.parseTemplate(notifierRequest)
			if err != nil {
				//
			}
			notification.Message = msg
		}

		if notifierRequest.IsForAndroid() {
			notification.Platform = domain.PlatFormAndroid

			if err := h.fcmAdapter.Send(&notification); err != nil {
				// TODO
			}
		}

		if notifierRequest.IsForIOS() {
			notification.Platform = domain.PlatFormIos
			if err := h.apnAdapter.Send(&notification); err != nil {
				// TODO
			}
		}
	}

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
	reqRes := notifier_dtos.NotifierResendReqDto{
		Req:     *req,
		Error:   req.Error.Error(),
		TimeReq: time.Now().Sub(req.TimeReqStart).String(),
	}
	// TODO: resend to error queue

	time.Sleep(time.Millisecond * 18)
	fmt.Println("[NotifierHandler] execute task: ", time.Now().Sub(req.TimeReqStart).String(), reqRes)
	return rabbitmq.Ack
}

func (h *NotifierHandler) parseTemplate(req *notifier_dtos.NotifierReqDto) (msg string, err error) {
	tml, err := template.ParseFiles(req.EmailSetting.Template)
	if err != nil {
		return "", errors.New("[NotifierHandler] Cannot parse template")
	}

	buf := new(bytes.Buffer)
	if err = tml.Execute(buf, req.Data); err != nil {
		return "", errors.New("[NotifierHandler] Cannot parse template")
	}

	return buf.String(), nil
}
