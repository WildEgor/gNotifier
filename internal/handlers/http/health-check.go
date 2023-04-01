package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (s *HealthCheckHandler) Handle(c *fiber.Ctx) error {
	return nil
}
