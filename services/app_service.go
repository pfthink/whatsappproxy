package services

import (
	"github.com/pfthink/whatsappproxy//structs"
)

type AppService interface {
	Login(c *fiber.Ctx) (response structs.LoginResponse, err error)
	Logout(c *fiber.Ctx) (err error)
	Reconnect(c *fiber.Ctx) (err error)
}
