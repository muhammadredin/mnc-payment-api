package entity

type RefreshToken struct {
	Id           string `json:"id"`
	CustomerId   string `json:"customer_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	IsValid      bool   `json:"is_valid"`
}
