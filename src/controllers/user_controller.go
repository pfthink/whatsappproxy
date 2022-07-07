package controllers

import (
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/services"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
	"whatsappproxy/validations"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) UserController {
	return UserController{Service: service}
}

func (controller *UserController) Route(app *fiber.App) {
	app.Post("/whatsappproxy/user/logout", controller.Logout)
	app.Post("/whatsappproxy/user/reconnect", controller.Reconnect)

	app.Post("/whatsappproxy/user/info", controller.UserInfo)
	app.Post("/whatsappproxy/user/avatar", controller.UserAvatar)
	app.Post("/whatsappproxy/user/my/privacy", controller.UserMyPrivacySetting)
	app.Post("/whatsappproxy/user/my/groups", controller.UserMyListGroups)
}

func (controller *UserController) Logout(c *fiber.Ctx) error {
	var request structs.UserRequest
	c.BodyParser(&request)
	err := controller.Service.Logout(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        nil,
	})
}

func (controller *UserController) Reconnect(c *fiber.Ctx) error {
	var request structs.UserRequest
	c.BodyParser(&request)
	err := controller.Service.Reconnect(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        nil,
	})
}

func (controller *UserController) UserInfo(c *fiber.Ctx) error {
	var request structs.UserRequest
	err := c.QueryParser(&request)
	utils.PanicIfNeeded(err)

	// add validation send message
	//validations.ValidateUserInfo(request)

	//request.Phone = request.Phone + "@s.whatsapp.net"
	response, err := controller.Service.UserInfo(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        response.Data[0],
	})
}

func (controller *UserController) UserAvatar(c *fiber.Ctx) error {
	var request structs.UserAvatarRequest
	err := c.QueryParser(&request)
	utils.PanicIfNeeded(err)

	// add validation send message
	validations.ValidateUserAvatar(request)

	request.Phone = request.Phone + "@s.whatsapp.net"
	response, err := controller.Service.UserAvatar(c, request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        response,
	})
}

func (controller *UserController) UserMyPrivacySetting(c *fiber.Ctx) error {
	response, err := controller.Service.UserMyPrivacySetting(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        response,
	})
}

func (controller *UserController) UserMyListGroups(c *fiber.Ctx) error {
	response, err := controller.Service.UserMyListGroups(c)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Code:         200,
		Succeeded:    true,
		ResponseCode: "SUCCESS",
		Value:        response,
	})
}
