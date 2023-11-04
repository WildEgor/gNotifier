package handlers

import (
	"encoding/json"
	"github.com/WildEgor/gNotifier/internal/repository/mongo"
	log "github.com/sirupsen/logrus"
	"time"

	dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	"github.com/gofiber/fiber/v2"
)

type UnsubTokenHandler struct {
	tokensRepo mongo.ITokensRepository
}

func NewUnsubTokenHandler(
	tokensRepo mongo.ITokensRepository,
) *UnsubTokenHandler {
	return &UnsubTokenHandler{
		tokensRepo,
	}
}

// Handle Unsubscribe token from subscriberID
func (h *UnsubTokenHandler) Handle(ctx *fiber.Ctx) error {
	log.Debug("[UnsubTokenHandler] consumed: %v\n", string(ctx.Body()))

	req := h.parseReq(ctx.Body())
	if req.HasError() {
		log.Error("[UnsubTokenHandler] error: ", req.Error.Error())
		ctx.Status(400).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "Validation error",
			},
		})
	}

	// TODO: delete token from sub

	return nil
}

func (h *UnsubTokenHandler) parseReq(b []byte) *dtos.UnsubTokenReqDto {
	req := dtos.UnsubTokenReqDto{
		TimeReqStart: time.Now(),
	}
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}

	return &req
}
