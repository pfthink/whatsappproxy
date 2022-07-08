package structs

type UserRequest struct {
	Jid string `json:"jid" query:"jid"`
}

type UserInfoResponseDataDevice struct {
	User   string
	Agent  uint8
	Device string
	Server string
	AD     bool
}

type UserInfoResponseData struct {
	VerifiedName string                       `json:"verified_name"`
	Status       string                       `json:"status"`
	PictureID    string                       `json:"picture_id"`
	Devices      []UserInfoResponseDataDevice `json:"devices"`
}

type UserInfoResponse struct {
	Data []UserInfoResponseData `json:"data"`
}

type UserAvatarRequest struct {
	Jid string `json:"jid" query:"jid"`
}

type UserAvatarResponse struct {
	Jid  string `json:"jid"`
	URL  string `json:"url"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

type UserOnlineStatusRequest struct {
	Jid       string `json:"jid" query:"jid"`
	TargetJid string `json:"targetJid" query:"targetJid"`
}

type UserOnlineStatusResponse struct {
	Jid          string `json:"jid"`
	TargetJid    string `json:"targetJid" query:"targetJid"`
	OnlineStatus int    `json:"onlineStatus"`
}
