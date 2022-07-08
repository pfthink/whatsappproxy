package routers

import (
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/services"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type SendController struct {
	Service services.SendService
}

func NewSendController(service services.SendService) SendController {
	return SendController{Service: service}
}

func (controller *SendController) Route(app *fiber.App) {
	app.Post("/whatsappproxy/send/text", controller.SendText)
}

func (controller *SendController) SendText(c *fiber.Ctx) error {
	var request structs.SendMessageRequest
	err := c.BodyParser(&request)
	utils.PanicIfNeeded(err)

	// validations.ValidateSendMessage(request)

	response, err := controller.Service.SendText(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        response,
	})
}
