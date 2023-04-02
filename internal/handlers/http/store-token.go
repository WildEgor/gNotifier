package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	"github.com/gofiber/fiber/v2"
)

type StoreTokenHandler struct {
}

func NewStoreTokenHandler() *StoreTokenHandler {
	return &StoreTokenHandler{}
}

// Store any ANDROID or IOS tokens (array of objects) with SubscriberID (could be unique userID for example)
func (h *StoreTokenHandler) Handle(ctx *fiber.Ctx) error {
	fmt.Printf("[StoreTokenHandler] consumed: %v\n", string(ctx.Body()))

	req := h.parseReq(ctx.Body())
	if req.HasError() {
		//
	}

	// TODO: implement upsert logic to any storage (MongoDB, Radis, Postgres ...)
	// Add new token to sub or update date if consume same token for sub

	return nil
}

func (h *StoreTokenHandler) parseReq(b []byte) *notifier_dtos.StoreTokenReqDto {
	req := notifier_dtos.StoreTokenReqDto{
		TimeReqStart: time.Now(),
	}
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}

	return &req
}
