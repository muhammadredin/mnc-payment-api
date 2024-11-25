package dto

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CustomerId   string `json:"customer_id"`
}
