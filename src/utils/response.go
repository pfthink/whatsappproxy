package utils

type ResponseData struct {
	Code         int         `json:"code"`
	Succeeded    bool        `json:"succeeded"`
	ResponseCode string      `json:"responseCode"`
	ResponseMsg  string      `json:"responseMsg"`
	Value        interface{} `json:"value"`
}
