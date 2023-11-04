package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"

	dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	"github.com/WildEgor/gNotifier/internal/models"
	"github.com/WildEgor/gNotifier/internal/repository/mongo"
	"github.com/gofiber/fiber/v2"
)

type StoreTokenHandler struct {
	tokensRepo mongo.ITokensRepository
}

func NewStoreTokenHandler(
	tokensRepo mongo.ITokensRepository,
) *StoreTokenHandler {
	return &StoreTokenHandler{
		tokensRepo: tokensRepo,
	}
}

// Handle Store any ANDROID or IOS tokens (array of objects) with SubscriberID (could be unique userID for example)
func (h *StoreTokenHandler) Handle(ctx *fiber.Ctx) error {
	log.Debug("[StoreTokenHandler] consumed: %v\n", string(ctx.Body()))

	req := h.parseReq(ctx.Body())
	if req.HasError() {
		log.Error("[StoreTokenHandler] error: ", req.Error.Error())
		ctx.Status(400).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "Validation error",
			},
		})
	}

	_, err := h.tokensRepo.UpsertToken(&models.SubTokenCreateModel{
		SubID: req.SubscriberID,
		Token: &models.TokenModel{
			Token:    req.Token,
			Platform: req.Platform,
		},
	})

	if err != nil {
		log.Error("[StoreTokenHandler] error: ", err.Error())
		ctx.Status(400).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "Cannot save token",
			},
		})
	}

	return nil
}

func (h *StoreTokenHandler) parseReq(b []byte) *dtos.StoreTokenReqDto {
	req := dtos.StoreTokenReqDto{
		TimeReqStart: time.Now(),
	}
	if err := json.Unmarshal(b, &req); err != nil {
		req.Error = err
	}

	return &req
}
