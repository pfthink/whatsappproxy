package structs

import (
	"time"
)

type LoginResponse struct {
	ImagePath string        `json:"image_path"`
	Duration  time.Duration `json:"duration"`
	Code      string        `json:"code"`
}

type ScanStatusResponse struct {
	LoginStatus int    `json:"loginStatus"`
	Jid         string `json:"jid"`
}
