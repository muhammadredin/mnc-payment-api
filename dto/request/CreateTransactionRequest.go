package dto

type CreateTransactionRequest struct {
	FromWalletId string  `json:"from_wallet_id"`
	ToWalletId   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
	Message      string  `json:"message"`
}
