package entity

type Blacklist struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
}
