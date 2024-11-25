package dto

type CustomerResponse struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	WalletId string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}
