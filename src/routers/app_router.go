package routers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/services"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type AppController struct {
	Service services.AppService
}

func NewAppController(service services.AppService) AppController {
	return AppController{Service: service}
}

func (controller *AppController) Route(app *fiber.App) {
	app.Post("/whatsappproxy/app/qrcode", controller.GetQrcode)

	app.Post("/whatsappproxy/app/qrcode/scan/status", controller.GetQrcodeScanStatus)
}

func (controller *AppController) GetQrcode(c *fiber.Ctx) error {
	response, err := controller.Service.GetQrcode(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value: map[string]interface{}{
			"qr_link":        fmt.Sprintf("%s://%s/%s", c.Protocol(), c.Hostname(), response.ImagePath),
			"qrcode":         response.Code,
			"qrcodeDuration": response.Duration,
		},
	})
}

func (controller *AppController) GetQrcodeScanStatus(c *fiber.Ctx) error {

	var request structs.QrcodeScanStatus
	c.BodyParser(&request)

	response, err := controller.Service.QrcodeScanStatus(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value: map[string]interface{}{
			"loginStatus": response.LoginStatus,
			"jid":         response.Jid,
		},
	})
}
