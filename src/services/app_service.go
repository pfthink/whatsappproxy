package services

import (
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/structs"
)

type AppService interface {
	GetQrcode(c *fiber.Ctx) (response structs.LoginResponse, err error)

	QrcodeScanStatus(c *fiber.Ctx, request structs.QrcodeScanStatus) (response structs.ScanStatusResponse, err error)
}
