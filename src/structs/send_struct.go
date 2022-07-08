package structs

import (
	"mime/multipart"
)

type SendType string

const TypeUser SendType = "user"
const TypeGroup SendType = "group"

type SendMessageRequest struct {
	FromJid string   `json:"fromJid" form:"fromJid"`
	ToJid   string   `json:"toJid" form:"toJid"`
	Message string   `json:"message" form:"message"`
	Type    SendType `json:"type" form:"type"`
}

type SendMessageResponse struct {
	MessageId string `json:"messageId"`
}

type SendImageRequest struct {
	FromJid  string   `json:"fromJid" form:"fromJid"`
	ToJid    string   `json:"toJid" form:"toJid"`
	Caption  string   `json:"caption" form:"caption"`
	ImageUrl string   `json:"imageUrl" form:"imageUrl"`
	ViewOnce bool     `json:"view_once" form:"view_once"`
	Type     SendType `json:"type" form:"type"`
	Compress bool     `json:"compress"`
}

type SendImageResponse struct {
	Status string `json:"status"`
}

type SendFileRequest struct {
	Phone string                `json:"phone" form:"phone"`
	File  *multipart.FileHeader `json:"file" form:"file"`
	Type  SendType              `json:"type" form:"type"`
}

type SendFileResponse struct {
	Status string `json:"status"`
}

type SendVideoRequest struct {
	Phone    string                `json:"phone" form:"phone"`
	Caption  string                `json:"caption" form:"caption"`
	Video    *multipart.FileHeader `json:"video" form:"video"`
	Type     SendType              `json:"type" form:"type"`
	ViewOnce bool                  `json:"view_once" form:"view_once"`
	Compress bool                  `json:"compress"`
}

type SendVideoResponse struct {
	Status string `json:"status"`
}
