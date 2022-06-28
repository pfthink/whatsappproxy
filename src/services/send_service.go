package services

import (
	"github.com/pfthink/whatsappproxy/src/structs"
)

type SendService interface {
	SendText(c *fiber.Ctx, request structs.SendMessageRequest) (response structs.SendMessageResponse, err error)
	SendImage(c *fiber.Ctx, request structs.SendImageRequest) (response structs.SendImageResponse, err error)
	SendFile(c *fiber.Ctx, request structs.SendFileRequest) (response structs.SendFileResponse, err error)
	SendVideo(c *fiber.Ctx, request structs.SendVideoRequest) (response structs.SendVideoResponse, err error)
}
