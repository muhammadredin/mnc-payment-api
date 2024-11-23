package entity

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	CustomerId   string `json:"customer_id"`
	ExpiresAt    string `json:"expires_at"`
}
