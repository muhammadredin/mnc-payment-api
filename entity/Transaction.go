package entity

import "PaymentAPI/enums"

type Transaction struct {
	Id           string                  `json:"id"`
	FromWalletId string                  `json:"from_wallet_id"`
	ToWalletId   string                  `json:"to_wallet_id"`
	Status       enums.TransactionStatus `json:"status"`
	CreatedAt    string                  `json:"created_at"`
	Message      string                  `json:"message"`
}
