package structs

type QrcodeScanStatus struct {
	NoiseKeyPub    string `json:"noiseKeyPub" query:"noiseKeyPub"`
	IdentityKeyPub string `json:"identityKeyPub" query:"identityKeyPub"`
	AdvSecret      string `json:"advSecret" query:"advSecret"`
}
