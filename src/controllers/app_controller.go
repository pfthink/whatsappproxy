package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/services"
	"whatsappproxy/utils"
)

type AppController struct {
	Service services.AppService
}

func NewAppController(service services.AppService) AppController {
	return AppController{Service: service}
}

func (controller *AppController) Route(app *fiber.App) {
	app.Get("/app/login", controller.Login)
	app.Get("/app/logout", controller.Logout)
	app.Get("/app/reconnect", controller.Reconnect)
}

func (controller *AppController) Login(c *fiber.Ctx) error {
	response, err := controller.Service.Login(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"qr_link":        fmt.Sprintf("%s://%s/%s", c.Protocol(), c.Hostname(), response.ImagePath),
			"qrcode":         response.Code,
			"qrcodeDuration": response.Duration,
		},
	})
}

func (controller *AppController) Logout(c *fiber.Ctx) error {
	err := controller.Service.Logout(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:    200,
		Message: "success",
		Data:    nil,
	})
}

func (controller *AppController) Reconnect(c *fiber.Ctx) error {
	err := controller.Service.Reconnect(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:    200,
		Message: "success",
		Data:    nil,
	})
}
