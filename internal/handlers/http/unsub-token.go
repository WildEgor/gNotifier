package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	"github.com/gofiber/fiber/v2"
)

type UnsubTokenHandler struct {
}

func NewUnsubTokenHandler() *UnsubTokenHandler {
	return &UnsubTokenHandler{}
}

// Unsubscribe token from subscriberID
func (h *UnsubTokenHandler) Handle(ctx *fiber.Ctx) error {
	fmt.Printf("[UnsubTokenHandler] consumed: %v\n", string(ctx.Body()))

	req := h.parseReq(ctx.Body())
	if req.HasError() {
		//
	}

	// TODO: delete token from sub

	return nil
}

func (h *UnsubTokenHandler) parseReq(b []byte) *notifier_dtos.UnsubTokenReqDto {
	req := notifier_dtos.UnsubTokenReqDto{
		TimeReqStart: time.Now(),
	}
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}

	return &req
}
