package services

import (
	"github.com/gofiber/fiber/v2"
	"whatsappproxy/structs"
)

type UserService interface {
	Logout(c *fiber.Ctx, request structs.UserRequest) (err error)
	Reconnect(c *fiber.Ctx, request structs.UserRequest) (err error)
	UserInfo(c *fiber.Ctx, request structs.UserRequest) (response structs.UserInfoResponse, err error)
	UserAvatar(c *fiber.Ctx, request structs.UserAvatarRequest) (response structs.UserAvatarResponse, err error)
	OnlineStatus(c *fiber.Ctx, request structs.UserOnlineStatusRequest) (response structs.UserOnlineStatusResponse, err error)
}
