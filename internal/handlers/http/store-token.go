package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	notifier_dtos "github.com/WildEgor/gNotifier/internal/dtos/notifier"
	models "github.com/WildEgor/gNotifier/internal/models"
	mongo "github.com/WildEgor/gNotifier/internal/repository/mongo"
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

// Store any ANDROID or IOS tokens (array of objects) with SubscriberID (could be unique userID for example)
func (h *StoreTokenHandler) Handle(ctx *fiber.Ctx) error {
	fmt.Printf("[StoreTokenHandler] consumed: %v\n", string(ctx.Body()))

	req := h.parseReq(ctx.Body())
	if req.HasError() {
		//
	}

	_, err := h.tokensRepo.UpsertToken(&models.SubTokenCreateModel{
		SubID: req.SubscriberID,
		Token: &models.TokenModel{
			Token:    req.Token,
			Platform: req.Platform,
		},
	})
	if err != nil {
		ctx.Status(400).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "Cannot save token",
			},
		})
	}

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
