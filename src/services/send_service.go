package services

import (
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/structs"
)

type SendService interface {
	SendText(c *fiber.Ctx, request structs.SendMessageRequest) (response structs.SendMessageResponse, err error)
}
