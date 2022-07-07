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
	UserMyListGroups(c *fiber.Ctx) (response structs.UserMyListGroupsResponse, err error)
	UserMyPrivacySetting(c *fiber.Ctx) (response structs.UserMyPrivacySettingResponse, err error)
}
